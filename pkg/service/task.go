package service

import (
	"context"
	"github.com/6547709/goct/pkg/adapter"
)

type TaskService struct{ c adapter.TaskOps }
func NewTask(c adapter.TaskOps) *TaskService { return &TaskService{c: c} }

func (s *TaskService) List(ctx context.Context, opts adapter.ListOpts) ([]adapter.Task, error) {
	return s.c.ListTasks(ctx, opts)
}

func (s *TaskService) Get(ctx context.Context, id string) (*adapter.Task, error) {
	return s.c.GetTask(ctx, id)
}
