module github.com/ydb-platform/ydb-go-yc

go 1.16

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/yandex-cloud/go-genproto v0.0.0-20211012081957-400ccab0fe15
	github.com/ydb-platform/ydb-go-sdk/v3 v3.0.1-release
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
)

retract v0.0.1
retract v0.0.2