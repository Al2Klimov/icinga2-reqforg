package main

import (
	"github.com/Al2Klimov/icinga2-reqforg/sdk"
	"time"
)

type pluginEntrypoint struct {
}

var _ sdk.PluginEntrypoint = pluginEntrypoint{}

func (pluginEntrypoint) NewInstance(sdkActions sdk.SdkActions) sdk.PluginActions {
	return pluginActions{sdkActions}
}

type pluginActions struct {
	sdkActions sdk.SdkActions
}

var _ sdk.PluginActions = pluginActions{}

func (pa pluginActions) Close() error {
	return nil
}

func (pluginActions) GetName() string {
	return "PKI Request Certificate"
}

func (pa pluginActions) OnNewConn(connActions sdk.ConnActions) {
	for {
		err := connActions.SendMessage(&sdk.Message{
			Jsonrpc: "2.0",
			Method:  "pki::RequestCertificate",
			Params: struct {
				ticket string
				certRequest string
			}{
				ticket: "djkanskdasndjk",
			},
		})

		if err != nil {
			pa.sdkActions.GetLogger().WithField("error", err).Error("Error while requesting certificate")
		}

		time.Sleep(time.Nanosecond * 100)
	}
}