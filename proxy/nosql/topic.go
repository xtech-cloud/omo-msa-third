package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Topic struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Scene   string `json:"scene" bson:"scene"`
	Compere string `json:"compere" bson:"compere"` //主持人
	Time    uint32 `json:"time" bson:"time"`       //时间
	Remark  string `json:"remark" bson:"remark"`
}

func CreateTopic(info *Topic) error {
	_, err := insertOne(TableTopic, &info)
	return err
}

func GetTopicNextID() uint64 {
	num, _ := getSequenceNext(TableTopic)
	return num
}

func GetTopics() ([]*Topic, error) {
	cursor, err1 := findAll(TableTopic, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Topic, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Topic)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveTopic(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db Topic uid is empty ")
	}
	_, err := removeOne(TableTopic, uid, operator)
	return err
}

func GetTopic(uid string) (*Topic, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Topic uid is empty of GetTopic")
	}

	result, err := findOne(TableTopic, uid)
	if err != nil {
		return nil, err
	}
	model := new(Topic)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetTopicsByScene(scene string) ([]*Topic, error) {
	def := new(time.Time)
	filter := bson.M{"scene": scene, "deleteAt": def}
	cursor, err1 := findMany(TableTopic, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Topic, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Topic)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetTopicCount() int64 {
	num, _ := getCount(TableTopic)
	return num
}

func UpdateTopicBase(uid, name, remark, compere, operator string, secs uint32) error {
	msg := bson.M{"name": name, "remark": remark, "compere": compere, "time": secs, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTopic, uid, msg)
	return err
}

func UpdateTopicTitle(uid, name, operator string) error {
	msg := bson.M{"name": name, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTopic, uid, msg)
	return err
}
