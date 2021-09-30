package v1alpha1

func (n *Nddoipam) GetIpam() *NddoipamIpam {
	if n.Ipam == nil {
		return &NddoipamIpam{
			Rir:       make([]*NddoipamIpamRir, 0),
			Aggregate: make([]*NddoipamIpamAggregate, 0),
			IpPrefix:  make([]*NddoipamIpamIpPrefix, 0),
			IpRange:   make([]*NddoipamIpamIpRange, 0),
			IpAddress: make([]*NddoipamIpamIpAddress, 0),
		}
	}
	return n.Ipam
}

func (n *Nddoipam) SetIpam(i *NddoipamIpam) {
	n.Ipam = i
}

func (n *Nddoipam) DeleteIpam() {
	n.Ipam.Rir = make([]*NddoipamIpamRir, 0)
	n.Ipam.Aggregate = make([]*NddoipamIpamAggregate, 0)
}

func (n *NddoipamIpam) GetRir() []*NddoipamIpamRir {
	return n.Rir
}

func (n *NddoipamIpam) GetAggregate() []*NddoipamIpamAggregate {
	return n.Aggregate
}

func (n *NddoipamIpam) GetIpPrefix() []*NddoipamIpamIpPrefix {
	return n.IpPrefix
}

func (n *NddoipamIpam) AppendRir(s []*NddoipamIpamRir) {
	for _, new := range s {
		found := false
		for _, orig := range n.Rir {
			if orig.GetName() == new.GetName() {
				//orig = new
				found = true
			}
		}
		if !found {
			n.Rir = append(n.Rir, new)
		}
	}
}

func (n *NddoipamIpam) DeleteRir(s []*NddoipamIpamRir) {
	for _, new := range s {
		found := false
		index := -1
		for i, orig := range n.Rir {
			if orig.GetName() == new.GetName() {
				//orig = new
				found = true
				index = i
			}
		}
		if found {
			// delete entry
			n.Rir = append(n.Rir[:index], n.Rir[index+1:]...)
		}
	}
}

func (n *NddoipamIpam) AppendAggregate(s []*NddoipamIpamAggregate) {
	for _, new := range s {
		found := false
		for _, orig := range n.Aggregate {
			if orig.GetTenant() == new.GetTenant() &&
				orig.GetNetworkInstance() == new.GetNetworkInstance() &&
				orig.GetPrefix() == new.GetPrefix() {
				//orig = new
				found = true
			}
		}
		if !found {
			n.Aggregate = append(n.Aggregate, new)
		}
	}
}

func (n *NddoipamIpam) DeleteAggregate(s []*NddoipamIpamAggregate) {
	for _, new := range s {
		found := false
		index := -1
		for i, orig := range n.Aggregate {
			if orig.GetTenant() == new.GetTenant() &&
				orig.GetNetworkInstance() == new.GetNetworkInstance() &&
				orig.GetPrefix() == new.GetPrefix() {
				//orig = new
				found = true
				index = i
			}
		}
		if found {
			// delete entry
			n.Aggregate = append(n.Aggregate[:index], n.Aggregate[index+1:]...)
		}
	}
}

func (n *NddoipamIpam) AppendPrefix(s []*NddoipamIpamIpPrefix) {
	for _, new := range s {
		found := false
		for _, orig := range n.IpPrefix {
			if orig.GetTenant() == new.GetTenant() &&
				orig.GetNetworkInstance() == new.GetNetworkInstance() &&
				orig.GetPrefix() == new.GetPrefix() {
				//orig = new
				found = true
			}
		}
		if !found {
			n.IpPrefix = append(n.IpPrefix, new)
		}
	}
}

func (n *NddoipamIpam) DeletePrefix(s []*NddoipamIpamIpPrefix) {
	for _, new := range s {
		found := false
		index := -1
		for i, orig := range n.IpPrefix {
			if orig.GetTenant() == new.GetTenant() &&
				orig.GetNetworkInstance() == new.GetNetworkInstance() &&
				orig.GetPrefix() == new.GetPrefix() {
				//orig = new
				found = true
				index = i
			}
		}
		if found {
			// delete entry
			n.IpPrefix = append(n.IpPrefix[:index], n.IpPrefix[index+1:]...)
		}
	}
}
