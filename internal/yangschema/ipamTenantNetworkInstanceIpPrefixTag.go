package yangschema

import (
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/yndd/ndd-yang/pkg/leafref"
	"github.com/yndd/ndd-yang/pkg/yentry"
)

func initIpamTenantNetworkInstanceIpPrefixTag(p *yentry.Entry, opts ...yentry.EntryOption) *yentry.Entry {
	children := map[string]yentry.EntryInitFunc{}
	e := &yentry.Entry{
		Name: "tag",
		Key: []string{
			"key",
		},
		Parent:           p,
		Children:         make(map[string]*yentry.Entry),
		ResourceBoundary: false,
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
