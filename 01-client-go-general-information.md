# How we can access the K8s cluster programatically using go

to talk to K8s cluster we need to use a set of libraries

- **client-go**
- **api**: provides all the API types
- **apimachinery**

We can perform any operation to the cluster using these libraries.

They are imported in the go project using the `k8s.io/client-go`, `k8s.io/api`, and `k8s.io/apimachinery`

# client-go

expose all the interfaces that can be used to interact or do certain operations with the K8s API-Server.

for example:

- `tools` from the `client-go` library is responsible for creating `clientset`
- `kubernetes` from the `client-go` library contains all types available

# api

All k8s resources (pods, deployments, services, ....) are maintained in the api repository

# apimachinery

Contains common code that can be used while developing any K8s API

for example, when we try to list all the pods, we can specify some kind of filter with the help of apimachinery

# DEMO: Writing an app (outside the cluster) to get all the pods from the default namespace 

we need to initialize a Go project:

`go mod init client-go-project`

this will create a `go.mod` file that will contains information about our dependencies, versions, ...

let's create a file called `main.go` that contains a **main** function.

from this function we want to talk to **K8s API-Server** using the `client-go`, `api`, and `apimachinery` libraries.

to be able to access (talk) to our cluster, we need to specify where our `kubeconfig` file is. Otherwise we will not be able to acces our cluster.

1. Get the `kubeconfig` file as an argument to our program using the `flag` module (predefined in go)

```go
	// this program is expecting a flag named 'kubeconfig'
	// and if the value is not provided it will get a default value (the second argument)
	kubeconfig := flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Location to your kubeconfig file")
	flag.Parse()
    fmt.Printf("%s\n", *kubeconfig)
```

2. Now we need to import the `k8s.io/client-go/tools/clientcmd` library to **build a config from the kubeconfig flag** 

```go
    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
		panic(err.Error())
	}
```

3. Now we have a well formed **config**, let's create our `clientset` from the `config` with the help of `k8s.io/client-go/kubernetes` library

```go
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
		panic(err.Error())
	}
```

this is called a clientset because it is **set of clients** that can be used to interact with the resources from different api-versions.

4. Let's try to get all the pods from the default namespace, to do this, we need a `context`  and a `ListOptions` from the apimachinery (we can import it `metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"`)

```go
    pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

    fmt.Println("pods from the default namespace")
	for _, pod := range pods.Items {
		fmt.Printf("Pod Name:%s\n", pod.Name)
	}
```

the full code:

```go
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

func main() {
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

}

```

## get all the deployments from the default namespace
 
we just use the `clientset` and go to the right `GroupVersion` (which is `AppV1` in the case of a deployment)

 ```go
    deployments, err := clientset.AppsV1().Deployments("default").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
		panic(err.Error())
	}

	fmt.Println("deployments from the default namespace")
	for _, deploy := range deployments.Items {
		fmt.Printf("Deployment Name:%s\n", deploy.Name)
	}
 ```

# Summary

- to interact with k8s api server we will have to use `clientset`.
- this is how we can query/interact with k8s resources: `clientset.GroupVersion().Resources(namespace).Verb()`

examples:

A Pod resource, it's part of `Core` **API group** and `V1` **version** ==> `clientset.corev1().Pods("namespace").List(....)`

A Deployments, it's part of `Apps` **API group** and `V1` **version** ==> `clientset.AppsV1().Deployments(namespace).List(....)`

you can use this command to get the API group and version of the ressource you're looking for 

`kubectl api-resources`


# DEMO: Writing an app (inside the cluster)

This example shows you how to configure a client with client-go to authenticate to the Kubernetes API from an application running inside the Kubernetes cluster.

client-go uses the **Service Account token** mounted inside the Pod at the `/var/run/secrets/kubernetes.io/serviceaccount` path when the `rest.InClusterConfig()` is used.

we need to change our previous code to use the `rest.InClusterConfig()`

to do that we just need to import `"k8s.io/client-go/rest"`

and then 

```go
    config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
```

and the rest of the code is the same

## Full code
```go
    package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

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

```


to run an application inside a K8s cluster, we need to containarize the application  

## Containarize an application

1. build the app:

`GOOS=linux go build -o ./app .`

2. create a Dockerfile

```Dockerfile
FROM debian
COPY ./app /app
ENTRYPOINT /app
```

3. build the image

If you are running a `Minikube` cluster, you can build this image directly on the Docker engine of the Minikube node without pushing it to a registry. To build the image on Minikube:

```bash
eval $(minikube docker-env)
docker build -t in-cluster .
```

If you are **not using Minikube**, you should build this image and push it to a registry that your Kubernetes cluster can pull from.

`docker build -t in-cluster .`

## RBAC Configuration

If you have RBAC enabled on your cluster, use the following snippet to create role binding which will grant the default service account view permissions.

`kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default`

## run the image in a Pod

`kubectl run --rm -i demo --image=in-cluster`

and then run

`kubectl logs demo`