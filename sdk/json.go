package sdk

type Request struct {
	Jsonrpc string      `json:"jsonrpc,omitempty"`
	Id      interface{} `json:"id,omitempty"`
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
}

type Response struct {
	Jsonrpc string      `json:"jsonrpc,omitempty"`
	Id      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}
