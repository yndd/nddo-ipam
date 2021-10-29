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

package server

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/openconfig/gnmi/proto/gnmi"
	ipamv1alpha1 "github.com/yndd/nddo-ipam/apis/ipam/v1alpha1"
	"github.com/yndd/nddo-ipam/internal/controllers/ipam"
)

func (s *Server) GetInitialState() error {
	nddoipam, err := s.client.ListNddoipam(s.ctx)
	if err != nil {
		return err
	}
	b, err := json.Marshal(nddoipam)
	if err != nil {
		return err
	}
	var x1 interface{}
	if err := json.Unmarshal(b, &x1); err != nil {
		return err
	}

	switch x := x1.(type) {
	case map[string]interface{}:
		x1 = x["nddo-ipam"]
	}

	rootPath := []*gnmi.Path{
		{
			Elem: []*gnmi.PathElem{
				{Name: "nddo-ipam"},
				{Name: "ipam"},
			},
		},
	}

	prefix := &gnmi.Path{Target: ipam.GnmiTarget, Origin: ipam.GnmiOrigin}

	updates := s.parser.GetUpdatesFromJSONDataGnmi(rootPath[0], s.parser.XpathToGnmiPath("/", 0), x1, resourceRefPaths)
	for _, u := range updates {
		s.log.Debug("Observe Fine Grane Updates X1", "Path", s.parser.GnmiPathToXPath(u.Path, true), "Value", u.GetVal())
		n, err := s.GetStateCache().GetNotificationFromUpdate(prefix, u)
		if err != nil {
			s.log.Debug("GetNotificationFromUpdate Error", "Notification", n, "Error", err)
			return err
		}
		s.log.Debug("Replace", "Notification", n)
		if n != nil {
			if err := s.GetStateCache().GnmiUpdate(ipam.GnmiTarget, n); err != nil {
				s.log.Debug("GnmiUpdate Error", "Notification", n, "Error", err)
				return err
			}
		}
	}
	return nil
}

func (s *Server) GetConfig() (*ipamv1alpha1.Ipam, error) {
	prefix := &gnmi.Path{Target: ipam.GnmiTarget, Origin: ipam.GnmiOrigin}
	x, err := s.GetConfigCache().GetJson(ipam.GnmiTarget, prefix, &gnmi.Path{Elem: []*gnmi.PathElem{}})
	if err != nil {
		return nil, err
	}
	s.log.Debug("GetConfig", "config", x)
	b, err := json.Marshal(x)
	if err != nil {
		return nil, err
	}
	n := ipamv1alpha1.Ipam{}
	if err := json.Unmarshal(b, &n); err != nil {
		return nil, err
	}
	return &n, nil
}

func (s *Server) GetState() (*ipamv1alpha1.Nddoipam, error) {
	prefix := &gnmi.Path{Target: ipam.GnmiTarget, Origin: ipam.GnmiOrigin}
	x, err := s.GetStateCache().GetJson(ipam.GnmiTarget, prefix, &gnmi.Path{Elem: []*gnmi.PathElem{}})
	if err != nil {
		return nil, err
	}
	s.log.Debug("GetState", "state", x)
	b, err := json.Marshal(x)
	if err != nil {
		return nil, err
	}
	n := ipamv1alpha1.Nddoipam{}
	if err := json.Unmarshal(b, &n); err != nil {
		return nil, err
	}
	return &n, nil
}

func getPath2process(p *gnmi.Path) (string, *gnmi.Path) {
	pathlen := len(p.GetElem())
	switch {
	case pathlen > 3:
		// ip-prefix, ip-range, ip-address
		pathElem := p.GetElem()[3].GetName()
		if pathElem == "ip-prefix" || pathElem == "ip-range" || pathElem == "ip-address" {
			return getKeyString(p.GetElem()[3].GetKey()), &gnmi.Path{
				Elem: p.GetElem()[:4], // we cut the path
			}
		}
	case pathlen > 2:
		// network-instance
		pathElem := p.GetElem()[2].GetName()
		if pathElem == "network-instance" {
			return getKeyString(p.GetElem()[2].GetKey()), &gnmi.Path{
				Elem: p.GetElem()[:3], // we cut the path
			}
		}
	case pathlen > 1:
		// rir or tenant
		pathElem := p.GetElem()[1].GetName()
		if pathElem == "rir" || pathElem == "tenant" {
			return getKeyString(p.GetElem()[1].GetKey()), &gnmi.Path{
				Elem: p.GetElem()[:2], // we cut the path
			}
		}
	case pathlen == 1:
		// ipam
		pathElem := p.GetElem()[0].GetName()
		if pathElem == "ipam" {
			return getKeyString(p.GetElem()[0].GetKey()), &gnmi.Path{
				Elem: p.GetElem()[:1], // we cut the path
			}
		}
	}
	return "", nil
}

