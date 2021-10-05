package server

/*
import (
	"encoding/json"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	ipamv1alpha1 "github.com/yndd/nddo-ipam/apis/ipam/v1alpha1"
)
*/

/*

func (s *Server) ProcessUpdate(resources map[string]map[string]*resource) error {
	ipam := s.GetSchema().GetIpam()
	for resourceName, r := range resources {
		log := s.log.WithValues("ProcessUpdate", resourceName)
		switch resourceName {
		case "rir":
			newRirs, err := s.handler.getNewRirStructs([]string{"name"}, r)
			if err != nil {
				return err
			}
			// this resource come in bulk and will update the complete state
			ipam.AppendRir(newRirs)
			s.GetSchema().SetIpam(ipam)
		case "aggregate":
			newAggregates, err := s.handler.getNewAggregateStructs([]string{"tenant", "network-instance", "prefix"}, r)
			if err != nil {
				return err
			}

			for _, o := range s.GetSchema().GetIpam().GetAggregate() {
				found := false
				for _, n := range newAggregates {
					found = n.Exists(o)
				}
				if !found {
					// we can delete as we do here or check if ip addresses were allocated or not
					log.Debug("will be deleted but not sure if this is correct")
				}
			}
			// get origin
			for _, n := range newAggregates {
				n.GetStatusFromOrigin(s.GetSchema().GetIpam().GetAggregate())
			}
			// Process overlap between newAggregates and original aggreagates
			for _, n := range newAggregates {
				err := n.ProcessOverlapWithAggregate(s.GetSchema().GetIpam().GetAggregate())
				if err != nil {
					return err
				}
			}
			// this resource come in bulk and will update the complete state
			ipam.AppendAggregate(newAggregates)
			s.GetSchema().SetIpam(ipam)
			// reevaluate the prefixes against prefixes and aggregates
			for _, n := range s.GetSchema().GetIpam().GetIpPrefix() {
				err := n.ProcessOverlapWithPrefix(s.GetSchema().GetIpam().GetIpPrefix())
				if err != nil {
					return err
				}
				log.Debug("Prefix Status aggregate", "Prefix", n.GetPrefix(), "Status", n.GetStatus())
				err = n.ProcessOverlapWithAggregate(s.GetSchema().GetIpam().GetAggregate())
				if err != nil {
					return err
				}
				log.Debug("Prefix Status aggregate", "Prefix", n.GetPrefix(), "Status", n.GetStatus())
			}
			s.GetSchema().SetIpam(ipam)
		case "ip-prefix":
			newPrefixes, err := s.handler.getNewPrefixStructs([]string{"tenant", "network-instance", "prefix"}, r)
			if err != nil {
				return err
			}
			// get latest status if it exists
			for _, n := range newPrefixes {
				n.GetStatusFromOrigin(s.GetSchema().GetIpam().GetIpPrefix())
			}
			// for each new prefix check if there is overlap within the prefix
			// afterwards evaluate the overlap with the aggregates
			for _, n := range newPrefixes {
				// Process overlap between newPrefixes and original prefix status
				err := n.ProcessOverlapWithPrefix(s.GetSchema().GetIpam().GetIpPrefix())
				if err != nil {
					return err
				}
				log.Debug("Prefix Status ip-prefix", "Prefix", n.GetPrefix(), "Status", n.GetStatus())
				// Process overlap between newPrefixes and the aggregate
				log.Debug("Prefix Status ip-prefix", "Aggregates", s.GetSchema().GetIpam().GetAggregate())
				err = n.ProcessOverlapWithAggregate(s.GetSchema().GetIpam().GetAggregate())
				if err != nil {
					return err
				}
				log.Debug("Prefix Status ip-prefix", "Prefix", n.GetPrefix(), "Status", n.GetStatus())
			}

			ipam.AppendPrefix(newPrefixes)
			s.GetSchema().SetIpam(ipam)
		case "ip-range":
			for _, resource := range r {
				log.Debug("Resource Validation", "Key", resource.GetKey(), "Data", resource.GetData())
			}
		case "ip-address":
			for _, resource := range r {
				log.Debug("Resource Validation", "Key", resource.GetKey(), "Data", resource.GetData())
			}
		}
	}
	return nil
}

func (s *Server) ProcessDelete(resources map[string]map[string]*resource) error {
	ipam := s.GetSchema().GetIpam()
	for resourceName, r := range resources {
		log := s.log.WithValues("ProcessDelete", resourceName)
		switch resourceName {
		case "ipam":
			s.GetSchema().DeleteIpam()
		case "rir":
			newRirs, err := s.handler.getNewRirStructs([]string{"name"}, r)
			if err != nil {
				return err
			}
			r, err := s.handler.ProcessRirDelete(s.schema.Ipam.Rir, newRirs)
			if err != nil {
				return err
			}
			// this resource come in bulk and will update the complete state
			ipam.DeleteRir(r)
			s.GetSchema().SetIpam(ipam)
		case "aggregate":
			newAggregates, err := s.handler.getNewAggregateStructs([]string{"tenant", "network-instance", "prefix"}, r)
			if err != nil {
				return err
			}
			a, err := s.handler.ProcessAggregateDelete(s.schema.Ipam.Aggregate, newAggregates)
			if err != nil {
				return err
			}
			// this resource come in bulk and will update the complete state
			ipam.DeleteAggregate(a)
			s.GetSchema().SetIpam(ipam)
		case "ip-prefix":
			newPrefixes, err := s.handler.getNewPrefixStructs([]string{"tenant", "network-instance", "prefix"}, r)
			if err != nil {
				return err
			}
			// TO BE UPDATED
			ipam.DeletePrefix(newPrefixes)
			s.GetSchema().SetIpam(ipam)
		case "ip-ranage":
			for _, resource := range r {
				log.Debug("Resource Validation", "Key", resource.GetKey(), "Data", resource.GetData())
			}
		case "ip-address":
			for _, resource := range r {
				log.Debug("Resource Validation", "Key", resource.GetKey(), "Data", resource.GetData())
			}
		}
	}
	return nil
}

func (h *Handler) ProcessRirDelete(origRir, newRir []*ipamv1alpha1.NddoipamIpamRir) ([]*ipamv1alpha1.NddoipamIpamRir, error) {
	var newR []*ipamv1alpha1.NddoipamIpamRir
	for _, orig := range origRir {
		for _, new := range newRir {
			// compare to see if the resource was found, if so delete it
			found := false
			if orig.GetName() == new.GetName() {
				// delete entry
				// TODO check if ip addresses were allocated
				found = true
			}
			if !found {
				// entry is not deleted, add it back to the list
				newR = append(newR, orig)
			}
		}
	}
	return newR, nil
}

func (h *Handler) ProcessAggregateDelete(origAggregate, newAggregate []*ipamv1alpha1.NddoipamIpamAggregate) ([]*ipamv1alpha1.NddoipamIpamAggregate, error) {
	var newA []*ipamv1alpha1.NddoipamIpamAggregate
	for _, orig := range origAggregate {
		for _, new := range newAggregate {

			// compare to see if the resource was found, if so delete it
			found := false
			if orig.GetTenant() == new.GetTenant() &&
				orig.GetNetworkInstance() == new.GetNetworkInstance() &&
				orig.GetPrefix() == new.GetPrefix() {
				// delete entry
				// TODO check if ip addresses were allocated
				found = true
			}
			if !found {
				// entry is not deleted, add it back to the list
				newA = append(newA, orig)
			}
		}
	}
	return newA, nil
}
*/

