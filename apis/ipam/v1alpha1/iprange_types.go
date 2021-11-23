package v1alpha1

import (
	"reflect"
)

// +k8s:deepcopy-gen=false

func (n *NddoipamIpamTenantNetworkInstanceIpRange) GetAdminState() string {
	if reflect.ValueOf(n.AdminState).IsZero() {
		return "enabale"
	}
	return *n.AdminState
}

func (n *NddoipamIpamTenantNetworkInstanceIpRange) GetDescription() string {
	if reflect.ValueOf(n.Description).IsZero() {
		return ""
	}
	return *n.Description
}

func (n *NddoipamIpamTenantNetworkInstanceIpRange) GetStatus() string {
	if reflect.ValueOf(n.Status).IsZero() {
		return ""
	}
	return *n.Status
}

func (n *NddoipamIpamTenantNetworkInstanceIpRange) SetAdminState(s string) {
	n.AdminState = &s
}

func (n *NddoipamIpamTenantNetworkInstanceIpRange) SetDescription(s string) {
	n.Description = &s
}

func (n *NddoipamIpamTenantNetworkInstanceIpRange) SetStatus(s string) {
	n.Status = &s
}

func (n *NddoipamIpamTenantNetworkInstanceIpRange) SetReason(s string) {
	n.Reason = &s
}

func (n *NddoipamIpamTenantNetworkInstanceIpRange) SetLastUpdate(s string) {
	n.LastUpdate = &s
}

func (n *NddoipamIpamTenantNetworkInstanceIpRange) GetTags() []*NddoipamIpamTenantNetworkInstanceIpRangeTag {
	return n.Tag
}
