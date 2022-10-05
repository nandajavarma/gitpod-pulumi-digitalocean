package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	// appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	// corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	// metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	do "github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
)

func createCluster(ctx *pulumi.Context) (*do.KubernetesCluster, error) {
	cluster, err := do.NewKubernetesCluster(ctx, "gitpod-cluster", &do.KubernetesClusterArgs{
		Region:  pulumi.String("fra1"),
		Version: pulumi.String("1.24.4-do.0"),
		NodePool: &do.KubernetesClusterNodePoolArgs{
			Name:      pulumi.String("services-nodepool"),
			Size:      pulumi.String("s-2vcpu-2gb"),
			NodeCount: pulumi.Int(3),
			Labels: pulumi.StringMap{
				"gitpod.io/workload_meta":               pulumi.String("true"),
				"gitpod.io/workload_ide":                pulumi.String("true"),
				"gitpod.io/workload_workspace_services": pulumi.String("true"),
			},
		},
	})

	if err != nil {
		return nil, err
	}
	_, err = do.NewKubernetesNodePool(ctx, "workspaces-nodepool", &do.KubernetesNodePoolArgs{
		ClusterId: cluster.ID(),
		Size:      pulumi.String("c-2"),
		NodeCount: pulumi.Int(2),
		Tags: pulumi.StringArray{
			pulumi.String("backend"),
		},
		Labels: pulumi.StringMap{
			"gitpod.io/workload_workspace_regular":  pulumi.String("true"),
			"gitpod.io/workload_workspace_headless": pulumi.String("true"),
		},
	})

	return cluster, nil
}

func createRegistry(ctx *pulumi.Context) (*do.ContainerRegistry, error) {
	args := do.ContainerRegistryArgs{
		SubscriptionTierSlug: pulumi.String("starter"),
	}

	registry, err := do.NewContainerRegistry(ctx, "gitpod-registry", &args)
	if err != nil {
		return nil, err
	}

	return registry, err
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		cluster, err := createCluster(ctx)
		if err != nil {
			return err
		}

		ctx.Export("kubeconfig", cluster.KubeConfigs)

		registry, err := createRegistry(ctx)
		if err != nil {
			return err
		}

		ctx.Export("registry", registry.ToContainerRegistryOutput())

		// TODO
		// create a k8s cluster
		// create a registry
		// create a storage
		// create a database
		return nil
	})
}