/*
func getResources2Update(u *gnmi.Update) (intentlogic.ResourceKind, string, *gnmi.Path) {
	// The order is important since we first want to check from the lowest level
	if len(u.GetPath().GetElem()) > 3 {
		// ip-address, ip-prefix, ip-range
		// the key is at place 4
		if _, ok := intentlogic.Dispatcher[u.GetPath().GetElem()[3].GetName()]; ok {
			return intentlogic.String2ResourceKind(u.GetPath().GetElem()[3].GetName()),
				getKeyString(u.GetPath().GetElem()[3].GetKey()),
				&gnmi.Path{
					Elem: u.GetPath().GetElem()[:4], // we cut the path
				}
		}
	}
	if len(u.GetPath().GetElem()) > 2 {
		// network-instance
		// the key is at place 3
		if _, ok := intentlogic.Dispatcher[u.GetPath().GetElem()[2].GetName()]; ok {
			return intentlogic.String2ResourceKind(u.GetPath().GetElem()[2].GetName()),
				getKeyString(u.GetPath().GetElem()[2].GetKey()),
				&gnmi.Path{
					Elem: u.GetPath().GetElem()[:3], // we cut the path
				}
		}
	}
	if len(u.GetPath().GetElem()) > 1 {
		// rir
		// the key is at place 2 of the pathElem for all resources
		// if there is a match in the dispatcher we can return a resource
		if _, ok := intentlogic.Dispatcher[u.GetPath().GetElem()[1].GetName()]; ok {
			return intentlogic.String2ResourceKind(u.GetPath().GetElem()[1].GetName()),
				getKeyString(u.GetPath().GetElem()[1].GetKey()),
				&gnmi.Path{
					Elem: u.GetPath().GetElem()[:2], // we cut the path
				}
		}
	}
	return "", "", nil
}
*/

/*
func getResources2Delete(p *gnmi.Path) (intentlogic.ResourceKind, string, *gnmi.Path) {
	if len(p.GetElem()) == 1 {
		// TODO delete all resources
	}
	if len(p.GetElem()) > 3 {
		// we do a more specific check first sicne instance would match if we check on level 2 first
		// ip-address, ip-prefix, ip-range
		// the key is at place 3
		if _, ok := intentlogic.Dispatcher[p.GetElem()[3].GetName()]; ok {
			return intentlogic.String2ResourceKind(p.GetElem()[3].GetName()),
				getKeyString(p.GetElem()[2].GetKey()),
				&gnmi.Path{
					Elem: p.GetElem()[:4], // we cut the path
				}
		}
	}
	if len(p.GetElem()) > 2 {
		// we do a more specific check first sicne instance would match if we check on level 2 first
		// ip-address, ip-prefix, ip-range
		// the key is at place 3
		if _, ok := intentlogic.Dispatcher[p.GetElem()[2].GetName()]; ok {
			return intentlogic.String2ResourceKind(p.GetElem()[2].GetName()),
				getKeyString(p.GetElem()[2].GetKey()),
				&gnmi.Path{
					Elem: p.GetElem()[:3], // we cut the path
				}
		}
	}
	if len(p.GetElem()) > 1 {
		// rir and instance
		// the key is at place 2 of the pathElem for all resources
		// if there is a match in the dispatcher we can return a resource
		if _, ok := intentlogic.Dispatcher[p.GetElem()[1].GetName()]; ok {
			return intentlogic.String2ResourceKind(p.GetElem()[1].GetName()),
				getKeyString(p.GetElem()[1].GetKey()),
				&gnmi.Path{
					Elem: p.GetElem()[:2], // we cut the path
				}
		}
	}
	return "", "", nil
}
*/

func getKeyString(key map[string]string) string {
	sb := strings.Builder{}
	i := 0
	type kv struct {
		Key   string
		Value string
	}
	var ss []kv
	for k, v := range key {
		ss = append(ss, kv{k, v})
	}
	// sort the slice keys so we have a determinsitic result
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Key > ss[j].Key
	})
	for _, kv := range ss {
		sb.WriteString(kv.Key)
		sb.WriteString("=")
		sb.WriteString(kv.Value)
		if i != len(ss)-1 {
			sb.WriteString(",")
		}
		i++
	}
	return sb.String()
}
