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

type iprange struct {
	dispatcher.Resource
	data   *ipamv1alpha1.NddoipamIpamTenantNetworkInstanceIpRange
	parent *networkinstance
}

func (r *iprange) WithLogging(log logging.Logger) {
	r.Log = log
}

func (r *iprange) WithStateCache(c *cache.Cache) {
	r.StateCache = c
}

func (r *iprange) WithConfigCache(c *cache.Cache) {
	r.ConfigCache = c
}

func (r *iprange) WithPrefix(p *gnmi.Path) {
	r.Prefix = p
}

func (r *iprange) WithPathElem(pe []*gnmi.PathElem) {
	r.PathElem = pe[0]
}

func (r *iprange) WithRootSchema(rs *yentry.Entry) {
	r.RootSchema = rs
}

func NewIpRange(n string, opts ...dispatcher.HandlerOption) dispatcher.Handler {
	x := &iprange{}
	x.Key = n

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
		dispatcher.WithPrefix(prefix),
		dispatcher.WithPathElem(pe),
		dispatcher.WithLogging(log),
		dispatcher.WithStateCache(sc),
		dispatcher.WithConfigCache(cc))
}

func (r *iprange) HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	log := r.Log.WithValues("Operation", o, "Path Elem", pe)

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

func (r *iprange) SetRootSchema(rs *yentry.Entry) {
	r.RootSchema = rs
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
		p = append(p, r.PathElem)
		return p, nil
	}
	return nil, nil
}

/*
func (r *iprange) UpdateStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	r.Log.Debug("iprange Update Cache", "PathElem", pe, "Prefix", r.Prefix, "data", r.data)
	if err := updateCache(r.Log, r.StateCache, r.Prefix, &gnmi.Path{Elem: pe}, r.data); err != nil {
		r.Log.Debug("iprange Update Error")
		return err
	}
	r.Log.Debug("iprange Update ok")
	return nil
}

func (r *iprange) DeleteStateCache() error {
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

func (r *iprange) UpdateStateCache() error {
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

func (r *iprange) DeleteStateCache() error {
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
