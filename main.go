package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main2() {
	// 1.
	// this program is expecting a flag named 'kubeconfig'
	// and if the value is not provided it will get a default value (the second argument)
	kubeconfig := flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Location to your kubeconfig file")
	flag.Parse()
	// fmt.Printf("%s\n", *kubeconfig)

	// 2. build a valid config from the *kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// config, err := rest.InClusterConfig()
	// if err != nil {
	// 	panic(err.Error())
	// }

	// runtime.Object

	// 3. create a clientset from the config
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 4. get all the pods from the default namespace
	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("pods from the default namespace")
	for _, pod := range pods.Items {
		fmt.Printf("Pod Name:%s\n", pod.Name)
	}

	// 5. get all the deployments from the default namespace
	deployments, err := clientset.AppsV1().Deployments("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("deployments from the default namespace")
	for _, deploy := range deployments.Items {
		fmt.Printf("Deployment Name:%s\n", deploy.Name)
	}

}
