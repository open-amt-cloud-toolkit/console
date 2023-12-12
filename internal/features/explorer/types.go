package explorer

import "github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/client"

type ConnectionParameters struct {
}

type Class struct {
	Name       string
	MethodList []Method
}

type Method struct {
	Name                   string
	Execute                func() (response client.Message, err error)
	ExecuteWithStringInput func(string) (response client.Message, err error)
	ExecuteWithIntInput    func(int) (response client.Message, err error)
}

type WsmanMethods struct {
	MethodList []Method
}
type Response struct {
	XMLInput  string
	XMLOutput string
}
