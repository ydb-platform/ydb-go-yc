package yc_test

import (
	"context"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	yc "github.com/ydb-platform/ydb-go-yc"
)

func ExampleNebius() {
	db, err := ydb.Open(context.TODO(), "grpcs://ydb.serverless.il.nebius.cloud:2135/il1/yc.mds-team.service-cloud/a3duh738rh48homerna57",
		ydb.WithCertificatesFromFile("./ydb/CA.pem"),
		yc.WithServiceAccountKeyFileCredentials("./nebius/sa.json",
			yc.WithEndpoint("iam.api.il.nebius.cloud:443"),
		),
	)
	if err != nil {
		panic(err)
	}
	_ = db.Close(context.TODO())
}

func ExampleYandexCloud() {
	db, err := ydb.Open(context.TODO(), "grpcs://ydb.serverless.yandexcloud.net:2135/?database=/ru-central1/yc.serverless.cloud/etne6yegf7346f0h71tpev2p",
		yc.WithInternalCA(),
		yc.WithServiceAccountKeyFileCredentials("./yc/sa.json"),
	)
	if err != nil {
		panic(err)
	}
	_ = db.Close(context.TODO())
}