///@@@@ above was here before

/*
// ProcessAggregateUpdate check for deleted objects and sets the status according to previous state if exists
// if the status is not active we validate overlap and set the status appropriately
func (h *Handler) ProcessAggregateUpdate(origAggregate, newAggregate []*ipamv1alpha1.NddoipamIpamAggregate) ([]*ipamv1alpha1.NddoipamIpamAggregate, error) {
	// we run over the original objects in the state datastore
	// we compare with the new objects
	// if found we update the status to the original status -> this is to retain the previous status when active
	// if not found this entry no longer exists as this is a bulk update
	// -> for these deleted entries do we need to check if their our pending ip addresses active or not
	// -> if so we could move the object into draining status
	deletes := make([]*resource, 0)
	for _, orig := range origAggregate {
		found := false
		for _, new := range newAggregate {
			if orig.GetTenant() == new.GetTenant() &&
				orig.GetNetworkInstance() == new.GetNetworkInstance() &&
				orig.GetPrefix() == new.GetPrefix() {
				// found
				found = true
				new.SetStatus(orig.GetStatus())
			}
		}
		if !found {
			// not found -> the entry will be deleted
			// TODO check if the deletes have active leases, if so we cannot delete the entry
			// We can update the status to draining Maybe
			key := make(map[string]string)
			key["tenant"] = orig.GetTenant()
			key["network-instance"] = orig.GetNetworkInstance()
			key["prefix"] = orig.GetPrefix()
			deletes = append(deletes, newResource(key))
		}
	}

	h.log.Debug("ProcessAggregateUpdate", "Deletes", deletes)

	// Process overlap
	// when a resource status was active, we retain this status -> could be triggered by a restart or update
	// for a resources with status NOT active we reevaluate the situation wrt overlap
	for _, newA := range newAggregate {
		// only evaluate the status that is not active
		if newA.GetStatus() != "active" {
			overlap := false
			for _, newB := range newAggregate {
				// evaluate against status that is active
				if newB.GetStatus() == "active" {
					// in order to check address overlap the tenant and network-instance should match
					if newA.GetTenant() == newB.GetTenant() &&
						newA.GetNetworkInstance() == newB.GetNetworkInstance() {
						// dont compare with yourself, should never happen anyhow since the status should be != active
						if newA.GetPrefix() != newB.GetPrefix() {
							// get the subnet
							_, subnet, err := net.ParseCIDR(newA.GetPrefix())
							if err != nil {
								return nil, err
							}
							// find the start address and end address of the subnet
							ipStart, ipEnd := cidr.AddressRange(subnet)
							// check if the start ip address is contained in the prefix
							cStart, err := newB.Contains(ipStart)
							if err != nil {
								return nil, err
							}
							cStop, err := newB.Contains(ipEnd)
							if err != nil {
								return nil, err
							}
							// there should not be any overlap between aggregates
							if cStart || cStop {
								overlap = true
								newA.SetStatus("inactive")
								// TODO set inactive reason to overlap
							}
						}
					}
				}
			}
			// when none overlap we can set the status to active
			if !overlap {
				newA.SetStatus("active")
			}
		}
	}
	return newAggregate, nil
}
*/

