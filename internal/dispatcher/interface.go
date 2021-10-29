package dispatcher

import "github.com/openconfig/gnmi/proto/gnmi"

type Handler interface {
	HandleConfigEvent(o Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (Handler, error)
	SetParent(interface{}) error
	GetChildren() map[string]string
	UpdateConfig(interface{}) error
	UpdateStateCache() error
	DeleteStateCache() error
}
