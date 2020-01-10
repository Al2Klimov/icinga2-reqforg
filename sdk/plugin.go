package sdk

import (
	"crypto/tls"
	log "github.com/sirupsen/logrus"
	"io"
)

type PluginEntrypoint interface {
	NewInstance(SdkActions) PluginActions
}

type SdkActions interface {
	io.Closer

	AddConn(*tls.Conn)
	GetLogger() log.Ext1FieldLogger
}

type PluginActions interface {
	io.Closer

	GetName() string
	OnNewConn(ConnActions)
}

type ConnActions interface {
	io.Closer

	GetConn() *tls.Conn
	OnClose(func(error))
	OnMessage(func(*Message))
	SendMessage(*Message) error
}
