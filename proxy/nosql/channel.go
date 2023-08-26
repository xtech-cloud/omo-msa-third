package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Channel struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Created     int64              `json:"created" bson:"created"`
	Updated     int64              `json:"updated" bson:"updated"`
	Deleted     int64              `json:"deleted" bson:"deleted"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Remark string   `json:"remark" bson:"remark"`
	Owner  string   `json:"owner" bson:"owner"`
	Type   uint8    `json:"type" bson:"type"`
	Bags   []string `json:"bags" bson:"bags"`
}

func CreateChannel(info *Channel) error {
	_, err := insertOne(TableChannel, &info)
	return err
}

func GetChannelNextID() uint64 {
	num, _ := getSequenceNext(TableChannel)
	return num
}

func GetChannelCount() int64 {
	num, _ := getCount(TableChannel)
	return num
}

func GetChannels() ([]*Channel, error) {
	cursor, err1 := findAllEnable(TableChannel, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Channel, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Channel)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveChannel(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db Channel uid is empty ")
	}
	_, err := removeOne(TableChannel, uid, operator)
	return err
}

func GetChannel(uid string) (*Channel, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Channel uid is empty of GetChannel")
	}

	result, err := findOne(TableChannel, uid)
	if err != nil {
		return nil, err
	}
	model := new(Channel)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetChannelByOwner(owner string) (*Channel, error) {
	filter := bson.M{"owner": owner, TimeDeleted: 0}
	result, err := findOneBy(TableChannel, filter)
	if err != nil {
		return nil, err
	}
	model := new(Channel)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func UpdateChannelBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableChannel, uid, msg)
	return err
}

func UpdateChannelBags(uid, operator string, list []string) error {
	msg := bson.M{"bags": list, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableChannel, uid, msg)
	return err
}

func AppendChannelBag(uid, prop string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"bags": prop}
	_, err := appendElement(TableChannel, uid, msg)
	return err
}

func SubtractChannelBag(uid, key string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"bags": key}
	_, err := removeElement(TableChannel, uid, msg)
	return err
}
