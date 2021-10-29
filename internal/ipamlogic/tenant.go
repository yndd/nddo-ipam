package ipamlogic

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-yang/pkg/cache"
	ipamv1alpha1 "github.com/yndd/nddo-ipam/apis/ipam/v1alpha1"
	"github.com/yndd/nddo-ipam/internal/dispatcher"
)

func init() {
	dispatcher.Register("tenant", []*dispatcher.EventHandler{
		{
			Kind: dispatcher.EventHandlerCreate,
			PathElem: []*gnmi.PathElem{
				{Name: "ipam"},
				{Name: "tenant", Key: map[string]string{"name": "*"}},
			},
			Handler: tenantCreate,
		},
	})
}

type Tenant interface {
	HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error)
}

type tenant struct {
	log              logging.Logger
	configCache      *cache.Cache
	stateCache       *cache.Cache
	pathElem         *gnmi.PathElem
	prefix           *gnmi.Path
	name             string
	parent           *ipam
	data             *ipamv1alpha1.NddoipamIpamTenant
	networkInstances map[string]dispatcher.Handler
}

type TenantOption func(*tenant)

// WithRirRirLogger initializes the logger.
func WithTenantLogger(log logging.Logger) TenantOption {
	return func(o *tenant) {
		o.log = log
	}
}

// WithRirRirCache initializes the cache.
func WithTenantStateCache(c *cache.Cache) TenantOption {
	return func(o *tenant) {
		o.stateCache = c
	}
}

func WithTenantConfigCache(c *cache.Cache) TenantOption {
	return func(o *tenant) {
		o.configCache = c
	}
}

func WithTenantPrefix(p *gnmi.Path) TenantOption {
	return func(o *tenant) {
		o.prefix = p
	}
}

func WithTenantPathElem(pe []*gnmi.PathElem) TenantOption {
	return func(o *tenant) {
		o.pathElem = pe[0]
	}
}

func NewTenant(n string, opts ...TenantOption) *tenant {
	x := &tenant{
		name:             n,
		networkInstances: make(map[string]dispatcher.Handler),
	}

	for _, opt := range opts {
		opt(x)
	}
	return x
}

func tenantGetKey(pe []*gnmi.PathElem) string {
	return pe[0].GetKey()["name"]
}

func tenantCreate(log logging.Logger, cc, sc *cache.Cache, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) dispatcher.Handler {
	tenantName := tenantGetKey(pe)
	return NewTenant(tenantName,
		WithTenantPrefix(prefix),
		WithTenantPathElem(pe),
		WithTenantLogger(log),
		WithTenantStateCache(sc),
		WithTenantConfigCache(cc))
}

func (r *tenant) HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	log := r.log.WithValues("Operation", o, "Path Elem", pe)

	log.Debug("tenant HandleConfigEvent")

	children := map[string]dispatcher.HandleConfigEventFunc{
		"network-instance": networkinstanceCreate,
	}

	// check path Element Name
	pathElemName := pe[0].GetName()
	if _, ok := children[pathElemName]; !ok {
		return nil, errors.Wrap(errors.New("unexpected pathElem"), fmt.Sprintf("ni HandleConfigEvent: %s", pathElemName))
	}

	if len(pe) == 1 {
		log.Debug("tenant Handle pathelem =1")
		// handle local
		switch o {
		case dispatcher.OperationUpdate:
			i, err := r.UpdateChild(children, pathElemName, prefix, pe, d)
			if err != nil {
				return nil, err
			}
			switch pathElemName {
			case "network-instance":
				r.networkInstances[networkinstanceGetKey(pe)] = i
			}
			return i, nil
		case dispatcher.OperationDelete:
			if err := r.DeleteChild(pathElemName, pe); err != nil {
				return nil, err
			}
			switch pathElemName {
			case "network-instance":
				delete(r.networkInstances, networkinstanceGetKey(pe))
			}
			return nil, nil
		}
	} else {
		log.Debug("tenant Handle pathelem >1")
		switch pathElemName {
		case "network-instance":
			if _, ok := r.networkInstances[rirGetKey(pe)]; !ok {
				// create resource with dummy data
				i, err := r.UpdateChild(children, pathElemName, prefix, pe[:1], nil)
				if err != nil {
					return nil, err
				}
				r.networkInstances[networkinstanceGetKey(pe)] = i
			}
			return r.networkInstances[networkinstanceGetKey(pe)].HandleConfigEvent(o, prefix, pe[1:], d)
		}
	}
	return nil, nil
}

func (r *tenant) UpdateChild(children map[string]dispatcher.HandleConfigEventFunc, pathElemName string, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	i := children[pathElemName](r.log, r.configCache, r.stateCache, prefix, pe, d)
	if err := i.SetParent(r); err != nil {
		return nil, err
	}
	if d != nil {
		if err := i.UpdateConfig(d); err != nil {
			return nil, err
		}
		if err := i.UpdateStateCache(); err != nil {
			return nil, err
		}
	}
	return i, nil
}

func (r *tenant) DeleteChild(pathElemName string, pe []*gnmi.PathElem) error {
	switch pathElemName {
	case "network-instance":
		if i, ok := r.networkInstances[networkinstanceGetKey(pe)]; ok {
			if err := i.DeleteStateCache(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *tenant) SetParent(parent interface{}) error {
	p, ok := parent.(*ipam)
	if !ok {
		return errors.New("wrong parent object")
	}
	r.parent = p
	return nil
}

func (r *tenant) GetChildren() map[string]string {
	x := make(map[string]string)
	for k := range r.networkInstances {
		x[k] = "network-instance"
	}
	return x
}

func (r *tenant) UpdateConfig(d interface{}) error {
	r.Copy(d)
	if r.data.GetAdminState() == "disable" {
		r.data.SetStatus("down")
		r.data.SetReason("admin disabled")
	} else {
		r.data.SetStatus("up")
		r.data.SetReason("")
	}
	r.data.SetLastUpdate(time.Now().String())
	// update the state cache
	if err := r.UpdateStateCache(); err != nil {
		return err
	}
	return nil
}

func (r *tenant) GetPathElem(p []*gnmi.PathElem, do_recursive bool) ([]*gnmi.PathElem, error) {
	r.log.Debug("GetPathElem", "PathElem tenant", r.pathElem)
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

func (r *tenant) UpdateStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	r.log.Debug("Tenant Update State Cache", "PathElem", pe, "Prefix", r.prefix, "data", r.data)
	if err := updateCache(r.log, r.stateCache, r.prefix, &gnmi.Path{Elem: pe}, r.data); err != nil {
		r.log.Debug("Tenant Update State Cache Error")
		return err
	}
	r.log.Debug("Tenant Update State Cache ok")
	return nil
}

func (r *tenant) DeleteStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	if err := deleteCache(r.log, r.stateCache, r.prefix, &gnmi.Path{Elem: pe}); err != nil {
		return err
	}
	return nil
}

func (r *tenant) Copy(d interface{}) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	x := ipamv1alpha1.NddoipamIpamTenant{}
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}
	r.data = (&x).DeepCopy()
	return nil
}