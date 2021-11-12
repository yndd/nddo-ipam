package yangschema

import (
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-yang/pkg/yentry"
	"github.com/yndd/ndd-yang/pkg/yparser"
)

type ipamTenantNetworkInstanceIpPrefixChild struct {
	*yentry.Entry
}

func initIpamTenantNetworkInstanceIpPrefixChild(p yentry.Handler, opts ...yentry.HandlerOption) yentry.Handler {
	children := map[string]yentry.HandleInitFunc{
		"ip-prefix": initIpamTenantNetworkInstanceIpPrefixChildIpPrefix,
	}
	e := &yentry.Entry{
		Name:     "child",
		Key:      []string{},
		Parent:   p,
		Children: make(map[string]yentry.Handler),
	}
	r := &ipamTenantNetworkInstanceIpPrefixChild{e}

	for _, opt := range opts {
		opt(r)
	}

	for name, initFunc := range children {
		r.Children[name] = initFunc(r, yentry.WithLogging(r.Log))
	}
	return r
}

func (r *ipamTenantNetworkInstanceIpPrefixChild) WithLogging(log logging.Logger) {
	r.Log = log
}

func (r *ipamTenantNetworkInstanceIpPrefixChild) GetKeys(p *gnmi.Path) []string {
	r.Log.Debug("Yangschema GetKeys", "Path", yparser.GnmiPath2XPath(p, true))
	if len(p.GetElem()) > 1 {
		return r.Children[p.GetElem()[0].GetName()].GetKeys(&gnmi.Path{Elem: p.GetElem()[1:]})
	} else {
		return r.GetKey()
	}
}