package main

import (
	"github.com/Al2Klimov/icinga2-reqforg/sdk"
)

type pluginEntrypoint struct {
}

var _ sdk.PluginEntrypoint = pluginEntrypoint{}

func (pluginEntrypoint) NewInstance(sdk.SdkActions) sdk.PluginActions {
	return pluginActions{}
}

type pluginActions struct {
}

var _ sdk.PluginActions = pluginActions{}

func (pa pluginActions) Close() error {
	return nil
}

func (pluginActions) GetName() string {
	return "Hello"
}

func (pa pluginActions) OnNewConn(connActions sdk.ConnActions) {
	connActions.SendMessage(&sdk.Message{
		Jsonrpc: "2.0",
		Method:  "icinga::Hello",
		Params:  struct{}{},
	})
}
