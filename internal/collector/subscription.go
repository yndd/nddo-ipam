/*
Copyright 2021 Wim Henderickx.

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

package collector

import (
	"context"
	"strings"

	"github.com/openconfig/gnmi/proto/gnmi"
)

// Subscription defines the parameters for the subscription
type Subscription struct {
	stopCh   chan bool
	cancelFn context.CancelFunc
	ctx      context.Context
}

// CreateSubscriptionRequest create a gnmi subscription
func CreateSubscriptionRequest(target, subName string, prefix *gnmi.Path, paths []*gnmi.Path) (*gnmi.SubscribeRequest, error) {
	// create subscription

	modeVal := gnmi.SubscriptionList_Mode_value[strings.ToUpper("STREAM")]
	qos := &gnmi.QOSMarking{Marking: 21}

	subscriptions := make([]*gnmi.Subscription, len(paths))
	for i, p := range paths {
		subscriptions[i] = &gnmi.Subscription{Path: p}
		switch gnmi.SubscriptionList_Mode(modeVal) {
		case gnmi.SubscriptionList_STREAM:
			mode := gnmi.SubscriptionMode_value[strings.Replace(strings.ToUpper("ON_CHANGE"), "-", "_", -1)]
			subscriptions[i].Mode = gnmi.SubscriptionMode(mode)
		}
	}
	req := &gnmi.SubscribeRequest{
		Request: &gnmi.SubscribeRequest_Subscribe{
			Subscribe: &gnmi.SubscriptionList{
				Prefix:       prefix,
				Mode:         gnmi.SubscriptionList_Mode(modeVal),
				Encoding:     gnmi.Encoding_JSON,
				Subscription: subscriptions,
				Qos:          qos,
			},
		},
	}
	return req, nil
}
