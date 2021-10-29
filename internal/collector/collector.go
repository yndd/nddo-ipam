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
	"sync"
	"time"

	"github.com/karimra/gnmic/target"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/pkg/errors"
	"github.com/yndd/ndd-runtime/pkg/logging"
)

const (
	defaultTargetReceivebuffer = 1000
	defaultLockRetry           = 5 * time.Second
	defaultRetryTimer          = 10 * time.Second

	// errors
	errCreateSubscriptionRequest = "cannot create subscription request"
)

// Collector defines the interfaces for the collector
type Collector interface {
	GetTarget() *target.Target
	Lock()
	Unlock()
	GetSubscription(subName string) bool
	StopSubscription(ctx context.Context, subName string) error
	StartSubscription(ctx context.Context, subName string, prefix *gnmi.Path, paths []*gnmi.Path) error
	StartSubscriptionHandler(ctx context.Context, subName string, prefix *gnmi.Path, paths []*gnmi.Path, handle func(resp *gnmi.SubscribeResponse))
}

// CollectorOption can be used to manipulate Options.
type CollectorOption func(*collector)

// WithCollectorLogger specifies how the collector logs messages.
func WithCollectorLogger(log logging.Logger) CollectorOption {
	return func(o *collector) {
		o.log = log
	}
}

// collector defines the parameters for the collector
type collector struct {
	targetReceiveBuffer uint
	retryTimer          time.Duration
	target              *target.Target
	//targetSubRespChan   chan *collector.SubscribeResponse
	//targetSubErrChan    chan *collector.TargetError
	subscriptions map[string]*Subscription
	mutex         sync.RWMutex
	log           logging.Logger
}

// NewCollector creates a new GNMI collector
func New(t *target.Target, opts ...CollectorOption) Collector {
	c := &collector{
		target:              t,
		subscriptions:       make(map[string]*Subscription),
		mutex:               sync.RWMutex{},
		targetReceiveBuffer: defaultTargetReceivebuffer,
		retryTimer:          defaultRetryTimer,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Lock locks a gnmi collector
func (c *collector) GetTarget() *target.Target {
	return c.target
}

// Lock locks a gnmi collector
func (c *collector) Lock() {
	c.mutex.RLock()
}

// Unlock unlocks a gnmi collector
func (c *collector) Unlock() {
	c.mutex.RUnlock()
}

// GetSubscription returns a bool based on a subscription name
func (c *collector) GetSubscription(subName string) bool {
	if _, ok := c.subscriptions[subName]; !ok {
		return true
	}
	return false
}

// StopSubscription stops a subscription
func (c *collector) StopSubscription(ctx context.Context, subName string) error {
	c.log.WithValues("subscription", subName)
	c.log.Debug("subscription stop...")
	c.subscriptions[subName].stopCh <- true // trigger quit

	c.log.Debug("subscription stopped")
	return nil
}

// StartSubscription starts a subscription
func (c *collector) StartSubscription(ctx context.Context, subName string, prefix *gnmi.Path, paths []*gnmi.Path) error {
	log := c.log.WithValues("subscription", subName, "Paths", paths)
	log.Debug("subscription start...")

	req, err := CreateSubscriptionRequest(c.target.Config.Name, subName, prefix, paths)
	if err != nil {
		c.log.Debug(errCreateSubscriptionRequest, "error", err)
		return errors.Wrap(err, errCreateSubscriptionRequest)
	}

	log.Debug("subscription request", "request", req)
	go func() {
		c.target.Subscribe(ctx, req, subName)
	}()
	log.Debug("subscription started ...")

	for {
		select {
		case <-c.subscriptions[subName].stopCh: // execute quit
			c.subscriptions[subName].cancelFn()
			c.mutex.Lock()
			delete(c.subscriptions, subName)
			c.mutex.Unlock()
			c.log.Debug("subscription cancelled")
			return nil
		}
	}
}

func (c *collector) StartSubscriptionHandler(ctx context.Context, subName string, prefix *gnmi.Path, paths []*gnmi.Path, handle func(resp *gnmi.SubscribeResponse)) {
	c.log.Debug("Starting subscription Handler...")

	// initialize new subscription
	ctx, cancel := context.WithCancel(ctx)

	c.subscriptions[subName] = &Subscription{
		stopCh:   make(chan bool),
		cancelFn: cancel,
		ctx:      ctx,
	}

	c.Lock()
	go c.StartSubscription(ctx, subName, prefix, paths)
	c.Unlock()

	chanSubResp, chanSubErr := c.target.ReadSubscriptions()

	for {
		select {
		case resp := <-chanSubResp:
			//c.log.Debug("subscription", "response", resp.Response)
			handle(resp.Response)
		case tErr := <-chanSubErr:
			c.log.Debug("subscription error", "subscriptionName", tErr.SubscriptionName, "error", tErr.Err)
			time.Sleep(60 * time.Second)
		case <-c.subscriptions[subName].stopCh:
			c.log.Debug("Stopping subscription process...")
			return
		}
	}
}
