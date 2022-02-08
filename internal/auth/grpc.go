package auth

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/yandex-cloud/go-genproto/yandex/cloud/iam/v1"
)

type grpcTransport struct {
	endpoint           string
	certPool           *x509.CertPool
	insecure           bool // Only for testing.
	insecureSkipVerify bool // Accept any TLS certificate from server.
}

func (t *grpcTransport) CreateToken(ctx context.Context, jwt string) (
	token string, expires time.Time, err error,
) {
	conn, err := t.conn(ctx)
	if err != nil {
		return
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
	if err == nil {
		token = res.IamToken
		expires = time.Unix(
			res.ExpiresAt.Seconds,
			int64(res.ExpiresAt.Nanos),
		)
	}
	return
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
	return grpc.DialContext(ctx, t.endpoint, opts...)
}
