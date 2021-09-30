package v1alpha1

import (
	"fmt"
	"net"
	"reflect"

	"github.com/apparentlymart/go-cidr/cidr"
)

// +k8s:deepcopy-gen=false

func (n *NddoipamIpamIpPrefix) GetAdminState() string {
	if reflect.ValueOf(n.AdminState).IsZero() {
		return "enabale"
	}
	return *n.AdminState
}

func (n *NddoipamIpamIpPrefix) GetDescription() string {
	if reflect.ValueOf(n.Description).IsZero() {
		return ""
	}
	return *n.Description
}

func (n *NddoipamIpamIpPrefix) GetNetworkInstance() string {
	if reflect.ValueOf(n.NetworkInstance).IsZero() {
		return "default"
	}
	return *n.NetworkInstance
}

func (n *NddoipamIpamIpPrefix) GetPrefix() string {
	if reflect.ValueOf(n.Prefix).IsZero() {
		return ""
	}
	return *n.Prefix
}

func (n *NddoipamIpamIpPrefix) GetAdresses() uint32 {
	if reflect.ValueOf(n.Adresses).IsZero() {
		return 0
	}
	return *n.Adresses
}

func (n *NddoipamIpamIpPrefix) GetStatus() string {
	if reflect.ValueOf(n.Status).IsZero() {
		return ""
	}
	return *n.Status
}

func (n *NddoipamIpamIpPrefix) GetTenant() string {
	if reflect.ValueOf(n.Tenant).IsZero() {
		return ""
	}
	return *n.Tenant
}

func (n *NddoipamIpamIpPrefix) GetPool() bool {
	if reflect.ValueOf(n.Pool).IsZero() {
		return false
	}
	return *n.Pool
}

func (n *NddoipamIpamIpPrefix) SetAdminState(s string) {
	n.AdminState = &s
}

func (n *NddoipamIpamIpPrefix) SetDescription(s string) {
	n.Description = &s
}

func (n *NddoipamIpamIpPrefix) SetNetworkInstance(s string) {
	n.NetworkInstance = &s
}

func (n *NddoipamIpamIpPrefix) SetPrefix(s string) {
	n.Prefix = &s
}

func (n *NddoipamIpamIpPrefix) SetAddresses(s uint32) {
	n.Adresses = &s
}

func (n *NddoipamIpamIpPrefix) SetStatus(s string) {
	n.Status = &s
}

func (n *NddoipamIpamIpPrefix) SetTenant(s string) {
	n.Tenant = &s
}

func (n *NddoipamIpamIpPrefix) SetPool(s bool) {
	n.Pool = &s
}

func (n *NddoipamIpamIpPrefix) GetTags() []*NddoipamIpamIpPrefixTag {
	return n.Tag
}

func (n *NddoipamIpamIpPrefix) GetTag(s string) *NddoipamIpamIpPrefixTag {
	for _, t := range n.GetTags() {
		if t.GetKey() == s {
			return t
		}
	}
	return nil
}

func (n *NddoipamIpamIpPrefix) AppendTag(s *NddoipamIpamIpPrefixTag) {
	found := false
	for _, a := range n.Tag {
		if a.GetKey() == s.GetKey() {
			a.SetValue(s.GetValue())
			found = true
		}
	}
	if !found {
		n.Tag = append(n.Tag, s)
	}
}

func (n *NddoipamIpamIpPrefixTag) GetKey() string {
	if reflect.ValueOf(n.Key).IsZero() {
		return ""
	}
	return *n.Key
}

func (n *NddoipamIpamIpPrefixTag) GetValue() string {
	if reflect.ValueOf(n.Value).IsZero() {
		return ""
	}
	return *n.Value
}

func (n *NddoipamIpamIpPrefixTag) SetValue(s string) {
	n.Value = &s
}

func (n *NddoipamIpamIpPrefix) GetAggregates() []*NddoipamIpamIpPrefixAggregate {
	return n.Aggregate
}

func (n *NddoipamIpamIpPrefix) GetAggregate(t, ni, p string) *NddoipamIpamIpPrefixAggregate {
	for _, a := range n.GetAggregates() {
		if a.GetTenant() == t && a.GetNetworkInstance() == ni && a.GetPrefix() == p {
			return a
		}
	}
	return nil
}

func (n *NddoipamIpamIpPrefix) AppendAggregate(s *NddoipamIpamIpPrefixAggregate) {
	found := false
	for _, a := range n.Aggregate {
		if a.GetNetworkInstance() == s.GetNetworkInstance() && a.GetPrefix() == s.GetPrefix() && a.GetTenant() == s.GetTenant() {
			found = true
		}
	}
	if !found {
		n.Aggregate = append(n.Aggregate, s)
	}
}

func (n *NddoipamIpamIpPrefixAggregate) GetTenant() string {
	if reflect.ValueOf(n.Tenant).IsZero() {
		return ""
	}
	return *n.Tenant
}

func (n *NddoipamIpamIpPrefixAggregate) GetNetworkInstance() string {
	if reflect.ValueOf(n.NetworkInstance).IsZero() {
		return ""
	}
	return *n.NetworkInstance
}

