package iam

import (
	"context"
	"fmt"
	"github.com/ydb-platform/ydb-go-sdk/v3/credentials"

	"github.com/ydb-platform/ydb-go-sdk/v3"
)

func WithMetadataCredentials(ctx context.Context) ydb.Option {
	return ydb.WithCredentials(
		InstanceServiceAccount(
			credentials.WithCredentialsSourceInfo(ctx, "connect.WithMetadataCredentials(ctx)"),
		),
	)
}

func WithServiceAccountKeyFileCredentials(serviceAccountKeyFile string) ydb.Option {
	return ydb.WithCreateCredentialsFunc(func(ctx context.Context) (credentials.Credentials, error) {
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
