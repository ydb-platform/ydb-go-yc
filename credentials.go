package yc

import (
	"github.com/ydb-platform/ydb-go-sdk/v3/credentials"

	"github.com/ydb-platform/ydb-go-yc/internal/auth"
)

func NewClient(opts ...ClientOption) (credentials.Credentials, error) {
	options := make([]auth.ClientOption, 0, len(opts))
	for _, option := range opts {
		options = append(options, auth.ClientOption(option))
	}
	return auth.NewClient(options...)
}
