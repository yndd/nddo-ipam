/*
Copyright 2021 NDDO.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

/*
func (s *Nddoipam) ValidateIPRangeOverlap() {
	for i1, r1 := range s.Ipam.IpRange {
		for i2, r2 := range s.Ipam.IpRange {
			if i1 != i2 {

			}
		}
	}
}

//VerifyNoOverlap takes a list subnets and supernet (CIDRBlock) and verifies
//none of the subnets overlap and all subnets are in the supernet
//it returns an error if any of those conditions are not satisfied
func (s *Nddoipam) VerifyIpPrefixOverlap(subnets []*net.IPNet, CIDRBlock *net.IPNet) error {
	firstLastIP := make([][]net.IP, len(subnets))
	for i, s := range subnets {
		first, last := AddressRange(s)
		firstLastIP[i] = []net.IP{first, last}
	}
	for i, s := range subnets {
		if !CIDRBlock.Contains(firstLastIP[i][0]) || !CIDRBlock.Contains(firstLastIP[i][1]) {
			return fmt.Errorf("%s does not fully contain %s", CIDRBlock.String(), s.String())
		}
		for j := 0; j < len(subnets); j++ {
			if i == j {
				continue
			}

			first := firstLastIP[j][0]
			last := firstLastIP[j][1]
			if s.Contains(first) || s.Contains(last) {
				return fmt.Errorf("%s overlaps with %s", subnets[j].String(), s.String())
			}
		}
	}
	return nil
}
*/
