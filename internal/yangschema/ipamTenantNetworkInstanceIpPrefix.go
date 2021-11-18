package yangschema

import (
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/yndd/ndd-yang/pkg/leafref"
	"github.com/yndd/ndd-yang/pkg/yentry"
)

func initIpamTenantNetworkInstanceIpPrefix(p *yentry.Entry, opts ...yentry.EntryOption) *yentry.Entry {
	children := map[string]yentry.EntryInitFunc{
		"child":  initIpamTenantNetworkInstanceIpPrefixChild,
		"parent": initIpamTenantNetworkInstanceIpPrefixParent,
		"tag":    initIpamTenantNetworkInstanceIpPrefixTag,
	}
	e := &yentry.Entry{
		Name: "ip-prefix",
		Key: []string{
			"prefix",
		},
		Parent:           p,
		Children:         make(map[string]*yentry.Entry),
		ResourceBoundary: true,
		LeafRefs: []*leafref.LeafRef{
			{
				LocalPath: &gnmi.Path{
					Elem: []*gnmi.PathElem{
						{Name: "rir-name"},
					},
				},
				RemotePath: &gnmi.Path{
					Elem: []*gnmi.PathElem{
						{Name: "ipam"},
						{Name: "rir", Key: map[string]string{"name": ""}},
					},
				},
			},
		},
	}

	for _, opt := range opts {
		opt(e)
	}

	for name, initFunc := range children {
		e.Children[name] = initFunc(e, yentry.WithLogging(e.Log))
	}
	return e
}
