package tfe

import (
	"context"
	"log"
)

// Compile-time proof of interface implementation.
var _ Registry = (*registry)(nil)

// Registry describes all the registry module related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs: https://www.terraform.io/docs/cloud/api/modules.html
type Registry interface {
	// Publish a module to the TFE private registry
	Publish(ctx context.Context, options ModulePublishOptions) (*Module, error)
}

// registry implements Registry.
type registry struct {
	client *Client
}

// ModulePublishOptions options for publishing a registry module
type ModulePublishOptions struct {
	// VCS connection information to import a module to the registry
	ModuleVCSOptions *ModuleVCSOptions `jsonapi:"attr,vcs-repo"`
}

// ModuleVCSOptions contains the configuration of a VCS integration.
type ModuleVCSOptions struct {
	Identifier        string `json:"identifier"`
	OAuthTokenID      string `json:"oauth-token-id"`
	DisplayIdentifier string `json:"display_identifier"`
}

// Module represents a registry module
type Module struct {
	ID   string `jsonapi:"primary,registry-modules"`
	Type string `json:"type"`

	Name      string `jsonapi:"attr,name"`
	Provider  string `jsonapi:"attr,provider"`
	Status    string `jsonapi:"attr,status"`
	CreatedAt string `jsonapi:"attr,created-at"`
	UpdatedAt string `jsonapi:"attr,updated-at"`

	// TODO:
	// version-statuses
	// permissions

	// Relations
	Organization *Organization `jsonapi:"relation,organization"`

	// Links
	// TODO
}

// Publish is used to publish a new module to the TFE private registry
func (r *registry) Publish(ctx context.Context, options ModulePublishOptions) (*Module, error) {
	req, err := r.client.newRequest("POST", "registry-modules", &options)
	if req == nil {
		log.Fatal(req)
	}
	if err != nil {
		return nil, err
	}

	m := &Module{}
	err = r.client.do(ctx, req, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
