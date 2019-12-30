package main

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func main() {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()

	configOverrides := &clientcmd.ConfigOverrides{}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	podInterface := clientset.CoreV1().Pods("default")

	w, err := podInterface.Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	podLabels := make(map[string]map[string]string)

	for {
		select {
			case e := <- w.ResultChan():
				pod := e.Object.(*corev1.Pod)
				switch e.Type {
				case watch.Added, watch.Modified:
					podLabels[pod.Name] = pod.Labels
				case watch.Deleted:
					delete(podLabels, pod.Name)
				}
				printLabels(podLabels)
		}
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
