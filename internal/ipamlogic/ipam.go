package ipamlogic

import (
	"encoding/json"
	"fmt"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/pkg/errors"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-yang/pkg/cache"
	ipamv1alpha1 "github.com/yndd/nddo-ipam/apis/ipam/v1alpha1"
	"github.com/yndd/nddo-ipam/internal/dispatcher"
)

const (
	ipamDummyName = "dummy"
)

func init() {
	dispatcher.Register("ipam", []*dispatcher.EventHandler{
		{
			Kind: dispatcher.EventHandlerCreate,
			PathElem: []*gnmi.PathElem{
				{Name: "ipam"},
			},
			Handler: ipamCreate,
		},
	})
}

type Ipam interface {
	HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error)
}

type ipam struct {
	log         logging.Logger
	configCache *cache.Cache
	stateCache  *cache.Cache
	pathElem    *gnmi.PathElem
	prefix      *gnmi.Path
	key         string
	data        *ipamv1alpha1.NddoipamIpam
	parent      Root
	rirs        map[string]dispatcher.Handler
	tenants     map[string]dispatcher.Handler
}

type IpamOption func(*ipam)

// WithIpamIpamLogger initializes the logger.
func WithIpamLogger(log logging.Logger) IpamOption {
	return func(o *ipam) {
		o.log = log
	}
}

func WithIpamStateCache(c *cache.Cache) IpamOption {
	return func(o *ipam) {
		o.stateCache = c
	}
}

func WithIpamConfigCache(c *cache.Cache) IpamOption {
	return func(o *ipam) {
		o.configCache = c
	}
}

func WithIpamPrefix(p *gnmi.Path) IpamOption {
	return func(o *ipam) {
		o.prefix = p
	}
}

func WithIpamPathElem(pe []*gnmi.PathElem) IpamOption {
	return func(o *ipam) {
		o.pathElem = pe[0]
	}
}

func NewIpam(n string, opts ...IpamOption) dispatcher.Handler {
	x := &ipam{
		key:     n,
		rirs:    make(map[string]dispatcher.Handler),
		tenants: make(map[string]dispatcher.Handler),
	}

	for _, opt := range opts {
		opt(x)
	}
	return x
}

func ipamGetKey(p []*gnmi.PathElem) string {
	return ipamDummyName
}

func ipamCreate(log logging.Logger, cc, sc *cache.Cache, prefix *gnmi.Path, p []*gnmi.PathElem, d interface{}) dispatcher.Handler {
	ipamName := ipamGetKey(p)
	return NewIpam(ipamName,
		WithIpamPrefix(prefix),
		WithIpamPathElem(p),
		WithIpamLogger(log),
		WithIpamStateCache(sc),
		WithIpamConfigCache(cc))
}

func (r *ipam) HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	log := r.log.WithValues("Operation", o, "Path Elem", pe)

	log.Debug("ipam Handle")

	children := map[string]dispatcher.HandleConfigEventFunc{
		"rir":    rirCreate,
		"tenant": tenantCreate,
	}

	// check path Element Name
	pathElemName := pe[0].GetName()
	if _, ok := children[pathElemName]; !ok {
		return nil, errors.Wrap(errors.New("unexpected pathElem"), fmt.Sprintf("ipam Handle: %s", pathElemName))
	}

	if len(pe) == 1 {
		log.Debug("ipam Handle pathelem =1")
		// handle local
		switch o {
		case dispatcher.OperationUpdate:
			i, err := r.CreateChild(children, pathElemName, prefix, pe, d)
			if err != nil {
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
		case dispatcher.OperationDelete:
			if err := r.DeleteChild(pathElemName, pe); err != nil {
				return nil, err
			}
			return nil, nil
		}
	} else {
		log.Debug("ipam Handle pathelem >1")
		i, err := r.CreateChild(children, pathElemName, prefix, pe[:1], nil)
		if err != nil {
			return nil, err
		}
		return i.HandleConfigEvent(o, prefix, pe[1:], d)
	}
	return nil, nil
}

func (r *ipam) CreateChild(children map[string]dispatcher.HandleConfigEventFunc, pathElemName string, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	switch pathElemName {
	case "rir":
		if i, ok := r.rirs[rirGetKey(pe)]; !ok {
			i = children[pathElemName](r.log, r.configCache, r.stateCache, prefix, pe, d)
			if err := i.SetParent(r); err != nil {
				return nil, err
			}
			r.rirs[rirGetKey(pe)] = i
			return i, nil
		} else {
			return i, nil
		}
	case "tenant":
		if i, ok := r.tenants[tenantGetKey(pe)]; !ok {
			i = children[pathElemName](r.log, r.configCache, r.stateCache, prefix, pe, d)
			if err := i.SetParent(r); err != nil {
				return nil, err
			}
			r.tenants[tenantGetKey(pe)] = i
			return i, nil
		} else {
			return i, nil
		}
	}
	return nil, nil
}

func (r *ipam) DeleteChild(pathElemName string, pe []*gnmi.PathElem) error {
	switch pathElemName {
	case "rir":
		if i, ok := r.rirs[rirGetKey(pe)]; ok {
			if err := i.DeleteStateCache(); err != nil {
				return err
			}
		}
	case "tenant":
		if i, ok := r.tenants[tenantGetKey(pe)]; ok {
			if err := i.DeleteStateCache(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *ipam) SetParent(parent interface{}) error {
	p, ok := parent.(*root)
	if !ok {
		return errors.New("wrong parent object")
	}
	r.parent = p
	return nil
}

func (r *ipam) GetChildren() map[string]string {
	x := make(map[string]string)
	for k := range r.rirs {
		x[k] = "rir"
	}
	for k := range r.tenants {
		x[k] = "tenant"
	}
	return x
}

func (r *ipam) UpdateConfig(d interface{}) error {
	// no updates required for ipam
	return nil
}

func (r *ipam) GetPathElem(p []*gnmi.PathElem, do_recursive bool) ([]*gnmi.PathElem, error) {
	r.log.Debug("GetPathElem", "PathElem ipam", r.pathElem)
	return []*gnmi.PathElem{r.pathElem}, nil
}

func (r *ipam) UpdateStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	r.log.Debug("Rir Update Cache", "PathElem", pe, "Prefix", r.prefix, "data", r.data)
	if err := updateCache(r.log, r.stateCache, r.prefix, &gnmi.Path{Elem: pe}, r.data); err != nil {
		r.log.Debug("Rir Update Error")
		return err
	}
	r.log.Debug("Rir Update ok")
	return nil
}

func (r *ipam) DeleteStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	if err := deleteCache(r.log, r.stateCache, r.prefix, &gnmi.Path{Elem: pe}); err != nil {
		return err
	}
	return nil
}

func (r *ipam) Copy(d interface{}) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	x := ipamv1alpha1.NddoipamIpam{}
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}
	r.data = (&x).DeepCopy()
	return nil
}