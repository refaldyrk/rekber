package config

import (
	"context"
	"fmt"
	"time"

	"github.com/qiniu/qmgo"
	"github.com/spf13/viper"
)

func ConnectMongo(ctx context.Context) *qmgo.Client {
	start := time.Now()
	cli, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: viper.GetString("MONGO_URL")})
	if err != nil {
		panic(err)
	}

	fmt.Println("Connect Mongo : ", time.Since(start).Milliseconds(), " ms")
	return cli
}
