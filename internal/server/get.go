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
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yndd/ndd-yang/pkg/yparser"
	"github.com/yndd/nddo-ipam/internal/controllers/ipam"
)

func (s *Server) Get(ctx context.Context, req *gnmi.GetRequest) (*gnmi.GetResponse, error) {
	ok := s.unaryRPCsem.TryAcquire(1)
	if !ok {
		return nil, status.Errorf(codes.ResourceExhausted, "max number of Unary RPC reached")
	}
	defer s.unaryRPCsem.Release(1)

	log := s.log.WithValues("Type", req.GetType())
	if req.GetPath() != nil {
		log.Debug("Get...", "Path", yparser.GnmiPath2XPath(req.GetPath()[0], true))
	} else {
		log.Debug("Get...")
	}

	// we overwrite the gnmi prefix for now
	prefix := &gnmi.Path{Target: ipam.GnmiTarget, Origin: ipam.GnmiOrigin}

	//log.Debug("GNMI GET...")
	updates, err := s.HandleGet(yparser.GetDataType(req.GetType()), prefix, req.GetPath())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error: %s", err))
	}

	return &gnmi.GetResponse{
		Notification: []*gnmi.Notification{
			{
				Timestamp: time.Now().UnixNano(),
				Prefix:    req.GetPrefix(),
				Update:    updates,
			},
		},
	}, nil
}

func (s *Server) HandleGet(cacheType string, prefix *gnmi.Path, reqPaths []*gnmi.Path) ([]*gnmi.Update, error) {
	//var err error
	updates := make([]*gnmi.Update, 0)

	cache := s.GetConfigCache()
	if cacheType == yparser.CacheTypeState {
		cache = s.GetStateCache()
	}

	//if reqPaths == nil || len(reqPaths[0].GetElem()) == 0 {
	//	reqPaths = append(reqPaths, &gnmi.Path{Elem: []*gnmi.PathElem{}})
	//}
	for _, path := range reqPaths {
		x, err := cache.GetJson(ipam.GnmiTarget, prefix, path)
		if err != nil {
			return nil, err
		}

		if updates, err = appendUpdateResponse(x, path, updates); err != nil {
			return nil, err
		}
	}

	/*
		if reqPaths == nil || len(reqPaths[0].GetElem()) == 0 {
			x, err := cache.GetJson(ipam.GnmiTarget, prefix, &gnmi.Path{Elem: []*gnmi.PathElem{}})
			if err != nil {
				return nil, err
			}
			if updates, err = appendUpdateResponse(x, &gnmi.Path{}, updates); err != nil {
				return nil, err
			}
		} else {
			for _, path := range reqPaths {
				xx, err := cache.GetJson(ipam.GnmiTarget, prefix, path)
				if err != nil {
					return nil, err
				}
				x, err := s.parser.DeepCopy(xx)
				if err != nil {
					if !strings.Contains(fmt.Sprint(err), "in cannot be nil") {
						return nil, err
					}
				}
				// prepareResponseData prepare the response data aligned with the controller
				// 1. the hierarchical elements are removed
				// 2. add the last element of the path back to the return data
				hElem, ok := hPathElements[*s.parser.GnmiPathToXPath(path, false)]
				if !ok {
					hElem = []string{}
				}
				s.log.Debug("prepareResponseData", "Path", s.parser.GnmiPathToXPath(path, true))
				newx, err := prepareResponseData(x, path, hElem)
				if err != nil {
					return nil, err
				}
				if updates, err = appendUpdateResponse(newx, &gnmi.Path{}, updates); err != nil {
					return nil, err
				}
			}
		}
	*/
	return updates, nil
}

// prepareResponseData prepare the response data aligned with the controller
// 1. the hierarchical elements within the resource should be removed
// 2. add the last element of the path back to the return data
func prepareResponseData(x interface{}, path *gnmi.Path, hElem interface{}) (interface{}, error) {
	// remove hierarchical elements
	switch x1 := x.(type) {
	case map[string]interface{}:
		switch he := hElem.(type) {
		case map[string]interface{}:
			x = deleteHierarchicalElements(x1, he)
		default:
		}
		//for _, elem := range hElem {
		//	delete(x1, elem)
		//}
	}
	// add last element of the path to the return data
	xx := make(map[string]interface{})
	xx[path.GetElem()[len(path.GetElem())-1].GetName()] = x
	return xx, nil
}

func deleteHierarchicalElements(x map[string]interface{}, he map[string]interface{}) map[string]interface{} {
	for element, heElements := range he {
		switch he := heElements.(type) {
		case nil:
			// this is the end of the hierarchical list
			delete(x, element)
		case map[string]interface{}:
			// there is still some elements in the hierarchical resource list
			if xx, ok := x[element]; ok {
				switch x1 := xx.(type) {
				case map[string]interface{}:
					deleteHierarchicalElements(x1, he)
				case []interface{}:
					for _, xxx := range x1 {
						switch x2 := xxx.(type) {
						case map[string]interface{}:
							deleteHierarchicalElements(x2, he)
						}
					}
				default:
					// it can be that no data is present, so we ignore this
				}
			}
		}
	}
	return x
}

func appendUpdateResponse(data interface{}, path *gnmi.Path, updates []*gnmi.Update) ([]*gnmi.Update, error) {
	var err error
	var d []byte
	if data != nil {
		d, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	upd := &gnmi.Update{
		Path: path,
		Val:  &gnmi.TypedValue{Value: &gnmi.TypedValue_JsonVal{JsonVal: d}},
	}
	updates = append(updates, upd)
	return updates, nil
}
