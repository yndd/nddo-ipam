package yangschema

import (
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-yang/pkg/yentry"
	"github.com/yndd/ndd-yang/pkg/yparser"
)

type ipam struct {
	*yentry.Entry
}

func initIpam(p yentry.Handler, opts ...yentry.HandlerOption) yentry.Handler {
	children := map[string]yentry.HandleInitFunc{
		"rir":    initIpamRir,
		"tenant": initIpamTenant,
	}
	e := &yentry.Entry{
		Name:     "ipam",
		Key:      []string{},
		Parent:   p,
		Children: make(map[string]yentry.Handler),
	}
	r := &ipam{e}

	for _, opt := range opts {
		opt(r)
	}

	for name, initFunc := range children {
		r.Children[name] = initFunc(r, yentry.WithLogging(r.Log))
	}
	return r
}

func (r *ipam) WithLogging(log logging.Logger) {
	r.Log = log
}

func (r *ipam) GetKeys(p *gnmi.Path) []string {
	r.Log.Debug("Yangschema GetKeys", "Path", yparser.GnmiPath2XPath(p, true), "Name", r.GetName(), "Key", r.GetKey())
	if len(p.GetElem()) >= 1 {
		return r.Children[p.GetElem()[0].GetName()].GetKeys(&gnmi.Path{Elem: p.GetElem()[1:]})
	} else {
		return r.GetKey()
	}
}
