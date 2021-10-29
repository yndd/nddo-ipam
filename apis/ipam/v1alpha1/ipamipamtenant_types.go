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
	// IpamTenantFinalizer is the name of the finalizer added to
	// IpamTenant to block delete operations until the physical node can be
	// deprovisioned.
	IpamTenantFinalizer string = "tenant.ipam.nddo.yndd.io"
)

// IpamTenant struct
type IpamTenant struct {
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
	Name *string          `json:"name,omitempty"`
	Tag  []*IpamTenantTag `json:"tag,omitempty"`
}

// IpamTenantTag struct
type IpamTenantTag struct {
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

// IpamTenantParameters are the parameter fields of a IpamTenant.
type IpamTenantParameters struct {
	IpamIpamTenant *IpamTenant `json:"tenant,omitempty"`
}

// IpamTenantObservation are the observable fields of a IpamTenant.
type IpamTenantObservation struct {
}

// A IpamTenantSpec defines the desired state of a IpamTenant.
type IpamTenantSpec struct {
	nddv1.ResourceSpec `json:",inline"`
	ForNetworkNode     IpamTenantParameters `json:"forNetworkNode"`
}

// A IpamTenantStatus represents the observed state of a IpamTenant.
type IpamTenantStatus struct {
	nddv1.ResourceStatus `json:",inline"`
	AtNetworkNode        IpamTenantObservation `json:"atNetworkNode,omitempty"`
}

// +kubebuilder:object:root=true

// IpamIpamTenant is the Schema for the IpamTenant API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="TARGET",type="string",JSONPath=".status.conditions[?(@.kind=='TargetFound')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="LOCALLEAFREF",type="string",JSONPath=".status.conditions[?(@.kind=='InternalLeafrefValidationSuccess')].status"
// +kubebuilder:printcolumn:name="EXTLEAFREF",type="string",JSONPath=".status.conditions[?(@.kind=='ExternalLeafrefValidationSuccess')].status"
// +kubebuilder:printcolumn:name="PARENTDEP",type="string",JSONPath=".status.conditions[?(@.kind=='ParentValidationSuccess')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={ndd,ipam}
type IpamIpamTenant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IpamTenantSpec   `json:"spec,omitempty"`
	Status IpamTenantStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IpamIpamTenantList contains a list of IpamTenants
type IpamIpamTenantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IpamIpamTenant `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IpamIpamTenant{}, &IpamIpamTenantList{})
}

// IpamTenant type metadata.
var (
	IpamTenantKindKind         = reflect.TypeOf(IpamIpamTenant{}).Name()
	IpamTenantGroupKind        = schema.GroupKind{Group: Group, Kind: IpamTenantKindKind}.String()
	IpamTenantKindAPIVersion   = IpamTenantKindKind + "." + GroupVersion.String()
	IpamTenantGroupVersionKind = GroupVersion.WithKind(IpamTenantKindKind)
)
