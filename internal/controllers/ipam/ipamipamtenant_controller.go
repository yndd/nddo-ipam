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
	"fmt"
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
	"github.com/yndd/ndd-yang/pkg/leafref"
	"github.com/yndd/ndd-yang/pkg/parser"
	"github.com/yndd/ndd-yang/pkg/yentry"
	"github.com/yndd/ndd-yang/pkg/yparser"
	"github.com/yndd/ndd-yang/pkg/yresource"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	cevent "sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	ipamv1alpha1 "github.com/yndd/nddo-ipam/apis/ipam/v1alpha1"
	"github.com/yndd/nddo-ipam/internal/shared"
)

const (
	// Errors
	errUnexpectedIpamTenant       = "the managed resource is not a IpamTenant resource"
	errKubeUpdateFailedIpamTenant = "cannot update IpamTenant"
	errReadIpamTenant             = "cannot read IpamTenant"
	errCreateIpamTenant           = "cannot create IpamTenant"
	erreUpdateIpamTenant          = "cannot update IpamTenant"
	errDeleteIpamTenant           = "cannot delete IpamTenant"

	// resource information
	// resourcePrefixIpamTenant = "ipam.nddo.yndd.io.v1alpha1.IpamTenant"
)

// SetupIpamTenant adds a controller that reconciles IpamTenants.
func SetupIpamTenant(mgr ctrl.Manager, o controller.Options, nddcopts *shared.NddControllerOptions) (string, chan cevent.GenericEvent, error) {

	name := managed.ControllerName(ipamv1alpha1.IpamTenantGroupKind)

	events := make(chan cevent.GenericEvent)

	y := initYangIpamTenant()

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(ipamv1alpha1.IpamTenantGroupVersionKind),
		managed.WithExternalConnecter(&connectorIpamTenant{
			log:              nddcopts.Logger,
			kube:             mgr.GetClient(),
			usage:            resource.NewNetworkNodeUsageTracker(mgr.GetClient(), &ndrv1.NetworkNodeUsage{}),
			rootSchema:       nddcopts.Yentry,
			y:                y,
			newClientFn:      target.NewTarget,
			grpcQueryAddress: nddcopts.GrpcQueryAddress,
		},
		),
		managed.WithParser(nddcopts.Logger),
		managed.WithValidator(&validatorIpamTenant{
			log:        nddcopts.Logger,
			rootSchema: nddcopts.Yentry,
			y:          y,
			parser:     *parser.NewParser(parser.WithLogger(nddcopts.Logger))},
		),
		managed.WithLogger(nddcopts.Logger.WithValues("controller", name)),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))))

	return ipamv1alpha1.IpamTenantGroupKind, events, ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o).
		For(&ipamv1alpha1.IpamIpamTenant{}).
		WithEventFilter(resource.IgnoreUpdateWithoutGenerationChangePredicate()).
		Watches(
			&source.Channel{Source: events},
			&handler.EnqueueRequestForObject{},
		).
		Complete(r)
}

type ipamTenant struct {
	*yresource.Resource
}

func initYangIpamTenant(opts ...yresource.Option) yresource.Handler {
	rr := &yresource.Resource{}
	r := &ipamTenant{rr}

	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *ipamTenant) GetRootPath(mg resource.Managed) []*gnmi.Path {
	o, ok := mg.(*ipamv1alpha1.IpamIpamTenant)
	if !ok {
		return nil
	}
	return []*gnmi.Path{
		{
			Elem: []*gnmi.PathElem{
				{Name: "ipam"},
				{Name: "tenant", Key: map[string]string{
					"name": *o.Spec.ForNetworkNode.IpamIpamTenant.Name,
				}},
			},
		},
	}
}

func (r *ipamTenant) GetParentDependency(mg resource.Managed) []*leafref.LeafRef {
	rootPath := r.GetRootPath(mg)
	// if the path is not bigger than 1 element there is no parent dependency
	if len(rootPath[0].GetElem()) < 2 {
		return []*leafref.LeafRef{}
	}
	// the dependency path is the rootPath except for the last element
	dependencyPathElem := rootPath[0].GetElem()[:(len(rootPath[0].GetElem()) - 1)]
	// check for keys present, if no keys present we return an empty list
	keysPresent := false
	for _, pathElem := range dependencyPathElem {
		if len(pathElem.GetKey()) != 0 {
			keysPresent = true
		}
	}
	if !keysPresent {
		return []*leafref.LeafRef{}
	}

	fmt.Printf("GetParentDependency tenant: %v\n", yparser.GnmiPath2XPath(&gnmi.Path{Elem: dependencyPathElem}, true))

	// return the rootPath except the last entry
	return []*leafref.LeafRef{{RemotePath: &gnmi.Path{Elem: dependencyPathElem}}}
}

