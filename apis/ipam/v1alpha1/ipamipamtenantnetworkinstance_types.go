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
	// IpamTenantNetworkinstanceFinalizer is the name of the finalizer added to
	// IpamTenantNetworkinstance to block delete operations until the physical node can be
	// deprovisioned.
	IpamTenantNetworkinstanceFinalizer string = "networkInstance.ipam.nddo.yndd.io"
)

// IpamTenantNetworkinstance struct
type IpamTenantNetworkinstance struct {
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
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="[A-Za-z0-9 !@#$^&()|+=`~.,'/_:;?-]*"
	// +kubebuilder:default:="default"
	Name *string                         `json:"name,omitempty"`
	Tag  []*IpamTenantNetworkinstanceTag `json:"tag,omitempty"`
}

// IpamTenantNetworkinstanceTag struct
type IpamTenantNetworkinstanceTag struct {
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

// IpamTenantNetworkinstanceParameters are the parameter fields of a IpamTenantNetworkinstance.
type IpamTenantNetworkinstanceParameters struct {
	TenantName                    *string                    `json:"tenant-name"`
	IpamIpamTenantNetworkinstance *IpamTenantNetworkinstance `json:"network-instance,omitempty"`
}

// IpamTenantNetworkinstanceObservation are the observable fields of a IpamTenantNetworkinstance.
type IpamTenantNetworkinstanceObservation struct {
}

// A IpamTenantNetworkinstanceSpec defines the desired state of a IpamTenantNetworkinstance.
type IpamTenantNetworkinstanceSpec struct {
	nddv1.ResourceSpec `json:",inline"`
	ForNetworkNode     IpamTenantNetworkinstanceParameters `json:"forNetworkNode"`
}

// A IpamTenantNetworkinstanceStatus represents the observed state of a IpamTenantNetworkinstance.
type IpamTenantNetworkinstanceStatus struct {
	nddv1.ResourceStatus `json:",inline"`
	AtNetworkNode        IpamTenantNetworkinstanceObservation `json:"atNetworkNode,omitempty"`
}

// +kubebuilder:object:root=true

// IpamIpamTenantNetworkinstance is the Schema for the IpamTenantNetworkinstance API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="TARGET",type="string",JSONPath=".status.conditions[?(@.kind=='TargetFound')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="LOCALLEAFREF",type="string",JSONPath=".status.conditions[?(@.kind=='InternalLeafrefValidationSuccess')].status"
// +kubebuilder:printcolumn:name="EXTLEAFREF",type="string",JSONPath=".status.conditions[?(@.kind=='ExternalLeafrefValidationSuccess')].status"
// +kubebuilder:printcolumn:name="PARENTDEP",type="string",JSONPath=".status.conditions[?(@.kind=='ParentValidationSuccess')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={ndd,ipam}
type IpamIpamTenantNetworkinstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IpamTenantNetworkinstanceSpec   `json:"spec,omitempty"`
	Status IpamTenantNetworkinstanceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IpamIpamTenantNetworkinstanceList contains a list of IpamTenantNetworkinstances
type IpamIpamTenantNetworkinstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IpamIpamTenantNetworkinstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IpamIpamTenantNetworkinstance{}, &IpamIpamTenantNetworkinstanceList{})
}

// IpamTenantNetworkinstance type metadata.
var (
	IpamTenantNetworkinstanceKindKind         = reflect.TypeOf(IpamIpamTenantNetworkinstance{}).Name()
	IpamTenantNetworkinstanceGroupKind        = schema.GroupKind{Group: Group, Kind: IpamTenantNetworkinstanceKindKind}.String()
	IpamTenantNetworkinstanceKindAPIVersion   = IpamTenantNetworkinstanceKindKind + "." + GroupVersion.String()
	IpamTenantNetworkinstanceGroupVersionKind = GroupVersion.WithKind(IpamTenantNetworkinstanceKindKind)
)
