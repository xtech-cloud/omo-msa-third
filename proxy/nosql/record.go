package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.third/proxy"
	"time"
)

type TopicRecord struct {
	UID      primitive.ObjectID `bson:"_id"`
	ID       uint64             `json:"id" bson:"id"`
	Name     string             `json:"name" bson:"name"`
	Created  int64              `json:"created" bson:"created"`
	Updated  int64              `json:"updated" bson:"updated"`
	Deleted  int64              `json:"deleted" bson:"deleted"`
	Creator  string             `json:"creator" bson:"creator"`
	Operator string             `json:"operator" bson:"operator"`

	Scene   string               `json:"scene" bson:"scene"`
	Compere string               `json:"compere" bson:"compere"` //主持人
	Date    int64                `json:"date" bson:"date"`       //时间
	Topic   string               `json:"topic" bson:"topic"`
	State   uint32               `json:"state" bson:"state"`
	SN      string               `json:"sn" bson:"sn"`
	Results []*proxy.TopicResult `json:"results" bson:"results"`
}

func CreateTopicRecord(info *TopicRecord) error {
	_, err := insertOne(TableRecords, &info)
	return err
}

func GetTopicRecordNextID() uint64 {
	num, _ := getSequenceNext(TableRecords)
	return num
}

func GetTopicRecords() ([]*TopicRecord, error) {
	cursor, err1 := findAllEnable(TableRecords, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*TopicRecord, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(TopicRecord)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveTopicRecord(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db TopicRecord uid is empty ")
	}
	_, err := removeOne(TableRecords, uid, operator)
	return err
}

func GetTopicRecord(uid string) (*TopicRecord, error) {
	if len(uid) < 2 {
		return nil, errors.New("db TopicRecord uid is empty of GetTopicRecord")
	}

	result, err := findOne(TableRecords, uid)
	if err != nil {
		return nil, err
	}
	model := new(TopicRecord)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetTopicRecordsByScene(scene string) ([]*TopicRecord, error) {
	filter := bson.M{"scene": scene, TimeDeleted: 0}
	cursor, err1 := findMany(TableRecords, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*TopicRecord, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(TopicRecord)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetTopicRecordsBySN(sn string) ([]*TopicRecord, error) {
	filter := bson.M{"sn": sn, TimeDeleted: 0}
	cursor, err1 := findMany(TableRecords, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*TopicRecord, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(TopicRecord)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetTopicRecordCount() int64 {
	num, _ := getCount(TableRecords)
	return num
}

func UpdateTopicRecordBase(uid, name, remark, compere, operator string, num uint32) error {
	msg := bson.M{"name": name, "remark": remark, "compere": compere, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableRecords, uid, msg)
	return err
}
