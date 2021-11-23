package v1alpha1

import (
	"reflect"
)

// +k8s:deepcopy-gen=false

func (n *NddoipamIpamTenantNetworkInstanceIpPrefix) GetAdminState() string {
	if reflect.ValueOf(n.AdminState).IsZero() {
		return "enabale"
	}
	return *n.AdminState
}

func (n *NddoipamIpamTenantNetworkInstanceIpPrefix) GetDescription() string {
	if reflect.ValueOf(n.Description).IsZero() {
		return ""
	}
	return *n.Description
}

func (n *NddoipamIpamTenantNetworkInstanceIpPrefix) GetPrefix() string {
	if reflect.ValueOf(n.Prefix).IsZero() {
		return ""
	}
	return *n.Prefix
}

func (n *NddoipamIpamTenantNetworkInstanceIpPrefix) GetStatus() string {
	if reflect.ValueOf(n.Status).IsZero() {
		return ""
	}
	return *n.Status
}

func (n *NddoipamIpamTenantNetworkInstanceIpPrefix) SetAdminState(s string) {
	n.AdminState = &s
}

func (n *NddoipamIpamTenantNetworkInstanceIpPrefix) SetDescription(s string) {
	n.Description = &s
}

func (n *NddoipamIpamTenantNetworkInstanceIpPrefix) SetPrefix(s string) {
	n.Prefix = &s
}

func (n *NddoipamIpamTenantNetworkInstanceIpPrefix) SetStatus(s string) {
	n.Status = &s
}

func (n *NddoipamIpamTenantNetworkInstanceIpPrefix) SetReason(s string) {
	n.Reason = &s
}

func (n *NddoipamIpamTenantNetworkInstanceIpPrefix) SetLastUpdate(s string) {
	n.LastUpdate = &s
}

func (n *NddoipamIpamTenantNetworkInstanceIpPrefix) GetTags() []*NddoipamIpamTenantNetworkInstanceIpPrefixTag {
	return n.Tag
}
