package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kercyn/crud_template/configs"
	"github.com/Kercyn/crud_template/internal/adapter/outbound/mongodb"
	"github.com/Kercyn/crud_template/internal/adapter/outbound/mongodb/dto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func Execute(config configs.Config) error {
	clientOptions := options.Client().ApplyURI(config.DataSourceURI)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(config.DataSourceTimeoutInMilliseconds)*time.Millisecond,
	)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	userRepository := mongodb.NewMongoUserRepository(client, ctx, config)
	userID, err := dto.NewMongoID("646556f77b0912000743ff46")
	if err != nil {
		return err
	}

	user, err := userRepository.GetByUserID(userID)
	if err != nil {
		return err
	}

	u, err := json.Marshal(user)
	fmt.Println(string(u))
	return nil
}
