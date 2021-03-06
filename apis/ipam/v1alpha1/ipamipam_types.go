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
	// IpamFinalizer is the name of the finalizer added to
	// Ipam to block delete operations until the physical node can be
	// deprovisioned.
	IpamFinalizer string = "ipam.ipam.nddo.yndd.io"
)

// Ipam struct
type Ipam struct {
	Rir []*IpamRir `json:"rir,omitempty"`
}

// IpamRir struct
type IpamRir struct {
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="[A-Za-z0-9 !@#$^&()|+=`~.,'/_:;?-]*"
	Description *string `json:"description,omitempty"`
	// +kubebuilder:validation:Enum=`afrinic`;`apnic`;`arin`;`lacnic`;`rfc1918`;`rfc6598`;`ripe`;`ula`
	Name *string       `json:"name"`
	Tag  []*IpamRirTag `json:"tag,omitempty"`
}

// IpamRirTag struct
type IpamRirTag struct {
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

// IpamParameters are the parameter fields of a Ipam.
type IpamParameters struct {
	IpamIpam *Ipam `json:"ipam,omitempty"`
}

// IpamObservation are the observable fields of a Ipam.
type IpamObservation struct {
	*Nddoipam `json:",inline"`
}

// A IpamSpec defines the desired state of a Ipam.
type IpamSpec struct {
	nddv1.ResourceSpec `json:",inline"`
	ForNetworkNode     IpamParameters `json:"forNetworkNode"`
}

// A IpamStatus represents the observed state of a Ipam.
type IpamStatus struct {
	nddv1.ResourceStatus `json:",inline"`
	AtNetworkNode        IpamObservation `json:"atNetworkNode,omitempty"`
}

// +kubebuilder:object:root=true

// IpamIpam is the Schema for the Ipam API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="TARGET",type="string",JSONPath=".status.conditions[?(@.kind=='TargetFound')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="LOCALLEAFREF",type="string",JSONPath=".status.conditions[?(@.kind=='InternalLeafrefValidationSuccess')].status"
// +kubebuilder:printcolumn:name="EXTLEAFREF",type="string",JSONPath=".status.conditions[?(@.kind=='ExternalLeafrefValidationSuccess')].status"
// +kubebuilder:printcolumn:name="PARENTDEP",type="string",JSONPath=".status.conditions[?(@.kind=='ParentValidationSuccess')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={ndd,ipam}
type IpamIpam struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IpamSpec   `json:"spec,omitempty"`
	Status IpamStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IpamIpamList contains a list of Ipams
type IpamIpamList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IpamIpam `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IpamIpam{}, &IpamIpamList{})
}

// Ipam type metadata.
var (
	IpamKindKind         = reflect.TypeOf(IpamIpam{}).Name()
	IpamGroupKind        = schema.GroupKind{Group: Group, Kind: IpamKindKind}.String()
	IpamKindAPIVersion   = IpamKindKind + "." + GroupVersion.String()
	IpamGroupVersionKind = GroupVersion.WithKind(IpamKindKind)
)
