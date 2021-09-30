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
	// IpamIpaddressFinalizer is the name of the finalizer added to
	// IpamIpaddress to block delete operations until the physical node can be
	// deprovisioned.
	IpamIpaddressFinalizer string = "ipAddress.ipam.nddo.yndd.io"
)

// IpamIpaddress struct
type IpamIpaddress struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])/(([0-9])|([1-2][0-9])|(3[0-2]))|((:|[0-9a-fA-F]{0,4}):)([0-9a-fA-F]{0,4}:){0,5}((([0-9a-fA-F]{0,4}:)?(:|[0-9a-fA-F]{0,4}))|(((25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])))(/(([0-9])|([0-9]{2})|(1[0-1][0-9])|(12[0-8])))`
	Address *string `json:"address"`
	// +kubebuilder:validation:Enum=`disable`;`enable`
	// +kubebuilder:default:="enable"
	AdminState *string `json:"admin-state,omitempty"`
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="[A-Za-z0-9 !@#$^&()|+=`~.,'/_:;?-]*"
	Description *string `json:"description,omitempty"`
	DnsName     *string `json:"dns-name,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])|((:|[0-9a-fA-F]{0,4}):)([0-9a-fA-F]{0,4}:){0,5}((([0-9a-fA-F]{0,4}:)?(:|[0-9a-fA-F]{0,4}))|(((25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])))`
	NatInside *string `json:"nat-inside,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])|((:|[0-9a-fA-F]{0,4}):)([0-9a-fA-F]{0,4}:){0,5}((([0-9a-fA-F]{0,4}:)?(:|[0-9a-fA-F]{0,4}))|(((25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])))`
	NatOutside *string `json:"nat-outside,omitempty"`
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="[A-Za-z0-9 !@#$^&()|+=`~.,'/_:;?-]*"
	// +kubebuilder:default:="default"
	NetworkInstance *string `json:"network-instance,omitempty"`
	// +kubebuilder:validation:Enum=`dhcp`;`slaac`;`static`
	Origin   *string                  `json:"origin,omitempty"`
	IpPrefix []*IpamIpaddressIpPrefix `json:"ip-prefix,omitempty"`
	IpRange  []*IpamIpaddressIpRange  `json:"ip-range,omitempty"`
	Tag      []*IpamIpaddressTag      `json:"tag,omitempty"`
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="[A-Za-z0-9 !@#$^&()|+=`~.,'/_:;?-]*"
	// +kubebuilder:default:="default"
	Tenant *string `json:"tenant,omitempty"`
}

// IpamIpaddressIpPrefix struct
type IpamIpaddressIpPrefix struct {
	NetworkInstance *string `json:"network-instance"`
	Prefix          *string `json:"prefix"`
	Tenant          *string `json:"tenant"`
}

// IpamIpaddressIpRange struct
type IpamIpaddressIpRange struct {
	End             *string `json:"end"`
	NetworkInstance *string `json:"network-instance"`
	Start           *string `json:"start"`
	Tenant          *string `json:"tenant"`
}

// IpamIpaddressTag struct
type IpamIpaddressTag struct {
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

// IpamIpaddressParameters are the parameter fields of a IpamIpaddress.
type IpamIpaddressParameters struct {
	IpamIpamIpaddress *IpamIpaddress `json:"ip-address,omitempty"`
}

// IpamIpaddressObservation are the observable fields of a IpamIpaddress.
type IpamIpaddressObservation struct {
}

// A IpamIpaddressSpec defines the desired state of a IpamIpaddress.
type IpamIpaddressSpec struct {
	nddv1.ResourceSpec `json:",inline"`
	ForNetworkNode     IpamIpaddressParameters `json:"forNetworkNode"`
}

// A IpamIpaddressStatus represents the observed state of a IpamIpaddress.
type IpamIpaddressStatus struct {
	nddv1.ResourceStatus `json:",inline"`
	AtNetworkNode        IpamIpaddressObservation `json:"atNetworkNode,omitempty"`
}

// +kubebuilder:object:root=true

// IpamIpamIpaddress is the Schema for the IpamIpaddress API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="TARGET",type="string",JSONPath=".status.conditions[?(@.kind=='TargetFound')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="LOCALLEAFREF",type="string",JSONPath=".status.conditions[?(@.kind=='InternalLeafrefValidationSuccess')].status"
// +kubebuilder:printcolumn:name="EXTLEAFREF",type="string",JSONPath=".status.conditions[?(@.kind=='ExternalLeafrefValidationSuccess')].status"
// +kubebuilder:printcolumn:name="PARENTDEP",type="string",JSONPath=".status.conditions[?(@.kind=='ParentValidationSuccess')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={ndd,ipam}
type IpamIpamIpaddress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IpamIpaddressSpec   `json:"spec,omitempty"`
	Status IpamIpaddressStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IpamIpamIpaddressList contains a list of IpamIpaddresss
type IpamIpamIpaddressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IpamIpamIpaddress `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IpamIpamIpaddress{}, &IpamIpamIpaddressList{})
}

// IpamIpaddress type metadata.
var (
	IpamIpaddressKindKind         = reflect.TypeOf(IpamIpamIpaddress{}).Name()
	IpamIpaddressGroupKind        = schema.GroupKind{Group: Group, Kind: IpamIpaddressKindKind}.String()
	IpamIpaddressKindAPIVersion   = IpamIpaddressKindKind + "." + GroupVersion.String()
	IpamIpaddressGroupVersionKind = GroupVersion.WithKind(IpamIpaddressKindKind)
)
