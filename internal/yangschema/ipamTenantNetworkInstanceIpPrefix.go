package yangschema

import (
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-yang/pkg/yentry"
	"github.com/yndd/ndd-yang/pkg/yparser"
)

type ipamTenantNetworkInstanceIpPrefix struct {
	*yentry.Entry
}

func initIpamTenantNetworkInstanceIpPrefix(p yentry.Handler, opts ...yentry.HandlerOption) yentry.Handler {
	children := map[string]yentry.HandleInitFunc{
		"child":  initIpamTenantNetworkInstanceIpPrefixChild,
		"parent": initIpamTenantNetworkInstanceIpPrefixParent,
		"tag":    initIpamTenantNetworkInstanceIpPrefixTag,
	}
	e := &yentry.Entry{
		Name: "ip-prefix",
		Key: []string{
			"prefix",
		},
		Parent:   p,
		Children: make(map[string]yentry.Handler),
	}
	r := &ipamTenantNetworkInstanceIpPrefix{e}

	for _, opt := range opts {
		opt(r)
	}

	for name, initFunc := range children {
		r.Children[name] = initFunc(r, yentry.WithLogging(r.Log))
	}
	return r
}

func (r *ipamTenantNetworkInstanceIpPrefix) WithLogging(log logging.Logger) {
	r.Log = log
}

func (r *ipamTenantNetworkInstanceIpPrefix) GetKeys(p *gnmi.Path) []string {
	r.Log.Debug("Yangschema GetKeys", "Path", yparser.GnmiPath2XPath(p, true))
	if len(p.GetElem()) > 1 {
		return r.Children[p.GetElem()[0].GetName()].GetKeys(&gnmi.Path{Elem: p.GetElem()[1:]})
	} else {
		return r.GetKey()
	}
}
