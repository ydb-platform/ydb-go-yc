package yc

import (
	"context"
	"crypto/x509"
	"fmt"
	"github.com/ydb-platform/ydb-go-sdk/v3/ydbsql"
	auth2 "github.com/ydb-platform/ydb-go-yc/internal/auth"
	"github.com/ydb-platform/ydb-go-yc/internal/pem"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/credentials"
)

func WithMetadataCredentials(ctx context.Context) ydb.Option {
	return ydb.WithCredentials(
		auth2.InstanceServiceAccount(
			credentials.WithCredentialsSourceInfo(ctx, "yc.WithMetadataCredentials(ctx)"),
		),
	)
}

func WithServiceAccountKeyFileCredentials(serviceAccountKeyFile string) ydb.Option {
	return ydb.WithCreateCredentialsFunc(func(ctx context.Context) (credentials.Credentials, error) {
		credentials, err := auth2.NewClient(
			auth2.WithServiceFile(serviceAccountKeyFile),
			auth2.WithDefaultEndpoint(),
			auth2.WithSystemCertPool(),
			auth2.WithSourceInfo("yc.WithServiceAccountKeyFileCredentials(\""+serviceAccountKeyFile+"\")"),
		)
		if err != nil {
			return nil, fmt.Errorf("configure credentials error: %w", err)
		}
		return credentials, nil
	})
}

func WithInternalCA() ydb.Option {
	certPool, err := x509.SystemCertPool()
	if err != nil {
		panic(err)
	}
	if !certPool.AppendCertsFromPEM(pem.YcPEM) {
		panic("cannot append yandex-cloud PEM")
	}
	return ydb.WithCertificates(certPool)
}

func WithYdbSqlInternalCA() ydbsql.ConnectorOption {
	return ydbsql.WithCertificatesFromPem(pem.YcPEM)
}
