package main

import (
	"context"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var cntext []string = []string{"docker-desktop"}

func buildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func main() {
	config, err := buildConfigFromFlags(cntext[0], "/Users/k_shi/.kube/config")
	if err != nil {
		panic(err.Error())
	}
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	scaledRes := schema.GroupVersionResource{Group: "keda.sh", Version: "v1alpha1", Resource: "scaledobjects"}
	NewScaledObject := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "keda.sh/v1alpha1",
			"kind":       "ScaledObject",
			"metadata": map[string]interface{}{
				"name": "kafka-scaler",
			},
			"spec": map[string]interface{}{
				"pollingInterval": 30,
				"scaleTargetRef": map[string]interface{}{
					"name": "nginx-deployment",
				},
				"triggers": []map[string]interface{}{
					{
						"type": "kafka",
						"metadata": map[string]interface{}{
							"bootstrapServers":  "localhost:9092",
							"consumerGroup":     "my-group",
							"topic":             "test-topic",
							"lagThreshold":      "50",
							"offsetResetPolicy": "latest",
						},
					},
				},
			},
		},
	}
	fmt.Println("Creating deployment...")
	result, err := clientset.Resource(scaledRes).Namespace(apiv1.NamespaceDefault).Create(context.TODO(), NewScaledObject, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetName())
}
