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
	// IpamTenantNetworkinstanceIpaddressFinalizer is the name of the finalizer added to
	// IpamTenantNetworkinstanceIpaddress to block delete operations until the physical node can be
	// deprovisioned.
	IpamTenantNetworkinstanceIpaddressFinalizer string = "ipAddress.ipam.nddo.yndd.io"
)

// IpamTenantNetworkinstanceIpaddress struct
type IpamTenantNetworkinstanceIpaddress struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])/(([0-9])|([1-2][0-9])|(3[0-2]))|((:|[0-9a-fA-F]{0,4}):)([0-9a-fA-F]{0,4}:){0,5}((([0-9a-fA-F]{0,4}:)?(:|[0-9a-fA-F]{0,4}))|(((25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])))(/(([0-9])|([0-9]{2})|(1[0-1][0-9])|(12[0-8])))`
	Address *string `json:"address"`
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
	DnsName     *string `json:"dns-name,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])|((:|[0-9a-fA-F]{0,4}):)([0-9a-fA-F]{0,4}:){0,5}((([0-9a-fA-F]{0,4}:)?(:|[0-9a-fA-F]{0,4}))|(((25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])))`
	NatInside *string `json:"nat-inside,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])|((:|[0-9a-fA-F]{0,4}):)([0-9a-fA-F]{0,4}:){0,5}((([0-9a-fA-F]{0,4}:)?(:|[0-9a-fA-F]{0,4}))|(((25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])))`
	NatOutside *string `json:"nat-outside,omitempty"`
	// +kubebuilder:validation:Enum=`dhcp`;`slaac`;`static`
	Origin *string                                  `json:"origin,omitempty"`
	Tag    []*IpamTenantNetworkinstanceIpaddressTag `json:"tag,omitempty"`
}

// IpamTenantNetworkinstanceIpaddressTag struct
type IpamTenantNetworkinstanceIpaddressTag struct {
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

// IpamTenantNetworkinstanceIpaddressParameters are the parameter fields of a IpamTenantNetworkinstanceIpaddress.
type IpamTenantNetworkinstanceIpaddressParameters struct {
	TenantName                             *string                             `json:"tenant-name"`
	NetworkInstanceName                    *string                             `json:"network-instance-name"`
	IpamIpamTenantNetworkinstanceIpaddress *IpamTenantNetworkinstanceIpaddress `json:"ip-address,omitempty"`
}

// IpamTenantNetworkinstanceIpaddressObservation are the observable fields of a IpamTenantNetworkinstanceIpaddress.
type IpamTenantNetworkinstanceIpaddressObservation struct {
}

// A IpamTenantNetworkinstanceIpaddressSpec defines the desired state of a IpamTenantNetworkinstanceIpaddress.
type IpamTenantNetworkinstanceIpaddressSpec struct {
	nddv1.ResourceSpec `json:",inline"`
	ForNetworkNode     IpamTenantNetworkinstanceIpaddressParameters `json:"forNetworkNode"`
}

// A IpamTenantNetworkinstanceIpaddressStatus represents the observed state of a IpamTenantNetworkinstanceIpaddress.
type IpamTenantNetworkinstanceIpaddressStatus struct {
	nddv1.ResourceStatus `json:",inline"`
	AtNetworkNode        IpamTenantNetworkinstanceIpaddressObservation `json:"atNetworkNode,omitempty"`
}

// +kubebuilder:object:root=true

// IpamIpamTenantNetworkinstanceIpaddress is the Schema for the IpamTenantNetworkinstanceIpaddress API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="TARGET",type="string",JSONPath=".status.conditions[?(@.kind=='TargetFound')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="LOCALLEAFREF",type="string",JSONPath=".status.conditions[?(@.kind=='InternalLeafrefValidationSuccess')].status"
// +kubebuilder:printcolumn:name="EXTLEAFREF",type="string",JSONPath=".status.conditions[?(@.kind=='ExternalLeafrefValidationSuccess')].status"
// +kubebuilder:printcolumn:name="PARENTDEP",type="string",JSONPath=".status.conditions[?(@.kind=='ParentValidationSuccess')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={ndd,ipam}
type IpamIpamTenantNetworkinstanceIpaddress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IpamTenantNetworkinstanceIpaddressSpec   `json:"spec,omitempty"`
	Status IpamTenantNetworkinstanceIpaddressStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IpamIpamTenantNetworkinstanceIpaddressList contains a list of IpamTenantNetworkinstanceIpaddresss
type IpamIpamTenantNetworkinstanceIpaddressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IpamIpamTenantNetworkinstanceIpaddress `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IpamIpamTenantNetworkinstanceIpaddress{}, &IpamIpamTenantNetworkinstanceIpaddressList{})
}

// IpamTenantNetworkinstanceIpaddress type metadata.
var (
	IpamTenantNetworkinstanceIpaddressKindKind         = reflect.TypeOf(IpamIpamTenantNetworkinstanceIpaddress{}).Name()
	IpamTenantNetworkinstanceIpaddressGroupKind        = schema.GroupKind{Group: Group, Kind: IpamTenantNetworkinstanceIpaddressKindKind}.String()
	IpamTenantNetworkinstanceIpaddressKindAPIVersion   = IpamTenantNetworkinstanceIpaddressKindKind + "." + GroupVersion.String()
	IpamTenantNetworkinstanceIpaddressGroupVersionKind = GroupVersion.WithKind(IpamTenantNetworkinstanceIpaddressKindKind)
)
