package v1alpha1

import "reflect"

// +k8s:deepcopy-gen=false

func (n *NddoipamIpamRir) GetDescription() string {
	if reflect.ValueOf(n.Description).IsZero() {
		return ""
	}
	return *n.Description
}

func (n *NddoipamIpamRir) GetName() string {
	if reflect.ValueOf(n.Name).IsZero() {
		return ""
	}
	return *n.Name
}

func (n *NddoipamIpamRir) GetPrivate() bool {
	if reflect.ValueOf(n.Private).IsZero() {
		return false
	}
	return *n.Private
}

func (n *NddoipamIpamRir) SetDescription(s string) {
	n.Description = &s
}

func (n *NddoipamIpamRir) SetName(s string) {
	n.Name = &s
}

func (n *NddoipamIpamRir) SetPrivate(s bool) {
	n.Private = &s
}

func (n *NddoipamIpamRir) GetTags() []*NddoipamIpamRirTag {
	return n.Tag
}

func (n *NddoipamIpamRir) GetTag(s string) *NddoipamIpamRirTag {
	for _, t := range n.GetTags() {
		if t.GetKey() == s {
			return t
		}
	}
	return nil
}

func (n *NddoipamIpamRir) AppendTag(s *NddoipamIpamRirTag) {
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

func (n *NddoipamIpamRirTag) GetKey() string {
	if reflect.ValueOf(n.Key).IsZero() {
		return ""
	}
	return *n.Key
}

func (n *NddoipamIpamRirTag) GetValue() string {
	if reflect.ValueOf(n.Value).IsZero() {
		return ""
	}
	return *n.Value
}

func (n *NddoipamIpamRirTag) SetValue(s string) {
	n.Value = &s
}
