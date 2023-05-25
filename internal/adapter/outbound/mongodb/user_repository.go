package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/Kercyn/crud_template/configs"
	"github.com/Kercyn/crud_template/internal/adapter/outbound/mongodb/dto"
	"github.com/Kercyn/crud_template/internal/core/domain"
	port "github.com/Kercyn/crud_template/internal/core/port/outbound"
	apperror "github.com/Kercyn/crud_template/pkg/error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	ctx    context.Context
	config configs.Config
	client *mongo.Client
}

func NewMongoUserRepository(
	client *mongo.Client,
	ctx context.Context,
	config configs.Config,
) MongoUserRepository {
	return MongoUserRepository{
		ctx:    ctx,
		config: config,
		client: client,
	}
}

func (r MongoUserRepository) GetHealthStats() error {
	// todo return more things like collection size etc.
	// https://www.mongodb.com/docs/manual/reference/command/collStats/#dbcmd.collStats
	return r.client.Ping(r.ctx, nil)
}

func (r MongoUserRepository) GetByUserID(userID port.DataSourceID) (domain.User, error) {
	collection, err := r.getUserCollection()
	if err != nil {
		return domain.User{}, err
	}

	userIDFilter := bson.M{
		"_id": userID.AsDataSourceID(),
	}
	result := collection.FindOne(r.ctx, userIDFilter)
	user, err := r.parseResult(result)
	if err != nil {
		return domain.User{}, err
	}

	return user.ToDomain()
}

func (r MongoUserRepository) getUserCollection() (*mongo.Collection, error) {
	database := r.client.Database(r.config.DataSourceDatabase)
	if database == nil {
		return nil, apperror.AppError{
			Err: errors.New(fmt.Sprintf(
				"Database %s is nil",
				r.config.DataSourceDatabase,
			)),
			UserFriendlyMessage: "Could not get database",
		}
	}

	collection := database.Collection(r.config.DataSourceCollection)
	if collection == nil {
		return nil, apperror.AppError{
			Err: errors.New(fmt.Sprintf(
				"Collection %s is nil",
				r.config.DataSourceDatabase,
			)),
			UserFriendlyMessage: "Could not get collection",
		}
	}

	return collection, nil
}

func (r MongoUserRepository) parseResult(result *mongo.SingleResult) (dto.User, error) {
	var user dto.User
	err := result.Decode(&user)
	if err == mongo.ErrNoDocuments {
		return dto.User{}, apperror.AppError{
			Err:                 err,
			UserFriendlyMessage: fmt.Sprintf("Could not find user"),
		}
	} else if err != nil {
		return dto.User{}, apperror.AppError{
			Err:                 err,
			UserFriendlyMessage: apperror.MessageUnknownError,
		}
	}

	return user, nil
}

func (r MongoUserRepository) Patch(request port.PatchData) error {
	collection, err := r.getUserCollection()
	if err != nil {
		return err
	}

	userIDFilter := bson.M{
		"_id": request.ID.AsDataSourceID(),
	}

	patchData, err := buildPatchData(request.Fields)
	if err != nil {
		return err
	}

	patchQuery := bson.M{"$set": patchData}

	_, err = collection.UpdateOne(r.ctx, userIDFilter, patchQuery)
	if err != nil {
		return err
	}

	return nil
}

func buildPatchData(fields map[string]interface{}) (bson.M, error) {
	var query bson.M

	for k, v := range fields {
		if isMap(v) {
			compositeValue, err := buildPatchData(v.(map[string]interface{}))
			if err != nil {
				return bson.M{}, err
			}
			query[k] = compositeValue
		} else {
			query[k] = v
		}
	}

	return query, nil
}

func isMap(valueToCheck interface{}) bool {
	_, ok := valueToCheck.(map[string]interface{})
	return ok
}
