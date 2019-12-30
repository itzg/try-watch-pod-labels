package main

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	// allow for all authentication types
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func main() {
	// Use clientcmd's standard behavior for both in-cluster and external config loading
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()

	// ...customize, if needed
	configOverrides := &clientcmd.ConfigOverrides{}

	// ...and load the config
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create a top-level clientset for accessing well-known kubernetes API resources
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// Create a clientset interface specifically for pod access in the default namespace
	podInterface := clientset.CoreV1().Pods("default")

	// Start a watch operation on all pod resources
	// ...list options could customize what label selectors, etc to use
	w, err := podInterface.Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// For this demo app, we'll maintain a mapping of pod name to pod labels
	podLabels := make(map[string]map[string]string)

	for {
		// Watch for events
		e := <-w.ResultChan()
		// can assume a pod resource since the watch was created for that resource type
		pod := e.Object.(*corev1.Pod)

		// Update tracking by event type
		switch e.Type {
		case watch.Added, watch.Modified:
			podLabels[pod.Name] = pod.Labels
		case watch.Deleted:
			delete(podLabels, pod.Name)
		}

		// ...and output latest state
		printLabels(podLabels)
	}
}

func printLabels(podLabels map[string]map[string]string) {
	fmt.Printf("=================================\n")
	for podName, labels := range podLabels {
		fmt.Printf("%s:\n", podName)
		for key, value := range labels {
			fmt.Printf("\t%s: %s\n", key, value)
		}
	}
}
