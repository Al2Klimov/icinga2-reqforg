package main

import (
	"crypto/tls"
	"github.com/Al2Klimov/icinga2-reqforg/sdk"
	log "github.com/sirupsen/logrus"
	"os"
)

type pluginEntrypoint struct {
}

var _ sdk.PluginEntrypoint = pluginEntrypoint{}

func (pluginEntrypoint) NewInstance(sdkActions sdk.SdkActions) sdk.PluginActions {
	go dialOut(sdkActions)
	return pluginActions{}
}

func dialOut(sdkActions sdk.SdkActions) {
	if doa, ok := os.LookupEnv("CONNECTION_ADDR"); ok {
		cert, err := tls.LoadX509KeyPair("certs/localhost.crt", "certs/localhost.key")
		if err != nil {
			log.Fatalf("server: loadkeys: %s", err)
		}
		config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

		conn, errDl := tls.Dial("tcp", doa, &config)
		if errDl != nil {
			sdkActions.GetLogger().WithFields(log.Fields{"error": errDl.Error()}).Error("Couldn't connect")
			return
		}

		sdkActions.AddConn(conn)
	} else {
		sdkActions.GetLogger().Error("Missing env var CONNECTION_ADDR")
	}
}

type pluginActions struct {
}

var _ sdk.PluginActions = pluginActions{}

func (pa pluginActions) Close() error {
	return nil
}

func (pluginActions) GetName() string {
	return "Dial out"
}

func (pa pluginActions) OnNewConn(sdk.ConnActions) {
}