package main

import "encoding/json"

func renderJson(in interface{}) string {
	if out, errJM := json.Marshal(in); errJM == nil {
		return string(out)
	} else {
		return errJM.Error()
	}
}
