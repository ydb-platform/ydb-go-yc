package yc

import (
	"github.com/ydb-platform/ydb-go-sdk/v3/credentials"

	"github.com/ydb-platform/ydb-go-yc/internal/auth"
)

func NewClient(opts ...ClientOption) (credentials.Credentials, error) {
	return auth.NewClient(opts...)
}
