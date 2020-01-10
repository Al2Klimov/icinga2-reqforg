package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"github.com/Al2Klimov/icinga2-reqforg/lib/base"
	"github.com/Al2Klimov/icinga2-reqforg/sdk"
	log "github.com/sirupsen/logrus"
	"sync"
)

type connInfo struct {
	sync.Mutex

	closed     bool
	closeError error
}

var conns = map[*tls.Conn]*connInfo{}
var connsLock sync.RWMutex

type sdkActions struct {
	plugin *myPlugin
}

var _ sdk.SdkActions = sdkActions{}

func (sa sdkActions) Close() error {
	pluginsLock.Lock()
	delete(plugins, sa.plugin)
	pluginsLock.Unlock()

	return nil
}

func (sa sdkActions) AddConn(conn *tls.Conn) {
	ci := &connInfo{}

	connsLock.Lock()
	conns[conn] = ci
	connsLock.Unlock()

	var wg sync.WaitGroup

	pluginsLock.RLock()

	for plugin := range plugins {
		wg.Add(1)
		go onNewConn(plugin, conn, ci, &wg)
	}

	pluginsLock.RUnlock()

	wg.Wait()

	go readLoop(conn, ci)
}

func (sa sdkActions) GetLogger() log.Ext1FieldLogger {
	// TODO: evaluate alternative
	return log.StandardLogger()
}

func onNewConn(plugin *myPlugin, conn *tls.Conn, ci *connInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	plugin.pluginActions.OnNewConn(connActions{plugin, conn, ci})
}

type connActions struct {
	plugin *myPlugin
	conn   *tls.Conn
	ci     *connInfo
}

var _ sdk.ConnActions = connActions{}

func (ca connActions) Close() error {
	ca.ci.Lock()
	defer ca.ci.Unlock()

	if !ca.ci.closed {
		ca.ci.closeError = ca.conn.Close()
		ca.ci.closed = true

		connsLock.Lock()
		delete(conns, ca.conn)
		connsLock.Unlock()
	}

	return ca.ci.closeError
}

func (ca connActions) GetConn() *tls.Conn {
	return ca.conn
}

func (ca connActions) OnClose(eventHandler func(error)) {
	ca.plugin.Lock()
	defer ca.plugin.Unlock()

	eventHandlers, ok := ca.plugin.connEventHandlers[ca.conn]
	if !ok {
		if eventHandler == nil {
			return
		}

		eventHandlers = &connEventHandlers{}
		ca.plugin.connEventHandlers[ca.conn] = eventHandlers
	}

	eventHandlers.onClose = eventHandler

	if eventHandlers.onClose == nil && eventHandlers.onMessage == nil {
		delete(ca.plugin.connEventHandlers, ca.conn)
	}
}

func (ca connActions) OnMessage(eventHandler func(*sdk.Message)) {
	ca.plugin.Lock()
	defer ca.plugin.Unlock()

	eventHandlers, ok := ca.plugin.connEventHandlers[ca.conn]
	if !ok {
		if eventHandler == nil {
			return
		}

		eventHandlers = &connEventHandlers{}
		ca.plugin.connEventHandlers[ca.conn] = eventHandlers
	}

	eventHandlers.onMessage = eventHandler

	if eventHandlers.onClose == nil && eventHandlers.onMessage == nil {
		delete(ca.plugin.connEventHandlers, ca.conn)
	}
}

func (ca connActions) SendMessage(message *sdk.Message) error {
	jsn, errJM := json.Marshal(message)
	if errJM != nil {
		return errJM
	}

	ca.ci.Lock()
	defer ca.ci.Unlock()

	if ca.ci.closed {
		return ca.ci.closeError
	} else {
		buf := bufio.NewWriter(ca.conn)

		if errWr := base.WriteNetStringToStream(buf, jsn); errWr != nil {
			return errWr
		}

		return buf.Flush()
	}
}

func readLoop(conn *tls.Conn, ci *connInfo) {
	buf := bufio.NewReader(conn)

	for {
		payload, errRd := base.ReadNetStringFromStream(buf, -1)
		if errRd != nil {
			pluginsLock.RLock()

			for plugin := range plugins {
				go onReadError(plugin, conn, errRd)
			}

			pluginsLock.RUnlock()

			connActions{nil, conn, ci}.Close()
			break
		}

		go handleNetstring(conn, payload)
	}
}

func handleNetstring(conn *tls.Conn, payload []byte) {
	message := &sdk.Message{}
	if errJU := json.Unmarshal(payload, message); errJU != nil {
		log.WithFields(log.Fields{
			"netstring": string(payload), "error": errJU.Error(),
		}).Error("Couldn't unJSONify netstring")
		return
	}

	pluginsLock.RLock()

	for plugin := range plugins {
		go onMessage(plugin, conn, message)
	}

	pluginsLock.RUnlock()
}

func onReadError(plugin *myPlugin, conn *tls.Conn, err error) {
	plugin.RLock()
	eventHandlers, ok := plugin.connEventHandlers[conn]
	plugin.RUnlock()

	if ok && eventHandlers.onClose != nil {
		eventHandlers.onClose(err)
	}
}

func onMessage(plugin *myPlugin, conn *tls.Conn, message *sdk.Message) {
	plugin.RLock()
	eventHandlers, ok := plugin.connEventHandlers[conn]
	plugin.RUnlock()

	if ok && eventHandlers.onMessage != nil {
		eventHandlers.onMessage(message)
	}
}
