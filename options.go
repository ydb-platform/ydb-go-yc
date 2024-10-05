package yc

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3"             //nolint:depguard
	"github.com/ydb-platform/ydb-go-sdk/v3/credentials" //nolint:depguard
	yc "github.com/ydb-platform/ydb-go-yc-metadata"     //nolint:depguard

	"github.com/ydb-platform/ydb-go-yc/internal/auth"
)

func WithMetadataCredentialsURL(url string) ydb.Option {
	return ydb.WithCredentials(
		NewInstanceServiceAccountURL(url),
	)
}

func WithMetadataCredentials(opts ...yc.InstanceServiceAccountCredentialsOption) ydb.Option {
	return ydb.WithCredentials(
		NewInstanceServiceAccount(opts...),
	)
}

func WithServiceAccountKeyFileCredentials(serviceAccountKeyFile string, opts ...ClientOption) ydb.Option {
	return WithAuthClientCredentials(
		append(
			[]ClientOption{auth.WithServiceFile(serviceAccountKeyFile)},
			opts...,
		)...,
	)
}

func WithServiceAccountKeyCredentials(serviceAccountKey string, opts ...ClientOption) ydb.Option {
	return WithAuthClientCredentials(
		append(
			[]ClientOption{auth.WithServiceKey(serviceAccountKey)},
			opts...,
		)...,
	)
}

func WithAuthClientCredentials(opts ...ClientOption) ydb.Option {
	return ydb.WithCreateCredentialsFunc(func(ctx context.Context) (credentials.Credentials, error) {
		c, err := auth.NewClient(opts...)
		if err != nil {
			return nil, fmt.Errorf("credentials configure error: %w", err)
		}
		return c, nil
	})
}

// WithInternalCA append internal yandex-cloud certs
func WithInternalCA() ydb.Option {
	return yc.WithInternalCA()
}

// WithFallbackCredentials makes fallback credentials if primary credentials are failed
func WithFallbackCredentials(fallback credentials.Credentials) ClientOption {
	return auth.WithFallbackCredentials(fallback)
}

// WithEndpoint set provided endpoint.
func WithEndpoint(endpoint string) ClientOption {
	return auth.WithEndpoint(endpoint)
}

// WithDefaultEndpoint set endpoint with default value.
func WithDefaultEndpoint() ClientOption {
	return auth.WithDefaultEndpoint()
}

// WithSourceInfo set sourceInfo
func WithSourceInfo(sourceInfo string) ClientOption {
	return auth.WithSourceInfo(sourceInfo)
}

// WithCertPool set provided certPool.
func WithCertPool(certPool *x509.CertPool) ClientOption {
	return auth.WithCertPool(certPool)
}

// WithCertPoolFile try set root certPool from provided cert file path.
func WithCertPoolFile(caFile string) ClientOption {
	return auth.WithCertPoolFile(caFile)
}

// WithSystemCertPool try set certPool with system root certificates.
func WithSystemCertPool() ClientOption {
	return auth.WithSystemCertPool()
}

// WithInsecureSkipVerify set insecureSkipVerify to true which force client accepts any TLS certificate
// presented by the iam server and any host name in that certificate.
//
// If insecureSkipVerify is set, then certPool field is not used.
//
// This should be used only for testing purposes.
func WithInsecureSkipVerify(insecure bool) ClientOption {
	return auth.WithInsecureSkipVerify(insecure)
}

// WithKeyID set provided keyID.
func WithKeyID(keyID string) ClientOption {
	return auth.WithKeyID(keyID)
}

// WithIssuer set provided issuer.
func WithIssuer(issuer string) ClientOption {
	return auth.WithIssuer(issuer)
}

// WithTokenTTL set provided tokenTTL duration.
func WithTokenTTL(tokenTTL time.Duration) ClientOption {
	return auth.WithTokenTTL(tokenTTL)
}

// WithAudience set provided audience.
func WithAudience(audience string) ClientOption {
	return auth.WithAudience(audience)
}

// WithPrivateKey set provided private key.
func WithPrivateKey(key *rsa.PrivateKey) ClientOption {
	return auth.WithPrivateKey(key)
}

// WithPrivateKeyFile try set key from provided private key file path
func WithPrivateKeyFile(path string) ClientOption {
	return auth.WithPrivateKeyFile(path)
}

// WithServiceFile try set key, keyID, issuer from provided service account file path.
//
// Do not mix this option with WithKeyID, WithIssuer and key options (WithPrivateKey, WithPrivateKeyFile, etc).
func WithServiceFile(path string) ClientOption {
	return auth.WithServiceFile(path)
}

// WithServiceKey try set key, keyID, issuer from provided service account key.
//
// Do not mix this option with WithKeyID, WithIssuer and key options (WithPrivateKey, WithPrivateKeyFile, etc).
func WithServiceKey(json string) ClientOption {
	return auth.WithServiceKey(json)
}
