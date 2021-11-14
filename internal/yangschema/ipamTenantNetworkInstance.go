package yangschema

import (
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-yang/pkg/yentry"
	"github.com/yndd/ndd-yang/pkg/yparser"
)

type ipamTenantNetworkInstance struct {
	*yentry.Entry
}

func initIpamTenantNetworkInstance(p yentry.Handler, opts ...yentry.HandlerOption) yentry.Handler {
	children := map[string]yentry.HandleInitFunc{
		"ip-address": initIpamTenantNetworkInstanceIpAddress,
		"ip-prefix":  initIpamTenantNetworkInstanceIpPrefix,
		"ip-range":   initIpamTenantNetworkInstanceIpRange,
		"tag":        initIpamTenantNetworkInstanceTag,
	}
	e := &yentry.Entry{
		Name: "network-instance",
		Key: []string{
			"name",
		},
		Parent:   p,
		Children: make(map[string]yentry.Handler),
	}
	r := &ipamTenantNetworkInstance{e}

	for _, opt := range opts {
		opt(r)
	}

	for name, initFunc := range children {
		r.Children[name] = initFunc(r, yentry.WithLogging(r.Log))
	}
	return r
}

func (r *ipamTenantNetworkInstance) WithLogging(log logging.Logger) {
	r.Log = log
}

func (r *ipamTenantNetworkInstance) GetKeys(p *gnmi.Path) []string {
	r.Log.Debug("Yangschema GetKeys", "Path", yparser.GnmiPath2XPath(p, true), "Name", r.GetName(), "Key", r.GetKey())
	if len(p.GetElem()) >= 1 {
		return r.Children[p.GetElem()[0].GetName()].GetKeys(&gnmi.Path{Elem: p.GetElem()[1:]})
	} else {
		return r.GetKey()
	}
}
