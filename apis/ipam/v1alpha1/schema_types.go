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

// Nddoipam struct
type Nddoipam struct {
	Ipam *NddoipamIpam `json:"ipam,omitempty"`
}

// NddoipamIpam struct
type NddoipamIpam struct {
	Rir    []*NddoipamIpamRir    `json:"rir,omitempty"`
	Tenant []*NddoipamIpamTenant `json:"tenant,omitempty"`
}

// NddoipamIpamRir struct
type NddoipamIpamRir struct {
	Description *string               `json:"description,omitempty"`
	Name        *string               `json:"name"`
	Prefixes    *uint32               `json:"prefixes,omitempty"`
	Private     *bool                 `json:"private,omitempty"`
	Tag         []*NddoipamIpamRirTag `json:"tag,omitempty"`
}

// NddoipamIpamRirTag struct
type NddoipamIpamRirTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// NddoipamIpamTenant struct
type NddoipamIpamTenant struct {
	AdminState      *string                              `json:"admin-state,omitempty"`
	Description     *string                              `json:"description,omitempty"`
	Name            *string                              `json:"name"`
	NetworkInstance []*NddoipamIpamTenantNetworkInstance `json:"network-instance,omitempty"`
	Tag             []*NddoipamIpamTenantTag             `json:"tag,omitempty"`
	Reason          *string                              `json:"reason,omitempty"`
	Status          *string                              `json:"status,omitempty"`
	LastUpdate      *string                              `json:"last-update,omitempty"`
}

// NddoipamIpamTenantNetworkInstance struct
type NddoipamIpamTenantNetworkInstance struct {
	AddressAllocationStrategy *string                                       `json:"address-allocation-strategy,omitempty"`
	AdminState                *string                                       `json:"admin-state,omitempty"`
	Description               *string                                       `json:"description,omitempty"`
	IpAddress                 []*NddoipamIpamTenantNetworkInstanceIpAddress `json:"ip-address,omitempty"`
	IpPrefix                  []*NddoipamIpamTenantNetworkInstanceIpPrefix  `json:"ip-prefix,omitempty"`
	IpRange                   []*NddoipamIpamTenantNetworkInstanceIpRange   `json:"ip-range,omitempty"`
	Name                      *string                                       `json:"name"`
	Tag                       []*NddoipamIpamTenantNetworkInstanceTag       `json:"tag,omitempty"`
	Reason                    *string                                       `json:"reason,omitempty"`
	Status                    *string                                       `json:"status,omitempty"`
	LastUpdate                *string                                       `json:"last-update,omitempty"`
}

// NddoipamIpamTenantNetworkInstanceIpAddress struct
type NddoipamIpamTenantNetworkInstanceIpAddress struct {
	Address                   *string                                               `json:"address"`
	AddressAllocationStrategy *string                                               `json:"address-allocation-strategy,omitempty"`
	AdminState                *string                                               `json:"admin-state,omitempty"`
	Description               *string                                               `json:"description,omitempty"`
	DnsName                   *string                                               `json:"dns-name,omitempty"`
	LastUpdate                *string                                               `json:"last-update,omitempty"`
	NatInside                 *string                                               `json:"nat-inside,omitempty"`
	NatOutside                *string                                               `json:"nat-outside,omitempty"`
	Origin                    *string                                               `json:"origin,omitempty"`
	IpPrefix                  []*NddoipamIpamTenantNetworkInstanceIpAddressIpPrefix `json:"ip-prefix,omitempty"`
	IpRange                   []*NddoipamIpamTenantNetworkInstanceIpAddressIpRange  `json:"ip-range,omitempty"`
	Reason                    *string                                               `json:"reason,omitempty"`
	Status                    *string                                               `json:"status,omitempty"`
	Tag                       []*NddoipamIpamTenantNetworkInstanceIpAddressTag      `json:"tag,omitempty"`
}

// NddoipamIpamTenantNetworkInstanceIpAddressIpPrefix struct
type NddoipamIpamTenantNetworkInstanceIpAddressIpPrefix struct {
	Prefix *string `json:"prefix"`
}

// NddoipamIpamTenantNetworkInstanceIpAddressIpRange struct
type NddoipamIpamTenantNetworkInstanceIpAddressIpRange struct {
	End   *string `json:"end"`
	Start *string `json:"start"`
}

