package ipamlogic

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/pkg/errors"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-yang/pkg/cache"
	"github.com/yndd/ndd-yang/pkg/yentry"
	"github.com/yndd/nddo-ipam/internal/dispatcher"
)

type root struct {
	dispatcher.Resource
	data  interface{}
	ipams map[string]dispatcher.Handler
}

func (r *root) WithLogging(log logging.Logger) {
	r.Log = log
}

func (r *root) WithStateCache(c *cache.Cache) {
	r.StateCache = c
}

func (r *root) WithConfigCache(c *cache.Cache) {
	r.ConfigCache = c
}

func (r *root) WithPrefix(p *gnmi.Path) {
	r.Prefix = p
}

func (r *root) WithPathElem(pe []*gnmi.PathElem) {
	r.PathElem = pe[0]
}

func (r *root) WithRootSchema(rs yentry.Handler) {
	r.RootSchema = rs
}

func NewRoot(opts ...dispatcher.HandlerOption) dispatcher.Handler {
	r := &root{
		ipams: make(map[string]dispatcher.Handler),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *root) HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	log := r.Log.WithValues("Operation", o, "Path Elem", pe)

	log.Debug("root Handle")

	children := map[string]dispatcher.HandleConfigEventFunc{
		"ipam": ipamCreate,
	}

	// check path Element Name
	pathElemName := pe[0].GetName()
	if _, ok := children[pathElemName]; !ok {
		return nil, errors.Wrap(errors.New("unexpected pathElem"), fmt.Sprintf("ipam Handle: %s", pathElemName))
	}

	if len(pe) == 1 {
		log.Debug("root Handle pathelem =1")
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
		log.Debug("root Handle pathelem >1")
		i, err := r.CreateChild(children, pathElemName, prefix, pe[:1], nil)
		if err != nil {
			return nil, err
		}
		return i.HandleConfigEvent(o, prefix, pe[1:], d)
	}
	return nil, nil
}

func (r *root) CreateChild(children map[string]dispatcher.HandleConfigEventFunc, pathElemName string, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	switch pathElemName {
	case "ipam":
		if i, ok := r.ipams[ipamGetKey(pe)]; !ok {
			i = children[pathElemName](r.Log, r.ConfigCache, r.StateCache, prefix, pe, d)
			i.SetRootSchema(r.RootSchema)
			if err := i.SetParent(r); err != nil {
				return nil, err
			}
			r.ipams[ipamGetKey(pe)] = i
			return i, nil
		} else {
			return i, nil
		}
	}
	return nil, errors.New("CreateChild unexpected pathElemName in root")
}

func (r *root) DeleteChild(pathElemName string, pe []*gnmi.PathElem) error {
	switch pathElemName {
	case "ipam":
		if i, ok := r.ipams[ipamGetKey(pe)]; ok {
			if err := i.DeleteStateCache(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *root) SetParent(parent interface{}) error {
	// no SetParent required for root
	return nil
}

func (r *root) SetRootSchema(rs yentry.Handler) {
	r.RootSchema = rs
}

func (r *root) GetChildren() map[string]string {
	x := make(map[string]string)
	for k := range r.ipams {
		x[k] = "ipam"
	}
	return x
}

func (r *root) UpdateConfig(d interface{}) error {
	// no updates required for root
	return nil
}

func (r *root) GetPathElem(p []*gnmi.PathElem, do_recursive bool) ([]*gnmi.PathElem, error) {
	return nil, nil
}

func (r *root) UpdateStateCache() error {
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

func (r *root) DeleteStateCache() error {
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
