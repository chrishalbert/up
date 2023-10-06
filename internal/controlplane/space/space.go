package space

import (
	"context"

	xpcommonv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"

	"github.com/upbound/up/internal/controlplane"
	"github.com/upbound/up/internal/resources"
)

var (
	resource = resources.ControlPlaneGVK.GroupVersion().WithResource("controlplanes")
)

// Client is the client used for interacting with the ControlPlanes API in an
// Upbound Space.
type Client struct {
	c dynamic.Interface
}

// New instantiates a new Client.
func New(c dynamic.Interface) *Client {
	return &Client{
		c: c,
	}
}

// Get the ControlPlane corresponding to the given ControlPlane name.
func (c *Client) Get(ctx context.Context, name string) (*controlplane.Response, error) {
	u, err := c.c.
		Resource(resource).
		Get(
			ctx,
			name,
			metav1.GetOptions{},
		)

	if err != nil {
		return nil, err
	}

	return convert(&resources.ControlPlane{Unstructured: *u}), nil
}

// List all ControlPlanes within the Space.
func (c *Client) List(ctx context.Context) ([]*controlplane.Response, error) {
	list, err := c.c.
		Resource(resource).
		List(
			ctx,
			metav1.ListOptions{},
		)
	if err != nil {
		return nil, err
	}

	resps := []*controlplane.Response{}
	for _, u := range list.Items {
		resps = append(resps, convert(&resources.ControlPlane{Unstructured: u}))
	}

	return resps, nil
}

// Create a new ControlPlane with the given name and the supplied Options.
func (c *Client) Create(ctx context.Context, name string, opts controlplane.Options) (*controlplane.Response, error) {
	ctp := &resources.ControlPlane{}
	ctp.SetName(name)
	ctp.SetGroupVersionKind(resources.ControlPlaneGVK)
	ctp.SetWriteConnectionSecretToReference(&xpcommonv1.SecretReference{
		Name:      opts.SecretName,
		Namespace: opts.SecretNamespace,
	})

	u, err := c.c.
		Resource(resource).
		Create(
			ctx,
			ctp.GetUnstructured(),
			metav1.CreateOptions{},
		)

	if err != nil {
		return nil, err
	}

	return convert(&resources.ControlPlane{Unstructured: *u}), nil
}

// Delete the ControlPlane corresponding to the given ControlPlane name.
func (c *Client) Delete(ctx context.Context, name string) error {
	return c.c.
		Resource(resource).
		Delete(
			context.Background(),
			name,
			metav1.DeleteOptions{},
		)
}

func convert(ctp *resources.ControlPlane) *controlplane.Response {
	cnd := ctp.GetCondition(xpcommonv1.TypeReady)
	ref := ctp.GetConnectionSecretToReference()
	if ref == nil {
		ref = &xpcommonv1.SecretReference{}
	}

	return &controlplane.Response{
		ID:            ctp.GetControlPlaneID(),
		Name:          ctp.GetName(),
		Message:       cnd.Message,
		Status:        string(cnd.Reason),
		ConnName:      ref.Name,
		ConnNamespace: ref.Namespace,
	}
}
