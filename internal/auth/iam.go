// Package auth provides interface for retrieving and caching iam tokens.
package auth

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jonboulle/clockwork"
	"github.com/ydb-platform/ydb-go-sdk/v3/credentials"
)

// Default client parameters.
const (
	DefaultAudience = "https://iam.api.cloud.yandex.net/iam/v1/tokens"
	DefaultEndpoint = "iam.api.cloud.yandex.net:443"
	DefaultTokenTTL = time.Hour
)

var (
	ErrServiceFileInvalid = errors.New("service account file is not valid")
	ErrKeyCannotBeParsed  = errors.New("private key can not be parsed")
)

// createTokenError contains reason of token creation failure.
type createTokenError struct {
	cause  error
	reason string
}

// Error implements error interface.
func (e *createTokenError) Error() string {
	return fmt.Sprintf("iam: create token error: %s", e.reason)
}

func (e *createTokenError) Unwrap() error {
	return e.cause
}

type transport interface {
	CreateToken(ctx context.Context, jwt string) (token string, expires time.Time, err error)
}

type ClientOption func(*client) error

// WithFallbackCredentials makes fallback credentials if primary credentials are failed
func WithFallbackCredentials(fallback credentials.Credentials) ClientOption {
	return func(c *client) error {
		c.fallback = fallback
		return nil
	}
}

// WithEndpoint set provided endpoint.
func WithEndpoint(endpoint string) ClientOption {
	return func(c *client) error {
		c.endpoint = endpoint
		return nil
	}
}

// WithDefaultEndpoint set endpoint with default value.
func WithDefaultEndpoint() ClientOption {
	return func(c *client) error {
		c.endpoint = DefaultEndpoint
		return nil
	}
}

// WithSourceInfo set sourceInfo
func WithSourceInfo(sourceInfo string) ClientOption {
	return func(c *client) error {
		c.sourceInfo = sourceInfo
		return nil
	}
}

// WithCertPool set provided certPool.
func WithCertPool(certPool *x509.CertPool) ClientOption {
	return func(c *client) error {
		c.certPool = certPool
		return nil
	}
}

// WithCertPoolFile try set root certPool from provided cert file path.
func WithCertPoolFile(caFile string) ClientOption {
	return func(c *client) error {
		if len(caFile) > 0 && caFile[0] == '~' {
			usr, err := user.Current()
			if err != nil {
				return err
			}
			caFile = filepath.Join(usr.HomeDir, caFile[1:])
		}
		bytes, err := os.ReadFile(caFile)
		if err != nil {
			return err
		}
		if !c.certPool.AppendCertsFromPEM(bytes) {
			return fmt.Errorf("cannot append certificates from file '%s' to certificates pool", caFile)
		}
		return nil
	}
}

// WithSystemCertPool try set certPool with system root certificates.
func WithSystemCertPool() ClientOption {
	return func(c *client) error {
		var err error
		c.certPool, err = x509.SystemCertPool()
		return err
	}
}

// WithInsecureSkipVerify set insecureSkipVerify to true which force client accepts any TLS certificate
// presented by the iam server and any host name in that certificate.
//
// If insecureSkipVerify is set, then certPool field is not used.
//
// This should be used only for testing purposes.
func WithInsecureSkipVerify(insecure bool) ClientOption {
	return func(c *client) error {
		c.insecureSkipVerify = insecure
		return nil
	}
}

// WithKeyID set provided keyID.
func WithKeyID(keyID string) ClientOption {
	return func(c *client) error {
		c.keyID = keyID
		return nil
	}
}

// WithIssuer set provided issuer.
func WithIssuer(issuer string) ClientOption {
	return func(c *client) error {
		c.issuer = issuer
		return nil
	}
}

// WithTokenTTL set provided tokenTTL duration.
func WithTokenTTL(tokenTTL time.Duration) ClientOption {
	return func(c *client) error {
		c.tokenTTL = tokenTTL
		return nil
	}
}

// WithAudience set provided audience.
func WithAudience(audience string) ClientOption {
	return func(c *client) error {
		c.audience = audience
		return nil
	}
}

// WithPrivateKey set provided private key.
func WithPrivateKey(key *rsa.PrivateKey) ClientOption {
	return func(c *client) error {
		c.key = key
		return nil
	}
}

// WithPrivateKeyFile try set key from provided private key file path
func WithPrivateKeyFile(path string) ClientOption {
	return func(c *client) error {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		key, err := parsePrivateKey(data)
		if err != nil {
			return err
		}
		c.key = key
		return nil
	}
}

// WithServiceFile try set key, keyID, issuer from provided service account file path.
//
// Do not mix this option with WithKeyID, WithIssuer and key options (WithPrivateKey, WithPrivateKeyFile, etc).
func WithServiceFile(path string) ClientOption {
	return func(c *client) error {
		if len(path) > 0 && path[0] == '~' {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			path = filepath.Join(home, path[1:])
		}
		data, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			return err
		}

		return parseAndApplyServiceAccountKeyData(c, data)
	}
}

// WithServiceKey try set key, keyID, issuer from provided service account data key.
//
// Do not mix this option with WithKeyID, WithIssuer and key options (WithPrivateKey, WithPrivateKeyFile, etc).
func WithServiceKey(key string) ClientOption {
	return func(c *client) error { return parseAndApplyServiceAccountKeyData(c, []byte(key)) }
}