func (n *NddoipamIpamIpPrefixAggregate) GetPrefix() string {
	if reflect.ValueOf(n.Prefix).IsZero() {
		return ""
	}
	return *n.Prefix
}

func (n *NddoipamIpamIpPrefix) Contains(ip net.IP) (bool, error) {
	_, subnet, err := net.ParseCIDR(n.GetPrefix())
	if err != nil {
		return false, err
	}
	return subnet.Contains(ip), nil
}

// GetPrefixStatusFromOrigin updates the status from the origin if it exists
func (n *NddoipamIpamIpPrefix) GetStatusFromOrigin(origPrefixes []*NddoipamIpamIpPrefix) {
	for _, o := range origPrefixes {
		if o.GetTenant() == n.GetTenant() &&
			o.GetNetworkInstance() == n.GetNetworkInstance() &&
			o.GetPrefix() == n.GetPrefix() {
			// update status if it existed
			n.SetStatus(o.GetStatus())
		}

	}
}

// ProcessOverlapWithPrefix updates the status if the status was NOT active
func (n *NddoipamIpamIpPrefix) ProcessOverlapWithPrefix(origPrefixes []*NddoipamIpamIpPrefix) error {
	// only evaluate the status that is not active
	if n.GetStatus() != "active" {
		overlap := false
		parent := false
		for _, o := range origPrefixes {
			// evaluate against status that is active
			if o.GetStatus() == "active" {
				// in order to check address overlap the tenant and network-instance should match
				if n.GetTenant() == o.GetTenant() &&
					n.GetNetworkInstance() == o.GetNetworkInstance() {
					// dont compare with yourself, should never happen anyhow since the status should be != active
					if n.GetPrefix() != o.GetPrefix() {
						// get the subnet
						_, subnet, err := net.ParseCIDR(n.GetPrefix())
						if err != nil {
							return err
						}
						// find the start address and end address of the subnet
						ipStart, ipEnd := cidr.AddressRange(subnet)
						// check if the start ip address is contained in the prefix
						cStart, err := o.Contains(ipStart)
						if err != nil {
							return err
						}
						cStop, err := o.Contains(ipEnd)
						if err != nil {
							return err
						}
						// there should not be any overlap between aggregates
						if cStart && cStop {
							overlap = true
							parent = true
							// complete overlap is ok
							// TODO add parents -> getParents from parent prefix and add the new parent to it
						} else {
							if cStart || cStop {
								overlap = true
								// TODO set inactive reason to overlap
							}
						}
					}
				}
			}
		}
		// when none overlap we can set the status to active
		switch {
		case overlap && parent:
			// complete overlap is ok
			n.SetStatus("active")
		case overlap && !parent:
			// half overlap is nok
			n.SetStatus("inactive")
		case !overlap:
			// no overlap means we need to check if an aggregate overlaps
			n.SetStatus("check parent aggregate")
		}
	}
	return nil
}

// ProcessOverlapWithPrefix updates the status if the status was NOT active
func (n *NddoipamIpamIpPrefix) ProcessOverlapWithAggregate(origAggregate []*NddoipamIpamAggregate) error {
	// only evaluate the status that is not active
	if n.GetStatus() != "active" {
		overlap := false
		parent := false
		for _, o := range origAggregate {
			// evaluate against status that is active
			if o.GetStatus() == "active" {
				// in order to check address overlap the tenant and network-instance should match
				if n.GetTenant() == o.GetTenant() &&
					n.GetNetworkInstance() == o.GetNetworkInstance() {
					// when comparing against aggregates we can have exact match of prefixes, which is ok, but one is an aggregate
					//, the other is a prefix so we dont have to check if they are the same or not
					// get the subnet
					_, subnet, err := net.ParseCIDR(n.GetPrefix())
					if err != nil {
						return err
					}
					// find the start address and end address of the subnet
					ipStart, ipEnd := cidr.AddressRange(subnet)
					// check if the start ip address is contained in the prefix
					cStart, err := o.Contains(ipStart)
					if err != nil {
						return err
					}
					cEnd, err := o.Contains(ipEnd)
					if err != nil {
						return err
					}
					fmt.Printf("ProcessOverlapWithAggregate Prefix: %s, Aggregate: %s ipStart: %s, ipEnd: %s, cStart: %t, cEnd: %t\n", n.GetPrefix(), o.GetPrefix(), ipStart, ipEnd, cStart, cEnd)
					// there should not be any overlap between aggregates
					if cStart && cEnd {
						overlap = true
						parent = true
						// complete overlap is ok
						// TODO add parents -> getParents from parent prefix and add the new parent to it
						n.SetStatus("active")
						return nil
					} else {
						if cStart || cEnd {
							overlap = true
							// TODO set inactive reason to overlap
							n.SetStatus("inactive")
							return nil
						}
					}
				}
			}
		}
		switch {
		case overlap && parent:
			// complete overlap is ok
			n.SetStatus("active")
		case overlap && !parent:
			// half overlap is nok
			n.SetStatus("inactive")
		case !overlap:
			// no overlap means we need to check if an aggregate overlaps
			n.SetStatus("inactive")
		}
	}
	return nil
}
