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
	Aggregate []*NddoipamIpamAggregate `json:"aggregate,omitempty"`
	IpAddress []*NddoipamIpamIpAddress `json:"ip-address,omitempty"`
	IpPrefix  []*NddoipamIpamIpPrefix  `json:"ip-prefix,omitempty"`
	IpRange   []*NddoipamIpamIpRange   `json:"ip-range,omitempty"`
	Rir       []*NddoipamIpamRir       `json:"rir,omitempty"`
}

// NddoipamIpamAggregate struct
type NddoipamIpamAggregate struct {
	AdminState      *string                     `json:"admin-state,omitempty"`
	Description     *string                     `json:"description,omitempty"`
	LastUpdate      *string                     `json:"last-update,omitempty"`
	NetworkInstance *string                     `json:"network-instance"`
	Prefix          *string                     `json:"prefix"`
	Prefixes        *uint32                     `json:"prefixes,omitempty"`
	Reason          *string                     `json:"reason,omitempty"`
	RirName         *string                     `json:"rir-name"`
	Status          *string                     `json:"status,omitempty"`
	Tag             []*NddoipamIpamAggregateTag `json:"tag,omitempty"`
	Tenant          *string                     `json:"tenant"`
}

// NddoipamIpamAggregateTag struct
type NddoipamIpamAggregateTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// NddoipamIpamIpAddress struct
type NddoipamIpamIpAddress struct {
	Address         *string                          `json:"address"`
	AdminState      *string                          `json:"admin-state,omitempty"`
	Description     *string                          `json:"description,omitempty"`
	DnsName         *string                          `json:"dns-name,omitempty"`
	LastUpdate      *string                          `json:"last-update,omitempty"`
	NatInside       *string                          `json:"nat-inside,omitempty"`
	NatOutside      *string                          `json:"nat-outside,omitempty"`
	NetworkInstance *string                          `json:"network-instance"`
	Origin          *string                          `json:"origin,omitempty"`
	IpPrefix        []*NddoipamIpamIpAddressIpPrefix `json:"ip-prefix,omitempty"`
	IpRange         []*NddoipamIpamIpAddressIpRange  `json:"ip-range,omitempty"`
	Reason          *string                          `json:"reason,omitempty"`
	Status          *string                          `json:"status,omitempty"`
	Tag             []*NddoipamIpamIpAddressTag      `json:"tag,omitempty"`
	Tenant          *string                          `json:"tenant"`
}

// NddoipamIpamIpAddressIpPrefix struct
type NddoipamIpamIpAddressIpPrefix struct {
	NetworkInstance *string `json:"network-instance"`
	Prefix          *string `json:"prefix"`
	Tenant          *string `json:"tenant"`
}

// NddoipamIpamIpAddressIpRange struct
type NddoipamIpamIpAddressIpRange struct {
	End             *string `json:"end"`
	NetworkInstance *string `json:"network-instance"`
	Start           *string `json:"start"`
	Tenant          *string `json:"tenant"`
}

// NddoipamIpamIpAddressTag struct
type NddoipamIpamIpAddressTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// NddoipamIpamIpPrefix struct
type NddoipamIpamIpPrefix struct {
	AdminState      *string                          `json:"admin-state,omitempty"`
	Adresses        *uint32                          `json:"adresses,omitempty"`
	Description     *string                          `json:"description,omitempty"`
	LastUpdate      *string                          `json:"last-update,omitempty"`
	NetworkInstance *string                          `json:"network-instance"`
	Aggregate       []*NddoipamIpamIpPrefixAggregate `json:"aggregate,omitempty"`
	IpPrefix        []*NddoipamIpamIpPrefixIpPrefix  `json:"ip-prefix,omitempty"`
	Pool            *bool                            `json:"pool,omitempty"`
	Prefix          *string                          `json:"prefix"`
	Reason          *string                          `json:"reason,omitempty"`
	Status          *string                          `json:"status,omitempty"`
	Tag             []*NddoipamIpamIpPrefixTag       `json:"tag,omitempty"`
	Tenant          *string                          `json:"tenant"`
}

// NddoipamIpamIpPrefixAggregate struct
type NddoipamIpamIpPrefixAggregate struct {
	NetworkInstance *string `json:"network-instance"`
	Prefix          *string `json:"prefix"`
	Tenant          *string `json:"tenant"`
}

// NddoipamIpamIpPrefixIpPrefix struct
type NddoipamIpamIpPrefixIpPrefix struct {
	NetworkInstance *string `json:"network-instance"`
	Prefix          *string `json:"prefix"`
	Tenant          *string `json:"tenant"`
}

// NddoipamIpamIpPrefixTag struct
type NddoipamIpamIpPrefixTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// NddoipamIpamIpRange struct
type NddoipamIpamIpRange struct {
	AdminState      *string                         `json:"admin-state,omitempty"`
	Description     *string                         `json:"description,omitempty"`
	End             *string                         `json:"end"`
	LastUpdate      *string                         `json:"last-update,omitempty"`
	NetworkInstance *string                         `json:"network-instance"`
	Aggregate       []*NddoipamIpamIpRangeAggregate `json:"aggregate,omitempty"`
	IpPrefix        []*NddoipamIpamIpRangeIpPrefix  `json:"ip-prefix,omitempty"`
	Reason          *string                         `json:"reason,omitempty"`
	Size            *uint32                         `json:"size,omitempty"`
	Start           *string                         `json:"start"`
	Status          *string                         `json:"status,omitempty"`
	Tag             []*NddoipamIpamIpRangeTag       `json:"tag,omitempty"`
	Tenant          *string                         `json:"tenant"`
}

// NddoipamIpamIpRangeAggregate struct
type NddoipamIpamIpRangeAggregate struct {
	NetworkInstance *string `json:"network-instance"`
	Prefix          *string `json:"prefix"`
	Tenant          *string `json:"tenant"`
}

// NddoipamIpamIpRangeIpPrefix struct
type NddoipamIpamIpRangeIpPrefix struct {
	NetworkInstance *string `json:"network-instance"`
	Prefix          *string `json:"prefix"`
	Tenant          *string `json:"tenant"`
}

// NddoipamIpamIpRangeTag struct
type NddoipamIpamIpRangeTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// NddoipamIpamRir struct
type NddoipamIpamRir struct {
	Aggregates  *uint32               `json:"aggregates,omitempty"`
	Description *string               `json:"description,omitempty"`
	Name        *string               `json:"name"`
	Private     *bool                 `json:"private,omitempty"`
	Tag         []*NddoipamIpamRirTag `json:"tag,omitempty"`
}

// NddoipamIpamRirTag struct
type NddoipamIpamRirTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value,omitempty"`
}

// Root is the root of the schema
type Root struct {
	IpamNddoipam *Nddoipam `json:"nddo-ipam,omitempty"`
}
