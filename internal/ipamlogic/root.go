package ipamlogic

import (
	"fmt"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/pkg/errors"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-yang/pkg/cache"
	"github.com/yndd/nddo-ipam/internal/dispatcher"
)

type Root interface {
	HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error)
}

type root struct {
	log         logging.Logger
	configCache *cache.Cache
	stateCache  *cache.Cache
	ipams       map[string]dispatcher.Handler
}

type RootOption func(*root)

// WithRootLogger initializes the logger.
func WithRootLogger(log logging.Logger) RootOption {
	return func(o *root) {
		o.log = log
	}
}

// WithRootStateCache initializes the state cache.
func WithRootStateCache(c *cache.Cache) RootOption {
	return func(o *root) {
		o.stateCache = c
	}
}

// WithRootConfigCache initializes the config cache.
func WithRootConfigCache(c *cache.Cache) RootOption {
	return func(o *root) {
		o.configCache = c
	}
}

func NewRoot(opts ...RootOption) (Root, error) {
	r := &root{
		ipams: make(map[string]dispatcher.Handler),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r, nil
}

func (r *root) HandleConfigEvent(o dispatcher.Operation, prefix *gnmi.Path, pe []*gnmi.PathElem, d interface{}) (dispatcher.Handler, error) {
	log := r.log.WithValues("Operation", o, "Path Elem", pe)

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
			i = children[pathElemName](r.log, r.configCache, r.stateCache, prefix, pe, d)
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
