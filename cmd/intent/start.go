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

package intent

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	pkgmetav1 "github.com/yndd/ndd-core/apis/pkg/meta/v1"
	"github.com/yndd/ndd-yang/pkg/cache"
	"github.com/yndd/ndd-yang/pkg/yentry"
	"github.com/yndd/ndd-yang/pkg/yparser"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-runtime/pkg/ratelimiter"

	"github.com/yndd/ndd-yang/pkg/dispatcher"
	"github.com/yndd/nddo-ipam/internal/applogic"
	"github.com/yndd/nddo-ipam/internal/controllers"
	"github.com/yndd/nddo-ipam/internal/controllers/ipam"
	"github.com/yndd/nddo-ipam/internal/gnmiserver"
	"github.com/yndd/nddo-ipam/internal/kapi"
	"github.com/yndd/nddo-ipam/internal/restconf"
	"github.com/yndd/nddo-ipam/internal/yangschema"
	//+kubebuilder:scaffold:imports
)

var (
	metricsAddr          string
	probeAddr            string
	enableLeaderElection bool
	concurrency          int
	pollInterval         time.Duration
	namespace            string
	podname              string
	grpcServerAddress    string
)

// startCmd represents the start command for the network device driver
var startCmd = &cobra.Command{
	Use:          "start",
	Short:        "start the ipam nddo intent manager",
	Long:         "start the ipam ndd0 intent manager",
	Aliases:      []string{"start"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		zlog := zap.New(zap.UseDevMode(debug), zap.JSONEncoder())
		if debug {
			// Only use a logr.Logger when debug is on
			ctrl.SetLogger(zlog)
		}
		zlog.Info("create manager")
		mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
			Scheme:                 scheme,
			MetricsBindAddress:     metricsAddr,
			Port:                   9443,
			HealthProbeBindAddress: probeAddr,
			//LeaderElection:         false,
			LeaderElection:   enableLeaderElection,
			LeaderElectionID: "c66ce353.ndd.yndd.io",
		})
		if err != nil {
			return errors.Wrap(err, "Cannot create manager")
		}

		// initialize the config and state caches
		configCache := cache.New(
			[]string{ipam.GnmiTarget},
			cache.WithLogging(logging.NewLogrLogger(zlog.WithName("configcache"))))

		stateCache := cache.New(
			[]string{ipam.GnmiTarget},
			cache.WithLogging(logging.NewLogrLogger(zlog.WithName("statecache"))))

		// initialize the root schema
		rootSchema := yangschema.InitRoot(
			nil,
			yentry.WithLogging(logging.NewLogrLogger(zlog.WithName("yangschema"))))

		for _, path := range rootSchema.Resources {
			zlog.Info("Resource", "Path", yparser.GnmiPath2XPath(path, true))
		}

		// intialize the root handler
		rootResource := applogic.NewRoot(
			dispatcher.WithLogging(logging.NewLogrLogger(zlog.WithName("yresource"))),
			dispatcher.WithConfigCache(configCache),
			dispatcher.WithStateCache(stateCache),
			dispatcher.WithRootSchema(rootSchema),
		)

		// initialize controllers
		eventChans, err := controllers.Setup(mgr, nddCtlrOptions(concurrency), logging.NewLogrLogger(zlog.WithName("ipam")), pollInterval, namespace, rootSchema)
		if err != nil {
			return errors.Wrap(err, "Cannot add nddo controllers to manager")
		}

		// initialize kubernetes api
		a, err := kapi.New(config.GetConfigOrDie(),
			kapi.WithScheme(scheme),
			kapi.WithLogger(logging.NewLogrLogger(zlog.WithName("grpcserver"))),
		)
		if err != nil {
			return errors.Wrap(err, "Cannot create kubernetes client")
		}

		// create and start restconf server
		rcs, err := restconf.New(
			restconf.WithLogger(logging.NewLogrLogger(zlog.WithName("restconfserver"))),
			restconf.WithStateCache(stateCache),
			restconf.WithConfigCache(configCache),
			restconf.WithRootResource(rootResource),
			restconf.WithRootSchema(rootSchema),
			restconf.WithConfig(
				restconf.Config{
					Address: ":" + "9998",
				},
			),
		)
		if err != nil {
			return errors.Wrap(err, "unable to initialize server")
		}
		if err := rcs.Run(context.Background()); err != nil {
			return errors.Wrap(err, "unable to start restconf server")
		}

		// create and start gnmi server
		gs, err := gnmiserver.New(
			gnmiserver.WithKapi(a),
			gnmiserver.WithEventChannels(eventChans),
			gnmiserver.WithLogger(logging.NewLogrLogger(zlog.WithName("grpcserver"))),
			gnmiserver.WithStateCache(stateCache),
			gnmiserver.WithConfigCache(configCache),
			gnmiserver.WithRootResource(rootResource),
			gnmiserver.WithRootSchema(rootSchema),
			gnmiserver.WithResources(rootSchema.Resources),
			gnmiserver.WithConfig(
				gnmiserver.Config{
					Address:    ":" + strconv.Itoa(pkgmetav1.GnmiServerPort),
					SkipVerify: true,
					InSecure:   true,
				},
			),
			//gnmiserver.WithParser(logging.NewLogrLogger(zlog.WithName("grpcserver"))),
			//server.WithConnecter(intentlogic.New(logging.NewLogrLogger(zlog.WithName("grpcserver")))),
		)
		if err != nil {
			return errors.Wrap(err, "unable to initialize server")
		}

		state, err := gs.GetState()
		if err != nil {
			return errors.Wrap(err, "unable to get state from cache")
		}
		zlog.Info("New Server", "State", state)

		if err := gs.Run(context.Background()); err != nil {
			return errors.Wrap(err, "unable to start grpc server")
		}

		// +kubebuilder:scaffold:builder

		if err := mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
			return errors.Wrap(err, "unable to set up health check")
		}
		if err := mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
			return errors.Wrap(err, "unable to set up ready check")
		}

		zlog.Info("starting manager")
		if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
			return errors.Wrap(err, "problem running manager")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&metricsAddr, "metrics-bind-address", "m", ":8080", "The address the metric endpoint binds to.")
	startCmd.Flags().StringVarP(&probeAddr, "health-probe-bind-address", "p", ":8081", "The address the probe endpoint binds to.")
	startCmd.Flags().BoolVarP(&enableLeaderElection, "leader-elect", "l", false, "Enable leader election for controller manager. "+
		"Enabling this will ensure there is only one active controller manager.")
	startCmd.Flags().IntVarP(&concurrency, "concurrency", "", 1, "Number of items to process simultaneously")
	startCmd.Flags().DurationVarP(&pollInterval, "poll-interval", "", 1*time.Minute, "Poll interval controls how often an individual resource should be checked for drift.")
	startCmd.Flags().StringVarP(&namespace, "namespace", "n", os.Getenv("POD_NAMESPACE"), "Namespace used to unpack and run packages.")
	startCmd.Flags().StringVarP(&podname, "podname", "", os.Getenv("POD_NAME"), "Name from the pod")
	startCmd.Flags().StringVarP(&grpcServerAddress, "grpc-server-address", "s", "", "The address of the grpc server binds to.")
}

func nddCtlrOptions(c int) controller.Options {
	return controller.Options{
		MaxConcurrentReconciles: c,
		RateLimiter:             ratelimiter.NewDefaultProviderRateLimiter(ratelimiter.DefaultProviderRPS),
	}
}