// parseAndApplyServiceAccountKeyData set key, keyID, issuer from provided service account data key,
// or form service account file path.
//
//	Do not mix this option with WithKeyID, WithIssuer and key options (WithPrivateKey, WithPrivateKeyFile, etc).
func parseAndApplyServiceAccountKeyData(c *client, data []byte) error {
	type keyFile struct {
		ID               string `json:"id"`
		ServiceAccountID string `json:"service_account_id"`
		PrivateKey       string `json:"private_key"`
		Endpoint         string `json:"endpoint,omitempty"`
	}
	var info keyFile
	if err := json.Unmarshal(data, &info); err != nil {
		return err
	}
	if info.ID == "" || info.ServiceAccountID == "" || info.PrivateKey == "" {
		return ErrServiceFileInvalid
	}

	key, err := parsePrivateKey([]byte(info.PrivateKey))
	if err != nil {
		return err
	}
	c.key = key
	c.keyID = info.ID
	c.issuer = info.ServiceAccountID
	if info.Endpoint != "" {
		c.endpoint = info.Endpoint
	}

	return nil
}

// NewClient creates IAM (jwt) authorized client from provided ClientOptions list.
//
// To create successfully at least one of endpoint options must be provided.
func NewClient(opts ...ClientOption) (_ credentials.Credentials, err error) {
	var (
		certPool *x509.CertPool
		issues   []error
	)
	certPool, err = x509.SystemCertPool()
	if err != nil {
		certPool = x509.NewCertPool()
	}

	c := &client{
		endpoint:           DefaultEndpoint,
		certPool:           certPool,
		insecureSkipVerify: true,
		tokenTTL:           DefaultTokenTTL,
		audience:           DefaultAudience,
		clock:              clockwork.NewRealClock(),
	}

	for _, opt := range opts {
		err = opt(c)
		if err != nil {
			issues = append(issues, err)
		}
	}

	if len(issues) > 0 {
		if c.fallback != nil {
			return c.fallback, nil
		}
		return nil, fmt.Errorf("cannot create IAM client: %v", issues)
	}

	c.transport = &grpcTransport{
		endpoint:           c.endpoint,
		certPool:           c.certPool,
		insecureSkipVerify: c.insecureSkipVerify,
	}

	return c, nil
}

// Client contains options for interaction with the iam.
type client struct {
	endpoint string
	certPool *x509.CertPool

	// If insecureSkipVerify is true, client accepts any TLS certificate
	// presented by the iam server and any host name in that certificate.
	//
	// If insecureSkipVerify is set, then certPool field is not used.
	//
	// This should be used only for testing.
	insecureSkipVerify bool

	key    *rsa.PrivateKey
	keyID  string
	issuer string

	tokenTTL time.Duration
	audience string

	once    sync.Once
	mu      sync.RWMutex
	err     error
	token   string
	expires time.Time

	// transport is a stub used for tests.
	transport transport

	sourceInfo string

	fallback credentials.Credentials

	clock clockwork.Clock
}

func (c *client) init() (err error) {
	c.once.Do(func() {
		if c.endpoint == "" {
			c.err = fmt.Errorf("iam: endpoint required")
			return
		}
		if c.audience == "" {
			c.audience = DefaultAudience
		}
		if c.tokenTTL == 0 {
			c.tokenTTL = DefaultTokenTTL
		}
		if c.transport == nil {
			c.transport = &grpcTransport{
				endpoint:           c.endpoint,
				certPool:           c.certPool,
				insecureSkipVerify: c.insecureSkipVerify,
			}
		}
	})
	return c.err
}

func (c *client) String() string {
	if c.sourceInfo == "" {
		return "iam.Client"
	}
	return "iam.Client created from " + c.sourceInfo
}

// Token returns cached token if no c.tokenTTL time has passed or no token
// expiration deadline from the last request exceeded. In other way, it makes
// request for a new one token.
func (c *client) Token(ctx context.Context) (token string, err error) {
	if err = c.init(); err != nil {
		return
	}
	c.mu.RLock()
	if !c.expired() {
		token = c.token
	}
	c.mu.RUnlock()
	if token != "" {
		return token, nil
	}
	now := c.clock.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.expired() {
		return c.token, nil
	}
	var expires time.Time
	jwtToken, err := c.jwt(now)
	if err != nil {
		return c.token, err
	}
	token, expires, err = c.transport.CreateToken(ctx, jwtToken)
	if err != nil {
		return "", &createTokenError{
			cause:  err,
			reason: err.Error(),
		}
	}
	c.token = token
	c.expires = now.Add(expires.Sub(now) / 2)
	return token, nil
}

func (c *client) expired() bool {
	return c.clock.Since(c.expires) > 0
}

// By default, Go RSA PSS uses PSSSaltLengthAuto, but RFC states that salt size
// must be equal to hash size.
//
// See https://tools.ietf.org/html/rfc7518#section-3.5
var ps256WithSaltLengthEqualsHash = &jwt.SigningMethodRSAPSS{
	SigningMethodRSA: jwt.SigningMethodPS256.SigningMethodRSA,
	Options: &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
	},
}

func (c *client) jwt(now time.Time) (string, error) {
	var (
		issued = jwt.NewNumericDate(now.UTC())
		expire = jwt.NewNumericDate(now.Add(c.tokenTTL).UTC())
		method = ps256WithSaltLengthEqualsHash
	)
	t := jwt.Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": method.Alg(),
			"kid": c.keyID,
		},
		Claims: jwt.RegisteredClaims{
			Issuer:    c.issuer,
			IssuedAt:  issued,
			Audience:  []string{c.audience},
			ExpiresAt: expire,
		},
		Method: method,
	}
	s, err := t.SignedString(c.key)
	if err != nil {
		return "", fmt.Errorf("iam: could not sign jwt token: %w", err)
	}
	return s, nil
}

func parsePrivateKey(raw []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, ErrKeyCannotBeParsed
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return key, err
	}

	x, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	if key, ok := x.(*rsa.PrivateKey); ok {
		return key, nil
	}
	return nil, ErrKeyCannotBeParsed
}
