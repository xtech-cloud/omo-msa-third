package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Motion struct {
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

	Count   uint32 `json:"count" bson:"count"`
	Scene   string `json:"scene" bson:"scene"`
	AppID   string `json:"app" bson:"app"`
	SN      string `json:"sn" bson:"sn"`
	UserID  string `json:"user" bson:"user"`
	EventID string `json:"event" bson:"event"`
	Content string `json:"content" bson:"content"`
}

func CreateMotion(info *Motion) error {
	_, err := insertOne(TableMotion, &info)
	return err
}

func GetMotionNextID() uint64 {
	num, _ := getSequenceNext(TableMotion)
	return num
}

func GetMotions() ([]*Motion, error) {
	cursor, err1 := findAllEnable(TableMotion, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Motion, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Motion)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveMotion(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db Motion uid is empty ")
	}
	_, err := removeOne(TableMotion, uid, operator)
	return err
}

func GetMotion(uid string) (*Motion, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Motion uid is empty of GetMotion")
	}

	result, err := findOne(TableMotion, uid)
	if err != nil {
		return nil, err
	}
	model := new(Motion)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetMotionsByEventID(scene, eve string) ([]*Motion, error) {
	filter := bson.M{"scene": scene, "event": eve, TimeDeleted: 0}
	cursor, err1 := findMany(TableMotion, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Motion, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Motion)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetMotionsBySNEvent(scene, sn, eve string) ([]*Motion, error) {
	filter := bson.M{"scene": scene, "sn": sn, "event": eve, TimeDeleted: 0}
	cursor, err1 := findMany(TableMotion, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Motion, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Motion)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetMotionsBySN(scene, sn string) ([]*Motion, error) {
	filter := bson.M{"scene": scene, "sn": sn, TimeDeleted: 0}
	cursor, err1 := findMany(TableMotion, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Motion, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Motion)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetMotionsBy(scene, sn, event, content string) ([]*Motion, error) {
	filter := bson.M{"scene": scene, "sn": sn, "event": event, "content": content, TimeDeleted: 0}
	cursor, err1 := findMany(TableMotion, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Motion, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Motion)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetMotionBy(scene, sn, event, content string) (*Motion, error) {
	filter := bson.M{"scene": scene, "sn": sn, "event": event, "content": content, TimeDeleted: 0}
	result, err := findOneBy(TableMotion, filter)
	if err != nil {
		return nil, err
	}
	model := new(Motion)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetMotionsByContent(scene, sn, content string) ([]*Motion, error) {
	filter := bson.M{"scene": scene, "sn": sn, "content": content, TimeDeleted: 0}
	cursor, err1 := findMany(TableMotion, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Motion, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Motion)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetMotionsByEventContent(scene, event, content string) ([]*Motion, error) {
	filter := bson.M{"scene": scene, "event": event, "content": content, TimeDeleted: 0}
	cursor, err1 := findMany(TableMotion, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Motion, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Motion)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetMotionCount() int64 {
	num, _ := getCount(TableMotion)
	return num
}

func UpdateMotionCount(uid, operator string, num uint32) error {
	msg := bson.M{"count": num, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableMotion, uid, msg)
	return err
}
