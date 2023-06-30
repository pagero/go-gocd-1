package gocd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllRoles(t *testing.T) {
	t.Parallel()
	client, server := newTestAPIClient("/go/api/admin/security/roles", serveFileAsJSON(t, "GET", "test-fixtures/get_all_roles.json", 0, DummyRequestBodyValidator))
	defer server.Close()
	roles, err := client.GetAllRoles()
	assert.NoError(t, err)
	assert.NotNil(t, roles)
	assert.Equal(t, 1, len(roles))

	role := roles[0]
	assert.NotNil(t, role)
	assert.Equal(t, "myrole", role.Name)

	assert.Equal(t, 1, len(role.Attributes.Users))
	u := role.Attributes.Users[0]
	assert.Equal(t, "auser", u)

	assert.Equal(t, 1, len(role.Policy))
	assert.Equal(t, role.Policy[0].Permission, "allow")
	assert.Equal(t, role.Policy[0].Action, "view")
	assert.Equal(t, role.Policy[0].Type, "environment")
	assert.Equal(t, role.Policy[0].Resource, "env1")
}

func TestBulkUpdateRole(t *testing.T) {
	t.Parallel()
	requestBodyValidator := func(body string) error {
		expectedBody := "{\"operations\":[{\"role\":\"role1\",\"users\":{\"add\":[\"user1\"],\"remove\":[\"user2\"]}}]}"
		if body != expectedBody {
			return fmt.Errorf("Request body (%s) didn't match the expected body (%s)", body, expectedBody)
		}
		return nil
	}

	client, server := newTestAPIClient("/go/api/admin/security/roles", serveFileAsJSON(t, "PATCH", "test-fixtures/patch_bulk_update_roles.json", 0, requestBodyValidator))
	defer server.Close()

	operations := []*RoleBulkOperation{
		{
			Role: "role1",
			Users: RoleBulkUsers{
				Add:    []string{"user1"},
				Remove: []string{"user2"},
			},
		},
	}
	updatedRoles, err := client.BulkUpdateRoles(operations)
	assert.NoError(t, err)
	assert.NotNil(t, updatedRoles)
	assert.Equal(t, 1, len(updatedRoles))
}
