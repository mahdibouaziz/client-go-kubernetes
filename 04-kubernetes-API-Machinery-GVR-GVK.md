# kubernetes API Machinery - GVR - GVK

# `Kind` in API Machinery

The term `Kind` in the API Machinery means a `Kubernetes Type`

Example of kinds: Pod, Deployment, ...

Kinds are `grouped` in some groups for example `apps` group or `authorization` group.

There are also `Versions` in each group.

For example: `Deployment` is from: `apps` group and `v1` version ==> `apps/v1`

==> This is how the term `Group Version Kind` (`GVK`) comes.

==> `GroupVersionKind` is a way to present a particular Kubernetes Object using Api Machinery

## HTTP endpoint 

These kinds are **not** mapped one to one with an HTTP endpoint (a Kind can be a part of many HTTP endpoints)


# `Resource` in API Machinery

Resources are similar to Kind but they are represented in plulars && lower case.

for example a **resource** for the **Kind** `Deployment` is `deployments`.

Resources are `grouped` in some groups for example `apps` group or `authorization` group.

There are also `Versions` in each group.

==> This is how the term `Group Version Resource` (`GVR`) comes.

==> `GroupVersionResource` is a way to present a Resource using Api Machinery

## HTTP endpoint 

resources are directly mapped one to one with an HTTP endpoint

using resources we can always figure out the HTTP endpoint

The link is in the format of: `api/[group]/[version]/namespaces/[namespace]/[resource]`

NOTE: there is an exception for historical reasons for the `coreapi` group

for example:

- the **resource** `deployments` are in the **Group** `apps` and **version** `v1`
==> `api/apps/v1/namespaces/default/depoyments`

- the **resource** `replicasets` are in the **Group** `apps` and **version** `v1` ==> `api/apps/v1/namespaces/default/replicasets`

# Conversion from `Kind` to `Resource`

How do we convert a `GroupVersionKind` to `GroupVersionResource`?

to do this we need `RESTMapping` method from the `RESTMapper` interface

## `RESTMapper` interface, `RESTMapping` method and `RESTMapping` struct

we need to call the `RESTMapping` **method** from the `RESTMapper` interface, it expects a `GroupKind` and will return a `RESTMapping` **struct** that has a `Resource` and a `GroupVersionKind`.

The `RESTMapper` interface also has methods like:
- `KindFor`: expect a `GroupVersionResource` and returns a `GroupVersionKind`. it return an error if there are multiple matches

- `KindsFor`: expect a `GroupVersionResource` and returns a list of `GroupVersionKind`. 

# Scheme to convert Go type to GVK

If we have a Go struct that has all the needed attributes t become a Kubernetes object, we can convert that object to `GroupVersionKind` with the help of `ObjectKinds` that is a method in the `Scheme` struct from the `apimachinery`

Note: this can work only if the Kubernetes object is already registred and we can do that using `AddKnownTypes` method in the `Scheme` struct

# Summary

1. we have a Go struct that has the needed method to be treated like a K8s object
2. register the struct using the `AddKnownTypes`
3. using `Scheme` we can convert the registred struct to a particular `GVK` 
4. using `RESTMapping`, we can convert it to `GVR`
5. once we have a `GVR`, we can easilly get the **HTTP path** to which we want to query.