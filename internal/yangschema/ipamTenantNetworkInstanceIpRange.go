package yangschema

import (
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/yndd/ndd-yang/pkg/leafref"
	"github.com/yndd/ndd-yang/pkg/yentry"
)

func initIpamTenantNetworkInstanceIpRange(p *yentry.Entry, opts ...yentry.EntryOption) *yentry.Entry {
	children := map[string]yentry.EntryInitFunc{
		"parent": initIpamTenantNetworkInstanceIpRangeParent,
		"tag":    initIpamTenantNetworkInstanceIpRangeTag,
	}
	e := &yentry.Entry{
		Name: "ip-range",
		Key: []string{
			"end",
			"start",
		},
		Parent:           p,
		Children:         make(map[string]*yentry.Entry),
		ResourceBoundary: true,
		LeafRefs:         []*leafref.LeafRef{},
	}

	for _, opt := range opts {
		opt(e)
	}

	for name, initFunc := range children {
		e.Children[name] = initFunc(e, yentry.WithLogging(e.Log))
	}
	if e.ResourceBoundary {
		e.Register(&gnmi.Path{Elem: []*gnmi.PathElem{}})
	}
	return e
}
