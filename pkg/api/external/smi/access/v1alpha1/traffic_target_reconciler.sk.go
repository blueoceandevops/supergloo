// Code generated by solo-kit. DO NOT EDIT.

package v1alpha1

import (
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/reconcile"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
)

// Option to copy anything from the original to the desired before writing. Return value of false means don't update
type TransitionTrafficTargetFunc func(original, desired *TrafficTarget) (bool, error)

type TrafficTargetReconciler interface {
	Reconcile(namespace string, desiredResources TrafficTargetList, transition TransitionTrafficTargetFunc, opts clients.ListOpts) error
}

func trafficTargetsToResources(list TrafficTargetList) resources.ResourceList {
	var resourceList resources.ResourceList
	for _, trafficTarget := range list {
		resourceList = append(resourceList, trafficTarget)
	}
	return resourceList
}

func NewTrafficTargetReconciler(client TrafficTargetClient) TrafficTargetReconciler {
	return &trafficTargetReconciler{
		base: reconcile.NewReconciler(client.BaseClient()),
	}
}

type trafficTargetReconciler struct {
	base reconcile.Reconciler
}

func (r *trafficTargetReconciler) Reconcile(namespace string, desiredResources TrafficTargetList, transition TransitionTrafficTargetFunc, opts clients.ListOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "trafficTarget_reconciler")
	var transitionResources reconcile.TransitionResourcesFunc
	if transition != nil {
		transitionResources = func(original, desired resources.Resource) (bool, error) {
			return transition(original.(*TrafficTarget), desired.(*TrafficTarget))
		}
	}
	return r.base.Reconcile(namespace, trafficTargetsToResources(desiredResources), transitionResources, opts)
}
