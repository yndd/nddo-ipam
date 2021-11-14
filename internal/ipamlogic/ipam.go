package ipamlogic

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/pkg/errors"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-yang/pkg/cache"
	"github.com/yndd/ndd-yang/pkg/yentry"
	"github.com/yndd/ndd-yang/pkg/yparser"
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

type ipam struct {
	dispatcher.Resource
	data    *ipamv1alpha1.NddoipamIpam
	parent  *root
	rirs    map[string]dispatcher.Handler
	tenants map[string]dispatcher.Handler
}

func (r *ipam) WithLogging(log logging.Logger) {
	r.Log = log
}

func (r *ipam) WithStateCache(c *cache.Cache) {
	r.StateCache = c
}

func (r *ipam) WithConfigCache(c *cache.Cache) {
	r.ConfigCache = c
}

func (r *ipam) WithPrefix(p *gnmi.Path) {
	r.Prefix = p
}

func (r *ipam) WithPathElem(pe []*gnmi.PathElem) {
	r.PathElem = pe[0]
}

func (r *ipam) WithRootSchema(rs yentry.Handler) {
	r.RootSchema = rs
}

func NewIpam(n string, opts ...dispatcher.HandlerOption) dispatcher.Handler {
	x := &ipam{
		rirs:    make(map[string]dispatcher.Handler),
		tenants: make(map[string]dispatcher.Handler),
	}
	x.Key = n

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
		dispatcher.WithPrefix(prefix),
		dispatcher.WithPathElem(p),
		dispatcher.WithLogging(log),
		dispatcher.WithStateCache(sc),
		dispatcher.WithConfigCache(cc))
}

func (r *ipam) HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	log := r.Log.WithValues("Operation", o, "Path Elem", pe)

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
			i = children[pathElemName](r.Log, r.ConfigCache, r.StateCache, prefix, pe, d)
			i.SetRootSchema(r.RootSchema)
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
			i = children[pathElemName](r.Log, r.ConfigCache, r.StateCache, prefix, pe, d)
			i.SetRootSchema(r.RootSchema)
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

func (r *ipam) SetRootSchema(rs yentry.Handler) {
	r.RootSchema = rs
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
	r.Log.Debug("GetPathElem", "PathElem ipam", r.PathElem)
	return []*gnmi.PathElem{r.PathElem}, nil
}

/*
func (r *ipam) UpdateStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	r.Log.Debug("Rir Update Cache", "PathElem", pe, "Prefix", r.Prefix, "data", r.data)
	if err := updateCache(r.Log, r.StateCache, r.Prefix, &gnmi.Path{Elem: pe}, r.data); err != nil {
		r.Log.Debug("Rir Update Error")
		return err
	}
	r.Log.Debug("Rir Update ok")
	return nil
}

func (r *ipam) DeleteStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	if err := deleteCache(r.Log, r.StateCache, r.Prefix, &gnmi.Path{Elem: pe}); err != nil {
		return err
	}
	return nil
}
*/

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

func (r *ipam) UpdateStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	b, err := json.Marshal(r.data)
	if err != nil {
		return err
	}
	var x interface{}
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}
	//log.Debug("Debug updateState", "refPaths", refPaths)
	r.Log.Debug("Debug updateState", "data", x)
	u, err := yparser.GetGranularUpdatesFromJSON(&gnmi.Path{Elem: pe}, x, r.RootSchema)
	n := &gnmi.Notification{
		Timestamp: time.Now().UnixNano(),
		Prefix:    r.Prefix,
		Update:    u,
	}
	//n, err := r.StateCache.GetNotificationFromJSON2(r.Prefix, &gnmi.Path{Elem: pe}, x, r.RootSchema)
	if err != nil {
		return err
	}
	if u != nil {
		if err := r.StateCache.GnmiUpdate(r.Prefix.Target, n); err != nil {
			if strings.Contains(fmt.Sprintf("%v", err), "stale") {
				return nil
			}
			return err
		}
	}
	return nil
}

func (r *ipam) DeleteStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	n := &gnmi.Notification{
		Timestamp: time.Now().UnixNano(),
		Prefix:    r.Prefix,
		Delete:    []*gnmi.Path{{Elem: pe}},
	}
	if err := r.StateCache.GnmiUpdate(r.Prefix.Target, n); err != nil {
		return err
	}
	return nil
}
