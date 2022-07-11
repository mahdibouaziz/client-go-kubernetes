# General Go informations

create a module

`go mod init [module name]`

to download all the required libraries based on the imports in the code

`go mod tidy`

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

# DEMO: Writing an app to get all the pods from the default namespace

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
