package main

import (
	"github.com/Al2Klimov/icinga2-reqforg/sdk"
	"github.com/Pallinder/go-randomdata"
	"math/rand"
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
	return "Check Result"
}

func (pa pluginActions) OnNewConn(connActions sdk.ConnActions) {
	for {
		go func() {
			timestamp := time.Now().Unix()
			cr := map[string]interface{}{
				"active":          true,
				"check_source":    "localhost",
				"command":         "dummy",
				"execution_start": timestamp,
				"execution_end":   timestamp,
				"exit_status":     0.0,
				"output":          "Check was successful.",
				"performance_data": generateRandomPerfData(),
				"schedule_start": timestamp,
				"schedule_end":   timestamp,
				"state":          rand.Intn(3),
				"ttl":            0.0,
				"type":           "CheckResult",
			}

			err := connActions.SendMessage(&sdk.Message{
				Jsonrpc: "2.0",
				Method:  "event::CheckResult",
				/**
				  host		String			Host name
				  service	String			Service name
				  cr		Serialized CR	Check result
				*/
				Params: map[string]interface{}{
					"host": "test-host",
					"cr":   cr,
				},
			})

			if err != nil {
				pa.sdkActions.GetLogger().WithField("error", err).Error("Error while sending checkresult")
			}
		}()

		time.Sleep(time.Millisecond * 1000)
	}
}

func generateRandomPerfData() *[]interface{} {
	count := rand.Intn(10000)
	var perfData []interface{}

	for i := 0; i < count; i++ {
		perfData = append(
			perfData,
			map[string]interface{}{
				"counter": false,
				"crit": nil,
				"label": randomdata.SillyName(),
				"max": nil,
				"min": nil,
				"type": "PerfdataValue",
				"unit": "",
				"value": rand.Float64(),
				"warn": nil,
			},
		)
	}

	return &perfData
}