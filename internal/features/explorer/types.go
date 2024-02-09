package explorer

import "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/client"

type ConnectionParameters struct {
}

type Class struct {
	Name       string
	MethodList []Method
}

type Method struct {
	Name    string
	Execute func(string) (response client.Message, err error)
}

type WsmanMethods struct {
	MethodList []Method
}
type Response struct {
	XMLInput  string
	XMLOutput string
}