/*
func (h *Handler) ProcessPrefixUpdate(origPrefixes, newPrefixes []*ipamv1alpha1.NddoipamIpamIpPrefix) ([]*ipamv1alpha1.NddoipamIpamIpPrefix, error) {
	// get latest status if it exists
	for _, n := range newPrefixes {
		n.GetPrefixStatusFromOrigin(origPrefixes)
	}

	// Process overlap
	// when a resource status was active, we retain this status -> could be triggered by a resart or update
	// for a resources with status NOT active we reevaluate the situation wrt overlap
	for _, n := range newPrefixes {
		err := n.ProcessOverlapWithPrefix(origPrefixes)
		if err != nil {
			return nil, err
		}
	}
	return newPrefixes, nil
}
*/

// @@ below was there before
/*

func (h *Handler) ProcessAggregateWithPrefix(newPrefixes []*ipamv1alpha1.NddoipamIpamIpPrefix, origAggregate []*ipamv1alpha1.NddoipamIpamAggregate) ([]*ipamv1alpha1.NddoipamIpamIpPrefix, error) {
	for _, new := range newPrefixes {
		if new.GetStatus() == "check parent aggregate" {
			overlap := false
			parent := false
			for _, orig := range origAggregate {
				// evaluate against status that is active
				if orig.GetStatus() == "active" {
					// in order to check address overlap the tenant and network-instance should match
					if new.GetTenant() == orig.GetTenant() &&
						new.GetNetworkInstance() == orig.GetNetworkInstance() {
						// dont compare with yourself, should never happen anyhow since the status should be != active
						// get the subnet
						_, subnet, err := net.ParseCIDR(new.GetPrefix())
						if err != nil {
							return nil, err
						}
						// find the start address and end address of the subnet
						ipStart, ipEnd := cidr.AddressRange(subnet)
						// check if the start ip address is contained in the prefix
						cStart, err := orig.Contains(ipStart)
						if err != nil {
							return nil, err
						}
						cStop, err := orig.Contains(ipEnd)
						if err != nil {
							return nil, err
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
						h.log.WithValues(
							"Aggregate Prefix", orig.GetPrefix(),
							"Prefix", new.GetPrefix(),
							"ipStart", ipStart,
							"ipEnd", ipEnd,
							"cStart", cStart,
							"cStop", cStop,
							"overlap", overlap,
							"parent", parent,
						).Debug("ProcessAggregatePrefix")
					}
				}
			}
			// when none overlap we can set the status to active
			switch {
			case overlap && parent:
				// overlap and parent found which is good
				new.SetStatus("active")
			case overlap && !parent:
				// half overlap is not good
				new.SetStatus("inactive")
			case !overlap:
				// no overlap means no parent axists
				new.SetStatus("inactive")
			}
			h.log.Debug("ProcessAggregatePrefix", "Prefix", new.GetPrefix(), "Status", new.GetStatus())
		}
	}
	return newPrefixes, nil
}

func (h *Handler) getNewRirStructs(keys []string, r map[string]*resource) ([]*ipamv1alpha1.NddoipamIpamRir, error) {
	newStruct := make([]*ipamv1alpha1.NddoipamIpamRir, 0)
	for keystring, new := range r {
		switch d := new.data.(type) {
		case map[string]interface{}:
			for _, key := range keys {
				d[key] = new.GetKey()[key]
			}
		case nil: // typically happens with a delete which has no data
			new.data = make(map[string]interface{})
			switch d := new.data.(type) {
			case map[string]interface{}:
				for _, key := range keys {
					d[key] = new.GetKey()[key]
				}
			}
		}
		d, err := json.Marshal(new.data)
		if err != nil {
			return nil, err
		}
		var a *ipamv1alpha1.NddoipamIpamRir
		if err := json.Unmarshal(d, &a); err != nil {
			return nil, err
		}
		log := h.log.WithValues("Key", keystring)
		log.Debug("getNewRirStructs", "Data", a)

		newStruct = append(newStruct, a)
	}
	return newStruct, nil
}

func (h *Handler) getNewAggregateStructs(keys []string, r map[string]*resource) ([]*ipamv1alpha1.NddoipamIpamAggregate, error) {
	h.log.Debug("getNewAggregateStructs", "Data", r)
	newStruct := make([]*ipamv1alpha1.NddoipamIpamAggregate, 0)
	for keystring, new := range r {
		switch d := new.data.(type) {
		case map[string]interface{}:
			for _, key := range keys {
				d[key] = new.GetKey()[key]
			}
		case nil:
			new.data = make(map[string]interface{})
			switch d := new.data.(type) {
			case map[string]interface{}:
				for _, key := range keys {
					d[key] = new.GetKey()[key]
				}
			}
		}
		d, err := json.Marshal(new.data)
		if err != nil {
			return nil, err
		}
		var a *ipamv1alpha1.NddoipamIpamAggregate
		if err := json.Unmarshal(d, &a); err != nil {
			return nil, err
		}
		log := h.log.WithValues("Key", keystring)
		log.Debug("getNewAggregateStructs", "Data", a)

		newStruct = append(newStruct, a)
	}
	return newStruct, nil
}

func (h *Handler) getNewPrefixStructs(keys []string, r map[string]*resource) ([]*ipamv1alpha1.NddoipamIpamIpPrefix, error) {
	newStruct := make([]*ipamv1alpha1.NddoipamIpamIpPrefix, 0)
	for keystring, new := range r {
		switch d := new.data.(type) {
		case map[string]interface{}:
			for _, key := range keys {
				d[key] = new.GetKey()[key]
			}
		case nil:
			new.data = make(map[string]interface{})
			switch d := new.data.(type) {
			case map[string]interface{}:
				for _, key := range keys {
					d[key] = new.GetKey()[key]
				}
			}
		}
		d, err := json.Marshal(new.data)
		if err != nil {
			return nil, err
		}
		var a *ipamv1alpha1.NddoipamIpamIpPrefix
		if err := json.Unmarshal(d, &a); err != nil {
			return nil, err
		}
		log := h.log.WithValues("Key", keystring)
		log.Debug("getNewIPPrefixStructs", "Data", a)

		newStruct = append(newStruct, a)
	}
	return newStruct, nil
}
*/

