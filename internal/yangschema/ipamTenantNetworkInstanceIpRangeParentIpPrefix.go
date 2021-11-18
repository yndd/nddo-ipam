package yangschema

import (
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/yndd/ndd-yang/pkg/leafref"
	"github.com/yndd/ndd-yang/pkg/yentry"
)

func initIpamTenantNetworkInstanceIpRangeParentIpPrefix(p *yentry.Entry, opts ...yentry.EntryOption) *yentry.Entry {
	children := map[string]yentry.EntryInitFunc{}
	e := &yentry.Entry{
		Name: "ip-prefix",
		Key: []string{
			"prefix",
		},
		Parent:           p,
		Children:         make(map[string]*yentry.Entry),
		ResourceBoundary: false,
		LeafRefs: []*leafref.LeafRef{
			{
				LocalPath: &gnmi.Path{
					Elem: []*gnmi.PathElem{
						{Name: "prefix"},
					},
				},
				RemotePath: &gnmi.Path{
					Elem: []*gnmi.PathElem{
						{Name: "ipam"},
						{Name: "tenant"},
						{Name: "network-instance"},
						{Name: "ip-prefix", Key: map[string]string{"prefix": ""}},
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
