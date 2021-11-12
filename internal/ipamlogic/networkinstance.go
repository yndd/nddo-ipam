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
			},
			Handler: networkinstanceCreate,
		},
	})
}

type networkinstance struct {
	dispatcher.Resource
	data        *ipamv1alpha1.NddoipamIpamTenantNetworkInstance
	parent      *tenant
	ipprefixes  map[string]dispatcher.Handler
	ipranges    map[string]dispatcher.Handler
	ipaddresses map[string]dispatcher.Handler
}

func (r *networkinstance) WithLogging(log logging.Logger) {
	r.Log = log
}

func (r *networkinstance) WithStateCache(c *cache.Cache) {
	r.StateCache = c
}

func (r *networkinstance) WithConfigCache(c *cache.Cache) {
	r.ConfigCache = c
}

func (r *networkinstance) WithPrefix(p *gnmi.Path) {
	r.Prefix = p
}

func (r *networkinstance) WithPathElem(pe []*gnmi.PathElem) {
	r.PathElem = pe[0]
}

func (r *networkinstance) WithRootSchema(rs yentry.Handler) {
	r.RootSchema = rs
}

func NewNetworkInstance(n string, opts ...dispatcher.HandlerOption) dispatcher.Handler {
	x := &networkinstance{
		ipprefixes:  make(map[string]dispatcher.Handler),
		ipranges:    make(map[string]dispatcher.Handler),
		ipaddresses: make(map[string]dispatcher.Handler),
	}
	x.Key = n

	for _, opt := range opts {
		opt(x)
	}
	return x
}

func networkinstanceGetKey(pe []*gnmi.PathElem) string {
	return pe[0].GetKey()["name"]
}

func networkinstanceCreate(log logging.Logger, cc, sc *cache.Cache, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) dispatcher.Handler {
	networkinstanceName := networkinstanceGetKey(pe)
	return NewNetworkInstance(networkinstanceName,
		dispatcher.WithPrefix(prefix),
		dispatcher.WithPathElem(pe),
		dispatcher.WithLogging(log),
		dispatcher.WithStateCache(sc),
		dispatcher.WithConfigCache(cc))
}

