# How Kubernetes Objects are represented in a Go code?

# Which Go structs are Kubernetes Objects

Any go `struct` that implement `runtime.Object` interface from the `apimachinery` package (`"k8s.io/apimachinery/pkg/runtime"`) is a **kubernetes object**

# the `runtime.Object` interface

```go
type Object interface {
	GetObjectKind() schema.ObjectKind
	DeepCopyObject() Object
}
```

this object has 2 methods:
- `GetObjectKind()` that return another interface called `ObjectKind` that has 2 methods
```go
    type ObjectKind interface {
	// SetGroupVersionKind sets or clears the intended serialized kind of an object. 
	// Passing kind nil should clear the current setting.
	SetGroupVersionKind(kind GroupVersionKind)
	// GroupVersionKind returns the stored group, version, and kind of an object, 
	// or return an empty struct if the object does not expose or provide these fields.
	GroupVersionKind() GroupVersionKind
}
```
- `DeepCopyObject()` to create a deep copy of the object

==> **if a particular go struct implements `SetGroupVersionKind(kind GroupVersionKind)`, `GroupVersionKind()`, and `DeepCopyObject()` it means that this struct can be considered as a kubernetes object**.

# Closer look into Pod Kubernetes Object

```go
type Pod struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of the desired behavior of the pod.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Spec PodSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of the pod.
	// This data may not be up to date.
	// Populated by the system.
	// Read-only.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Status PodStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}
```

## How this Pod struct implement the 3 methods: `SetGroupVersionKind(kind GroupVersionKind)`, `GroupVersionKind()`, and `DeepCopyObject()`

The 2 methods `SetGroupVersionKind(kind GroupVersionKind)`, `GroupVersionKind()` are implements in the `metav1.TypeMeta` struct in this way:

```go
// SetGroupVersionKind satisfies the ObjectKind interface for all objects that embed TypeMeta
func (obj *TypeMeta) SetGroupVersionKind(gvk schema.GroupVersionKind) {
	obj.APIVersion, obj.Kind = gvk.ToAPIVersionAndKind()
}

// GroupVersionKind satisfies the ObjectKind interface for all objects that embed TypeMeta
func (obj *TypeMeta) GroupVersionKind() schema.GroupVersionKind {
	return schema.FromAPIVersionAndKind(obj.APIVersion, obj.Kind)
}
```

the `DeepCopyObject()` is implemented in the Pod struct
```go
func (in *Pod) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
```

==> **Pod implemets `runtime.Object` interface with the help of `metav1.TypeMeta` (implement the 2 methods) and it directly implement `DeepCopyObject()`**

# Closer look into `TypeMeta` struct 


Any kubernetes object will have `TypeMeta` (that have the 2 methods of the `runtime.Object` interface) and also have 2 attributes `Kind` and `APIVersion`


```go
type TypeMeta struct {
	// Kind is a string value representing the REST resource this object represents.
	// Servers may infer this from the endpoint the client submits requests to.
	// Cannot be updated.
	// In CamelCase.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	Kind string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`

	// APIVersion defines the versioned schema of this representation of an object.
	// Servers should convert recognized schemas to the latest internal value, and
	// may reject unrecognized values.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	// +optional
	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`
}
```
==> this is used to specify the Kind and the APIVersion (like the 2 first lines in the yaml file)

# Closer look into `ObjectMeta` struct 

used to specify any other **metadata** details about the Kubernetes object.

```go
type ObjectMeta struct {
	// Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically.
	// Cannot be updated.
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`

	// GenerateName is an optional prefix, used by the server, to generate a unique  name ONLY IF the Name field has not been provided.
	GenerateName string `json:"generateName,omitempty" protobuf:"bytes,2,opt,name=generateName"`

	Namespace string `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`

	// UID is the unique in time and space value for this object. It is typically generated by the server on successful creation of a resource and is not allowed to change on PUT operations.
	UID types.UID `json:"uid,omitempty" protobuf:"bytes,5,opt,name=uid,casttype=k8s.io/kubernetes/pkg/types.UID"`

	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	Labels map[string]string `json:"labels,omitempty" protobuf:"bytes,11,rep,name=labels"`

	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	Annotations map[string]string `json:"annotations,omitempty" protobuf:"bytes,12,rep,name=annotations"`

	... and more check definition file
}
```

# Closer look into `Spec` and `Status` structs

- `Spec` is the desired State
- `Status` is the current State

they change for each specific object (each object has its own implementation)


# Closer look into `kubernetes.Clientset` structs

`Clientset` contains the clients for groups. 

Each group has exactly one  version included in a Clientset. 

We don't use the clientset directly, because it calls the API-server for each request, we need a **cache mechanism** ==> **Informer**

# Closer look into the config that we get from `config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)`

This particular config is used to get the `clientset`.

A config is a struct and we can also specify custom configuration:

```go
type Config struct {
	Host string
	APIPath string
	ContentConfig
	Username string
	Password string `datapolicy:"password"`
	BearerToken string `datapolicy:"token"`
	BearerTokenFile string
	Impersonate ImpersonationConfig
	AuthProvider *clientcmdapi.AuthProviderConfig
	AuthConfigPersister AuthProviderConfigPersister
	ExecProvider *clientcmdapi.ExecConfig
	TLSClientConfig
	UserAgent string
	DisableCompression bool
	Transport http.RoundTripper
	WrapTransport transport.WrapperFunc
	QPS float32
	Burst int
	RateLimiter flowcontrol.RateLimiter
	WarningHandler WarningHandler
	Timeout time.Duration
	Dial func(ctx context.Context, network, address string) (net.Conn, error)
	Proxy func(*http.Request) (*url.URL, error)	
}
```

