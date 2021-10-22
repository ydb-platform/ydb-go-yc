package yc

import (
	"context"
	"fmt"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/credentials"

	"github.com/ydb-platform/ydb-go-yc/internal/auth"
	"github.com/ydb-platform/ydb-go-yc/internal/pem"
)

func WithMetadataCredentials(ctx context.Context) ydb.Option {
	return ydb.WithCredentials(
		auth.InstanceServiceAccount(
			credentials.WithCredentialsSourceInfo(ctx, "yc.WithMetadataCredentials(ctx)"),
		),
	)
}

func WithServiceAccountKeyFileCredentials(serviceAccountKeyFile string) ydb.Option {
	return ydb.WithCreateCredentialsFunc(func(ctx context.Context) (credentials.Credentials, error) {
		credentials, err := auth.NewClient(
			auth.WithServiceFile(serviceAccountKeyFile),
			auth.WithDefaultEndpoint(),
			auth.WithSystemCertPool(),
			auth.WithSourceInfo("yc.WithServiceAccountKeyFileCredentials(\""+serviceAccountKeyFile+"\")"),
		)
		if err != nil {
			return nil, fmt.Errorf("configure credentials error: %w", err)
		}
		return credentials, nil
	})
}

func WithInternalCA() ydb.Option {
	return ydb.WithCertificatesFromPem(pem.YcPEM)
}
