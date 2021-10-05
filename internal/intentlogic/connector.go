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

package intentlogic

import (
	"encoding/json"

	"github.com/yndd/ndd-runtime/pkg/logging"

	ipamv1alpha1 "github.com/yndd/nddo-ipam/apis/ipam/v1alpha1"
	"github.com/yndd/nddo-ipam/internal/connector"
	"github.com/yndd/nddo-ipam/internal/server"
)

// A Operation represents a crud operation
type Operation string

// Operations Kinds.
const (
	// create
	OperationCreate Operation = "Create"
	// update
	OperationUpdate Operation = "Update"
	// delete
	OperationDelete Operation = "Delete"
	// read
	OperationRead Operation = "Read"
)

type Connector struct {
	log logging.Logger
}

func New(l logging.Logger) *Connector {
	return &Connector{
		log: l,
	}
}

func (c *Connector) Read(o *connector.Object) (*connector.ConnectorObservation, error) {
	if err := c.dispatchObject(OperationRead, o); err != nil {
		return nil, err
	}
	return &connector.ConnectorObservation{}, nil
}
func (c *Connector) Create(o *connector.Object) (*connector.ConnectorCreation, error) {
	if err := c.dispatchObject(OperationCreate, o); err != nil {
		return nil, err
	}
	return &connector.ConnectorCreation{}, nil
}
func (c *Connector) Update(o *connector.Object) (*connector.ConnectorUpdate, error) {
	if err := c.dispatchObject(OperationUpdate, o); err != nil {
		return nil, err
	}
	return &connector.ConnectorUpdate{}, nil
}
func (c *Connector) Delete(o *connector.Object) error {
	if err := c.dispatchObject(OperationDelete, o); err != nil {
		return err
	}
	return nil
}

func (c *Connector) dispatchObject(a Operation, o *connector.Object) error {
	log := c.log.WithValues("Kind", o.Kind, "Key", o.Key, "Path", o.Path)
	b, err := json.Marshal(o.Data)
	if err != nil {
		return err
	}
	switch o.Kind {
	case server.ResourceNameRir:
		x := ipamv1alpha1.NddoipamIpamRir{}
		if err := json.Unmarshal(b, &x); err != nil {
			return err
		}
		switch a {
		case OperationCreate:
			log.Debug("Rir Create", "Data", x)
		case OperationUpdate:
			log.Debug("Rir Update", "Data", x)
		case OperationDelete:
			log.Debug("Rir Delete")
		case OperationRead:
			log.Debug("Rir Read")
		}
		return nil

	case server.ResourceNameAggregate:
		x := ipamv1alpha1.NddoipamIpamAggregate{}
		if err := json.Unmarshal(b, &x); err != nil {
			return err
		}
		switch a {
		case OperationCreate:
			log.Debug("Aggregate Create", "Data", x)
		case OperationUpdate:
			log.Debug("Aggregate Update", "Data", x)
		case OperationDelete:
			log.Debug("Aggregate Delete")
		case OperationRead:
			log.Debug("Aggregate Read")
		}
		return nil
	case server.ResourceNameIpPrefix:
		x := ipamv1alpha1.NddoipamIpamIpPrefix{}
		if err := json.Unmarshal(b, &x); err != nil {
			return err
		}
		switch a {
		case OperationCreate:
			log.Debug("IpPrefix Create", "Data", x)
		case OperationUpdate:
			log.Debug("IpPrefix Update", "Data", x)
		case OperationDelete:
			log.Debug("IpPrefix Delete")
		case OperationRead:
			log.Debug("IpPrefix Read")
		}
		return nil
	case server.ResourceNameIpRange:
		x := ipamv1alpha1.NddoipamIpamIpRange{}
		if err := json.Unmarshal(b, &x); err != nil {
			return err
		}
		switch a {
		case OperationCreate:
			log.Debug("IpRange Create", "Data", x)
		case OperationUpdate:
			log.Debug("IpRange Update", "Data", x)
		case OperationDelete:
			log.Debug("IpRange Delete")
		case OperationRead:
			log.Debug("IpRange Read")
		}
		return nil
	case server.ResourceNameIpAddress:
		x := ipamv1alpha1.NddoipamIpamIpAddress{}
		if err := json.Unmarshal(b, &x); err != nil {
			return err
		}
		switch a {
		case OperationCreate:
			log.Debug("IpAddress Create", "Data", x)
		case OperationUpdate:
			log.Debug("IpAddress Update", "Data", x)
		case OperationDelete:
			log.Debug("IpAddress Delete")
		case OperationRead:
			log.Debug("IpAddress Read")
		}
		return nil
	}
	return nil
}
