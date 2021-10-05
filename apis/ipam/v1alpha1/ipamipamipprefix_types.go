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
	// IpamIpprefixFinalizer is the name of the finalizer added to
	// IpamIpprefix to block delete operations until the physical node can be
	// deprovisioned.
	IpamIpprefixFinalizer string = "ipPrefix.ipam.nddo.yndd.io"
)

// IpamIpprefix struct
type IpamIpprefix struct {
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
	NetworkInstance *string                  `json:"network-instance,omitempty"`
	Aggregate       []*IpamIpprefixAggregate `json:"aggregate,omitempty"`
	IpPrefix        []*IpamIpprefixIpPrefix  `json:"ip-prefix,omitempty"`
	Pool            *bool                    `json:"pool,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])/(([0-9])|([1-2][0-9])|(3[0-2]))|((:|[0-9a-fA-F]{0,4}):)([0-9a-fA-F]{0,4}:){0,5}((([0-9a-fA-F]{0,4}:)?(:|[0-9a-fA-F]{0,4}))|(((25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])))(/(([0-9])|([0-9]{2})|(1[0-1][0-9])|(12[0-8])))`
	Prefix *string            `json:"prefix"`
	Tag    []*IpamIpprefixTag `json:"tag,omitempty"`
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="[A-Za-z0-9 !@#$^&()|+=`~.,'/_:;?-]*"
	// +kubebuilder:default:="default"
	Tenant *string `json:"tenant,omitempty"`
}

// IpamIpprefixAggregate struct
type IpamIpprefixAggregate struct {
	NetworkInstance *string `json:"network-instance"`
	Prefix          *string `json:"prefix"`
	Tenant          *string `json:"tenant"`
}

// IpamIpprefixIpPrefix struct
type IpamIpprefixIpPrefix struct {
	NetworkInstance *string `json:"network-instance"`
	Prefix          *string `json:"prefix"`
	Tenant          *string `json:"tenant"`
}

// IpamIpprefixTag struct
type IpamIpprefixTag struct {
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

// IpamIpprefixParameters are the parameter fields of a IpamIpprefix.
type IpamIpprefixParameters struct {
	IpamIpamIpprefix *IpamIpprefix `json:"ip-prefix,omitempty"`
}

// IpamIpprefixObservation are the observable fields of a IpamIpprefix.
type IpamIpprefixObservation struct {
}

// A IpamIpprefixSpec defines the desired state of a IpamIpprefix.
type IpamIpprefixSpec struct {
	nddv1.ResourceSpec `json:",inline"`
	ForNetworkNode     IpamIpprefixParameters `json:"forNetworkNode"`
}

// A IpamIpprefixStatus represents the observed state of a IpamIpprefix.
type IpamIpprefixStatus struct {
	nddv1.ResourceStatus `json:",inline"`
	AtNetworkNode        IpamIpprefixObservation `json:"atNetworkNode,omitempty"`
}

// +kubebuilder:object:root=true

// IpamIpamIpprefix is the Schema for the IpamIpprefix API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="TARGET",type="string",JSONPath=".status.conditions[?(@.kind=='TargetFound')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="LOCALLEAFREF",type="string",JSONPath=".status.conditions[?(@.kind=='InternalLeafrefValidationSuccess')].status"
// +kubebuilder:printcolumn:name="EXTLEAFREF",type="string",JSONPath=".status.conditions[?(@.kind=='ExternalLeafrefValidationSuccess')].status"
// +kubebuilder:printcolumn:name="PARENTDEP",type="string",JSONPath=".status.conditions[?(@.kind=='ParentValidationSuccess')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={ndd,ipam}
type IpamIpamIpprefix struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IpamIpprefixSpec   `json:"spec,omitempty"`
	Status IpamIpprefixStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IpamIpamIpprefixList contains a list of IpamIpprefixs
type IpamIpamIpprefixList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IpamIpamIpprefix `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IpamIpamIpprefix{}, &IpamIpamIpprefixList{})
}

// IpamIpprefix type metadata.
var (
	IpamIpprefixKindKind         = reflect.TypeOf(IpamIpamIpprefix{}).Name()
	IpamIpprefixGroupKind        = schema.GroupKind{Group: Group, Kind: IpamIpprefixKindKind}.String()
	IpamIpprefixKindAPIVersion   = IpamIpprefixKindKind + "." + GroupVersion.String()
	IpamIpprefixGroupVersionKind = GroupVersion.WithKind(IpamIpprefixKindKind)
)
