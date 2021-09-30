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

package kapi

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ipamv1alpha1 "github.com/yndd/nddo-ipam/apis/ipam/v1alpha1"
)

const (
	// errors
	errCreateK8sClient = "cannot create k8s client"
	errListIpam        = "cannot list ipam"
)

// Kapi is a struct to hold the information to talk to the K8s api server
type Kapi struct {
	client.Client
	Scheme *runtime.Scheme

	log logging.Logger
}

// Option can be used to manipulate Options.
type KapiOption func(*Kapi)

// WithLogger specifies how the object should log messages.
func WithLogger(l logging.Logger) KapiOption {
	return func(o *Kapi) {
		o.log = l
	}
}

// WithDeviceName initializes the device name in the device driver
func WithScheme(s *runtime.Scheme) KapiOption {
	return func(o *Kapi) {
		o.Scheme = s
	}
}

func New(config *rest.Config, opts ...KapiOption) (*Kapi, error) {
	a := &Kapi{}
	for _, opt := range opts {
		opt(a)
	}

	cl, err := getClient(config, a.Scheme)
	if err != nil {
		return nil, err
	}
	a.Client = cl
	return a, nil
}

// getClient gets the client to interact with the k8s apiserver
func getClient(config *rest.Config, scheme *runtime.Scheme) (client.Client, error) {
	k8sclopts := client.Options{
		Scheme: scheme,
	}
	c, err := client.New(config, k8sclopts)
	if err != nil {
		return nil, errors.Wrap(err, errCreateK8sClient)
	}
	return c, nil
}

func (a *Kapi) ListNddoipam(ctx context.Context) (*ipamv1alpha1.Nddoipam, error) {
	nddoipams := &ipamv1alpha1.IpamIpamList{}
	if err := a.List(ctx, nddoipams); err != nil {
		return nil, errors.Wrap(err, errListIpam)
	}
	a.log.Debug("ListNddoipam", "nddoipams", nddoipams)
	for _, nddoipam := range nddoipams.Items {
		// we expect only 1 ipam to exists for now
		// TODO handling extra nddoipams
		return nddoipam.Status.AtNetworkNode.Nddoipam, nil
	}
	return &ipamv1alpha1.Nddoipam{}, nil
}
