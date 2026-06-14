package auth

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"time"

	v1 "github.com/yandex-cloud/go-genproto/yandex/cloud/iam/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcTransport struct {
	endpoint           string
	certPool           *x509.CertPool
	insecure           bool // Only for testing.
	insecureSkipVerify bool // Accept any TLS certificate from server.
}

func (t *grpcTransport) CreateToken(ctx context.Context, jwt string) (string, time.Time, error) {
	conn, err := t.conn(ctx)
	if err != nil {
		return "", time.Time{}, err
	}
	defer func() {
		_ = conn.Close()
	}()

	client := v1.NewIamTokenServiceClient(conn)
	res, err := client.Create(ctx, &v1.CreateIamTokenRequest{
		Identity: &v1.CreateIamTokenRequest_Jwt{
			Jwt: jwt,
		},
	})
	if err != nil {
		return "", time.Time{}, err
	}

	expiresAt := res.GetExpiresAt()

	return res.GetIamToken(), time.Unix(
		expiresAt.GetSeconds(),
		int64(expiresAt.GetNanos()),
	), nil
}

func (t *grpcTransport) conn(ctx context.Context) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	switch {
	case t.insecure:
		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
	case t.insecureSkipVerify:
		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(
				credentials.NewTLS(&tls.Config{
					//nolint: gosec
					InsecureSkipVerify: true,
				}),
			),
		}
	case t.certPool != nil:
		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(
				credentials.NewClientTLSFromCert(t.certPool, ""),
			),
		}
	}

	//nolint:staticcheck,nolintlint // grpc.NewClient migration is out of scope for this change.
	return grpc.DialContext(ctx, t.endpoint, opts...)
}
