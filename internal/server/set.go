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
	"fmt"
	"time"

	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yndd/nddo-ipam/internal/connector"
	"github.com/yndd/nddo-ipam/internal/controllers/ipam"
)

func (s *Server) Set(ctx context.Context, req *gnmi.SetRequest) (*gnmi.SetResponse, error) {

	ok := s.unaryRPCsem.TryAcquire(1)
	if !ok {
		return nil, status.Errorf(codes.ResourceExhausted, errMaxNbrOfUnaryRPCReached)
	}
	defer s.unaryRPCsem.Release(1)

	numUpdates := len(req.GetUpdate())
	numReplaces := len(req.GetReplace())
	numDeletes := len(req.GetDelete())
	if numUpdates+numReplaces+numDeletes == 0 {
		return nil, status.Errorf(codes.InvalidArgument, errMissingPathsInGNMISet)
	}

	log := s.log.WithValues("numUpdates", numUpdates, "numReplaces", numReplaces, "numDeletes", numDeletes)
	log.Debug("Set")

	if numReplaces > 0 {
		updateObjects := make(map[string]*connector.Object)

		for _, u := range req.GetReplace() {
			//log.Debug("Replace", "Update", u)
			n, err := s.GetConfigCache().GetNotificationFromUpdate(ipam.GnmiTarget, ipam.GnmiOrigin, u)
			if err != nil {
				log.Debug("GetNotificationFromUpdate Error", "Notification", n, "Error", err)
				return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error: %v", err))
			}
			//log.Debug("Replace", "Notification", n)
			if n != nil {
				if err := s.GetConfigCache().GnmiUpdate(ipam.GnmiTarget, n); err != nil {
					log.Debug("GnmiUpdate Error", "Notification", n, "Error", err)
					return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error: %v", err))
				}
			}
			// gather which resources need to be updated
			resourceName, resourceKey, path := getResources2Update(u)
			if path != nil {
				if _, ok := updateObjects[resourceKey]; !ok {
					updateObjects[resourceKey] = &connector.Object{
						Kind: resourceName,
						Key:  resourceKey,
						Path: path,
					}
				}
			}
		}
		// Update the connected processing engine with the updated resources
		for _, o := range updateObjects {
			d, err := s.GetConfigCache().GetJson(ipam.GnmiTarget, o.Path)
			if err != nil {
				return nil, err
			}
			o.Data = d
			if _, err := s.connector.Create(o); err != nil {
				return nil, err
			}
		}
	}

	if numUpdates > 0 {
		updateObjects := make(map[string]*connector.Object)

		for _, u := range req.GetUpdate() {
			n, err := s.GetConfigCache().GetNotificationFromUpdate(ipam.GnmiTarget, ipam.GnmiOrigin, u)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error: %v", err))
			}
			if n != nil {
				if err := s.GetConfigCache().GnmiUpdate(ipam.GnmiTarget, n); err != nil {
					return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error: %v", err))
				}
			}
			// gather which resources need to be updated
			resourceName, resourceKey, path := getResources2Update(u)
			if path != nil {
				if _, ok := updateObjects[resourceKey]; !ok {
					updateObjects[resourceKey] = &connector.Object{
						Kind: resourceName,
						Key:  resourceKey,
						Path: path,
					}
				}
			}
		}
		// Update the connected processing engine with the updated resources
		for _, o := range updateObjects {
			d, err := s.GetConfigCache().GetJson(ipam.GnmiTarget, o.Path)
			if err != nil {
				return nil, err
			}
			o.Data = d
			if _, err := s.connector.Update(o); err != nil {
				return nil, err
			}
		}
	}

	if numDeletes > 0 {
		deleteObjects := make(map[string]*connector.Object)
		for _, p := range req.GetDelete() {
			n, err := s.GetConfigCache().GetNotificationFromDelete(ipam.GnmiTarget, ipam.GnmiOrigin, p)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error: %v", err))
			}
			if err := s.GetConfigCache().GnmiUpdate(ipam.GnmiTarget, n); err != nil {
				return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error: %v", err))
			}

			// gather which resources need to be deleted
			resourceName, resourceKey, path := getResources2Delete(p)
			if path != nil {
				if _, ok := deleteObjects[resourceKey]; !ok {
					deleteObjects[resourceKey] = &connector.Object{
						Kind: resourceName,
						Key:  resourceKey,
						Path: path,
					}
				}
			}
		}
		// Update the connected processing engine with the deleted resources
		for _, o := range deleteObjects {
			if err := s.connector.Delete(o); err != nil {
				return nil, err
			}
		}
	}

	// TODO process updatePaths, deletePaths
	// get JSON blobs
	// process the updates -> update the status

	cfg, _ := s.GetConfig()
	state, _ := s.GetState()
	log.Debug("Set Result Config Data", "schema", cfg)
	log.Debug("Set Result Status Data", "schema", state)

	return &gnmi.SetResponse{
		Response: []*gnmi.UpdateResult{
			{
				Timestamp: time.Now().UnixNano(),
			},
		}}, nil
}
