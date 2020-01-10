package main

import (
	"github.com/Al2Klimov/icinga2-reqforg/sdk"
	log "github.com/sirupsen/logrus"
	"unsafe"
)

type pluginEntrypoint struct {
}

var _ sdk.PluginEntrypoint = pluginEntrypoint{}

func (pluginEntrypoint) NewInstance(sdkActions sdk.SdkActions) sdk.PluginActions {
	sdkActions.GetLogger().Info("Starting")
	return pluginActions{sdkActions}
}

type pluginActions struct {
	sdkActions sdk.SdkActions
}

var _ sdk.PluginActions = pluginActions{}

func (pa pluginActions) Close() error {
	pa.sdkActions.GetLogger().Info("Stopping")
	return nil
}

func (pluginActions) GetName() string {
	return "FYI"
}

func (pa pluginActions) OnNewConn(connActions sdk.ConnActions) {
	pa.sdkActions.GetLogger().WithFields(log.Fields{
		"conn": uintptr(unsafe.Pointer(connActions.GetConn())),
	}).Info("New connection")

	connActions.OnClose(func(err error) {
		pa.sdkActions.GetLogger().WithFields(log.Fields{
			"conn": uintptr(unsafe.Pointer(connActions.GetConn())), "error": err.Error(),
		}).Info("Connection closed")
	})

	connActions.OnMessage(func(message *sdk.Message) {
		pa.sdkActions.GetLogger().WithFields(log.Fields{
			"conn": uintptr(unsafe.Pointer(connActions.GetConn())), "message": renderJson(message),
		}).Info("Got message")
	})
}