// @@ above was there

/*
	// check the origin
	b, err := json.Marshal(s.schemaRaw)
	if err != nil {
		return err
	}
	var newSchema *ipamv1alpha1.Nddoipam
	if err := json.Unmarshal(b, &newSchema); err != nil {
		return err
	}

	/*
		// check object deletes
		deleteRirs := make([]string, 0)
		for _, rir := range s.schema.Ipam.Rir {
			found := false
			for _, newrir := range newSchema.Ipam.Rir {
				if rir.GetName() == newrir.GetName() {
					found = true
				}
			}
			if !found {
				deleteRirs = append(deleteRirs, rir.GetName())
			}
		}
		// update object or create the new object
		for _, newrir := range newSchema.Ipam.Rir {
			found := false
			for i, rir := range s.schema.Ipam.Rir {
				if rir.GetName() == newrir.GetName() {
					// update the object
					found = true
					s.schema.Ipam.Rir[i] = newrir.DeepCopy()
				}
			}
			if !found {
				// create the new object
				if s.schema.Ipam.Rir != nil {
					s.schema.Ipam.Rir = make([]*ipamv1alpha1.NddoipamIpamRir, 0)
				}
				s.schema.Ipam.Rir = append(s.schema.Ipam.Rir, newrir)
			}
		}
		for _, delete := range deleteRirs {

		}

	// check deletes

	// check additions

	return nil
}
*/

