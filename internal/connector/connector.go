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

package connector

type Connector interface {
	Read(*Object) (*ConnectorObservation, error)
	Create(*Object) (*ConnectorCreation, error)
	Update(*Object) (*ConnectorUpdate, error)
	Delete(*Object) error
}

type ConnectorFns struct {
	ReadFn   func(*Object) (*ConnectorObservation, error)
	CreateFn func(*Object) (*ConnectorCreation, error)
	UpdateFn func(*Object) (*ConnectorUpdate, error)
	DeleteFn func(*Object) error
}

func (e ConnectorFns) Read(o *Object) (*ConnectorObservation, error) {
	return e.ReadFn(o)
}

func (e ConnectorFns) Create(o *Object) (*ConnectorCreation, error) {
	return e.CreateFn(o)
}

func (e ConnectorFns) Update(o *Object) (*ConnectorUpdate, error) {
	return e.UpdateFn(o)
}

// Delete the external resource upon deletion of its associated Managed
// resource.
func (e ConnectorFns) Delete(o *Object) error {
	return e.DeleteFn(o)
}

// A NopClient does nothing.
type NopConnecter struct{}

// Observe does nothing. It returns an empty ExternalObservation and no error.
func (c *NopConnecter) Read(o *Object) (*ConnectorObservation, error) {
	return &ConnectorObservation{}, nil
}

// Create does nothing. It returns an empty ExternalCreation and no error.
func (c *NopConnecter) Create(o *Object) (*ConnectorCreation, error) {
	return &ConnectorCreation{}, nil
}

// Update does nothing. It returns an empty ExternalUpdate and no error.
func (c *NopConnecter) Update(o *Object) (*ConnectorUpdate, error) {
	return &ConnectorUpdate{}, nil
}

// Delete does nothing. It never returns an error.
func (c *NopConnecter) Delete(o *Object) error { return nil }

type ConnectorObservation struct{}

type ConnectorCreation struct{}

type ConnectorUpdate struct{}
