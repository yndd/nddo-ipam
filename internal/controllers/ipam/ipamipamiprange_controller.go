/*
Copyright 2021 NDDO.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ipam

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/karimra/gnmic/target"
	gnmitypes "github.com/karimra/gnmic/types"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/pkg/errors"
	ndrv1 "github.com/yndd/ndd-core/apis/dvr/v1"
	pkgmetav1 "github.com/yndd/ndd-core/apis/pkg/meta/v1"
	"github.com/yndd/ndd-runtime/pkg/event"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-runtime/pkg/reconciler/managed"
	"github.com/yndd/ndd-runtime/pkg/resource"
	"github.com/yndd/ndd-runtime/pkg/utils"
	"github.com/yndd/ndd-yang/pkg/parser"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	cevent "sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	ipamv1alpha1 "github.com/yndd/nddo-ipam/apis/ipam/v1alpha1"
)

const (
	// Errors
	errUnexpectedIpamIprange       = "the managed resource is not a IpamIprange resource"
	errKubeUpdateFailedIpamIprange = "cannot update IpamIprange"
	errReadIpamIprange             = "cannot read IpamIprange"
	errCreateIpamIprange           = "cannot create IpamIprange"
	erreUpdateIpamIprange          = "cannot update IpamIprange"
	errDeleteIpamIprange           = "cannot delete IpamIprange"

	// resource information
	// resourcePrefixIpamIprange = "ipam.nddo.yndd.io.v1alpha1.IpamIprange"
)

var resourceRefPathsIpamIprange = []*gnmi.Path{
	{
		Elem: []*gnmi.PathElem{
			{Name: "ip-range", Key: map[string]string{
				"end":              "",
				"network-instance": "",
				"start":            "",
				"tenant":           "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "ip-range", Key: map[string]string{
				"end":              "",
				"network-instance": "",
				"start":            "",
				"tenant":           "",
			}},
			{Name: "aggregate", Key: map[string]string{
				"network-instance": "",
				"prefix":           "",
				"tenant":           "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "ip-range", Key: map[string]string{
				"end":              "",
				"network-instance": "",
				"start":            "",
				"tenant":           "",
			}},
			{Name: "ip-prefix", Key: map[string]string{
				"network-instance": "",
				"prefix":           "",
				"tenant":           "",
			}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "ip-range", Key: map[string]string{
				"end":              "",
				"network-instance": "",
				"start":            "",
				"tenant":           "",
			}},
			{Name: "tag", Key: map[string]string{
				"key": "",
			}},
		},
	},
}
var localleafRefIpamIprange = []*parser.LeafRefGnmi{}
var externalLeafRefIpamIprange = []*parser.LeafRefGnmi{}

// SetupIpamIprange adds a controller that reconciles IpamIpranges.
func SetupIpamIprange(mgr ctrl.Manager, o controller.Options, l logging.Logger, poll time.Duration, namespace string) (string, chan cevent.GenericEvent, error) {

	name := managed.ControllerName(ipamv1alpha1.IpamIprangeGroupKind)

	events := make(chan cevent.GenericEvent)

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(ipamv1alpha1.IpamIprangeGroupVersionKind),
		managed.WithExternalConnecter(&connectorIpamIprange{
			log:         l,
			kube:        mgr.GetClient(),
			usage:       resource.NewNetworkNodeUsageTracker(mgr.GetClient(), &ndrv1.NetworkNodeUsage{}),
			newClientFn: target.NewTarget},
		),
		managed.WithParser(l),
		managed.WithValidator(&validatorIpamIprange{log: l, parser: *parser.NewParser(parser.WithLogger(l))}),
		managed.WithLogger(l.WithValues("controller", name)),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))))

	return ipamv1alpha1.IpamIprangeGroupKind, events, ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o).
		For(&ipamv1alpha1.IpamIpamIprange{}).
		WithEventFilter(resource.IgnoreUpdateWithoutGenerationChangePredicate()).
		Watches(
			&source.Channel{Source: events},
			&handler.EnqueueRequestForObject{},
		).
		Complete(r)
}

type validatorIpamIprange struct {
	log    logging.Logger
	parser parser.Parser
}

func (v *validatorIpamIprange) ValidateLocalleafRef(ctx context.Context, mg resource.Managed) (managed.ValidateLocalleafRefObservation, error) {
	log := v.log.WithValues("resource", mg.GetName())
	log.Debug("ValidateLocalleafRef...")

	// json unmarshal the resource
	o, ok := mg.(*ipamv1alpha1.IpamIpamIprange)
	if !ok {
		return managed.ValidateLocalleafRefObservation{}, errors.New(errUnexpectedIpamIprange)
	}
	d, err := json.Marshal(&o.Spec.ForNetworkNode)
	if err != nil {
		return managed.ValidateLocalleafRefObservation{}, errors.Wrap(err, errJSONMarshal)
	}
	var x1 interface{}
	json.Unmarshal(d, &x1)

	// For local leafref validation we dont need to supply the external data so we use nil
	success, resultleafRefValidation, err := v.parser.ValidateLeafRefGnmi(
		parser.LeafRefValidationLocal, x1, nil, localleafRefIpamIprange, log)
	if err != nil {
		return managed.ValidateLocalleafRefObservation{
			Success: false,
		}, nil
	}
	if !success {
		log.Debug("ValidateLocalleafRef failed", "resultleafRefValidation", resultleafRefValidation)
		return managed.ValidateLocalleafRefObservation{
			Success:          false,
			ResolvedLeafRefs: resultleafRefValidation}, nil
	}
	log.Debug("ValidateLocalleafRef success", "resultleafRefValidation", resultleafRefValidation)
	return managed.ValidateLocalleafRefObservation{
		Success:          true,
		ResolvedLeafRefs: resultleafRefValidation}, nil
}

func (v *validatorIpamIprange) ValidateExternalleafRef(ctx context.Context, mg resource.Managed, cfg []byte) (managed.ValidateExternalleafRefObservation, error) {
	log := v.log.WithValues("resource", mg.GetName())
	log.Debug("ValidateExternalleafRef...")

	// json unmarshal the resource
	o, ok := mg.(*ipamv1alpha1.IpamIpamIprange)
	if !ok {
		return managed.ValidateExternalleafRefObservation{}, errors.New(errUnexpectedIpamIprange)
	}
	d, err := json.Marshal(&o.Spec.ForNetworkNode)
	if err != nil {
		return managed.ValidateExternalleafRefObservation{}, errors.Wrap(err, errJSONMarshal)
	}
	var x1 interface{}
	json.Unmarshal(d, &x1)

	// json unmarshal the external data
	var x2 interface{}
	json.Unmarshal(cfg, &x2)

	// For local external leafref validation we need to supply the external
	// data to validate the remote leafref, we use x2 for this
	success, resultleafRefValidation, err := v.parser.ValidateLeafRefGnmi(
		parser.LeafRefValidationExternal, x1, x2, externalLeafRefIpamIprange, log)
	if err != nil {
		return managed.ValidateExternalleafRefObservation{
			Success: false,
		}, nil
	}
	if !success {
		log.Debug("ValidateExternalleafRef failed", "resultleafRefValidation", resultleafRefValidation)
		return managed.ValidateExternalleafRefObservation{
			Success:          false,
			ResolvedLeafRefs: resultleafRefValidation}, nil
	}
	log.Debug("ValidateExternalleafRef success", "resultleafRefValidation", resultleafRefValidation)
	return managed.ValidateExternalleafRefObservation{
		Success:          true,
		ResolvedLeafRefs: resultleafRefValidation}, nil
}

func (v *validatorIpamIprange) ValidateParentDependency(ctx context.Context, mg resource.Managed, cfg []byte) (managed.ValidateParentDependencyObservation, error) {
	log := v.log.WithValues("resource", mg.GetName())
	log.Debug("ValidateParentDependency...")

	// we initialize a global list for finer information on the resolution
	resultleafRefValidation := make([]*parser.ResolvedLeafRefGnmi, 0)
	log.Debug("ValidateParentDependency success", "resultParentValidation", resultleafRefValidation)
	return managed.ValidateParentDependencyObservation{
		Success:          true,
		ResolvedLeafRefs: resultleafRefValidation}, nil
}

// ValidateResourceIndexes validates if the indexes of a resource got changed
// if so we need to delete the original resource, because it will be dangling if we dont delete it
func (v *validatorIpamIprange) ValidateResourceIndexes(ctx context.Context, mg resource.Managed) (managed.ValidateResourceIndexesObservation, error) {
	log := v.log.WithValues("resource", mg.GetName())

	// json unmarshal the resource
	o, ok := mg.(*ipamv1alpha1.IpamIpamIprange)
	if !ok {
		return managed.ValidateResourceIndexesObservation{}, errors.New(errUnexpectedIpamIprange)
	}
	log.Debug("ValidateResourceIndexes", "Spec", o.Spec)

	rootPath := []*gnmi.Path{
		{
			Elem: []*gnmi.PathElem{
				{Name: "ipam"},
				{Name: "ip-range", Key: map[string]string{
					"tenant":           *o.Spec.ForNetworkNode.IpamIpamIprange.Tenant,
					"network-instance": *o.Spec.ForNetworkNode.IpamIpamIprange.NetworkInstance,
					"start":            *o.Spec.ForNetworkNode.IpamIpamIprange.Start,
					"end":              *o.Spec.ForNetworkNode.IpamIpamIprange.End,
				}},
			},
		},
	}

	origResourceIndex := mg.GetResourceIndexes()
	// we call the CompareConfigPathsWithResourceKeys irrespective is the get resource index returns nil
	changed, deletPaths, newResourceIndex := v.parser.CompareGnmiPathsWithResourceKeys(rootPath[0], origResourceIndex)
	if changed {
		log.Debug("ValidateResourceIndexes changed", "deletPaths", deletPaths[0])
		return managed.ValidateResourceIndexesObservation{Changed: true, ResourceDeletes: deletPaths, ResourceIndexes: newResourceIndex}, nil
	}

	log.Debug("ValidateResourceIndexes success")
	return managed.ValidateResourceIndexesObservation{Changed: false, ResourceIndexes: newResourceIndex}, nil
}

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type connectorIpamIprange struct {
	log         logging.Logger
	kube        client.Client
	usage       resource.Tracker
	newClientFn func(c *gnmitypes.TargetConfig) *target.Target
	//newClientFn func(ctx context.Context, cfg ndd.Config) (config.ConfigurationClient, error)
}

// Connect produces an ExternalClient by:
// 1. Tracking that the managed resource is using a NetworkNode.
// 2. Getting the managed resource's NetworkNode with connection details
// A resource is mapped to a single target
func (c *connectorIpamIprange) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	log := c.log.WithValues("resource", mg.GetName())
	log.Debug("Connect")

	cfg := &gnmitypes.TargetConfig{
		Name:       "dummy",
		Address:    pkgmetav1.PrefixService + "-" + "ipam" + "." + pkgmetav1.NamespaceLocalK8sDNS + strconv.Itoa((pkgmetav1.GnmiServerPort)),
		Username:   utils.StringPtr("admin"),
		Password:   utils.StringPtr("admin"),
		Timeout:    10 * time.Second,
		SkipVerify: utils.BoolPtr(true),
		Insecure:   utils.BoolPtr(true),
		TLSCA:      utils.StringPtr(""), //TODO TLS
		TLSCert:    utils.StringPtr(""), //TODO TLS
		TLSKey:     utils.StringPtr(""),
		Gzip:       utils.BoolPtr(false),
	}

	cl := target.NewTarget(cfg)
	if err := cl.CreateGNMIClient(ctx); err != nil {
		return nil, errors.Wrap(err, errNewClient)
	}

	// we make a string here since we use a trick in registration to go to multiple targets
	// while here the object is mapped to a single target/network node
	tns := []string{"localGNMIServer"}

	return &externalIpamIprange{client: cl, targets: tns, log: log, parser: *parser.NewParser(parser.WithLogger(log))}, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type externalIpamIprange struct {
	//client  config.ConfigurationClient
	client  *target.Target
	targets []string
	log     logging.Logger
	parser  parser.Parser
}

func (e *externalIpamIprange) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	o, ok := mg.(*ipamv1alpha1.IpamIpamIprange)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedIpamIprange)
	}
	log := e.log.WithValues("Resource", o.GetName())
	log.Debug("Observing ...")

	// rootpath of the resource
	rootPath := []*gnmi.Path{
		{
			Elem: []*gnmi.PathElem{
				{Name: "ipam"},
				{Name: "ip-range", Key: map[string]string{
					"tenant":           *o.Spec.ForNetworkNode.IpamIpamIprange.Tenant,
					"network-instance": *o.Spec.ForNetworkNode.IpamIpamIprange.NetworkInstance,
					"start":            *o.Spec.ForNetworkNode.IpamIpamIprange.Start,
					"end":              *o.Spec.ForNetworkNode.IpamIpamIprange.End,
				}},
			},
		},
	}

	// gnmi get request
	req := &gnmi.GetRequest{
		Path:     rootPath,
		Encoding: gnmi.Encoding_JSON,
	}

	// gnmi get response
	resp, err := e.client.Get(ctx, req)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errReadIpamIprange)
	}

	// prepare the input data to compare against the response data
	d, err := json.Marshal(&o.Spec.ForNetworkNode)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errJSONMarshal)
	}
	var x1 interface{}
	if err := json.Unmarshal(d, &x1); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errJSONUnMarshal)
	}

	// remove the hierarchical elements for data processing, comparison, etc
	// they are used in the provider for parent dependency resolution
	// but are not relevant in the data, they are referenced in the rootPath
	// when interacting with the device driver
	//hids := make([]string, 0)
	//
	//
	//
	//x1 = e.parser.RemoveLeafsFromJSONData(x1, hids)

	// validate gnmi resp information
	var exists bool
	var x2 interface{}
	if len(resp.GetNotification()) != 0 {
		if len(resp.GetNotification()[0].GetUpdate()) != 0 {
			// get value from gnmi get response
			x2, err = e.parser.GetValue(resp.GetNotification()[0].GetUpdate()[0].Val)
			if err != nil {
				log.Debug("Observe response get value issue")
				return managed.ExternalObservation{}, errors.Wrap(err, errJSONMarshal)
			}
			switch x := x2.(type) {
			case map[string]interface{}:
				if x, ok := x["ip-range"]; ok {
					if x != nil {
						exists = true
					}
				}
			}
		}
	}

	// logging information that will be used to provide the response
	log.Debug("Spec Data", "X1", x1)
	log.Debug("Resp Data", "X2", x2)

	// if the cache is not ready we back off and return
	if !exists {
		// Resource Does not Exists
		log.Debug("Observing Response:", "Exists", false, "HasData", false, "UpToDate", false, "Response", resp)
		return managed.ExternalObservation{
			Ready:            true,
			ResourceExists:   false,
			ResourceHasData:  false,
			ResourceUpToDate: false,
		}, nil
	}

	// data is present
	// for lists with keys we need to create a list before calulating the paths since this is what
	// the object eventually happens to be based upon. We avoid having multiple entries in a list object
	// and hence we have to add this step
	x1, err = e.parser.AddJSONDataToList(x1)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errWrongInputdata)
	}

	updatesx1 := e.parser.GetUpdatesFromJSONDataGnmi(rootPath[0], e.parser.XpathToGnmiPath("/", 0), x1, resourceRefPathsIpamIprange)
	for _, update := range updatesx1 {
		log.Debug("Observe Fine Grane Updates X1", "Path", e.parser.GnmiPathToXPath(update.Path, true), "Value", update.GetVal())
	}
	// for lists with keys we need to create a list before calulating the paths since this is what
	// the object eventually happens to be based upon. We avoid having multiple entries in a list object
	// and hence we have to add this step
	x2, err = e.parser.AddJSONDataToList(x2)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errWrongInputdata)
	}
	updatesx2 := e.parser.GetUpdatesFromJSONDataGnmi(rootPath[0], e.parser.XpathToGnmiPath("/", 0), x2, resourceRefPathsIpamIprange)
	for _, update := range updatesx2 {
		log.Debug("Observe Fine Grane Updates X2", "Path", e.parser.GnmiPathToXPath(update.Path, true), "Value", update.GetVal())
	}

	deletes, updates, err := e.parser.FindResourceDeltaGnmi(updatesx1, updatesx2, log)
	if err != nil {
		return managed.ExternalObservation{}, err
	}
	// resource is NOT up to date
	if len(deletes) != 0 || len(updates) != 0 {
		// resource is NOT up to date
		log.Debug("Observing Response: resource NOT up to date", "Exists", true, "HasData", true, "UpToDate", false, "Response", resp, "Updates", updates, "Deletes", deletes)
		for _, del := range deletes {
			log.Debug("Observing Response: resource NOT up to date, deletes", "path", e.parser.GnmiPathToXPath(del, true))
		}
		for _, upd := range updates {
			val, _ := e.parser.GetValue(upd.GetVal())
			log.Debug("Observing Response: resource NOT up to date, updates", "path", e.parser.GnmiPathToXPath(upd.GetPath(), true), "data", val)
		}
		return managed.ExternalObservation{
			Ready:            true,
			ResourceExists:   true,
			ResourceHasData:  true,
			ResourceUpToDate: false,
			ResourceDeletes:  deletes,
			ResourceUpdates:  updates,
		}, nil
	}
	// resource is up to date
	log.Debug("Observing Response: resource up to date", "Exists", true, "HasData", true, "UpToDate", true, "Response", resp)
	return managed.ExternalObservation{
		Ready:            true,
		ResourceExists:   true,
		ResourceHasData:  true,
		ResourceUpToDate: true,
	}, nil

}

func (e *externalIpamIprange) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	o, ok := mg.(*ipamv1alpha1.IpamIpamIprange)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errUnexpectedIpamIprange)
	}
	log := e.log.WithValues("Resource", o.GetName())
	log.Debug("Creating ...")

	rootPath := []*gnmi.Path{
		{
			Elem: []*gnmi.PathElem{
				{Name: "ipam"},
				{Name: "ip-range", Key: map[string]string{
					"tenant":           *o.Spec.ForNetworkNode.IpamIpamIprange.Tenant,
					"network-instance": *o.Spec.ForNetworkNode.IpamIpamIprange.NetworkInstance,
					"start":            *o.Spec.ForNetworkNode.IpamIpamIprange.Start,
					"end":              *o.Spec.ForNetworkNode.IpamIpamIprange.End,
				}},
			},
		},
	}

	d, err := json.Marshal(&o.Spec.ForNetworkNode)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errJSONMarshal)
	}

	var x1 interface{}
	if err := json.Unmarshal(d, &x1); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errJSONUnMarshal)
	}

	// remove the hierarchical elements for data processing, comparison, etc
	// they are used in the provider for parent dependency resolution
	// but are not relevant in the data, they are referenced in the rootPath
	// when interacting with the device driver
	hids := make([]string, 0)
	x1 = e.parser.RemoveLeafsFromJSONData(x1, hids)
	// for lists with keys we need to create a list before calulating the paths since this is what
	// the object eventually happens to be based upon. We avoid having multiple entries in a list object
	// and hence we have to add this step
	x1, err = e.parser.AddJSONDataToList(x1)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errWrongInputdata)
	}

	updates := e.parser.GetUpdatesFromJSONDataGnmi(rootPath[0], e.parser.XpathToGnmiPath("/", 0), x1, resourceRefPathsIpamIprange)
	for _, update := range updates {
		log.Debug("Create Fine Grane Updates", "Path", update.Path, "Value", update.GetVal())
	}

	if len(updates) == 0 {
		log.Debug("cannot create object since there are no updates present")
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateObject)
	}

	req := &gnmi.SetRequest{
		Replace: updates,
	}

	_, err = e.client.Set(ctx, req)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateIpamIprange)
	}

	return managed.ExternalCreation{}, nil
}

func (e *externalIpamIprange) Update(ctx context.Context, mg resource.Managed, obs managed.ExternalObservation) (managed.ExternalUpdate, error) {
	o, ok := mg.(*ipamv1alpha1.IpamIpamIprange)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errUnexpectedIpamIprange)
	}
	log := e.log.WithValues("Resource", o.GetName())
	log.Debug("Updating ...")

	for _, u := range obs.ResourceUpdates {
		log.Debug("Update -> Update", "Path", u.Path, "Value", u.GetVal())
	}
	for _, d := range obs.ResourceDeletes {
		log.Debug("Update -> Delete", "Path", d)
	}

	req := &gnmi.SetRequest{
		Update: obs.ResourceUpdates,
		Delete: obs.ResourceDeletes,
	}

	_, err := e.client.Set(ctx, req)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errReadIpamIprange)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *externalIpamIprange) Delete(ctx context.Context, mg resource.Managed) error {
	o, ok := mg.(*ipamv1alpha1.IpamIpamIprange)
	if !ok {
		return errors.New(errUnexpectedIpamIprange)
	}
	log := e.log.WithValues("Resource", o.GetName())
	log.Debug("Deleting ...")

	rootPath := []*gnmi.Path{
		{
			Elem: []*gnmi.PathElem{
				{Name: "ipam"},
				{Name: "ip-range", Key: map[string]string{
					"tenant":           *o.Spec.ForNetworkNode.IpamIpamIprange.Tenant,
					"network-instance": *o.Spec.ForNetworkNode.IpamIpamIprange.NetworkInstance,
					"start":            *o.Spec.ForNetworkNode.IpamIpamIprange.Start,
					"end":              *o.Spec.ForNetworkNode.IpamIpamIprange.End,
				}},
			},
		},
	}

	req := gnmi.SetRequest{
		Delete: rootPath,
	}

	_, err := e.client.Set(ctx, &req)
	if err != nil {
		return errors.Wrap(err, errDeleteIpamIprange)
	}

	return nil
}

func (e *externalIpamIprange) GetTarget() []string {
	return e.targets
}

func (e *externalIpamIprange) GetConfig(ctx context.Context) ([]byte, error) {
	e.log.Debug("Get Config ...")
	req := &gnmi.GetRequest{
		Path:     []*gnmi.Path{},
		Encoding: gnmi.Encoding_JSON,
	}

	resp, err := e.client.Get(ctx, req)
	if err != nil {
		return make([]byte, 0), errors.Wrap(err, errGetConfig)
	}

	if len(resp.GetNotification()) != 0 {
		if len(resp.GetNotification()[0].GetUpdate()) != 0 {
			x2, err := e.parser.GetValue(resp.GetNotification()[0].GetUpdate()[0].Val)
			if err != nil {
				return make([]byte, 0), errors.Wrap(err, errGetConfig)
			}

			data, err := json.Marshal(x2)
			if err != nil {
				return make([]byte, 0), errors.Wrap(err, errJSONMarshal)
			}
			return data, nil
		}
	}
	e.log.Debug("Get Config Empty response")
	return nil, nil
}

func (e *externalIpamIprange) GetResourceName(ctx context.Context, path []*gnmi.Path) (string, error) {
	return "", nil
}
