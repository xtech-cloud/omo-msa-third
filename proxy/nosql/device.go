package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Invite struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Created     int64              `json:"created" bson:"created"`
	Updated     int64              `json:"updated" bson:"updated"`
	Deleted     int64              `json:"deleted" bson:"deleted"`

	Creator     string `json:"creator" bson:"creator"`
	Operator    string `json:"operator" bson:"operator"`
	Scene       string `json:"scene" bson:"scene"` // 所属场景
	Type        uint8  `json:"type" bson:"type"`   //类型
	Status      uint8  `json:"status" bson:"status"`
	Remark      string `json:"remark" bson:"remark"`           //备注
	SN          string `json:"sn" bson:"sn"`                   //设备SN或者邀请码
	OS          string `json:"os" bson:"os"`                   //操作系统
	ExpiryTime  uint32 `json:"expiry" bson:"expiry"`           //有效时长
	ActiveTime  int64  `json:"activated" bson:"activated"`     //激活时间
	Quote       string `json:"quote" bson:"quote"`             //
	Certificate string `json:"certificate" bson:"certificate"` //激活证书
}

func GetAllDevices() ([]*Invite, error) {
	cursor, err1 := findAllOld(TableDevice, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Invite, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Invite)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func findAllOld(collection string, limit int64) (*mongo.Cursor, error) {
	if len(collection) < 1 {
		return nil, errors.New("the collection is empty")
	}
	c := noSql.Collection(collection)
	if c == nil {
		return nil, errors.New("can not found the collection of" + collection)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	filter := bson.M{"deleteAt": new(time.Time)}
	var cursor *mongo.Cursor
	var err error
	if limit > 0 {
		cursor, err = c.Find(ctx, filter, options.Find().SetSort(bson.M{TimeCreated: -1}).SetLimit(limit))
	} else {
		cursor, err = c.Find(ctx, filter)
	}
	if err != nil {
		return nil, err
	}
	return cursor, nil
}
