package applogic

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/pkg/errors"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-yang/pkg/cache"
	"github.com/yndd/ndd-yang/pkg/dispatcher"
	"github.com/yndd/ndd-yang/pkg/yentry"
	"github.com/yndd/ndd-yang/pkg/yparser"
	ipamv1alpha1 "github.com/yndd/nddo-ipam/apis/ipam/v1alpha1"
	"inet.af/netaddr"
)

/*
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
*/

type ipprefix struct {
	dispatcher.Resource
	data   *ipamv1alpha1.NddoipamIpamTenantNetworkInstanceIpPrefix
	parent *networkinstance
}

func (r *ipprefix) WithLogging(log logging.Logger) {
	r.Log = log
}

func (r *ipprefix) WithStateCache(c *cache.Cache) {
	r.StateCache = c
}

func (r *ipprefix) WithConfigCache(c *cache.Cache) {
	r.ConfigCache = c
}

func (r *ipprefix) WithPrefix(p *gnmi.Path) {
	r.Prefix = p
}

func (r *ipprefix) WithPathElem(pe []*gnmi.PathElem) {
	r.PathElem = pe[0]
}

func (r *ipprefix) WithRootSchema(rs *yentry.Entry) {
	r.RootSchema = rs
}

func NewIpPrefix(n string, opts ...dispatcher.Option) dispatcher.Handler {
	x := &ipprefix{}
	x.Key = n

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
		dispatcher.WithPrefix(prefix),
		dispatcher.WithPathElem(pe),
		dispatcher.WithLogging(log),
		dispatcher.WithStateCache(sc),
		dispatcher.WithConfigCache(cc))
}

func (r *ipprefix) HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	log := r.Log.WithValues("Operation", o, "Path Elem", pe)

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

func (r *ipprefix) SetRootSchema(rs *yentry.Entry) {
	r.RootSchema = rs
}

func (r *ipprefix) GetChildren() map[string]string {
	var x map[string]string
	return x
}

func (r *ipprefix) GetPathElem(p []*gnmi.PathElem, do_recursive bool) ([]*gnmi.PathElem, error) {
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
func (r *ipprefix) UpdateStateCache() error {
	pe, err := r.GetPathElem(nil, true)
	if err != nil {
		return err
	}
	r.Log.Debug("ipprefix Update Cache", "PathElem", pe, "Prefix", r.Prefix, "data", r.data)
	if err := updateCache(r.Log, r.StateCache, r.Prefix, &gnmi.Path{Elem: pe}, r.data); err != nil {
		r.Log.Debug("ipprefix Update Error")
		return err
	}
	r.Log.Debug("ipprefix Update ok")
	return nil
}

func (r *ipprefix) DeleteStateCache() error {
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

func (r *ipprefix) UpdateStateCache() error {
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

func (r *ipprefix) DeleteStateCache() error {
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

func (r *ipprefix) GetTags() map[string]string {
	tags := make(map[string]string)
	for _, t := range r.data.Tag {
		tags[*t.DeepCopy().Key] = *t.DeepCopy().Value
	}
	return tags
}

func (r *ipprefix) UpdateConfig(d interface{}) error {
	r.Copy(d)
	if r.data.Prefix != nil {
		log := r.Log.WithValues("Prefix", *r.data.Prefix)
		log.Debug("UpdateConfig iprepfix")
		fmt.Printf("UpdateConfig iprepfix: %s\n", *r.data.Prefix)
	}

	p, err := netaddr.ParseIPPrefix(*r.data.Prefix)
	if err != nil {
		r.Log.Debug("UpdateConfig ParseIPPrefix", "Error", err)
		return err
	}

	if err := r.parent.AddRoute(p, r.GetTags()); err != nil {
		r.data.SetStatus("down")
		r.data.SetReason(fmt.Sprintf("%v", err))
		fmt.Printf("UpdateConfig iprepfix AddRoute error %v\n", err)
		//log.Debug("UpdateConfig AddRoute", "Error", err)
	} else {
		r.data.SetStatus("up")
		r.data.SetReason("")
	}
	r.data.SetLastUpdate(time.Now().String())

	if err := r.UpdateStateCache(); err != nil {
		return err
	}
	return nil
}
