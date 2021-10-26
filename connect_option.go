package yc

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"time"

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

func WithMetadataCredentialsURL(ctx context.Context, url string) ydb.Option {
	return ydb.WithCredentials(
		auth.InstanceServiceAccountURL(
			credentials.WithCredentialsSourceInfo(ctx, "yc.WithMetadataCredentialsURL(ctx, "+url+")"),
			url,
		),
	)
}

func WithServiceAccountKeyFileCredentials(serviceAccountKeyFile string, opts ...auth.ClientOption) ydb.Option {
	return ydb.WithCreateCredentialsFunc(func(ctx context.Context) (credentials.Credentials, error) {
		credentials, err := auth.NewClient(
			append(
				[]auth.ClientOption{
					auth.WithServiceFile(serviceAccountKeyFile),
					auth.WithDefaultEndpoint(),
					auth.WithSystemCertPool(),
					auth.WithSourceInfo("yc.WithServiceAccountKeyFileCredentials(\"" + serviceAccountKeyFile + "\")"),
				},
				opts...,
			)...,
		)
		if err != nil {
			return nil, fmt.Errorf("configure credentials error: %w", err)
		}
		return credentials, nil
	})
}

// WithInternalCA append internal yandex-cloud certs
func WithInternalCA() ydb.Option {
	return ydb.WithCertificatesFromPem(pem.YcPEM)
}

// WithEndpoint set provided endpoint.
func WithEndpoint(endpoint string) auth.ClientOption {
	return auth.WithEndpoint(endpoint)
}

// WithDefaultEndpoint set endpoint with default value.
func WithDefaultEndpoint() auth.ClientOption {
	return auth.WithDefaultEndpoint()
}

// WithSourceInfo set sourceInfo
func WithSourceInfo(sourceInfo string) auth.ClientOption {
	return auth.WithSourceInfo(sourceInfo)
}

// WithCertPool set provided certPool.
func WithCertPool(certPool *x509.CertPool) auth.ClientOption {
	return auth.WithCertPool(certPool)
}

// WithCertPoolFile try set root certPool from provided cert file path.
func WithCertPoolFile(caFile string) auth.ClientOption {
	return auth.WithCertPoolFile(caFile)
}

// WithSystemCertPool try set certPool with system root certificates.
func WithSystemCertPool() auth.ClientOption {
	return auth.WithSystemCertPool()
}

// WithInsecureSkipVerify set insecureSkipVerify to true which force client accepts any TLS certificate
// presented by the iam server and any host name in that certificate.
//
// If InsecureSkipVerify is set, then certPool field is not used.
//
// This should be used only for testing purposes.
func WithInsecureSkipVerify(insecure bool) auth.ClientOption {
	return auth.WithInsecureSkipVerify(insecure)
}

// WithKeyID set provided keyID.
func WithKeyID(keyID string) auth.ClientOption {
	return auth.WithKeyID(keyID)
}

// WithIssuer set provided issuer.
func WithIssuer(issuer string) auth.ClientOption {
	return auth.WithIssuer(issuer)
}

// WithTokenTTL set provided tokenTTL duration.
func WithTokenTTL(tokenTTL time.Duration) auth.ClientOption {
	return auth.WithTokenTTL(tokenTTL)
}

// WithAudience set provided audience.
func WithAudience(audience string) auth.ClientOption {
	return auth.WithAudience(audience)
}

// WithPrivateKey set provided private key.
func WithPrivateKey(key *rsa.PrivateKey) auth.ClientOption {
	return auth.WithPrivateKey(key)
}

// WithPrivateKeyFile try set key from provided private key file path
func WithPrivateKeyFile(path string) auth.ClientOption {
	return auth.WithPrivateKeyFile(path)
}

// WithServiceFile try set key, keyID, issuer from provided service account file path.
//
// Do not mix this option with WithKeyID, WithIssuer and key options (WithPrivateKey, WithPrivateKeyFile, etc).
func WithServiceFile(path string) auth.ClientOption {
	return auth.WithServiceFile(path)
}
