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
	// IpamTenantNetworkinstanceIprangeFinalizer is the name of the finalizer added to
	// IpamTenantNetworkinstanceIprange to block delete operations until the physical node can be
	// deprovisioned.
	IpamTenantNetworkinstanceIprangeFinalizer string = "ipRange.ipam.nddo.yndd.io"
)

// IpamTenantNetworkinstanceIprange struct
type IpamTenantNetworkinstanceIprange struct {
	// +kubebuilder:validation:Enum=`first-address`;`last-address`
	// +kubebuilder:default:="first-address"
	AddressAllocationStrategy *string `json:"address-allocation-strategy,omitempty"`
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
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])|((:|[0-9a-fA-F]{0,4}):)([0-9a-fA-F]{0,4}:){0,5}((([0-9a-fA-F]{0,4}:)?(:|[0-9a-fA-F]{0,4}))|(((25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])))`
	Start *string                                `json:"start"`
	Tag   []*IpamTenantNetworkinstanceIprangeTag `json:"tag,omitempty"`
}

// IpamTenantNetworkinstanceIprangeTag struct
type IpamTenantNetworkinstanceIprangeTag struct {
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

// IpamTenantNetworkinstanceIprangeParameters are the parameter fields of a IpamTenantNetworkinstanceIprange.
type IpamTenantNetworkinstanceIprangeParameters struct {
	TenantName                           *string                           `json:"tenant-name"`
	NetworkInstanceName                  *string                           `json:"network-instance-name"`
	IpamIpamTenantNetworkinstanceIprange *IpamTenantNetworkinstanceIprange `json:"ip-range,omitempty"`
}

// IpamTenantNetworkinstanceIprangeObservation are the observable fields of a IpamTenantNetworkinstanceIprange.
type IpamTenantNetworkinstanceIprangeObservation struct {
}

// A IpamTenantNetworkinstanceIprangeSpec defines the desired state of a IpamTenantNetworkinstanceIprange.
type IpamTenantNetworkinstanceIprangeSpec struct {
	nddv1.ResourceSpec `json:",inline"`
	ForNetworkNode     IpamTenantNetworkinstanceIprangeParameters `json:"forNetworkNode"`
}

// A IpamTenantNetworkinstanceIprangeStatus represents the observed state of a IpamTenantNetworkinstanceIprange.
type IpamTenantNetworkinstanceIprangeStatus struct {
	nddv1.ResourceStatus `json:",inline"`
	AtNetworkNode        IpamTenantNetworkinstanceIprangeObservation `json:"atNetworkNode,omitempty"`
}

// +kubebuilder:object:root=true

// IpamIpamTenantNetworkinstanceIprange is the Schema for the IpamTenantNetworkinstanceIprange API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="TARGET",type="string",JSONPath=".status.conditions[?(@.kind=='TargetFound')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="LOCALLEAFREF",type="string",JSONPath=".status.conditions[?(@.kind=='InternalLeafrefValidationSuccess')].status"
// +kubebuilder:printcolumn:name="EXTLEAFREF",type="string",JSONPath=".status.conditions[?(@.kind=='ExternalLeafrefValidationSuccess')].status"
// +kubebuilder:printcolumn:name="PARENTDEP",type="string",JSONPath=".status.conditions[?(@.kind=='ParentValidationSuccess')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={ndd,ipam}
type IpamIpamTenantNetworkinstanceIprange struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IpamTenantNetworkinstanceIprangeSpec   `json:"spec,omitempty"`
	Status IpamTenantNetworkinstanceIprangeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IpamIpamTenantNetworkinstanceIprangeList contains a list of IpamTenantNetworkinstanceIpranges
type IpamIpamTenantNetworkinstanceIprangeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IpamIpamTenantNetworkinstanceIprange `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IpamIpamTenantNetworkinstanceIprange{}, &IpamIpamTenantNetworkinstanceIprangeList{})
}

// IpamTenantNetworkinstanceIprange type metadata.
var (
	IpamTenantNetworkinstanceIprangeKindKind         = reflect.TypeOf(IpamIpamTenantNetworkinstanceIprange{}).Name()
	IpamTenantNetworkinstanceIprangeGroupKind        = schema.GroupKind{Group: Group, Kind: IpamTenantNetworkinstanceIprangeKindKind}.String()
	IpamTenantNetworkinstanceIprangeKindAPIVersion   = IpamTenantNetworkinstanceIprangeKindKind + "." + GroupVersion.String()
	IpamTenantNetworkinstanceIprangeGroupVersionKind = GroupVersion.WithKind(IpamTenantNetworkinstanceIprangeKindKind)
)
