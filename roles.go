package gocd

import (
	multierror "github.com/hashicorp/go-multierror"
)

type RoleAttributeProperty struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type RoleAttributes struct {
	Users        []string                `json:"users,omitempty"`
	AuthConfigId string                  `json:"auth_config_id,omitempty"`
	Properties   []RoleAttributeProperty `json:"properties,omitempty"`
}

type RolePolicy struct {
	Permission string `json:"permission,omitempty"`
	Action     string `json:"action,omitempty"`
	Type       string `json:"type,omitempty"`
	Resource   string `json:"resource,omitempty"`
}

type Role struct {
	Name       string         `json:"name,omitempty"`
	Type       string         `json:"type,omitempty"`
	Attributes RoleAttributes `json:"attributes,omitempty"`
	Policy     []RolePolicy   `json:"policy,omitempty"`
}

type RoleBulkUsers struct {
	Add    []string `json:"add,omitempty"`
	Remove []string `json:"remove,omitempty"`
}

type RoleBulkOperation struct {
	Role  string        `json:"role,omitempty"`
	Users RoleBulkUsers `json:"users,omitempty"`
}

func (c *DefaultClient) GetAllRoles() ([]*Role, error) {
	type EmbeddedObj struct {
		Role []*Role `json:"roles"`
	}
	type AllRolesResponse struct {
		Embedded EmbeddedObj `json:"_embedded"`
	}
	res := new(AllRolesResponse)
	headers := map[string]string{"Accept": "application/vnd.go.cd+json"}
	err := c.getJSON("/go/api/admin/security/roles", headers, res)
	if err != nil {
		return []*Role{}, err
	}

	var errors *multierror.Error

	return res.Embedded.Role, errors.ErrorOrNil()
}

func (c *DefaultClient) BulkUpdateRoles(bulkOperations []*RoleBulkOperation) ([]*Role, error) {
	headers := map[string]string{
		"Accept":       "application/vnd.go.cd+json",
		"Content-Type": "application/json",
	}
	type EmbeddedObj struct {
		Role []*Role `json:"roles"`
	}
	type RolesResponse struct {
		Embedded EmbeddedObj `json:"_embedded"`
	}
	res := new(RolesResponse)

	type BulkRequest struct {
		Operations []*RoleBulkOperation `json:"operations"`
	}

	req := BulkRequest{Operations: bulkOperations}
	err := c.patchJSON("/go/api/admin/security/roles", headers, req, res)
	if err != nil {
		return []*Role{}, err
	}
	return res.Embedded.Role, nil
}
