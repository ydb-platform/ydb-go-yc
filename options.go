package yc

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/credentials"

	yc "github.com/ydb-platform/ydb-go-yc-metadata"
	"github.com/ydb-platform/ydb-go-yc/internal/auth"
)

type ClientOption auth.ClientOption

func NewInstanceServiceAccount() *yc.InstanceServiceAccountCredentials {
	return yc.NewInstanceServiceAccount()
}

func NewInstanceServiceAccountURL(url string) *yc.InstanceServiceAccountCredentials {
	return yc.NewInstanceServiceAccount(yc.WithURL(url))
}

func WithMetadataCredentialsURL(url string) ydb.Option {
	return ydb.WithCredentials(
		NewInstanceServiceAccountURL(url),
	)
}

func WithMetadataCredentials() ydb.Option {
	return ydb.WithCredentials(
		NewInstanceServiceAccount(),
	)
}

func WithServiceAccountKeyFileCredentials(serviceAccountKeyFile string, opts ...ClientOption) ydb.Option {
	return WithAuthClientCredentials(
		append(
			[]ClientOption{ClientOption(auth.WithServiceFile(serviceAccountKeyFile))},
			opts...,
		)...,
	)
}

func WithServiceAccountKeyCredentials(serviceAccountKey string, opts ...ClientOption) ydb.Option {
	return WithAuthClientCredentials(
		append(
			[]ClientOption{ClientOption(auth.WithServiceKey(serviceAccountKey))},
			opts...,
		)...,
	)
}

func WithAuthClientCredentials(opts ...ClientOption) ydb.Option {
	return ydb.WithCreateCredentialsFunc(func(ctx context.Context) (credentials.Credentials, error) {
		options := make([]auth.ClientOption, 0, len(opts))
		for _, option := range opts {
			options = append(options, auth.ClientOption(option))
		}
		c, err := auth.NewClient(options...)
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
	return ClientOption(auth.WithFallbackCredentials(fallback))
}

// WithEndpoint set provided endpoint.
func WithEndpoint(endpoint string) ClientOption {
	return ClientOption(auth.WithEndpoint(endpoint))
}

// WithDefaultEndpoint set endpoint with default value.
func WithDefaultEndpoint() ClientOption {
	return ClientOption(auth.WithDefaultEndpoint())
}

// WithSourceInfo set sourceInfo
func WithSourceInfo(sourceInfo string) ClientOption {
	return ClientOption(auth.WithSourceInfo(sourceInfo))
}

// WithCertPool set provided certPool.
func WithCertPool(certPool *x509.CertPool) ClientOption {
	return ClientOption(auth.WithCertPool(certPool))
}

// WithCertPoolFile try set root certPool from provided cert file path.
func WithCertPoolFile(caFile string) ClientOption {
	return ClientOption(auth.WithCertPoolFile(caFile))
}

// WithSystemCertPool try set certPool with system root certificates.
func WithSystemCertPool() ClientOption {
	return ClientOption(auth.WithSystemCertPool())
}

// WithInsecureSkipVerify set insecureSkipVerify to true which force client accepts any TLS certificate
// presented by the iam server and any host name in that certificate.
//
// If insecureSkipVerify is set, then certPool field is not used.
//
// This should be used only for testing purposes.
func WithInsecureSkipVerify(insecure bool) ClientOption {
	return ClientOption(auth.WithInsecureSkipVerify(insecure))
}

// WithKeyID set provided keyID.
func WithKeyID(keyID string) ClientOption {
	return ClientOption(auth.WithKeyID(keyID))
}

// WithIssuer set provided issuer.
func WithIssuer(issuer string) ClientOption {
	return ClientOption(auth.WithIssuer(issuer))
}

// WithTokenTTL set provided tokenTTL duration.
func WithTokenTTL(tokenTTL time.Duration) ClientOption {
	return ClientOption(auth.WithTokenTTL(tokenTTL))
}

// WithAudience set provided audience.
func WithAudience(audience string) ClientOption {
	return ClientOption(auth.WithAudience(audience))
}

// WithPrivateKey set provided private key.
func WithPrivateKey(key *rsa.PrivateKey) ClientOption {
	return ClientOption(auth.WithPrivateKey(key))
}

// WithPrivateKeyFile try set key from provided private key file path
func WithPrivateKeyFile(path string) ClientOption {
	return ClientOption(auth.WithPrivateKeyFile(path))
}

// WithServiceFile try set key, keyID, issuer from provided service account file path.
//
// Do not mix this option with WithKeyID, WithIssuer and key options (WithPrivateKey, WithPrivateKeyFile, etc).
func WithServiceFile(path string) ClientOption {
	return ClientOption(auth.WithServiceFile(path))
}
