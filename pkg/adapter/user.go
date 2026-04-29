package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	sdkuser "github.com/smartxworks/cloudtower-go-sdk/v2/client/user"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// UserOps 定义用户操作。
type UserOps interface {
	ListUsers(ctx context.Context, opts ListOpts) ([]User, error)
	GetUser(ctx context.Context, id string) (*User, error)
	CreateUser(ctx context.Context, spec UserCreateSpec) (TaskRef, error)
	DeleteUser(ctx context.Context, id string) (TaskRef, error)
}

func (c *defaultClient) ListUsers(ctx context.Context, opts ListOpts) ([]User, error) {
	params := sdkuser.NewGetUsersParams()
	params.SetContext(ctx)
	body := &models.GetUsersRequestBody{}
	where := &models.UserWhereInput{}
	hasWhere := false
	if opts.NameContains != "" {
		where.NameContains = pointy.String(opts.NameContains)
		hasWhere = true
	}
	if hasWhere { body.Where = where }
	if opts.Limit > 0 { body.First = pointy.Int32(opts.Limit) }
	params.SetRequestBody(body)
	resp, err := c.api.User.GetUsers(params)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	out := make([]User, 0, len(resp.Payload))
	for _, u := range resp.Payload {
		out = append(out, toUser(u))
	}
	return out, nil
}

func (c *defaultClient) GetUser(ctx context.Context, id string) (*User, error) {
	params := sdkuser.NewGetUsersParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.GetUsersRequestBody{
		Where: &models.UserWhereInput{ID: pointy.String(id)},
	})
	resp, err := c.api.User.GetUsers(params)
	if err != nil {
		return nil, fmt.Errorf("get user %s: %w", id, err)
	}
	if len(resp.Payload) == 0 {
		return nil, fmt.Errorf("get user %s: %w", id, ErrNotFound)
	}
	u := toUser(resp.Payload[0])
	return &u, nil
}

func (c *defaultClient) CreateUser(ctx context.Context, spec UserCreateSpec) (TaskRef, error) {
	p := &models.UserCreationParams{
		Name:     pointy.String(spec.Name),
		Username: pointy.String(spec.Username),
		Password: pointy.String(spec.Password),
		RoleID:   pointy.String(spec.RoleID),
	}
	if spec.Email != "" {
		p.EmailAddress = pointy.String(spec.Email)
	}
	params := sdkuser.NewCreateUserParams()
	params.SetContext(ctx)
	params.SetRequestBody([]*models.UserCreationParams{p})
	resp, err := c.api.User.CreateUser(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("create user %s: %w", spec.Username, err)
	}
	return firstUserTaskRef(resp.Payload), nil
}

func (c *defaultClient) DeleteUser(ctx context.Context, id string) (TaskRef, error) {
	params := sdkuser.NewDeleteUserParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.UserDeletionParams{
		Where: &models.UserWhereInput{ID: pointy.String(id)},
	})
	resp, err := c.api.User.DeleteUser(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("delete user %s: %w", id, err)
	}
	return firstDeleteUserTaskRef(resp.Payload), nil
}

func toUser(u *models.User) User {
	out := User{}
	if u.ID != nil { out.ID = *u.ID }
	if u.Name != nil { out.Name = *u.Name }
	if u.Username != nil { out.Username = *u.Username }
	if u.Source != nil { out.Source = string(*u.Source) }
	if u.Role != nil { out.Role = string(*u.Role) }
	if u.EmailAddress != nil { out.Email = *u.EmailAddress }
	return out
}

func firstUserTaskRef(items []*models.WithTaskUser) TaskRef {
	if len(items) == 0 { return TaskRef{} }
	ref := TaskRef{EntityKind: "User"}
	if items[0].TaskID != nil { ref.ID = *items[0].TaskID }
	if items[0].Data != nil && items[0].Data.ID != nil { ref.EntityID = *items[0].Data.ID }
	return ref
}

func firstDeleteUserTaskRef(items []*models.WithTaskDeleteUser) TaskRef {
	if len(items) == 0 { return TaskRef{} }
	ref := TaskRef{EntityKind: "User"}
	if items[0].TaskID != nil { ref.ID = *items[0].TaskID }
	if items[0].Data != nil && items[0].Data.ID != nil { ref.EntityID = *items[0].Data.ID }
	return ref
}
