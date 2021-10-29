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
	dispatcher.Register("rir", []*dispatcher.EventHandler{
		{
			Kind: dispatcher.EventHandlerCreate,
			PathElem: []*gnmi.PathElem{
				{Name: "ipam"},
				{Name: "rir", Key: map[string]string{"name": "*"}},
			},
			Handler: rirCreate,
		},
	})
}

type Rir interface {
	HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error)
}

type rir struct {
	log         logging.Logger
	configCache *cache.Cache
	stateCache  *cache.Cache
	pathElem    *gnmi.PathElem
	prefix      *gnmi.Path
	key         string
	parent      *ipam
	data        *ipamv1alpha1.NddoipamIpamRir
}

type RirOption func(*rir)

// WithRirRirLogger initializes the logger.
func WithRirLogger(log logging.Logger) RirOption {
	return func(o *rir) {
		o.log = log
	}
}

// WithRirRirCache initializes the cache.
func WithRirStateCache(c *cache.Cache) RirOption {
	return func(o *rir) {
		o.stateCache = c
	}
}

func WithRirConfigCache(c *cache.Cache) RirOption {
	return func(o *rir) {
		o.configCache = c
	}
}

func WithRirPrefix(p *gnmi.Path) RirOption {
	return func(o *rir) {
		o.prefix = p
	}
}

func WithRirPathElem(pe []*gnmi.PathElem) RirOption {
	return func(o *rir) {
		o.pathElem = pe[0]
	}
}

func NewRir(n string, opts ...RirOption) *rir {
	x := &rir{
		key: n,
	}

	for _, opt := range opts {
		opt(x)
	}
	return x
}

func rirGetKey(pe []*gnmi.PathElem) string {
	return pe[0].GetKey()["name"]
}

func rirCreate(log logging.Logger, cc, sc *cache.Cache, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) dispatcher.Handler {
	rirName := rirGetKey(pe)
	return NewRir(rirName,
		WithRirPrefix(prefix),
		WithRirPathElem(pe),
		WithRirLogger(log),
		WithRirStateCache(sc),
		WithRirConfigCache(cc))
}

func (r *rir) HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	log := r.log.WithValues("Operation", o, "Path Elem", pe)

	log.Debug("rir Handle")

	if len(pe) == 1 {
		return nil, errors.New("the handle should have been terminated in the parent")
	} else {
		return nil, errors.New("there is no children in the rir resource ")
	}
}

func (r *rir) UpdateChild(children map[string]dispatcher.HandleConfigEventFunc, pathElemName string, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	return nil, errors.New("there is no children in the rir resource ")
}

func (r *rir) DeleteChild(pathElemName string, pe []*gnmi.PathElem) error {
	return errors.New("there is no children in the rir resource ")
}

func (r *rir) SetParent(parent interface{}) error {
	p, ok := parent.(*ipam)
	if !ok {
		return errors.New("wrong parent object")
	}
	r.parent = p
	return nil
}

func (r *rir) GetChildren() map[string]string {
	var x map[string]string
	return x
}

func (r *rir) UpdateConfig(d interface{}) error {
	r.Copy(d)
	switch r.data.GetName() {
	case "rfc1918", "rfc6598", "ula":
		r.data.SetPrivate(true)
	default:
		r.data.SetPrivate(false)
	}
	return nil
}

func (r *rir) GetPathElem(p []*gnmi.PathElem, do_recursive bool) ([]*gnmi.PathElem, error) {
	r.log.Debug("GetPathElem", "PathElem rir", r.pathElem)
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

func (r *rir) UpdateStateCache() error {
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

func (r *rir) DeleteStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	if err := deleteCache(r.log, r.stateCache, r.prefix, &gnmi.Path{Elem: pe}); err != nil {
		return err
	}
	return nil
}

func (r *rir) Copy(d interface{}) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	x := ipamv1alpha1.NddoipamIpamRir{}
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}
	r.data = (&x).DeepCopy()
	return nil
}
