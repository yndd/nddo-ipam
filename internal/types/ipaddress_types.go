package types

import "reflect"

func (n *NddoipamIpamIpAddress) GetAdminState() string {
	if reflect.ValueOf(n.AdminState).IsZero() {
		return "enabale"
	}
	return *n.AdminState
}

func (n *NddoipamIpamIpAddress) GetDescription() string {
	if reflect.ValueOf(n.Description).IsZero() {
		return ""
	}
	return *n.Description
}

func (n *NddoipamIpamIpAddress) GetNetworkInstance() string {
	if reflect.ValueOf(n.NetworkInstance).IsZero() {
		return "default"
	}
	return *n.NetworkInstance
}

func (n *NddoipamIpamIpAddress) GetAddress() string {
	if reflect.ValueOf(n.Address).IsZero() {
		return ""
	}
	return *n.Address
}

func (n *NddoipamIpamIpAddress) GetOrigin() string {
	if reflect.ValueOf(n.Origin).IsZero() {
		return ""
	}
	return *n.Origin
}

func (n *NddoipamIpamIpAddress) GetStatus() string {
	if reflect.ValueOf(n.Status).IsZero() {
		return ""
	}
	return *n.Status
}

func (n *NddoipamIpamIpAddress) GetTenant() string {
	if reflect.ValueOf(n.Tenant).IsZero() {
		return ""
	}
	return *n.Tenant
}

func (n *NddoipamIpamIpAddress) GetDnsName() string {
	if reflect.ValueOf(n.DnsName).IsZero() {
		return ""
	}
	return *n.DnsName
}

func (n *NddoipamIpamIpAddress) GetNatInside() string {
	if reflect.ValueOf(n.NatInside).IsZero() {
		return ""
	}
	return *n.NatInside
}

func (n *NddoipamIpamIpAddress) GetNatOutside() string {
	if reflect.ValueOf(n.NatOutside).IsZero() {
		return ""
	}
	return *n.NatOutside
}

func (n *NddoipamIpamIpAddress) SetAdminState(s string) {
	n.AdminState = &s
}

func (n *NddoipamIpamIpAddress) SetDescription(s string) {
	n.Description = &s
}

func (n *NddoipamIpamIpAddress) SetNetworkInstance(s string) {
	n.NetworkInstance = &s
}

func (n *NddoipamIpamIpAddress) SetAddress(s string) {
	n.Address = &s
}

func (n *NddoipamIpamIpAddress) SetOrigin(s string) {
	n.Origin = &s
}

func (n *NddoipamIpamIpAddress) SetStatus(s string) {
	n.Status = &s
}

func (n *NddoipamIpamIpAddress) SetTenant(s string) {
	n.Tenant = &s
}

func (n *NddoipamIpamIpAddress) SetDnsName(s string) {
	n.DnsName = &s
}

func (n *NddoipamIpamIpAddress) SetNatInside(s string) {
	n.NatInside = &s
}

func (n *NddoipamIpamIpAddress) SetNatOutside(s string) {
	n.NatOutside = &s
}

func (n *NddoipamIpamIpAddress) GetTags() []*NddoipamIpamIpAddressTag {
	return n.Tag
}

func (n *NddoipamIpamIpAddress) GetTag(s string) *NddoipamIpamIpAddressTag {
	for _, t := range n.GetTags() {
		if t.GetKey() == s {
			return t
		}
	}
	return nil
}

func (n *NddoipamIpamIpAddress) AppendTag(s *NddoipamIpamIpAddressTag) {
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

func (n *NddoipamIpamIpAddressTag) GetKey() string {
	if reflect.ValueOf(n.Key).IsZero() {
		return ""
	}
	return *n.Key
}

func (n *NddoipamIpamIpAddressTag) GetValue() string {
	if reflect.ValueOf(n.Value).IsZero() {
		return ""
	}
	return *n.Value
}

func (n *NddoipamIpamIpAddressTag) SetValue(s string) {
	n.Value = &s
}

func (n *NddoipamIpamIpAddress) GetIPRanges() []*NddoipamIpamIpAddressIpRange {
	return n.IpRange
}

func (n *NddoipamIpamIpAddress) GetIPRange(t, ni, s, e string) *NddoipamIpamIpAddressIpRange {
	for _, a := range n.GetIPRanges() {
		if a.GetTenant() == t && a.GetNetworkInstance() == ni && a.GetStart() == s && a.GetEnd() == e {
			return a
		}
	}
	return nil
}

func (n *NddoipamIpamIpAddress) AppendIPRange(s *NddoipamIpamIpAddressIpRange) {
	found := false
	for _, a := range n.IpRange {
		if a.GetNetworkInstance() == s.GetNetworkInstance() && a.GetStart() == s.GetStart() && a.GetEnd() == s.GetEnd() && a.GetTenant() == s.GetTenant() {
			found = true
		}
	}
	if !found {
		n.IpRange = append(n.IpRange, s)
	}
}

func (n *NddoipamIpamIpAddressIpRange) GetTenant() string {
	if reflect.ValueOf(n.Tenant).IsZero() {
		return ""
	}
	return *n.Tenant
}

func (n *NddoipamIpamIpAddressIpRange) GetNetworkInstance() string {
	if reflect.ValueOf(n.NetworkInstance).IsZero() {
		return ""
	}
	return *n.NetworkInstance
}

func (n *NddoipamIpamIpAddressIpRange) GetStart() string {
	if reflect.ValueOf(n.Start).IsZero() {
		return ""
	}
	return *n.Start
}

func (n *NddoipamIpamIpAddressIpRange) GetEnd() string {
	if reflect.ValueOf(n.End).IsZero() {
		return ""
	}
	return *n.End
}

func (n *NddoipamIpamIpAddress) GetIPPrefixes() []*NddoipamIpamIpAddressIpPrefix {
	return n.IpPrefix
}

func (n *NddoipamIpamIpAddress) GetIPPrefix(t, ni, p string) *NddoipamIpamIpAddressIpPrefix {
	for _, a := range n.GetIPPrefixes() {
		if a.GetTenant() == t && a.GetNetworkInstance() == ni && a.GetPrefix() == p {
			return a
		}
	}
	return nil
}

func (n *NddoipamIpamIpAddress) AppendIPPrefix(s *NddoipamIpamIpAddressIpPrefix) {
	found := false
	for _, a := range n.IpPrefix {
		if a.GetNetworkInstance() == s.GetNetworkInstance() && a.GetPrefix() == s.GetPrefix() && a.GetTenant() == s.GetTenant() {
			found = true
		}
	}
	if !found {
		n.IpPrefix = append(n.IpPrefix, s)
	}
}

func (n *NddoipamIpamIpAddressIpPrefix) GetTenant() string {
	if reflect.ValueOf(n.Tenant).IsZero() {
		return ""
	}
	return *n.Tenant
}

func (n *NddoipamIpamIpAddressIpPrefix) GetNetworkInstance() string {
	if reflect.ValueOf(n.NetworkInstance).IsZero() {
		return ""
	}
	return *n.NetworkInstance
}

func (n *NddoipamIpamIpAddressIpPrefix) GetPrefix() string {
	if reflect.ValueOf(n.Prefix).IsZero() {
		return ""
	}
	return *n.Prefix
}
