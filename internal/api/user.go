package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

type User struct {
	ID          string `json:"id" tf:"id"`
	Email       string `json:"email" tf:"email"`
	FirstName   string `json:"first_name" tf:"first_name"`
	LastName    string `json:"last_name" tf:"last_name"`
	OrgRoleID   string `json:"role_id" tf:"org_role_id"`
	OrgRoleName string `json:"role_name" tf:"org_role_name"`
}

func (s *User) Encode() (map[string]interface{}, error) {
	m, err := tfutils.Encode(s)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (client *Client) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	path := fmt.Sprintf("/users?email=%s", email)
	return Request[any, User](http.MethodGet, path, client, ctx, nil)
}
