package v1alpha1

import "reflect"

// +k8s:deepcopy-gen=false

func (n *NddoipamIpamTenantNetworkInstance) GetAddressAllocationStrategy() string {
	if reflect.ValueOf(n.AddressAllocationStrategy).IsZero() {
		return "enabale"
	}
	return *n.AddressAllocationStrategy
}

func (n *NddoipamIpamTenantNetworkInstance) GetAdminState() string {
	if reflect.ValueOf(n.AdminState).IsZero() {
		return "enabale"
	}
	return *n.AdminState
}

func (n *NddoipamIpamTenantNetworkInstance) GetDescription() string {
	if reflect.ValueOf(n.Description).IsZero() {
		return ""
	}
	return *n.Description
}

func (n *NddoipamIpamTenantNetworkInstance) GetName() string {
	if reflect.ValueOf(n.Name).IsZero() {
		return ""
	}
	return *n.Name
}

func (n *NddoipamIpamTenantNetworkInstance) GetStatus() string {
	if reflect.ValueOf(n.Status).IsZero() {
		return ""
	}
	return *n.Status
}

func (n *NddoipamIpamTenantNetworkInstance) GetReason() string {
	if reflect.ValueOf(n.Reason).IsZero() {
		return ""
	}
	return *n.Reason
}

func (n *NddoipamIpamTenantNetworkInstance) GetLastUpdate() string {
	if reflect.ValueOf(n.LastUpdate).IsZero() {
		return ""
	}
	return *n.LastUpdate
}

func (n *NddoipamIpamTenantNetworkInstance) SetAdminState(s string) {
	n.AdminState = &s
}

func (n *NddoipamIpamTenantNetworkInstance) SetDescription(s string) {
	n.Description = &s
}

func (n *NddoipamIpamTenantNetworkInstance) SetName(s string) {
	n.Name = &s
}

func (n *NddoipamIpamTenantNetworkInstance) SetStatus(s string) {
	n.Status = &s
}

func (n *NddoipamIpamTenantNetworkInstance) SetReason(s string) {
	n.Reason = &s
}

func (n *NddoipamIpamTenantNetworkInstance) SetLastUpdate(s string) {
	n.LastUpdate = &s
}

func (n *NddoipamIpamTenantNetworkInstance) GetTags() []*NddoipamIpamTenantNetworkInstanceTag {
	return n.Tag
}

func (n *NddoipamIpamTenantNetworkInstance) GetTag(s string) *NddoipamIpamTenantNetworkInstanceTag {
	for _, t := range n.GetTags() {
		if t.GetKey() == s {
			return t
		}
	}
	return nil
}

func (n *NddoipamIpamTenantNetworkInstance) AppendTag(s *NddoipamIpamTenantNetworkInstanceTag) {
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

func (n *NddoipamIpamTenantNetworkInstanceTag) GetKey() string {
	if reflect.ValueOf(n.Key).IsZero() {
		return ""
	}
	return *n.Key
}

func (n *NddoipamIpamTenantNetworkInstanceTag) GetValue() string {
	if reflect.ValueOf(n.Value).IsZero() {
		return ""
	}
	return *n.Value
}

func (n *NddoipamIpamTenantNetworkInstanceTag) SetValue(s string) {
	n.Value = &s
}
