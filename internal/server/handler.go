package server

import (
	"fmt"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/pkg/errors"
	"github.com/yndd/ndd-runtime/pkg/logging"
	ynddparser "github.com/yndd/ndd-yang/pkg/parser"
)

// Option can be used to manipulate Options.
type HandlerOption func(*Handler)

// WithLogger specifies how the Reconciler should log messages.
func WithHandlerLogger(log logging.Logger) HandlerOption {
	return func(h *Handler) {
		h.log = log
	}
}

func WithHandlerParser(log logging.Logger) HandlerOption {
	return func(h *Handler) {
		h.parser = ynddparser.NewParser(ynddparser.WithLogger(log))
	}
}

type Handler struct {
	log    logging.Logger
	parser *ynddparser.Parser
}

func NewHandler(opts ...HandlerOption) (*Handler, error) {
	h := &Handler{}

	for _, opt := range opts {
		opt(h)
	}

	return h, nil
}
func (h *Handler) Replace(d interface{}, update *gnmi.Update) (interface{}, error) {
	log := h.log.WithValues("Path", update.GetPath())
	log.Debug("ReplaceHandler", "Value", update.GetVal())
	val, err := h.parser.GetValue(update.GetVal())
	if err != nil {
		return nil, errors.Wrap(err, errGetValue)
	}
	tc := &ynddparser.TraceCtxtGnmi{
		Path:   update.GetPath(),
		Value:  val,
		Idx:    0,
		Msg:    make([]string, 0),
		Action: ynddparser.ConfigTreeActionCreate,
	}
	return h.parser.ParseTreeWithActionGnmi(d, tc, 0, 0), nil

}

func (h *Handler) Update(d interface{}, update *gnmi.Update) (interface{}, error) {
	log := h.log.WithValues("Path", update.GetPath())
	log.Debug("Update")
	val, err := h.parser.GetValue(update.GetVal())
	if err != nil {
		return nil, errors.Wrap(err, errGetValue)
	}
	tc := &ynddparser.TraceCtxtGnmi{
		Path:   update.GetPath(),
		Value:  val,
		Idx:    0,
		Msg:    make([]string, 0),
		Action: ynddparser.ConfigTreeActionUpdate,
	}
	return h.parser.ParseTreeWithActionGnmi(d, tc, 0, 0), nil
}

func (h *Handler) Delete(d interface{}, path *gnmi.Path) (interface{}, error) {
	log := h.log.WithValues("Path", path)
	log.Debug("Delete")
	tc := &ynddparser.TraceCtxtGnmi{
		Path:   path,
		Idx:    0,
		Msg:    make([]string, 0),
		Action: ynddparser.ConfigTreeActionDelete,
	}
	return h.parser.ParseTreeWithActionGnmi(d, tc, 0, 0), nil
}

func (h *Handler) Get(d interface{}, path *gnmi.Path) (interface{}, error) {
	log := h.log.WithValues("Path", path)
	log.Debug("Get")
	tc := &ynddparser.TraceCtxtGnmi{
		Path:   path,
		Idx:    0,
		Msg:    make([]string, 0),
		Action: ynddparser.ConfigTreeActionGet,
	}
	return h.parser.ParseTreeWithActionGnmi(d, tc, 0, 0), nil
}

// A ConditionKind represents a condition kind for a resource
type ResourceAction string

// Condition Kinds.
const (
	// replace
	ResourceActionReplace ResourceAction = "Replace"
	// update
	ResourceActionUpdate ResourceAction = "Update"
	// delete
	ResourceActionDelete ResourceAction = "Delete"
)

func (h *Handler) GetResources(a ResourceAction, upd []*gnmi.Update) (map[string]map[string]*resource, error) {
	resources := make(map[string]map[string]*resource)
	for _, u := range upd {
		log := h.log.WithValues("Path", u.GetPath().GetElem(), "Value", u.GetVal())
		log.Debug(fmt.Sprintf("GetResource%s", a))
		if len(u.GetPath().GetElem()) > 1 {
			switch {
			case u.GetPath().GetElem()[1].GetName() == "rir" ||
				u.GetPath().GetElem()[1].GetName() == "aggregate" ||
				u.GetPath().GetElem()[1].GetName() == "ip-prefix" ||
				u.GetPath().GetElem()[1].GetName() == "ip-range" ||
				u.GetPath().GetElem()[1].GetName() == "ip-address":

				if _, ok := resources[u.GetPath().GetElem()[1].GetName()]; !ok {
					resources[u.GetPath().GetElem()[1].GetName()] = make(map[string]*resource)
				}

				res := newResource(u.GetPath().GetElem()[1].GetKey()) // they key is at place 2 of the pathElem for all resources
				resKey := res.GetKeyString()

				if _, ok := resources[u.GetPath().GetElem()[1].GetName()][resKey]; !ok {
					// resource does not exist, initialize with the new data
					resources[u.GetPath().GetElem()[1].GetName()][resKey] = res
				}

				var err error
				switch a {
				case ResourceActionReplace:
					resources[u.GetPath().GetElem()[1].GetName()][resKey].data, err = h.Replace(resources[u.GetPath().GetElem()[1].GetName()][resKey].data, u)
					if err != nil {
						return nil, err
					}
				case ResourceActionUpdate:
					resources[u.GetPath().GetElem()[1].GetName()][resKey].data, err = h.Update(resources[u.GetPath().GetElem()[1].GetName()][resKey].data, u)
					if err != nil {
						return nil, err
					}
				}
				resources[u.GetPath().GetElem()[1].GetName()][resKey].data, err = h.Get(resources[u.GetPath().GetElem()[1].GetName()][resKey].data, u.GetPath())
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return resources, nil
}

func (h *Handler) GetResources2Delete(path []*gnmi.Path) (map[string]map[string]*resource, error) {
	resources := make(map[string]map[string]*resource)
	for _, p := range path {
		if len(p.GetElem()) == 1 {
			resources["ipam"] = make(map[string]*resource)
			resources["ipam"]["all"] = newResource(map[string]string{"name": "all"})
		}
		if len(p.GetElem()) > 1 {
			switch {
			case p.GetElem()[1].GetName() == "rir" ||
				p.GetElem()[1].GetName() == "aggregate" ||
				p.GetElem()[1].GetName() == "ip-prefix" ||
				p.GetElem()[1].GetName() == "ip-range" ||
				p.GetElem()[1].GetName() == "ip-address":

				if _, ok := resources[p.GetElem()[1].GetName()]; !ok {
					resources[p.GetElem()[1].GetName()] = make(map[string]*resource)
				}

				res := newResource(p.GetElem()[1].GetKey()) // they key is at place 2 of the pathElem for all resources
				resKey := res.GetKeyString()

				if _, ok := resources[p.GetElem()[1].GetName()][resKey]; !ok {
					// resource does not exist, initialize with the new data
					resources[p.GetElem()[1].GetName()][resKey] = res
				}
			}
		}
	}
	return resources, nil
}
