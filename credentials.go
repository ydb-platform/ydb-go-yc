package yc

import (
	yc "github.com/ydb-platform/ydb-go-yc-metadata" //nolint:depguard

	"github.com/ydb-platform/ydb-go-yc/internal/auth"
)

type ClientOption = auth.ClientOption

func mustNewClient(opts ...ClientOption) *auth.Client {
	client, err := auth.NewClient(opts...)
	if err != nil {
		panic(err)
	}
	return client
}

func NewClient(opts ...ClientOption) (*auth.Client, error) {
	return auth.NewClient(opts...)
}

func NewInstanceServiceAccount(
	opts ...yc.InstanceServiceAccountCredentialsOption,
) *yc.InstanceServiceAccountCredentials {
	return yc.NewInstanceServiceAccount(opts...)
}

func NewInstanceServiceAccountURL(url string) *yc.InstanceServiceAccountCredentials {
	return yc.NewInstanceServiceAccount(yc.WithURL(url))
}
