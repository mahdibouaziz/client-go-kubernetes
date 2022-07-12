# Informer, Cache and Queue -Kubernetes Informers vs Watch 

# What is custom controller

A `custom controller` is an **application** that is going to listen on a particular Kubernetes resource (object - for example Pod) and as soon as this object is created, the custom controller is going to do certain things (for example as soon as a pod resource is created, the controller creates a respective service to expose this pod)

## How the custom controller is going to know if a particular resource has been created, deleted or updated?

we should have a mechanism to let this application (custom controller) know if an object has been created, updated or deleted.

by default the custom controller use the `watch` to look to this events

## Why we don't use `watch`

`watch` is going to query the API-Server a lot to get to know if a particular resource was created, deleted or updated and that is going to **increase load on the API-Server** and the processes are going to get slow because of these number of requests to the API-Server

==> the solution is the `Informers`

## What is an `Informer`

`informer` use `watch` internally, but it uses the in **memory local cache** and because of that we don't need to query the API-Server a lot.

If we create an `informer`, it is going to call the `API-Server` with the `List()` and it is going to return all the details in the informer as a **local cache memory**. Now if we need anything we get it from the **Informer local cache memory** 

## What is `Shared Informer Factory`

We can create an `Informer` for every group version resource, but in this case for each group version resouce we will call the API-Server and this will increase the load on the API-Server.

==> The solution is the `Shared Informer Factory`

# Demo: how to use shared informer factory

We suppose that you already have your clientset configured
```go
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
```

# Don't update informer cache resources

All these resources are from in memory cache and that cache is maintained by informer ==> we should not change these resources directly.

If we want to change we should create `DeepCopy()` of the object before changing 

# The `NewFilteredSharedInformerFactory(...)` method

- The `NewSharedInformerFactory(....)` is going to return informer factory for all the resources and all the namespaces

- We can use `NewFilteredSharedInformerFactory` to filter with some options like the namespace, and other metadata

```go
    // Example using the NewFilteredSharedInformerFactory
	informerfactory = informers.NewFilteredSharedInformerFactory(clientset, 30*time.Second, "kube-system", func(lo *metav1.ListOptions) {
		// specify the filter here (metadata filters of type ListOptions)
		lo.LabelSelector = ""
		lo.APIVersion = "v1"
	})
```

# Queues

with informers, we can **register some functions** like we did in this part of code

```go
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
```

When a function is called, it can **fail**, in this case we need to **retry** the call of this function.

To do this we have the `Queue`, whenever we need to call a function:
1. we `enqueue` the function call to the queue 
2. it will stay in the queue and run in another go routine 
3. if the function call fails, we `enqueue` the function call again
4. if the function call succeed, we remove the function call from the queue
