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

package controllers

import (
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"

	"github.com/yndd/ndd-runtime/pkg/logging"

	"github.com/yndd/nddo-ipam/internal/controllers/ipam"
)

// Setup package controllers.
func Setup(mgr ctrl.Manager, option controller.Options, l logging.Logger, poll time.Duration, namespace string) (map[string]chan event.GenericEvent, error) {
	eventChans := make(map[string]chan event.GenericEvent)
	for _, setup := range []func(ctrl.Manager, controller.Options, logging.Logger, time.Duration, string) (string, chan event.GenericEvent, error){
		ipam.SetupIpam,
		ipam.SetupIpamIprange,
		ipam.SetupIpamIpaddress,
		ipam.SetupIpamIpprefix,
	} {
		gvk, eventChan, err := setup(mgr, option, l, poll, namespace)
		if err != nil {
			return nil, err
		}
		eventChans[gvk] = eventChan
	}

	return eventChans, nil
}