func (r *networkinstance) HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	log := r.Log.WithValues("Operation", o, "Path Elem", pe)

	log.Debug("networkinstance HandleConfigEvent")

	children := map[string]dispatcher.HandleConfigEventFunc{
		"ip-prefix":  ipprefixCreate,
		"ip-range":   iprangeCreate,
		"ip-address": ipaddressCreate,
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

func (r *networkinstance) CreateChild(children map[string]dispatcher.HandleConfigEventFunc, pathElemName string, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	switch pathElemName {
	case "ip-prefix":
		if i, ok := r.ipprefixes[ipprefixGetKey(pe)]; !ok {
			i = children[pathElemName](r.Log, r.ConfigCache, r.StateCache, prefix, pe, d)
			i.SetRootSchema(r.RootSchema)
			if err := i.SetParent(r); err != nil {
				return nil, err
			}
			r.ipprefixes[ipprefixGetKey(pe)] = i
			return i, nil
		} else {
			return i, nil
		}
	case "ip-range":
		if i, ok := r.ipranges[iprangeGetKey(pe)]; !ok {
			i = children[pathElemName](r.Log, r.ConfigCache, r.StateCache, prefix, pe, d)
			i.SetRootSchema(r.RootSchema)
			if err := i.SetParent(r); err != nil {
				return nil, err
			}
			r.ipranges[iprangeGetKey(pe)] = i
			return i, nil
		} else {
			return i, nil
		}
	case "ip-address":
		if i, ok := r.ipaddresses[ipaddressGetKey(pe)]; !ok {
			i = children[pathElemName](r.Log, r.ConfigCache, r.StateCache, prefix, pe, d)
			i.SetRootSchema(r.RootSchema)
			if err := i.SetParent(r); err != nil {
				return nil, err
			}
			r.ipaddresses[ipaddressGetKey(pe)] = i
			return i, nil
		} else {
			return i, nil
		}
	}
	return nil, nil
}

func (r *networkinstance) DeleteChild(pathElemName string, pe []*gnmi.PathElem) error {
	switch pathElemName {
	case "ip-prefix":
		if i, ok := r.ipprefixes[ipprefixGetKey(pe)]; ok {
			if err := i.DeleteStateCache(); err != nil {
				return err
			}
		}
	case "ip-range":
		if i, ok := r.ipranges[iprangeGetKey(pe)]; ok {
			if err := i.DeleteStateCache(); err != nil {
				return err
			}
		}
	case "ip-address":
		if i, ok := r.ipaddresses[ipaddressGetKey(pe)]; ok {
			if err := i.DeleteStateCache(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *networkinstance) SetParent(parent interface{}) error {
	p, ok := parent.(*tenant)
	if !ok {
		return errors.New("wrong parent object")
	}
	r.parent = p
	return nil
}

func (r *networkinstance) SetRootSchema(rs yentry.Handler) {
	r.RootSchema = rs
}

func (r *networkinstance) GetChildren() map[string]string {
	x := make(map[string]string)
	for k := range r.ipprefixes {
		x[k] = "ip-prefix"
	}
	for k := range r.ipranges {
		x[k] = "ip-range"
	}
	for k := range r.ipaddresses {
		x[k] = "ip-address"
	}
	return x
}

func (r *networkinstance) UpdateConfig(d interface{}) error {
	r.Copy(d)
	if r.parent.data.GetStatus() == "down" {
		r.data.SetStatus("down")
		r.data.SetReason("parent status disabled")
	} else {
		if r.data.GetAdminState() == "down" {
			r.data.SetStatus("down")
			r.data.SetReason("admin disable")
		} else {
			r.data.SetStatus("up")
			r.data.SetReason("")
		}
	}
	r.data.SetLastUpdate(time.Now().String())
	// update the state cache
	if err := r.UpdateStateCache(); err != nil {
		return err
	}
	return nil
}

func (r *networkinstance) GetPathElem(p []*gnmi.PathElem, do_recursive bool) ([]*gnmi.PathElem, error) {
	r.Log.Debug("GetPathElem", "PathElem networkinstance", r.PathElem)
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
func (r *networkinstance) UpdateStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	r.Log.Debug("NetworkInstance Update Cache", "PathElem", pe, "Prefix", r.Prefix, "data", r.data)
	if err := updateCache(r.Log, r.StateCache, r.Prefix, &gnmi.Path{Elem: pe}, r.data); err != nil {
		r.Log.Debug("NetworkInstance Update Error")
		return err
	}
	r.Log.Debug("NetworkInstance Update ok")
	return nil
}

func (r *networkinstance) DeleteStateCache() error {
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

func (r *networkinstance) Copy(d interface{}) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	x := ipamv1alpha1.NddoipamIpamTenantNetworkInstance{}
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}
	r.data = (&x).DeepCopy()
	return nil
}

func (r *networkinstance) UpdateStateCache() error {
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
	n, err := r.StateCache.GetNotificationFromJSON2(r.Prefix, &gnmi.Path{Elem: pe}, x, r.RootSchema)
	if err != nil {
		return err
	}

	//printNotification(log, n)
	if n != nil {
		if err := r.StateCache.GnmiUpdate(r.Prefix.Target, n); err != nil {
			if strings.Contains(fmt.Sprintf("%v", err), "stale") {
				return nil
			}
			return err
		}
	}
	return nil
}

func (r *networkinstance) DeleteStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	n, err := r.StateCache.GetNotificationFromDelete(r.Prefix, &gnmi.Path{Elem: pe})
	if err != nil {
		return err
	}
	if err := r.StateCache.GnmiUpdate(r.Prefix.Target, n); err != nil {
		return err
	}

	return nil
}
