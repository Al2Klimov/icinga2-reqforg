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
	return "Set Next Notification"
}

func (pa pluginActions) OnNewConn(connActions sdk.ConnActions) {
	for {
		go func() {
			err := connActions.SendMessage(&sdk.Message{
				Jsonrpc: "2.0",
				Method:  "event::SetNextNotification",
				/**
				  	host				String		Host name
					service				String		Service name
					notification		String		Notification name
					next_notification	Timestamp	Next scheduled notification time as UNIX timestamp.
				*/
				Params: map[string]interface{}{
					"host": "test-host",
					"notification": "test-notification",
					"next_notification": time.Now().Unix() + 60,
				},
			})

			if err != nil {
				pa.sdkActions.GetLogger().WithField("error", err).Error("Error while sending checkresult")
			}
		}()

		time.Sleep(time.Millisecond * 100)
	}
}