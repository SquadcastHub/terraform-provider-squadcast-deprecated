package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-provider-squadcast/internal/tfutils"
)

type Contact struct {
	DialCode    string `json:"dial_code" tf:"-"`
	PhoneNumber string `json:"phone_number" tf:"-"`
}

type Ability struct {
	ID string `json:"id" tf:"-"`
	// Name    string `json:"name" tf:"-"`
	Slug string `json:"slug" tf:"-"`
	// Default bool   `json:"default" tf:"-"`
}

type PersonalNotificationRule struct {
	Type         string `json:"type" tf:"type"`
	DelayMinutes int    `json:"time" tf:"delay_minutes"`
}

func (r *PersonalNotificationRule) Encode() (map[string]interface{}, error) {
	return tfutils.Encode(r)
}

type OncallReminderRule struct {
	Type         string `json:"type" tf:"type"`
	DelayMinutes int    `json:"time" tf:"delay_minutes"`
}

func (r *OncallReminderRule) Encode() (map[string]interface{}, error) {
	return tfutils.Encode(r)
}

type DataSourceUser struct {
	AbilitiesSlugs            []string                    `json:"-" tf:"abilities"`
	Name                      string                      `json:"-" tf:"name"`
	PhoneNumber               string                      `json:"-" tf:"phone"`
	ID                        string                      `json:"id" tf:"id"`
	Abilities                 []*Ability                  `json:"abilities" tf:"-"`
	Bio                       string                      `json:"bio" tf:"-"`
	Contact                   Contact                     `json:"contact" tf:"-"`
	Email                     string                      `json:"email" tf:"email"`
	FirstName                 string                      `json:"first_name" tf:"first_name"`
	IsEmailVerified           bool                        `json:"email_verified" tf:"is_email_verified"`
	IsInGracePeriod           bool                        `json:"in_grace_period" tf:"-"`
	IsOverrideDnDEnabled      bool                        `json:"is_override_dnd_enabled" tf:"is_override_dnd_enabled"`
	IsPhoneVerified           bool                        `json:"phone_verified" tf:"is_phone_verified"`
	IsTrialSignup             bool                        `json:"is_trial_signup" tf:"-"`
	LastName                  string                      `json:"last_name" tf:"last_name"`
	OncallReminderRules       []*OncallReminderRule       `json:"oncall_reminder_rules" tf:"-"`
	PersonalNotificationRules []*PersonalNotificationRule `json:"notification_rules" tf:"-"`
	Role                      string                      `json:"role" tf:"role"`
	TimeZone                  string                      `json:"time_zone" tf:"time_zone"`
	Title                     string                      `json:"title" tf:"-"`
}

type User struct {
	ID        string `json:"id" tf:"id"`
	Email     string `json:"email" tf:"email"`
	FirstName string `json:"first_name" tf:"first_name"`
	LastName  string `json:"last_name" tf:"last_name"`
	Role      string `json:"role" tf:"role"`
}

func (u *DataSourceUser) Encode() (map[string]interface{}, error) {
	u.Name = u.FirstName + " " + u.LastName

	if u.Contact.DialCode != "" && u.Contact.PhoneNumber != "" {
		u.PhoneNumber = u.Contact.DialCode + u.Contact.PhoneNumber
	}

	for _, v := range u.Abilities {
		u.AbilitiesSlugs = append(u.AbilitiesSlugs, v.Slug)
	}

	m, err := tfutils.Encode(u)
	if err != nil {
		return nil, err
	}

	rules, err := tfutils.EncodeSlice(u.OncallReminderRules)
	if err != nil {
		return nil, err
	}
	m["oncall_reminder_rules"] = rules

	rules, err = tfutils.EncodeSlice(u.PersonalNotificationRules)
	if err != nil {
		return nil, err
	}
	m["notification_rules"] = rules

	return m, nil
}

func (u *User) Encode() (map[string]interface{}, error) {
	m, err := tfutils.Encode(u)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (client *Client) GetUserById(ctx context.Context, id string) (*User, error) {
	url := fmt.Sprintf("%s/users/%s", client.BaseURLV3, id)

	return Request[any, User](http.MethodGet, url, client, ctx, nil)
}

func (client *Client) GetUserByEmail(ctx context.Context, email string) (*DataSourceUser, error) {
	url := fmt.Sprintf("%s/users?email=%s", client.BaseURLV3, url.QueryEscape(email))

	users, err := RequestSlice[any, DataSourceUser](http.MethodGet, url, client, ctx, nil)
	if err != nil {
		return nil, err
	}

	if users[0] == nil {
		return nil, fmt.Errorf("cannot find user with email `%s`", email)
	}

	return users[0], nil
}

func (client *Client) ListUsers(ctx context.Context) ([]*User, error) {
	url := fmt.Sprintf("%s/users", client.BaseURLV3)

	return RequestSlice[any, User](http.MethodGet, url, client, ctx, nil)
}

type CreateUserReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

type UpdateUserReq struct {
	Role string `json:"role"`
}

func (client *Client) CreateUser(ctx context.Context, req *CreateUserReq) (*User, error) {
	url := fmt.Sprintf("%s/users", client.BaseURLV3)
	return Request[CreateUserReq, User](http.MethodPost, url, client, ctx, req)
}

func (client *Client) UpdateUser(ctx context.Context, id string, req *UpdateUserReq) (*User, error) {
	url := fmt.Sprintf("%s/users/%s", client.BaseURLV3, id)
	return Request[UpdateUserReq, User](http.MethodPut, url, client, ctx, req)
}

func (client *Client) DeleteUser(ctx context.Context, id string) (*any, error) {
	url := fmt.Sprintf("%s/users/%s", client.BaseURLV3, id)
	return Request[any, any](http.MethodDelete, url, client, ctx, nil)
}
