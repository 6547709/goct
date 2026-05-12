package service

import (
	"context"

	"github.com/6547709/goct/pkg/adapter"
)

type NTPService struct{ c adapter.NtpOps }
func NewNTP(c adapter.NtpOps) *NTPService { return &NTPService{c: c} }

func (s *NTPService) GetSettings(ctx context.Context) (*adapter.NtpSettings, error) {
	return s.c.GetNtpServiceURL(ctx)
}