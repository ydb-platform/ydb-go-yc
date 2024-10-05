package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ydb-platform/ydb-go-sdk/v3" //nolint:depguard

	yc "github.com/ydb-platform/ydb-go-yc" //nolint:depguard
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := ydb.Open(
		ctx,
		"grpcs://ydb.serverless.yandexcloud.net:2135/ru-central1/b1g8skpblkos03malf3s/etnaeujopcre7mubi9lj",

		// credentials to access YDB outside yandex-cloud
		yc.WithServiceAccountKeyFileCredentials(os.Getenv("YDB_SERVICE_ACCOUNT_KEY_FILE_CREDENTIALS")),
		// credentials to access YDB inside yandex-cloud (yandex function, yandex cloud virtual machine)
		// yc.WithMetadataCredentials(ctx),

		// certificates for access to yandex-cloud
		yc.WithInternalCA(),
		// or append certificates from file directly
		// ydb.WithCertificatesFromFile(os.Getenv("YDB_SSL_ROOT_CERTIFICATES_FILE")),
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = db.Close(ctx)
	}()
	whoAmI, err := db.Discovery().WhoAmI(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(whoAmI.String())
}
