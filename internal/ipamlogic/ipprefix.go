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
				{Name: "ip-prefix", Key: map[string]string{"prefix": "*"}},
			},
			Handler: networkinstanceCreate,
		},
	})
}

type IpPrefix interface {
	HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error)
}

type ipprefix struct {
	log         logging.Logger
	configCache *cache.Cache
	stateCache  *cache.Cache
	pathElem    *gnmi.PathElem
	prefix      *gnmi.Path
	key         string
	parent      *networkinstance
	data        *ipamv1alpha1.NddoipamIpamTenantNetworkInstanceIpPrefix
}

type IpPrefixOption func(*ipprefix)

// WithRirRirLogger initializes the logger.
func WithIpPrefixLogger(log logging.Logger) IpPrefixOption {
	return func(o *ipprefix) {
		o.log = log
	}
}

// WithRirRirCache initializes the cache.
func WithIpPrefixStateCache(c *cache.Cache) IpPrefixOption {
	return func(o *ipprefix) {
		o.stateCache = c
	}
}

func WithIpPrefixConfigCache(c *cache.Cache) IpPrefixOption {
	return func(o *ipprefix) {
		o.configCache = c
	}
}

func WithIpPrefixPrefix(p *gnmi.Path) IpPrefixOption {
	return func(o *ipprefix) {
		o.prefix = p
	}
}

func WithIpPrefixPathElem(pe []*gnmi.PathElem) IpPrefixOption {
	return func(o *ipprefix) {
		o.pathElem = pe[0]
	}
}

func NewIpPrefix(n string, opts ...IpPrefixOption) *ipprefix {
	x := &ipprefix{
		key: n,
	}

	for _, opt := range opts {
		opt(x)
	}
	return x
}

func ipprefixGetKey(pe []*gnmi.PathElem) string {
	return pe[0].GetKey()["ip-prefix"]
}

func ipprefixCreate(log logging.Logger, cc, sc *cache.Cache, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) dispatcher.Handler {
	niName := ipprefixGetKey(pe)
	return NewIpPrefix(niName,
		WithIpPrefixPrefix(prefix),
		WithIpPrefixPathElem(pe),
		WithIpPrefixLogger(log),
		WithIpPrefixStateCache(sc),
		WithIpPrefixConfigCache(cc))
}

func (r *ipprefix) HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	log := r.log.WithValues("Operation", o, "Path Elem", pe)

	log.Debug("ipprefix Handle")

	if len(pe) == 1 {
		return nil, errors.New("the handle should have been terminated in the parent")
	} else {
		return nil, errors.New("there is no children in the ipprefix resource ")
	}
}

func (r *ipprefix) UpdateChild(children map[string]dispatcher.HandleConfigEventFunc, pathElemName string, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	return nil, errors.New("there is no children in the ipprefix resource ")
}

func (r *ipprefix) DeleteChild(pathElemName string, pe []*gnmi.PathElem) error {
	return errors.New("there is no children in the ipprefix resource ")
}

func (r *ipprefix) SetParent(parent interface{}) error {
	p, ok := parent.(*networkinstance)
	if !ok {
		return errors.New("wrong parent object")
	}
	r.parent = p
	return nil
}

func (r *ipprefix) GetChildren() map[string]string {
	var x map[string]string
	return x
}

func (r *ipprefix) UpdateConfig(d interface{}) error {
	r.Copy(d)

	// TBD

	if err := r.UpdateStateCache(); err != nil {
		return err
	}
	return nil
}

func (r *ipprefix) GetPathElem(p []*gnmi.PathElem, do_recursive bool) ([]*gnmi.PathElem, error) {
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

func (r *ipprefix) UpdateStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	r.log.Debug("ipprefix Update Cache", "PathElem", pe, "Prefix", r.prefix, "data", r.data)
	if err := updateCache(r.log, r.stateCache, r.prefix, &gnmi.Path{Elem: pe}, r.data); err != nil {
		r.log.Debug("ipprefix Update Error")
		return err
	}
	r.log.Debug("ipprefix Update ok")
	return nil
}

func (r *ipprefix) DeleteStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	if err := deleteCache(r.log, r.stateCache, r.prefix, &gnmi.Path{Elem: pe}); err != nil {
		return err
	}
	return nil
}

func (r *ipprefix) Copy(d interface{}) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	x := ipamv1alpha1.NddoipamIpamTenantNetworkInstanceIpPrefix{}
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}
	r.data = (&x).DeepCopy()
	return nil
}
