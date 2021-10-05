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
	"net"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/pkg/errors"
	"github.com/yndd/ndd-runtime/pkg/logging"
	ynddparser "github.com/yndd/ndd-yang/pkg/parser"
	"golang.org/x/sync/semaphore"
	"google.golang.org/grpc"
	"sigs.k8s.io/controller-runtime/pkg/event"

	"github.com/yndd/ndd-yang/pkg/cache"

	"github.com/yndd/nddo-ipam/internal/connector"
	"github.com/yndd/nddo-ipam/internal/controllers/ipam"
	"github.com/yndd/nddo-ipam/internal/kapi"
)

const (
	// defaults
	defaultMaxSubscriptions = 64
	defaultMaxGetRPC        = 64
)

type Config struct {
	// Address
	GrpcServerAddress string
	// Generic
	MaxSubscriptions int64
	MaxUnaryRPC      int64
	// TLS
	InSecure   bool
	SkipVerify bool
	CaFile     string
	CertFile   string
	KeyFile    string
	// observability
	EnableMetrics bool
	Debug         bool
}

// Option can be used to manipulate Options.
type ServerOption func(*Server)

// WithLogger specifies how the Reconciler should log messages.
func WithServerLogger(log logging.Logger) ServerOption {
	return func(s *Server) {
		s.log = log
	}
}

func WithServerConfig(cfg Config) ServerOption {
	return func(s *Server) {
		s.cfg = cfg
	}
}

func WithParser(log logging.Logger) ServerOption {
	return func(s *Server) {
		s.parser = ynddparser.NewParser(ynddparser.WithLogger(log))
	}
}

func WithKapi(a *kapi.Kapi) ServerOption {
	return func(s *Server) {
		s.client = a
	}
}

func WithEventChannels(e map[string]chan event.GenericEvent) ServerOption {
	return func(s *Server) {
		s.EventChannels = e
	}
}

func WithConnecter(c connector.Connector) ServerOption {
	return func(s *Server) {
		s.connector.Connector = c
	}
}

type mrConnector struct {
	connector.Connector
}

type Server struct {
	gnmi.UnimplementedGNMIServer
	cfg Config

	// kubernetes
	client        *kapi.Kapi
	EventChannels map[string]chan event.GenericEvent

	// schema
	configCache *cache.Cache
	stateCache  *cache.Cache
	//schemaRaw interface{}
	//schema    *ipamv1alpha1.Nddoipam
	// gnmi calls
	subscribeRPCsem *semaphore.Weighted
	unaryRPCsem     *semaphore.Weighted
	// logging and parsing
	parser *ynddparser.Parser
	//handler *Handler
	log logging.Logger

	connector mrConnector

	// context
	ctx context.Context
}

func NewServer(opts ...ServerOption) (*Server, error) {
	s := &Server{
		configCache: cache.New([]string{ipam.GnmiTarget}),
		stateCache:  cache.New([]string{ipam.GnmiTarget}),

		//schemaRaw: nil,
		// check if the resource has a key or not
		//schema: &ipamv1alpha1.Nddoipam{},
		//	Topology: make([]*NddotopologyTopology, 0),
		//},
	}

	for _, opt := range opts {
		opt(s)
	}

	s.ctx = context.Background()

	// get the original status from k8s
	if err := s.GetInitialState(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) GetConfigCache() *cache.Cache {
	return s.configCache
}

func (s *Server) GetStateCache() *cache.Cache {
	return s.stateCache
}

func (s *Server) Run(ctx context.Context) error {
	log := s.log.WithValues("grpcServerAddress", s.cfg.GrpcServerAddress)
	log.Debug("grpc server run...")
	errChannel := make(chan error)
	go func() {
		if err := s.Start(); err != nil {
			errChannel <- errors.Wrap(err, errStartGRPCServer)
		}
		errChannel <- nil
	}()
	return nil
}

// Start GRPC Server
func (s *Server) Start() error {
	s.subscribeRPCsem = semaphore.NewWeighted(defaultMaxSubscriptions)
	s.unaryRPCsem = semaphore.NewWeighted(defaultMaxGetRPC)
	log := s.log.WithValues("grpcServerAddress", s.cfg.GrpcServerAddress)
	log.Debug("grpc server start...")

	// create a listener on a specific address:port
	l, err := net.Listen("tcp", s.cfg.GrpcServerAddress)
	if err != nil {
		return errors.Wrap(err, errCreateTcpListener)
	}

	// TODO, proper handling of the certificates with CERT Manager
	/*
		opts, err := s.serverOpts()
		if err != nil {
			return err
		}
	*/
	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach the gnmi service to the grpc server
	gnmi.RegisterGNMIServer(grpcServer, s)

	// start the server
	log.Debug("grpc server serve...")
	if err := grpcServer.Serve(l); err != nil {
		s.log.Debug("Errors", "error", err)
		return errors.Wrap(err, errGrpcServer)
	}
	return nil
}
