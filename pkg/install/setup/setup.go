package setup

import (
	"context"
	"time"

	"github.com/solo-io/go-utils/installutils/kuberesource"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/solo-io/supergloo/pkg/install/linkerd"

	"github.com/solo-io/go-utils/installutils/kubeinstall"

	"github.com/solo-io/supergloo/pkg/install/gloo"
	"github.com/solo-io/supergloo/pkg/install/istio"

	"github.com/solo-io/supergloo/pkg/api/clientset"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	v1 "github.com/solo-io/supergloo/pkg/api/v1"
)

func RunInstallEventLoop(ctx context.Context, cs *clientset.Clientset, customErrHandler func(error)) error {
	ctx = contextutils.WithLogger(ctx, "install-event-loop")
	logger := contextutils.LoggerFrom(ctx)

	errHandler := func(err error) {
		if err == nil {
			return
		}
		logger.Errorf("install error: %v", err)
		if customErrHandler != nil {
			customErrHandler(err)
		}
	}

	installCache := kubeinstall.NewCache()
	go func() {
		logger.Infof("beginning install cache sync, this may take a while...")
		started := time.Now()
		if err := installCache.Init(ctx, cs.RestConfig, append(kubeinstall.DefaultFilters, cacheFilters...)...); err != nil {
			logger.Fatalf("failed to initialize installation cache: %v", err)
		}
		logger.Infof("finished install cache sync. took %v", time.Now().Sub(started))
	}()

	kubeInstaller, err := kubeinstall.NewKubeInstaller(cs.RestConfig, installCache, nil)
	if err != nil {
		return err
	}

	installSyncers := createInstallSyncers(cs, kubeInstaller)

	if err := startEventLoop(ctx, errHandler, cs, installSyncers); err != nil {
		return err
	}

	return nil
}

// Add install syncers here
func createInstallSyncers(clientset *clientset.Clientset, installer kubeinstall.Installer) v1.InstallSyncers {
	return v1.InstallSyncers{
		istio.NewInstallSyncer(
			installer,
			clientset.Supergloo.Mesh,
			reporter.NewReporter("istio-install-reporter", clientset.Supergloo.Install.BaseClient()),
		),
		linkerd.NewInstallSyncer(
			installer,
			clientset.Kube,
			clientset.Supergloo.Mesh,
			reporter.NewReporter("linkerd-install-reporter", clientset.Supergloo.Install.BaseClient()),
		),
		gloo.NewInstallSyncer(
			installer,
			clientset.Supergloo.MeshIngress,
			reporter.NewReporter("gloo-install-reporter", clientset.Supergloo.Install.BaseClient()),
		),
	}
}

// start the install event loop
func startEventLoop(ctx context.Context, errHandler func(err error), c *clientset.Clientset, syncers v1.InstallSyncers) error {
	installEmitter := v1.NewInstallEmitter(c.Supergloo.Install, c.Supergloo.Mesh, c.Supergloo.MeshIngress)
	installEventLoop := v1.NewInstallEventLoop(installEmitter, syncers)

	watchOpts := clients.WatchOpts{
		Ctx:         ctx,
		RefreshRate: time.Second * 1,
	}

	installEventLoopErrs, err := installEventLoop.Run(nil, watchOpts)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case err := <-installEventLoopErrs:
				errHandler(err)
			case <-ctx.Done():
			}
		}
	}()
	return nil
}

/*
to speed up the cache init, filter out resource types
*/
var cacheFilters = []kuberesource.FilterResource{
	func(resource schema.GroupVersionResource) bool {
		for _, ignoredType := range ignoreTypesForInstall {
			if resource.String() == ignoredType.String() {
				return true
			}
		}
		return false
	},
}

// types the installer should ignore and the cache should skip
var ignoreTypesForInstall = []schema.GroupVersionResource{
	{Resource: "certificatesigningrequests", Version: "v1beta1", Group: "certificates.k8s.io"},
}
