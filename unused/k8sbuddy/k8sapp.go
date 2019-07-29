package k8sbuddy

import (
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
)

// K8sApp holds info on a
type K8sApp struct {
	Kind string
	SS   *appsv1.StatefulSet
	DP   *appsv1.Deployment
	DS   *appsv1.DaemonSet
	NS   *apiv1.Namespace
}
