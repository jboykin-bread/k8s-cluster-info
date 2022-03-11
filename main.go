package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubecontext := "my-k8s-cluster-context"
	clientset, err := connectToK8s(kubecontext)
	if err != nil {
		log.Fatal(err)
	}

	_, err = clientset.CoreV1().Services("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to connect to cluster and fetch services in default ns", err)
	}
}

// connectToK8s returns a kubernetes clientset if a kube config can be created from the given kube context.
func connectToK8s(kubeContext string) (*kubernetes.Clientset, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".kube", "config")
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: configPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: kubeContext,
		}).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create K8s config: %s", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s clientset: %s", err)
	}

	return clientset, nil
}
