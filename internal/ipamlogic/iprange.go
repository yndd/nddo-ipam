package ipamlogic

import (
	"encoding/json"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/pkg/errors"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-yang/pkg/cache"
	ipamv1alpha1 "github.com/yndd/nddo-ipam/apis/ipam/v1alpha1"
	"github.com/yndd/nddo-ipam/internal/dispatcher"
)

func init() {
	dispatcher.Register("network-instance", []*dispatcher.EventHandler{
		{
			Kind: dispatcher.EventHandlerCreate,
			PathElem: []*gnmi.PathElem{
				{Name: "ipam"},
				{Name: "tenant", Key: map[string]string{"name": "*"}},
				{Name: "network-instance", Key: map[string]string{"name": "*"}},
				{Name: "ip-range", Key: map[string]string{"range": "*"}},
			},
			Handler: networkinstanceCreate,
		},
	})
}

type IpRange interface {
	HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error)
}

type iprange struct {
	log         logging.Logger
	configCache *cache.Cache
	stateCache  *cache.Cache
	pathElem    *gnmi.PathElem
	prefix      *gnmi.Path
	key         string
	parent      *networkinstance
	data        *ipamv1alpha1.NddoipamIpamTenantNetworkInstanceIpRange
}

type IpRangeOption func(*iprange)

// WithRirRirLogger initializes the logger.
func WithIpRangeLogger(log logging.Logger) IpRangeOption {
	return func(o *iprange) {
		o.log = log
	}
}

// WithRirRirCache initializes the cache.
func WithIpRangeStateCache(c *cache.Cache) IpRangeOption {
	return func(o *iprange) {
		o.stateCache = c
	}
}

func WithIpRangeConfigCache(c *cache.Cache) IpRangeOption {
	return func(o *iprange) {
		o.configCache = c
	}
}

func WithIpRangePrefix(p *gnmi.Path) IpRangeOption {
	return func(o *iprange) {
		o.prefix = p
	}
}

func WithIpRangePathElem(pe []*gnmi.PathElem) IpRangeOption {
	return func(o *iprange) {
		o.pathElem = pe[0]
	}
}

func NewIpRange(n string, opts ...IpRangeOption) *iprange {
	x := &iprange{
		key: n,
	}

	for _, opt := range opts {
		opt(x)
	}
	return x
}

func iprangeGetKey(pe []*gnmi.PathElem) string {
	return pe[0].GetKey()["start"] + "." + pe[0].GetKey()["end"]
}

func iprangeCreate(log logging.Logger, cc, sc *cache.Cache, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) dispatcher.Handler {
	niName := iprangeGetKey(pe)
	return NewIpRange(niName,
		WithIpRangePrefix(prefix),
		WithIpRangePathElem(pe),
		WithIpRangeLogger(log),
		WithIpRangeStateCache(sc),
		WithIpRangeConfigCache(cc))
}

func (r *iprange) HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	log := r.log.WithValues("Operation", o, "Path Elem", pe)

	log.Debug("iprange Handle")

	if len(pe) == 1 {
		return nil, errors.New("the handle should have been terminated in the parent")
	} else {
		return nil, errors.New("there is no children in the iprange resource ")
	}
}

func (r *iprange) UpdateChild(children map[string]dispatcher.HandleConfigEventFunc, pathElemName string, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	return nil, errors.New("there is no children in the iprange resource ")
}

func (r *iprange) DeleteChild(pathElemName string, pe []*gnmi.PathElem) error {
	return errors.New("there is no children in the iprange resource ")
}

func (r *iprange) SetParent(parent interface{}) error {
	p, ok := parent.(*networkinstance)
	if !ok {
		return errors.New("wrong parent object")
	}
	r.parent = p
	return nil
}

func (r *iprange) GetChildren() map[string]string {
	var x map[string]string
	return x
}

func (r *iprange) UpdateConfig(d interface{}) error {
	r.Copy(d)

	// TBD

	if err := r.UpdateStateCache(); err != nil {
		return err
	}

	return nil
}

func (r *iprange) GetPathElem(p []*gnmi.PathElem, do_recursive bool) ([]*gnmi.PathElem, error) {
	if r.parent != nil {
		p, err := r.parent.GetPathElem(p, true)
		if err != nil {
			return nil, err
		}
		p = append(p, r.pathElem)
		return p, nil
	}
	return nil, nil
}

func (r *iprange) UpdateStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	r.log.Debug("iprange Update Cache", "PathElem", pe, "Prefix", r.prefix, "data", r.data)
	if err := updateCache(r.log, r.stateCache, r.prefix, &gnmi.Path{Elem: pe}, r.data); err != nil {
		r.log.Debug("iprange Update Error")
		return err
	}
	r.log.Debug("iprange Update ok")
	return nil
}

func (r *iprange) DeleteStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	if err := deleteCache(r.log, r.stateCache, r.prefix, &gnmi.Path{Elem: pe}); err != nil {
		return err
	}
	return nil
}

func (r *iprange) Copy(d interface{}) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	x := ipamv1alpha1.NddoipamIpamTenantNetworkInstanceIpRange{}
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}
	r.data = (&x).DeepCopy()
	return nil
}
