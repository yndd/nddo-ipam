package v1alpha1

import "reflect"

// +k8s:deepcopy-gen=false

func (n *NddoipamIpamTenant) GetAdminState() string {
	if reflect.ValueOf(n.AdminState).IsZero() {
		return "enable"
	}
	return *n.AdminState
}

func (n *NddoipamIpamTenant) GetDescription() string {
	if reflect.ValueOf(n.Description).IsZero() {
		return ""
	}
	return *n.Description
}

func (n *NddoipamIpamTenant) GetName() string {
	if reflect.ValueOf(n.Name).IsZero() {
		return ""
	}
	return *n.Name
}

func (n *NddoipamIpamTenant) GetStatus() string {
	if reflect.ValueOf(n.Status).IsZero() {
		return ""
	}
	return *n.Status
}

func (n *NddoipamIpamTenant) GetReason() string {
	if reflect.ValueOf(n.Reason).IsZero() {
		return ""
	}
	return *n.Reason
}

func (n *NddoipamIpamTenant) GetLastUpdate() string {
	if reflect.ValueOf(n.LastUpdate).IsZero() {
		return ""
	}
	return *n.LastUpdate
}

func (n *NddoipamIpamTenant) SetAdminState(s string) {
	n.AdminState = &s
}

func (n *NddoipamIpamTenant) SetDescription(s string) {
	n.Description = &s
}

func (n *NddoipamIpamTenant) SetName(s string) {
	n.Name = &s
}

func (n *NddoipamIpamTenant) SetStatus(s string) {
	n.Status = &s
}

func (n *NddoipamIpamTenant) SetReason(s string) {
	n.Reason = &s
}

func (n *NddoipamIpamTenant) SetLastUpdate(s string) {
	n.LastUpdate = &s
}

func (n *NddoipamIpamTenant) GetTags() []*NddoipamIpamTenantTag {
	return n.Tag
}

func (n *NddoipamIpamTenant) GetTag(s string) *NddoipamIpamTenantTag {
	for _, t := range n.GetTags() {
		if t.GetKey() == s {
			return t
		}
	}
	return nil
}

func (n *NddoipamIpamTenant) AppendTag(s *NddoipamIpamTenantTag) {
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

func (n *NddoipamIpamTenantTag) GetKey() string {
	if reflect.ValueOf(n.Key).IsZero() {
		return ""
	}
	return *n.Key
}

func (n *NddoipamIpamTenantTag) GetValue() string {
	if reflect.ValueOf(n.Value).IsZero() {
		return ""
	}
	return *n.Value
}

func (n *NddoipamIpamTenantTag) SetValue(s string) {
	n.Value = &s
}