// NddoipamIpamTenantNetworkInstanceIpAddressTag struct
type NddoipamIpamTenantNetworkInstanceIpAddressTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// NddoipamIpamTenantNetworkInstanceIpPrefix struct
type NddoipamIpamTenantNetworkInstanceIpPrefix struct {
	AddressAllocationStrategy *string                                          `json:"address-allocation-strategy,omitempty"`
	AdminState                *string                                          `json:"admin-state,omitempty"`
	Adresses                  *uint32                                          `json:"adresses,omitempty"`
	Child                     *NddoipamIpamTenantNetworkInstanceIpPrefixChild  `json:"child,omitempty"`
	Description               *string                                          `json:"description,omitempty"`
	LastUpdate                *string                                          `json:"last-update,omitempty"`
	Parent                    *NddoipamIpamTenantNetworkInstanceIpPrefixParent `json:"parent,omitempty"`
	Pool                      *bool                                            `json:"pool,omitempty"`
	Prefix                    *string                                          `json:"prefix"`
	Reason                    *string                                          `json:"reason,omitempty"`
	RirName                   *string                                          `json:"rir-name,omitempty"`
	Status                    *string                                          `json:"status,omitempty"`
	Tag                       []*NddoipamIpamTenantNetworkInstanceIpPrefixTag  `json:"tag,omitempty"`
}

// NddoipamIpamTenantNetworkInstanceIpPrefixChild struct
type NddoipamIpamTenantNetworkInstanceIpPrefixChild struct {
	IpPrefix []*NddoipamIpamTenantNetworkInstanceIpPrefixChildIpPrefix `json:"ip-prefix,omitempty"`
}

// NddoipamIpamTenantNetworkInstanceIpPrefixChildIpPrefix struct
type NddoipamIpamTenantNetworkInstanceIpPrefixChildIpPrefix struct {
	Prefix *string `json:"prefix"`
}

// NddoipamIpamTenantNetworkInstanceIpPrefixParent struct
type NddoipamIpamTenantNetworkInstanceIpPrefixParent struct {
	IpPrefix []*NddoipamIpamTenantNetworkInstanceIpPrefixParentIpPrefix `json:"ip-prefix,omitempty"`
}

// NddoipamIpamTenantNetworkInstanceIpPrefixParentIpPrefix struct
type NddoipamIpamTenantNetworkInstanceIpPrefixParentIpPrefix struct {
	Prefix *string `json:"prefix"`
}

// NddoipamIpamTenantNetworkInstanceIpPrefixTag struct
type NddoipamIpamTenantNetworkInstanceIpPrefixTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// NddoipamIpamTenantNetworkInstanceIpRange struct
type NddoipamIpamTenantNetworkInstanceIpRange struct {
	AddressAllocationStrategy *string                                         `json:"address-allocation-strategy,omitempty"`
	AdminState                *string                                         `json:"admin-state,omitempty"`
	Description               *string                                         `json:"description,omitempty"`
	End                       *string                                         `json:"end"`
	LastUpdate                *string                                         `json:"last-update,omitempty"`
	Parent                    *NddoipamIpamTenantNetworkInstanceIpRangeParent `json:"parent,omitempty"`
	Reason                    *string                                         `json:"reason,omitempty"`
	Size                      *uint32                                         `json:"size,omitempty"`
	Start                     *string                                         `json:"start"`
	Status                    *string                                         `json:"status,omitempty"`
	Tag                       []*NddoipamIpamTenantNetworkInstanceIpRangeTag  `json:"tag,omitempty"`
}

// NddoipamIpamTenantNetworkInstanceIpRangeParent struct
type NddoipamIpamTenantNetworkInstanceIpRangeParent struct {
	IpPrefix []*NddoipamIpamTenantNetworkInstanceIpRangeParentIpPrefix `json:"ip-prefix,omitempty"`
}

// NddoipamIpamTenantNetworkInstanceIpRangeParentIpPrefix struct
type NddoipamIpamTenantNetworkInstanceIpRangeParentIpPrefix struct {
	Prefix *string `json:"prefix"`
}

// NddoipamIpamTenantNetworkInstanceIpRangeTag struct
type NddoipamIpamTenantNetworkInstanceIpRangeTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// NddoipamIpamTenantNetworkInstanceTag struct
type NddoipamIpamTenantNetworkInstanceTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// NddoipamIpamTenantTag struct
type NddoipamIpamTenantTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// Root is the root of the schema
type Root struct {
	IpamNddoipam *Nddoipam `json:"nddo-ipam,omitempty"`
}