type validatorIpamTenant struct {
	log        logging.Logger
	parser     parser.Parser
	rootSchema *yentry.Entry
	y          yresource.Handler
}

func (v *validatorIpamTenant) ValidateLocalleafRef(ctx context.Context, mg resource.Managed) (managed.ValidateLocalleafRefObservation, error) {
	return managed.ValidateLocalleafRefObservation{
		Success:          true,
		ResolvedLeafRefs: []*leafref.ResolvedLeafRef{}}, nil
}

func (v *validatorIpamTenant) ValidateExternalleafRef(ctx context.Context, mg resource.Managed, cfg []byte) (managed.ValidateExternalleafRefObservation, error) {
	log := v.log.WithValues("resource", mg.GetName(), "kind", mg.GetObjectKind())
	log.Debug("ValidateLeafRef...")

	// json unmarshal the resource
	o, ok := mg.(*ipamv1alpha1.IpamIpamTenant)
	if !ok {
		return managed.ValidateExternalleafRefObservation{}, errors.New(errUnexpectedIpamTenant)
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

	rootPath := v.y.GetRootPath(o)

	leafRefs := v.rootSchema.GetLeafRefsLocal(true, rootPath[0], &gnmi.Path{}, make([]*leafref.LeafRef, 0))
	log.Debug("Validate leafRefs ...", "Path", yparser.GnmiPath2XPath(rootPath[0], false), "leafRefs", leafRefs)

	// For local external leafref validation we need to supply the external
	// data to validate the remote leafref, we use x2 for this
	success, resultValidation, err := yparser.ValidateLeafRef(
		rootPath[0], x1, x2, leafRefs, v.rootSchema)

	if err != nil {
		return managed.ValidateExternalleafRefObservation{
			Success: false,
		}, nil
	}
	if !success {
		log.Debug("ValidateExternalleafRef failed", "resultleafRefValidation", resultValidation)
		return managed.ValidateExternalleafRefObservation{
			Success:          false,
			ResolvedLeafRefs: resultValidation}, nil
	}
	log.Debug("ValidateExternalleafRef success", "resultleafRefValidation", resultValidation)
	return managed.ValidateExternalleafRefObservation{
		Success:          true,
		ResolvedLeafRefs: resultValidation}, nil
}

func (v *validatorIpamTenant) ValidateParentDependency(ctx context.Context, mg resource.Managed, cfg []byte) (managed.ValidateParentDependencyObservation, error) {
	log := v.log.WithValues("resource", mg.GetName(), "kind", mg.GetObjectKind().GroupVersionKind().Kind)
	log.Debug("ValidateParentDependency...")

	// we initialize a global list for finer information on the resolution
	// json unmarshal the resource
	o, ok := mg.(*ipamv1alpha1.IpamIpamTenant)
	if !ok {
		return managed.ValidateParentDependencyObservation{}, errors.New(errUnexpectedIpamTenant)
	}

	dependencyLeafRef := v.y.GetParentDependency(o)

	// unmarshal the config
	var x1 interface{}
	json.Unmarshal(cfg, &x1)
	log.Debug("Latest Config", "data", x1)

	success, resultValidation, err := yparser.ValidateParentDependency(
		x1, dependencyLeafRef, v.rootSchema)
	if err != nil {
		return managed.ValidateParentDependencyObservation{
			Success: false,
		}, nil
	}
	if !success {
		log.Debug("ValidateParentDependency failed", "resultParentValidation", resultValidation)
		return managed.ValidateParentDependencyObservation{
			Success:          false,
			ResolvedLeafRefs: resultValidation}, nil
	}
	log.Debug("ValidateParentDependency success", "resultParentValidation", resultValidation)
	return managed.ValidateParentDependencyObservation{
		Success:          true,
		ResolvedLeafRefs: resultValidation}, nil
}

// ValidateResourceIndexes validates if the indexes of a resource got changed
// if so we need to delete the original resource, because it will be dangling if we dont delete it
func (v *validatorIpamTenant) ValidateResourceIndexes(ctx context.Context, mg resource.Managed) (managed.ValidateResourceIndexesObservation, error) {
	log := v.log.WithValues("resource", mg.GetName(), "kind", mg.GetObjectKind())

	// json unmarshal the resource
	o, ok := mg.(*ipamv1alpha1.IpamIpamTenant)
	if !ok {
		return managed.ValidateResourceIndexesObservation{}, errors.New(errUnexpectedIpamTenant)
	}
	log.Debug("ValidateResourceIndexes", "Spec", o.Spec)

	rootPath := v.y.GetRootPath(o)

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
type connectorIpamTenant struct {
	log         logging.Logger
	kube        client.Client
	usage       resource.Tracker
	rootSchema  *yentry.Entry
	y           yresource.Handler
	newClientFn func(c *gnmitypes.TargetConfig) *target.Target
	//newClientFn func(ctx context.Context, cfg ndd.Config) (config.ConfigurationClient, error)
	grpcQueryAddress string
}

// Connect produces an ExternalClient by:
// 1. Tracking that the managed resource is using a NetworkNode.
// 2. Getting the managed resource's NetworkNode with connection details
// A resource is mapped to a single target
func (c *connectorIpamTenant) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	log := c.log.WithValues("resource", mg.GetName())
	log.Debug("Connect")

	query_address := pkgmetav1.PrefixService + "-" + "ipam" + "." + pkgmetav1.NamespaceLocalK8sDNS + strconv.Itoa((pkgmetav1.GnmiServerPort))
	if c.grpcQueryAddress != "" {
		query_address = c.grpcQueryAddress
	}

	cfg := &gnmitypes.TargetConfig{
		Name:       "dummy",
		Address:    query_address,
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

	return &externalIpamTenant{client: cl, targets: tns, log: log, parser: *parser.NewParser(parser.WithLogger(log)), rootSchema: c.rootSchema, y: c.y}, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type externalIpamTenant struct {
	//client  config.ConfigurationClient
	client     *target.Target
	targets    []string
	log        logging.Logger
	parser     parser.Parser
	rootSchema *yentry.Entry
	y          yresource.Handler
}

func (e *externalIpamTenant) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	o, ok := mg.(*ipamv1alpha1.IpamIpamTenant)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedIpamTenant)
	}
	log := e.log.WithValues("Resource", o.GetName())
	log.Debug("Observing ...")

	// rootpath of the resource
	rootPath := e.y.GetRootPath(o)
	hierElements := e.rootSchema.GetHierarchicalResourcesLocal(true, rootPath[0], &gnmi.Path{}, make([]*gnmi.Path, 0))
	log.Debug("Observing hierElements ...", "Path", yparser.GnmiPath2XPath(rootPath[0], false), "hierElements", hierElements)

	// gnmi get request
	req := &gnmi.GetRequest{
		Prefix:   &gnmi.Path{Target: GnmiTarget, Origin: GnmiOrigin},
		Path:     rootPath,
		Encoding: gnmi.Encoding_JSON,
		Type:     gnmi.GetRequest_DataType(gnmi.GetRequest_CONFIG),
	}

	// gnmi get response
	resp, err := e.client.Get(ctx, req)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errReadIpamTenant)
	}

	// processObserve
	// o. marshal/unmarshal data
	// 1. check if resource exists
	// 2. remove parent hierarchical elements from spec
	// 3. remove resource hierarchicaal elements from gnmi response
	// 4. transform the data in gnmi to process the delta
	// 5. find the resource delta: updates and/or deletes in gnmi
	exists, deletes, updates, err := processObserve(rootPath[0], hierElements, &o.Spec.ForNetworkNode, resp, e.rootSchema)
	if err != nil {
		return managed.ExternalObservation{}, err
	}
	if !exists {
		// No Data exists -> Create
		log.Debug("Observing Response:", "Exists", false, "HasData", false, "UpToDate", false, "Response", resp)
		return managed.ExternalObservation{
			Ready:            true,
			ResourceExists:   false,
			ResourceHasData:  false,
			ResourceUpToDate: false,
		}, nil
	}
	// Data exists
	if len(deletes) != 0 || len(updates) != 0 {
		// resource is NOT up to date
		log.Debug("Observing Response: resource NOT up to date", "Exists", true, "HasData", true, "UpToDate", false, "Response", resp, "Updates", updates, "Deletes", deletes)
		for _, del := range deletes {
			log.Debug("Observing Response: resource NOT up to date, deletes", "path", yparser.GnmiPath2XPath(del, true))
		}
		for _, upd := range updates {
			val, _ := e.parser.GetValue(upd.GetVal())
			log.Debug("Observing Response: resource NOT up to date, updates", "path", yparser.GnmiPath2XPath(upd.GetPath(), true), "data", val)
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

func (e *externalIpamTenant) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	o, ok := mg.(*ipamv1alpha1.IpamIpamTenant)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errUnexpectedIpamTenant)
	}
	log := e.log.WithValues("Resource", o.GetName())
	log.Debug("Creating ...")

	// get the rootpath of the resource
	rootPath := e.y.GetRootPath(o)

	// processCreate
	// 0. marshal/unmarshal data
	// 1. transform the spec data to gnmi updates
	updates, err := processCreate(rootPath[0], &o.Spec.ForNetworkNode, e.rootSchema)
	for _, update := range updates {
		log.Debug("Create Fine Grane Updates", "Path", yparser.GnmiPath2XPath(update.Path, true), "Value", update.GetVal())
	}

	if len(updates) == 0 {
		log.Debug("cannot create object since there are no updates present")
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateObject)
	}

	req := &gnmi.SetRequest{
		Prefix:  &gnmi.Path{Target: GnmiTarget, Origin: GnmiOrigin},
		Replace: updates,
	}

	_, err = e.client.Set(ctx, req)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateIpamTenantNetworkinstanceIpprefix)
	}

	return managed.ExternalCreation{}, nil
}

