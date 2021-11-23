package yangschema

import (
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/yndd/ndd-yang/pkg/leafref"
	"github.com/yndd/ndd-yang/pkg/yentry"
)

func initIpamTenantNetworkInstanceIpAddressIpRange(p *yentry.Entry, opts ...yentry.EntryOption) *yentry.Entry {
	children := map[string]yentry.EntryInitFunc{}
	e := &yentry.Entry{
		Name: "ip-range",
		Key: []string{
			"end",
			"start",
		},
		Parent:           p,
		Children:         make(map[string]*yentry.Entry),
		ResourceBoundary: false,
		LeafRefs: []*leafref.LeafRef{
			{
				LocalPath: &gnmi.Path{
					Elem: []*gnmi.PathElem{
						{Name: "end"},
					},
				},
				RemotePath: &gnmi.Path{
					Elem: []*gnmi.PathElem{
						{Name: "ipam"},
						{Name: "tenant"},
						{Name: "ip-range", Key: map[string]string{"start": ""}},
						{Name: "start]", Key: map[string]string{"end": ""}},
					},
				},
			},
			{
				LocalPath: &gnmi.Path{
					Elem: []*gnmi.PathElem{
						{Name: "start"},
					},
				},
				RemotePath: &gnmi.Path{
					Elem: []*gnmi.PathElem{
						{Name: "ipam"},
						{Name: "tenant"},
						{Name: "ip-range", Key: map[string]string{"end": ""}},
						{Name: "end]", Key: map[string]string{"start": ""}},
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
	if e.ResourceBoundary {
		e.Register(&gnmi.Path{Elem: []*gnmi.PathElem{}})
	}
	return e
}
