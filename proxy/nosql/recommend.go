package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Recommend struct {
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

	Type  uint8  `json:"type" bson:"type"`
	Owner string `json:"owner" bson:"owner"`
	Quote string `json:"quote" bson:"quote"`

	Targets []string `json:"targets" bson:"targets"`
}

func CreateRecommend(info *Recommend) error {
	_, err := insertOne(TableRecommend, &info)
	return err
}

func GetRecommendNextID() uint64 {
	num, _ := getSequenceNext(TableRecommend)
	return num
}

func RemoveRecommend(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db Recommend uid is empty ")
	}
	_, err := removeOne(TableRecommend, uid, operator)
	return err
}

func GetRecommend(uid string) (*Recommend, error) {
	if len(uid) < 2 {
		return nil, errors.New("db activity uid is empty of GetRecommend")
	}
	result, err := findOne(TableRecommend, uid)
	if err != nil {
		return nil, err
	}
	model := new(Recommend)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetRecommendCount() int64 {
	num, _ := getTotalCount(TableRecommend)
	return num
}

func GetRecommendBy(owner string, tp uint8) (*Recommend, error) {
	if len(owner) < 2 {
		return nil, errors.New("db owner is empty of GetRecommendBy")
	}
	filter := bson.M{"owner": owner, "type": tp, TimeDeleted: 0}
	result, err := findOneBy(TableRecommend, filter)
	if err != nil {
		return nil, err
	}
	model := new(Recommend)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetRecommendByOwnerQuote(owner, quote string) ([]*Recommend, error) {
	filter := bson.M{"owner": owner, "quote": quote, TimeDeleted: 0}
	cursor, err1 := findMany(TableRecommend, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Recommend, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Recommend)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetRecommendByOwnerTarget(owner, target string) ([]*Recommend, error) {
	filter := bson.M{"owner": owner, "targets": target, TimeDeleted: 0}
	cursor, err1 := findMany(TableRecommend, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Recommend, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Recommend)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetRecommendsByQuote(quote string) ([]*Recommend, error) {
	filter := bson.M{"quote": quote, TimeDeleted: 0}
	cursor, err1 := findMany(TableRecommend, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Recommend, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Recommend)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetRecommendByT(owner string) (*Recommend, error) {
	if len(owner) < 2 {
		return nil, errors.New("db owner is empty of GetRecommendBy")
	}
	filter := bson.M{"owner": owner, TimeDeleted: 0}
	result, err := findOneBy(TableRecommend, filter)
	if err != nil {
		return nil, err
	}
	model := new(Recommend)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}
func GetRecommendByOwner(owner string) ([]*Recommend, error) {
	filter := bson.M{"owner": owner, TimeDeleted: 0}
	cursor, err1 := findMany(TableRecommend, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Recommend, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Recommend)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}
func GetRecommendsByType(owner string, tp uint32) ([]*Recommend, error) {
	filter := bson.M{"owner": owner, "type": tp, TimeDeleted: 0}
	cursor, err1 := findMany(TableRecommend, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Recommend, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Recommend)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateRecommendBase(uid, name, sub, body, operator string) error {
	msg := bson.M{"name": name, "body": body, "subtitle": sub, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableRecommend, uid, msg)
	return err
}

func UpdateRecommendTargets(uid, operator string, list []string) error {
	msg := bson.M{"targets": list, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableRecommend, uid, msg)
	return err
}
