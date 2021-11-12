package yangschema

import (
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-yang/pkg/yentry"
	"github.com/yndd/ndd-yang/pkg/yparser"
)

type ipamTenantNetworkInstanceIpPrefixParent struct {
	*yentry.Entry
}

func initIpamTenantNetworkInstanceIpPrefixParent(p yentry.Handler, opts ...yentry.HandlerOption) yentry.Handler {
	children := map[string]yentry.HandleInitFunc{
		"ip-prefix": initIpamTenantNetworkInstanceIpPrefixParentIpPrefix,
	}
	e := &yentry.Entry{
		Name:     "parent",
		Key:      []string{},
		Parent:   p,
		Children: make(map[string]yentry.Handler),
	}
	r := &ipamTenantNetworkInstanceIpPrefixParent{e}

	for _, opt := range opts {
		opt(r)
	}

	for name, initFunc := range children {
		r.Children[name] = initFunc(r, yentry.WithLogging(r.Log))
	}
	return r
}

func (r *ipamTenantNetworkInstanceIpPrefixParent) WithLogging(log logging.Logger) {
	r.Log = log
}

func (r *ipamTenantNetworkInstanceIpPrefixParent) GetKeys(p *gnmi.Path) []string {
	r.Log.Debug("Yangschema GetKeys", "Path", yparser.GnmiPath2XPath(p, true))
	if len(p.GetElem()) > 1 {
		return r.Children[p.GetElem()[0].GetName()].GetKeys(&gnmi.Path{Elem: p.GetElem()[1:]})
	} else {
		return r.GetKey()
	}
}