/*
func (s *Server) Reconcile() error {
	b, err := json.Marshal(s.schemaRaw)
	if err != nil {
		return err
	}
	var schema *types.Nddoipam
	if err := json.Unmarshal(b, &schema); err != nil {
		return err
	}

	// set the private status to the rir
	for _, x := range schema.Ipam.Rir {
		switch x.GetName() {
		case "rfc1918", "rfc6598", "ula":
			x.SetPrivate(true)
		default:
			x.SetPrivate(false)
		}
	}

	// set status for all aggregates
	for _, x := range schema.Ipam.Aggregate {
		if x.GetAdminState() == "disable" {
			x.SetStatus("inactive")
		} else {
			x.SetStatus("container")
		}
	}

	// parse the ip ranges and fill in the size
	for _, r := range schema.Ipam.IpRange {
		if err := r.Parse(); err != nil {
			return err
		}
	}

	// sort the ranges
	schema.Ipam.SortRanges()
	schema.Ipam.SortPrefixes()

	// Process overlap of ip ranges
	for i1, r1 := range schema.Ipam.IpRange {
		for i2, r2 := range schema.Ipam.IpRange {
			// Process the tenancy and network-instance match otherwise we cannot compare the ip ranges
			if r1.GetTenant() == r2.GetTenant() && r1.GetNetworkInstance() == r2.GetNetworkInstance() {
				// dont compare the same range
				if i1 != i2 {
					c1, err := r1.Contains(r2.GetStart())
					if err != nil {
						return err
					}
					c2, err := r1.Contains(r2.GetEnd())
					if err != nil {
						return err
					}
					if c1 || c2 {
						// overlap detected
						// the status will not get active unless there was no overlap so we can safely
						// use the status to Process
						switch {
						case r1.GetStatus() == "active":
							r2.SetStatus("inactive overlap")
						case r2.GetStatus() == "active":
							r1.SetStatus("inactive overlap")
						default:
							r1.SetStatus("inactive overlap")
							r2.SetStatus("inactive overlap")
						}
					}
				}
			}
		}
	}

	// Process overlap of ip prefixes
	for i1, p1 := range schema.Ipam.IpPrefix {
		for i2, p2 := range schema.Ipam.IpPrefix {
			// Process the tenancy and network-instance
			if p1.GetTenant() == p2.GetTenant() && p1.GetNetworkInstance() == p2.GetNetworkInstance() {
				// dont compare the same prefix
				if i1 != i2 {
					// given we sort from small to big we always check if the smaller one fits in the bigger prefix
					// we check if p1 is contained in p2
					_, subnet, err := net.ParseCIDR(p1.GetPrefix())
					if err != nil {
						return err
					}
					ipStart, ipEnd := cidr.AddressRange(subnet)
					c1, err := p2.Contains(ipStart)
					if err != nil {
						return err
					}
					c2, err := p2.Contains(ipEnd)
					if err != nil {
						return err
					}
					switch {
					case c1 && c2:
						// p1 prefix is contained in the p2 prefix, which is ok
						// we dont fill out the parent hierarchy yet since there could be other overlaps that makes this
						// obsolete
					case c1 || c2:
						// this is not good the prefixes intermingle
						// the status will not get active unless there was no overlap so we can safely
						// use the status to Process
						switch {
						case p1.GetStatus() == "active":
							p2.SetStatus("inactive intermingle overlap")
						case p2.GetStatus() == "active":
							p1.SetStatus("inactive intermingle overlap")
						default:
							p1.SetStatus("inactive intermingle overlap")
							p2.SetStatus("inactive intermingle overlap")
						}
					default:
						// 2 indepent prefixes
					}
				}
			}
		}
	}

	// Process Overlap and hierarchy
	// -> we need to take into account he previous status, if not set this is a new prefix
	// -> overlap should use the empty status for validation
	// -> if there is overlap we dont allow the prefix to be used for ip address allocation
	// range can overlap with other range -> if empty status we set both to invalid, if one was active we keep it active and stop the other one
	// range can overlap with prefix, but start and end should be contained in the prefix, this is ok, if not we have a problem

	// check if prefixes have a parent aggregate
	for _, a := range schema.Ipam.Aggregate {
		for _, p := range schema.Ipam.IpPrefix {
			found := false
			if a.GetTenant() == p.GetTenant() && a.GetNetworkInstance() == p.GetNetworkInstance() {
				if a.ContainsPrefix(p.GetPrefix()) {
					found = true

					p.AppendAggregate(&types.NddoipamIpamIpPrefixAggregate{
						Tenant:          utils.StringPtr(a.GetTenant()),
						NetworkInstance: utils.StringPtr(a.GetNetworkInstance()),
						Prefix:          utils.StringPtr(a.GetPrefix()),
					})
				}
			}
			if !found {
				// if not found set status to inactive
				p.SetStatus("inactive")
			} else {
				// found
				if a.GetStatus() != "container" {
					// if parent status is nok set status to inactive
					p.SetStatus("inactive")
				} else {
					// found and parent status is ok so we can acitvate the status
					p.SetStatus("active")
				}
			}
			if p.GetAdminState() == "disable" {
				p.SetStatus("inactive")
			}
		}
		for _, r := range schema.Ipam.IpRange {
			found := false
			if a.GetTenant() == r.GetTenant() && a.GetNetworkInstance() == r.GetNetworkInstance() {
				if a.ContainsAddress(r.GetStart()) && a.ContainsAddress(r.GetEnd()) {
					found = true

					r.AppendAggregate(&types.NddoipamIpamIpRangeAggregate{
						Tenant:          utils.StringPtr(a.GetTenant()),
						NetworkInstance: utils.StringPtr(a.GetNetworkInstance()),
						Prefix:          utils.StringPtr(a.GetPrefix()),
					})
				}
			}
			if !found {
				// if not found set status to inactive
				r.SetStatus("inactive")
			} else {
				// found
				if a.GetStatus() != "container" {
					// if parent status is nok set status to inactive
					r.SetStatus("inactive")
				} else {
					// found and parent status is ok so we can acitvate the status
					r.SetStatus("active")
				}
			}
			if r.GetAdminState() == "disable" {
				r.SetStatus("inactive")
			}
		}
	}
	// check
	for _, x := range schema.Ipam.IpPrefix {
		for _, a := range schema.Ipam.Aggregate {

		}
		if x.GetAdminState() == "disable" {
			x.SetStatus("inactive")
		} else {
			x.SetStatus("reserved")
		}
	}

	// ranges
	for _, x := range schema.Ipam.IpRange {
		if x.GetAdminState() == "disable" {
			x.SetStatus("inactive")
		} else {
			x.SetStatus("reserved")
		}
	}

	// ip address
	for _, a := range schema.Ipam.IpAddress {
		found := false
		// check if it belongs to an ip range
		for _, r := range schema.Ipam.IpRange {
			// check if they belong to the same tenant and network-instance
			if r.GetTenant() == a.GetTenant() && r.GetNetworkInstance() == a.GetNetworkInstance() {
				if r.Contains(a.GetAddress()) {
					a.AppendIPRange(&types.NddoipamIpamIpAddressIpRange{
						Tenant:          utils.StringPtr(r.GetTenant()),
						NetworkInstance: utils.StringPtr(r.GetNetworkInstance()),
						Start:           utils.StringPtr(r.GetStart()),
						End:             utils.StringPtr(r.GetEnd()),
					})
					found = true

					if a.GetStatus() != "active" {
						a.SetStatus("inactive")
					} else {
						a.SetStatus("active")
					}
				}
			}
		}
		// check if it belongs to an ip prefix
		for _, p := range schema.Ipam.IpPrefix {
			// check if they belong to the same tenant and network-instance
			if p.GetTenant() == a.GetTenant() && p.GetNetworkInstance() == a.GetNetworkInstance() {
				if p.Contains(a.GetAddress()) {
					a.AppendIPPrefix(&types.NddoipamIpamIpAddressIpPrefix{
						Tenant:          utils.StringPtr(p.GetTenant()),
						NetworkInstance: utils.StringPtr(p.GetNetworkInstance()),
						Prefix:          utils.StringPtr(p.GetPrefix()),
					})
					found = true

					if p.GetStatus() != "active" {
						p.SetStatus("inactive")
					} else {
						p.SetStatus("active")
					}
				}
			}
		}
		if !found {
			// TODO to see if this is ok
			a.SetStatus("reserved")
		}
		if a.GetAdminState() == "disable" {
			a.SetStatus("inactive")
		}
	}

	return nil
}
*/
