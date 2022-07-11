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
