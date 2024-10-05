package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"       //nolint:depguard
	"github.com/jonboulle/clockwork"     //nolint:depguard
	"github.com/stretchr/testify/assert" //nolint:depguard
	"github.com/stretchr/testify/require"
)

type TransportFunc func(context.Context, string) (string, time.Time, error)

func (f TransportFunc) CreateToken(ctx context.Context, jwt string) (string, time.Time, error) {
	return f(ctx, jwt)
}

func TestClientToken(t *testing.T) {
	const (
		keyID    = "key-id"
		issuer   = "issuer"
		audience = "audience"
		endpoint = "endpoint"

		ttl = 12 * time.Hour
	)
	fakeTime := clockwork.NewFakeClock()

	key, err := rsa.GenerateKey(rand.Reader, 4096)
	assert.NoError(t, err)

	prevTimeFunc := jwt.TimeFunc
	jwt.TimeFunc = fakeTime.Now
	defer func() {
		jwt.TimeFunc = prevTimeFunc
	}()

	var (
		i int

		results = [...]struct {
			token   string
			expires time.Duration
		}{
			{"foo", ttl},
			{"bar", time.Second},
			{"baz", 0},
		}
	)
	c := Client{
		clock:    fakeTime,
		endpoint: endpoint,
		key:      key,
		keyID:    keyID,
		issuer:   issuer,

		audience: audience,
		tokenTTL: ttl,

		// Stub the real transport logic to check jwt token for correctness.
		transport: TransportFunc(func(ctx context.Context, jwtString string) (
			string, time.Time, error,
		) {
			var claims jwt.RegisteredClaims
			keyFunc := func(t *jwt.Token) (interface{}, error) {
				// Use the public part of our key as IAM service will.
				return key.Public(), nil
			}
			token, err := jwt.ParseWithClaims(jwtString, &claims, keyFunc)
			assert.NoError(t, err, "parse token error")
			assert.Equal(t, keyID, token.Header["kid"], `unexpected "kid" header`)

			// Get the "now" moment. Note that this is the same as for sourceInfo –
			// we stubbed time above.
			now := fakeTime.Now()

			iat := jwt.NewNumericDate(now.Local())
			exp := jwt.NewNumericDate(now.Add(ttl).Local())

			assert.Equal(t, issuer, claims.Issuer, "unexpected claims.issuer field")
			assert.Contains(t, claims.Audience, audience, "unexpected claims.audience field")
			assert.Equal(t, iat, claims.IssuedAt, "unexpected claims.IssuedAt field")
			assert.Equal(t, exp, claims.ExpiresAt, "unexpected claims.ExpiresAt field")

			tokenString := results[i].token
			e := results[i].expires
			i++

			return tokenString, now.Add(e), nil
		}),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var attempt int
	getToken := func(expResult int) {
		t1, err := c.Token(ctx)
		require.NoError(t, err)
		assert.Equal(t, results[expResult].token, t1, " %d Token(): unexpected token", attempt)
		attempt++
	}

	getToken(0)

	fakeTime.Advance(time.Second)
	getToken(0)

	fakeTime.Advance(ttl) // time.Minute
	getToken(1)

	// Now server respond with time.Second expiration time.
	// Thus, we expect Token() request server again after second, not after
	// ttl (which is time.Minute).
	fakeTime.Advance(time.Second)
	getToken(2)
}

func TestOptionsConfig(t *testing.T) {
	const (
		keyID    = "key-id"
		issuer   = "issuer"
		audience = "audience"
		endpoint = "endpoint"

		ttl = time.Minute
	)
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	assert.NoError(t, err)

	client, err := NewClient(
		WithKeyID(keyID),
		WithIssuer(issuer),
		WithAudience(audience),
		WithEndpoint(endpoint),
		WithTokenTTL(ttl),
		WithPrivateKey(key),
	)

	assert.NoError(t, err)
	assert.Equal(t, keyID, client.keyID)
	assert.Equal(t, issuer, client.issuer)
	assert.Equal(t, audience, client.audience)
	assert.Equal(t, endpoint, client.endpoint)
	assert.Equal(t, ttl, client.tokenTTL)
}