func (e *externalIpamTenant) Update(ctx context.Context, mg resource.Managed, obs managed.ExternalObservation) (managed.ExternalUpdate, error) {
	o, ok := mg.(*ipamv1alpha1.IpamIpamTenant)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errUnexpectedIpamTenant)
	}
	log := e.log.WithValues("Resource", o.GetName())
	log.Debug("Updating ...")

	for _, u := range obs.ResourceUpdates {
		log.Debug("Update -> Update", "Path", yparser.GnmiPath2XPath(u.GetPath(), true), "Value", u.GetVal())
	}
	for _, d := range obs.ResourceDeletes {
		log.Debug("Update -> Delete", "Path", yparser.GnmiPath2XPath(d, true))
	}

	req := &gnmi.SetRequest{
		Prefix: &gnmi.Path{Target: GnmiTarget, Origin: GnmiOrigin},
		Update: obs.ResourceUpdates,
		Delete: obs.ResourceDeletes,
	}

	_, err := e.client.Set(ctx, req)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errReadIpamTenant)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *externalIpamTenant) Delete(ctx context.Context, mg resource.Managed) error {
	o, ok := mg.(*ipamv1alpha1.IpamIpamTenant)
	if !ok {
		return errors.New(errUnexpectedIpamTenant)
	}
	log := e.log.WithValues("Resource", o.GetName())
	log.Debug("Deleting ...")

	// get the rootpath of the resource
	rootPath := e.y.GetRootPath(o)

	req := gnmi.SetRequest{
		Prefix: &gnmi.Path{Target: GnmiTarget, Origin: GnmiOrigin},
		Delete: rootPath,
	}

	_, err := e.client.Set(ctx, &req)
	if err != nil {
		return errors.Wrap(err, errDeleteIpamTenant)
	}

	return nil
}

func (e *externalIpamTenant) GetTarget() []string {
	return e.targets
}

func (e *externalIpamTenant) GetConfig(ctx context.Context) ([]byte, error) {
	e.log.Debug("Get Config ...")
	req := &gnmi.GetRequest{
		Prefix:   &gnmi.Path{Target: GnmiTarget, Origin: GnmiOrigin},
		Path:     []*gnmi.Path{},
		Encoding: gnmi.Encoding_JSON,
		Type:     gnmi.GetRequest_DataType(gnmi.GetRequest_CONFIG),
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

func (e *externalIpamTenant) GetResourceName(ctx context.Context, path []*gnmi.Path) (string, error) {
	return "", nil
}
