package main

import (
	"crypto/tls"
	"github.com/Al2Klimov/icinga2-reqforg/sdk"
	log "github.com/sirupsen/logrus"
	"plugin"
	"reflect"
	"sync"
)

type myPlugin struct {
	sync.RWMutex

	pluginActions     sdk.PluginActions
	connEventHandlers map[*tls.Conn]*connEventHandlers
}

type connEventHandlers struct {
	onClose   func(error)
	onMessage func(*sdk.Message)
}

var plugins = map[*myPlugin]struct{}{}
var pluginsLock sync.RWMutex

func loadPlugin(file string) {
	plug, errPO := plugin.Open(file)
	if errPO != nil {
		log.WithFields(log.Fields{"plugin": file, "error": errPO.Error()}).Error("Couldn't open plugin")
		return
	}

	sym, errPL := plug.Lookup("PluginEntrypoint")
	if errPL != nil {
		log.WithFields(log.Fields{
			"plugin": file, "symbol": "PluginEntrypoint", "error": errPL.Error(),
		}).Warn("Couldn't lookup symbol")
		return
	}

	ep, ok := sym.(*sdk.PluginEntrypoint)
	if !ok || ep == nil || *ep == nil {
		log.WithFields(log.Fields{
			"plugin": file, "expected": "PluginEntrypoint", "actual": reflect.TypeOf(sym),
		}).Warn("Bad plugin entrypoint")
		return
	}

	log.WithFields(log.Fields{"plugin": file}).Info("Instanciating plugin")

	newPlugin := &myPlugin{connEventHandlers: map[*tls.Conn]*connEventHandlers{}}

	newPlugin.Lock()

	newPlugin.pluginActions = (*ep).NewInstance(sdkActions{newPlugin})

	if newPlugin.pluginActions == nil {
		log.WithFields(log.Fields{"plugin": file}).Warn("Bad plugin instance (nil)")
		return
	}

	pluginsLock.Lock()
	plugins[newPlugin] = struct{}{}
	pluginsLock.Unlock()

	newPlugin.Unlock()

	log.WithFields(log.Fields{"plugin": file, "name": newPlugin.pluginActions.GetName()}).Info("Instanciated plugin")
}
