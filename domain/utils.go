package domain

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StrToObjectID - Convert type of ID between string and primitive.ObjectID
type StrToObjectID string

// MarshalBSONValue - Help ID type conversion
func (id StrToObjectID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	p, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return bsontype.Null, nil, err
	}
	return bson.MarshalValue(p)
}

// CookieConfig -
type CookieConfig interface {
	GetAccessTokenMaxAge() int
	GetRefreshTokenMaxAge() int
	GetDomain() string
	GetSecure() bool
	GetHTTPOnly() bool
}

/*
 */
