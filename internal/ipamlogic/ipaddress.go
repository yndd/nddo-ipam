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
				{Name: "ip-address", Key: map[string]string{"address": "*"}},
			},
			Handler: networkinstanceCreate,
		},
	})
}

type ipaddress struct {
	dispatcher.Resource
	data   *ipamv1alpha1.NddoipamIpamTenantNetworkInstanceIpAddress
	parent *networkinstance
}

func (r *ipaddress) WithLogging(log logging.Logger) {
	r.Log = log
}

func (r *ipaddress) WithStateCache(c *cache.Cache) {
	r.StateCache = c
}

func (r *ipaddress) WithConfigCache(c *cache.Cache) {
	r.ConfigCache = c
}

func (r *ipaddress) WithPrefix(p *gnmi.Path) {
	r.Prefix = p
}

func (r *ipaddress) WithPathElem(pe []*gnmi.PathElem) {
	r.PathElem = pe[0]
}

func NewIpAddress(n string, opts ...dispatcher.HandlerOption) dispatcher.Handler {
	x := &ipaddress{}
	x.Key = n

	for _, opt := range opts {
		opt(x)
	}
	return x
}

func ipaddressGetKey(pe []*gnmi.PathElem) string {
	return pe[0].GetKey()["ip-address"]
}

func ipaddressCreate(log logging.Logger, cc, sc *cache.Cache, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) dispatcher.Handler {
	n := ipaddressGetKey(pe)
	return NewIpAddress(n,
		dispatcher.WithPrefix(prefix),
		dispatcher.WithPathElem(pe),
		dispatcher.WithLogging(log),
		dispatcher.WithStateCache(sc),
		dispatcher.WithConfigCache(cc))
}

func (r *ipaddress) HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	log := r.Log.WithValues("Operation", o, "Path Elem", pe)

	log.Debug("ipaddress Handle")

	if len(pe) == 1 {
		return nil, errors.New("the handle should have been terminated in the parent")
	} else {
		return nil, errors.New("there is no children in the ipaddress resource ")
	}
}

func (r *ipaddress) UpdateChild(children map[string]dispatcher.HandleConfigEventFunc, pathElemName string, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	return nil, errors.New("there is no children in the ipaddress resource ")
}

func (r *ipaddress) DeleteChild(pathElemName string, pe []*gnmi.PathElem) error {
	return errors.New("there is no children in the ipaddress resource ")
}

func (r *ipaddress) SetParent(parent interface{}) error {
	p, ok := parent.(*networkinstance)
	if !ok {
		return errors.New("wrong parent object")
	}
	r.parent = p
	return nil
}

func (r *ipaddress) GetChildren() map[string]string {
	var x map[string]string
	return x
}

func (r *ipaddress) UpdateConfig(d interface{}) error {
	r.Copy(d)

	// TBD

	if err := r.UpdateStateCache(); err != nil {
		return err
	}
	return nil
}

func (r *ipaddress) GetPathElem(p []*gnmi.PathElem, do_recursive bool) ([]*gnmi.PathElem, error) {
	if r.parent != nil {
		p, err := r.parent.GetPathElem(p, true)
		if err != nil {
			return nil, err
		}
		p = append(p, r.PathElem)
		return p, nil
	}
	return nil, nil
}

func (r *ipaddress) UpdateStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	r.Log.Debug("ipaddress Update Cache", "PathElem", pe, "Prefix", r.Prefix, "data", r.data)
	if err := updateCache(r.Log, r.StateCache, r.Prefix, &gnmi.Path{Elem: pe}, r.data); err != nil {
		r.Log.Debug("ipaddress Update Error")
		return err
	}
	r.Log.Debug("ipaddress Update ok")
	return nil
}

func (r *ipaddress) DeleteStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	if err := deleteCache(r.Log, r.StateCache, r.Prefix, &gnmi.Path{Elem: pe}); err != nil {
		return err
	}
	return nil
}

func (r *ipaddress) Copy(d interface{}) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	x := ipamv1alpha1.NddoipamIpamTenantNetworkInstanceIpAddress{}
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}
	r.data = (&x).DeepCopy()
	return nil
}
