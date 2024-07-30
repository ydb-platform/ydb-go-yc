package yc_test

import (
	"context"

	ydb "github.com/ydb-platform/ydb-go-sdk/v3"
	yc "github.com/ydb-platform/ydb-go-yc"
)

func Example_withMetadataCredentials() {
	db, err := ydb.Open(context.TODO(), "grpc://localhost:2136/local",
		yc.WithMetadataCredentials(),
		yc.WithInternalCA(),
	)
	if err != nil {
		panic(err)
	}
	_ = db.Close(context.TODO())
}

func Example_withServiceAccountKeyFileCredentials() {
	db, err := ydb.Open(context.TODO(), "grpc://localhost:2136/local",
		yc.WithServiceAccountKeyFileCredentials("~/.ydb/sa.json"),
		yc.WithInternalCA(),
	)
	if err != nil {
		panic(err)
	}
	_ = db.Close(context.TODO())
}
