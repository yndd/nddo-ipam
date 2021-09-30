package v1alpha1

import (
	"net"
	"reflect"

	"github.com/apparentlymart/go-cidr/cidr"
)

// +k8s:deepcopy-gen=false

func (n *NddoipamIpamAggregate) GetAdminState() string {
	if reflect.ValueOf(n.AdminState).IsZero() {
		return "enabale"
	}
	return *n.AdminState
}

func (n *NddoipamIpamAggregate) GetDescription() string {
	if reflect.ValueOf(n.Description).IsZero() {
		return ""
	}
	return *n.Description
}

func (n *NddoipamIpamAggregate) GetNetworkInstance() string {
	if reflect.ValueOf(n.NetworkInstance).IsZero() {
		return "default"
	}
	return *n.NetworkInstance
}

func (n *NddoipamIpamAggregate) GetPrefix() string {
	if reflect.ValueOf(n.Prefix).IsZero() {
		return ""
	}
	return *n.Prefix
}

func (n *NddoipamIpamAggregate) GetPrefixes() uint32 {
	if reflect.ValueOf(n.Prefixes).IsZero() {
		return 0
	}
	return *n.Prefixes
}

func (n *NddoipamIpamAggregate) GetRirName() string {
	if reflect.ValueOf(n.RirName).IsZero() {
		return ""
	}
	return *n.RirName
}

func (n *NddoipamIpamAggregate) GetStatus() string {
	if reflect.ValueOf(n.Status).IsZero() {
		return ""
	}
	return *n.Status
}

func (n *NddoipamIpamAggregate) GetTenant() string {
	if reflect.ValueOf(n.Tenant).IsZero() {
		return ""
	}
	return *n.Tenant
}

func (n *NddoipamIpamAggregate) SetAdminState(s string) {
	n.AdminState = &s
}

func (n *NddoipamIpamAggregate) SetDescription(s string) {
	n.Description = &s
}

func (n *NddoipamIpamAggregate) SetNetworkInstance(s string) {
	n.NetworkInstance = &s
}

func (n *NddoipamIpamAggregate) SetPrefix(s string) {
	n.Prefix = &s
}

func (n *NddoipamIpamAggregate) SetPrefixes(s uint32) {
	n.Prefixes = &s
}

func (n *NddoipamIpamAggregate) SetStatus(s string) {
	n.Status = &s
}

func (n *NddoipamIpamAggregate) SetTenant(s string) {
	n.Tenant = &s
}

func (n *NddoipamIpamAggregate) GetTags() []*NddoipamIpamAggregateTag {
	return n.Tag
}

func (n *NddoipamIpamAggregate) GetTag(s string) *NddoipamIpamAggregateTag {
	for _, t := range n.GetTags() {
		if t.GetKey() == s {
			return t
		}
	}
	return nil
}

func (n *NddoipamIpamAggregate) AppendTag(s *NddoipamIpamAggregateTag) {
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

func (n *NddoipamIpamAggregateTag) GetKey() string {
	if reflect.ValueOf(n.Key).IsZero() {
		return ""
	}
	return *n.Key
}

func (n *NddoipamIpamAggregateTag) GetValue() string {
	if reflect.ValueOf(n.Value).IsZero() {
		return ""
	}
	return *n.Value
}

func (n *NddoipamIpamAggregateTag) SetValue(s string) {
	n.Value = &s
}

func (n *NddoipamIpamAggregate) Contains(ip net.IP) (bool, error) {
	_, subnet, err := net.ParseCIDR(n.GetPrefix())
	if err != nil {
		return false, err
	}
	return subnet.Contains(ip), nil
}

// GetAggregateStatusFromOrigin updates the status from the origin if it exists
func (n *NddoipamIpamAggregate) GetStatusFromOrigin(origAggregate []*NddoipamIpamAggregate) {
	for _, o := range origAggregate {
		if o.GetTenant() == n.GetTenant() &&
			o.GetNetworkInstance() == n.GetNetworkInstance() &&
			o.GetPrefix() == n.GetPrefix() {
			// update status if it existed
			n.SetStatus(o.GetStatus())
		}
	}
}

// GetAggregateStatusFromOrigin updates the status from the origin if it exists
func (n *NddoipamIpamAggregate) Exists(o *NddoipamIpamAggregate) bool {
	if o.GetTenant() == n.GetTenant() &&
		o.GetNetworkInstance() == n.GetNetworkInstance() &&
		o.GetPrefix() == n.GetPrefix() {
		// found
		return true
	}

	return false
}

func (n *NddoipamIpamAggregate) ProcessOverlapWithAggregate(origAggregate []*NddoipamIpamAggregate) error {
	// only evaluate the status that is not active
	if n.GetStatus() != "active" {
		overlap := false
		for _, o := range origAggregate {
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
						if cStart || cStop {
							overlap = true
							n.SetStatus("inactive")
							// TODO set inactive reason to overlap
							return nil
						}
					}
				}
			}
		}
		// when none overlap we can set the status to active
		if !overlap {
			n.SetStatus("active")
		}
	}
	return nil
}
