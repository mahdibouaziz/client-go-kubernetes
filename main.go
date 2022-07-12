package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
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

	// Using Shared Informer factory
	// create a shared informer factory from the clientset
	// the second argument is the period of the default resync time (to keep the cache updated)
	informerfactory := informers.NewSharedInformerFactory(clientset, 30*time.Second)

	// get a specif informer for pods
	podInformer := informerfactory.Core().V1().Pods()

	// we can register some functions to listen on some events
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(new interface{}) {
			fmt.Println("Add was called")
		},
		UpdateFunc: func(old, new interface{}) {
			// This function is called also in every resync period
			// before doing anything we need to check the `resource version`
			// of the new object to figure out if the resource has changed
			// NOTE: `resourceVersion` change every time a particular resource was edited
			fmt.Println("update was called")
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("delete was called")
		},
	})

	// Once we have the informer, we need to start it
	// start == initialize the Informer local cache memory from the API-Server
	informerfactory.Start(wait.NeverStop)

	// wait for the first `List()` call (sync) to be completed
	informerfactory.WaitForCacheSync(wait.NeverStop)

	// from the pod informer we call the lister to get the pods from the local cache memory of the informer
	pod, err := podInformer.Lister().Pods("default").Get("nginx")
	fmt.Println(pod.Name)

	// // Example using the NewFilteredSharedInformerFactory
	// informerfactory = informers.NewFilteredSharedInformerFactory(clientset, 30*time.Second, "kube-system", func(lo *metav1.ListOptions) {
	// 	// specify the filter here (metadata filters of type ListOptions)
	// 	lo.LabelSelector = ""
	// 	lo.APIVersion = "v1"
	// })

}
