package main

import (
	"context"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// kubectl apply -f https://k8s.io/examples/pods/simple-pod.yaml

// 1. RestClient
// 2. ClientSet

func TestRestClient(t *testing.T) {
	// config
	// create client
	// get data

	// 第一个参数： '' 默认 ./kube/config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	config.GroupVersion = &v1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs
	config.APIPath = "/api"

	client, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}
	result := client.Get().Namespace("default").Resource("pods").Name("nginx").Do(context.TODO())

	pod := v1.Pod{}
	if err := result.Into(&pod); err != nil {
		t.Error(err)
		return
	}

	t.Log(pod.Name) // nginx
}

func TestClientSet(t *testing.T) {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	clinetset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	coreV1 := clinetset.CoreV1()
	pod, err := coreV1.Pods("default").Get(context.TODO(), "nginx", metav1.GetOptions{})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(pod.Name)
	}
}
