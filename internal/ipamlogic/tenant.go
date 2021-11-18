package ipamlogic

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-yang/pkg/cache"
	"github.com/yndd/ndd-yang/pkg/yentry"
	"github.com/yndd/ndd-yang/pkg/yparser"
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

//type Tenant interface {
//	HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error)
//}

type tenant struct {
	dispatcher.Resource
	//log         logging.Logger
	//configCache *cache.Cache
	//stateCache  *cache.Cache
	//pathElem    *gnmi.PathElem
	//prefix      *gnmi.Path
	//key         string
	data             *ipamv1alpha1.NddoipamIpamTenant
	parent           *ipam
	networkInstances map[string]dispatcher.Handler
}

func (r *tenant) WithLogging(log logging.Logger) {
	r.Log = log
}

func (r *tenant) WithStateCache(c *cache.Cache) {
	r.StateCache = c
}

func (r *tenant) WithConfigCache(c *cache.Cache) {
	r.ConfigCache = c
}

func (r *tenant) WithPrefix(p *gnmi.Path) {
	r.Prefix = p
}

func (r *tenant) WithPathElem(pe []*gnmi.PathElem) {
	r.PathElem = pe[0]
}

func (r *tenant) WithRootSchema(rs *yentry.Entry) {
	r.RootSchema = rs
}

func NewTenant(n string, opts ...dispatcher.HandlerOption) dispatcher.Handler {
	x := &tenant{
		networkInstances: make(map[string]dispatcher.Handler),
	}
	x.Key = n

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
		dispatcher.WithPrefix(prefix),
		dispatcher.WithPathElem(pe),
		dispatcher.WithLogging(log),
		dispatcher.WithStateCache(sc),
		dispatcher.WithConfigCache(cc))
}

func (r *tenant) HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	log := r.Log.WithValues("Operation", o, "Path Elem", pe)

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

func (r *tenant) CreateChild(children map[string]dispatcher.HandleConfigEventFunc, pathElemName string, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	switch pathElemName {
	case "network-instance":
		if i, ok := r.networkInstances[rirGetKey(pe)]; !ok {
			i = children[pathElemName](r.Log, r.ConfigCache, r.StateCache, prefix, pe, d)
			i.SetRootSchema(r.RootSchema)
			if err := i.SetParent(r); err != nil {
				return nil, err
			}
			r.networkInstances[rirGetKey(pe)] = i
			return i, nil
		} else {
			return i, nil
		}
	}
	return nil, nil
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

func (r *tenant) SetRootSchema(rs *yentry.Entry) {
	r.RootSchema = rs
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
	r.Log.Debug("GetPathElem", "PathElem tenant", r.PathElem)
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

/*
func (r *tenant) UpdateStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	r.Log.Debug("Tenant Update State Cache", "PathElem", pe, "Prefix", r.Prefix, "data", r.data)
	if err := updateCache(r.Log, r.StateCache, r.Prefix, &gnmi.Path{Elem: pe}, r.data); err != nil {
		r.Log.Debug("Tenant Update State Cache Error")
		return err
	}
	r.Log.Debug("Tenant Update State Cache ok")
	return nil
}

func (r *tenant) DeleteStateCache() error {
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

func (r *tenant) UpdateStateCache() error {
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

func (r *tenant) DeleteStateCache() error {
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
