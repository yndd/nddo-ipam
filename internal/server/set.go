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
		resReplace, err := s.handler.GetResources(ResourceActionReplace, req.GetReplace())
		s.ProcessUpdate(resReplace)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error: %v", err))
		}
		for _, update := range req.GetReplace() {
			var err error
			s.schemaRaw, err = s.handler.Replace(s.schemaRaw, update)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error: %v", err))
			}
		}
	}

	if numUpdates > 0 {
		resUpdate, err := s.handler.GetResources(ResourceActionUpdate, req.GetUpdate())
		s.ProcessUpdate(resUpdate)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error: %v", err))
		}
		for _, update := range req.GetUpdate() {
			var err error
			s.schemaRaw, err = s.handler.Update(s.schemaRaw, update)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error: %v", err))
			}
		}
	}

	if numDeletes > 0 {
		resDelete, err := s.handler.GetResources2Delete(req.GetDelete())
		s.ProcessDelete(resDelete)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error: %v", err))
		}
		for _, path := range req.GetDelete() {
			var err error
			s.schemaRaw, err = s.handler.Delete(s.schemaRaw, path)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error: %v", err))
			}
		}
	}

	log.Debug("Set Result Config Data", "schema", s.schemaRaw)
	log.Debug("Set Result Status Data", "schema", s.GetSchema())

	return &gnmi.SetResponse{
		Response: []*gnmi.UpdateResult{
			{
				Timestamp: time.Now().UnixNano(),
			},
		}}, nil
}
