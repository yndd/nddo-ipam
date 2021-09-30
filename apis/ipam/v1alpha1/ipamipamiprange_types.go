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

package v1alpha1

import (
	"reflect"

	nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// IpamIprangeFinalizer is the name of the finalizer added to
	// IpamIprange to block delete operations until the physical node can be
	// deprovisioned.
	IpamIprangeFinalizer string = "ipRange.ipam.nddo.yndd.io"
)

// IpamIprange struct
type IpamIprange struct {
	// +kubebuilder:validation:Enum=`disable`;`enable`
	// +kubebuilder:default:="enable"
	AdminState *string `json:"admin-state,omitempty"`
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="[A-Za-z0-9 !@#$^&()|+=`~.,'/_:;?-]*"
	Description *string `json:"description,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])|((:|[0-9a-fA-F]{0,4}):)([0-9a-fA-F]{0,4}:){0,5}((([0-9a-fA-F]{0,4}:)?(:|[0-9a-fA-F]{0,4}))|(((25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])))`
	End *string `json:"end"`
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="[A-Za-z0-9 !@#$^&()|+=`~.,'/_:;?-]*"
	// +kubebuilder:default:="default"
	NetworkInstance *string                 `json:"network-instance,omitempty"`
	Aggregate       []*IpamIprangeAggregate `json:"aggregate,omitempty"`
	IpPrefix        []*IpamIprangeIpPrefix  `json:"ip-prefix,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])|((:|[0-9a-fA-F]{0,4}):)([0-9a-fA-F]{0,4}:){0,5}((([0-9a-fA-F]{0,4}:)?(:|[0-9a-fA-F]{0,4}))|(((25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])))`
	Start *string           `json:"start"`
	Tag   []*IpamIprangeTag `json:"tag,omitempty"`
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="[A-Za-z0-9 !@#$^&()|+=`~.,'/_:;?-]*"
	// +kubebuilder:default:="default"
	Tenant *string `json:"tenant,omitempty"`
}

// IpamIprangeAggregate struct
type IpamIprangeAggregate struct {
	NetworkInstance *string `json:"network-instance"`
	Prefix          *string `json:"prefix"`
	Tenant          *string `json:"tenant"`
}

// IpamIprangeIpPrefix struct
type IpamIprangeIpPrefix struct {
	NetworkInstance *string `json:"network-instance"`
	Prefix          *string `json:"prefix"`
	Tenant          *string `json:"tenant"`
}

// IpamIprangeTag struct
type IpamIprangeTag struct {
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="[A-Za-z0-9 !@#$^&()|+=`~.,'/_:;?-]*"
	Key *string `json:"key"`
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="[A-Za-z0-9 !@#$^&()|+=`~.,'/_:;?-]*"
	Value *string `json:"value,omitempty"`
}

// IpamIprangeParameters are the parameter fields of a IpamIprange.
type IpamIprangeParameters struct {
	IpamIpamIprange *IpamIprange `json:"ip-range,omitempty"`
}

// IpamIprangeObservation are the observable fields of a IpamIprange.
type IpamIprangeObservation struct {
}

// A IpamIprangeSpec defines the desired state of a IpamIprange.
type IpamIprangeSpec struct {
	nddv1.ResourceSpec `json:",inline"`
	ForNetworkNode     IpamIprangeParameters `json:"forNetworkNode"`
}

// A IpamIprangeStatus represents the observed state of a IpamIprange.
type IpamIprangeStatus struct {
	nddv1.ResourceStatus `json:",inline"`
	AtNetworkNode        IpamIprangeObservation `json:"atNetworkNode,omitempty"`
}

// +kubebuilder:object:root=true

// IpamIpamIprange is the Schema for the IpamIprange API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="TARGET",type="string",JSONPath=".status.conditions[?(@.kind=='TargetFound')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="LOCALLEAFREF",type="string",JSONPath=".status.conditions[?(@.kind=='InternalLeafrefValidationSuccess')].status"
// +kubebuilder:printcolumn:name="EXTLEAFREF",type="string",JSONPath=".status.conditions[?(@.kind=='ExternalLeafrefValidationSuccess')].status"
// +kubebuilder:printcolumn:name="PARENTDEP",type="string",JSONPath=".status.conditions[?(@.kind=='ParentValidationSuccess')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={ndd,ipam}
type IpamIpamIprange struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IpamIprangeSpec   `json:"spec,omitempty"`
	Status IpamIprangeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IpamIpamIprangeList contains a list of IpamIpranges
type IpamIpamIprangeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IpamIpamIprange `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IpamIpamIprange{}, &IpamIpamIprangeList{})
}

// IpamIprange type metadata.
var (
	IpamIprangeKindKind         = reflect.TypeOf(IpamIpamIprange{}).Name()
	IpamIprangeGroupKind        = schema.GroupKind{Group: Group, Kind: IpamIprangeKindKind}.String()
	IpamIprangeKindAPIVersion   = IpamIprangeKindKind + "." + GroupVersion.String()
	IpamIprangeGroupVersionKind = GroupVersion.WithKind(IpamIprangeKindKind)
)
