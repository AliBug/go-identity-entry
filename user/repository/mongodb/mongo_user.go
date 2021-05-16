package mongorepo

import (
	"context"
	"time"

	"github.com/alibug/go-identity-utils/status"
	"github.com/alibug/go-identity/domain"
	"github.com/alibug/go-identity/user/repository/body"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoUserRepository struct {
	userColl *mongo.Collection
}

// NewMongoUserRepository will create an object that represent the user.Repository interface
func NewMongoUserRepository(coll *mongo.Collection) domain.UserRepository {
	return &mongoUserRepository{coll}
}

func (m *mongoUserRepository) RegisterUser(ctx context.Context, register domain.Register) error {
	// 1、 GetByUsername, if user existed, throw err
	// 	  如果已有同名用户，抛出错误
	_, err := m.GetByAccount(ctx, register.GetAccount())
	if err == nil {
		return status.ErrConflict
	}

	// 2、SetCreatedTime
	//	 如果没有错误，设置用户注册日期
	now := time.Now()
	register.SetCreatedTime(&now)

	// 3、SetCyptPass
	err = register.SetCryptPass()
	if err != nil {
		return err
	}

	// 3、Insert User
	//	 将用户数据入库
	_, err = m.userColl.InsertOne(ctx, register)
	return err
}

func (m *mongoUserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, status.ErrBadParamInput
	}

	var u body.UserBody
	err = m.userColl.FindOne(ctx, bson.M{"_id": objectID}).Decode(&u)
	return &u, err
}

func (m *mongoUserRepository) GetByAccount(ctx context.Context, account string) (domain.User, error) {
	var u body.UserBody
	err := m.userColl.FindOne(ctx, bson.M{"account": account}).Decode(&u)
	return &u, err
}
