package iam

import (
	"context"
	"fmt"

	"github.com/YandexDatabase/ydb-go-sdk/v3"
	"github.com/YandexDatabase/ydb-go-sdk/v3/connect"
)

func WithMetadataCredentials(ctx context.Context) connect.ConnectOption {
	return connect.WithCredentials(
		InstanceServiceAccount(
			ydb.WithCredentialsSourceInfo(ctx, "connect.WithMetadataCredentials(ctx)"),
		),
	)
}

func WithServiceAccountKeyFileCredentials(serviceAccountKeyFile string) connect.ConnectOption {
	return connect.WithCreateCredentialsFunc(func() (ydb.Credentials, error) {
		credentials, err := NewClient(
			WithServiceFile(serviceAccountKeyFile),
			WithDefaultEndpoint(),
			WithSystemCertPool(),
			WithSourceInfo("connect.WithServiceAccountKeyFileCredentials(\""+serviceAccountKeyFile+"\")"),
		)
		if err != nil {
			return nil, fmt.Errorf("configure credentials error: %w", err)
		}
		return credentials, nil
	})
}
