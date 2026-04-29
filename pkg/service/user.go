package service

import (
	"context"
	"github.com/6547709/goct/pkg/adapter"
)

type UserService struct{ c adapter.UserOps }
func NewUser(c adapter.UserOps) *UserService { return &UserService{c: c} }

func (s *UserService) List(ctx context.Context, opts adapter.ListOpts) ([]adapter.User, error) {
	return s.c.ListUsers(ctx, opts)
}

func (s *UserService) Resolve(ctx context.Context, idOrName string) (*adapter.User, error) {
	return Resolve(ctx, s.c.ListUsers, s.c.GetUser,
		func(u adapter.User) (string, string) { return u.ID, u.Name },
		idOrName)
}

func (s *UserService) Create(ctx context.Context, spec adapter.UserCreateSpec) (adapter.TaskRef, error) {
	return s.c.CreateUser(ctx, spec)
}

func (s *UserService) Delete(ctx context.Context, idOrName string) (adapter.TaskRef, error) {
	u, err := s.Resolve(ctx, idOrName)
	if err != nil { return adapter.TaskRef{}, err }
	return s.c.DeleteUser(ctx, u.ID)
}
