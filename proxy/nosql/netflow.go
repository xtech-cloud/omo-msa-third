package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.third/proxy"
)

type Netflow struct {
	UID     primitive.ObjectID `bson:"_id"`
	ID      uint64             `json:"id" bson:"id"`
	Name    string             `json:"name" bson:"name"`
	Created int64              `json:"created" bson:"created"`
	Creator string             `json:"creator" bson:"creator"`

	Type     uint8                `json:"type" bson:"type"`
	Remark   string               `json:"remark" bson:"remark"`
	Owner    string               `json:"owner" bson:"owner"`
	Quote    string               `json:"quote" bson:"quote"` //引用对象，可能是播放列表或者区域
	Size     uint64               `json:"size" bson:"size"`
	Template string               `json:"template" bson:"template"`
	Target   string               `json:"target" bson:"target"` //目标终端设备area
	Contents []*proxy.ContentInfo `json:"contents" bson:"contents"`
}

func CreateNetflow(info *Netflow) error {
	_, err := insertOne(TableNetflow, &info)
	return err
}

func GetNetflowNextID() uint64 {
	num, _ := getSequenceNext(TableNetflow)
	return num
}

func GetNetflows() ([]*Netflow, error) {
	cursor, err1 := findAllEnable(TableNetflow, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Netflow, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Netflow)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveNetflow(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db Netflow uid is empty ")
	}
	_, err := removeOne(TableNetflow, uid, operator)
	return err
}

func GetNetflow(uid string) (*Netflow, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Netflow uid is empty of GetNetflow")
	}

	result, err := findOne(TableNetflow, uid)
	if err != nil {
		return nil, err
	}
	model := new(Netflow)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetNetflowsByOwner(owner string) ([]*Netflow, error) {
	filter := bson.M{"owner": owner}
	cursor, err1 := findMany(TableNetflow, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Netflow, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Netflow)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetNetflowsByDuration(owner string, from, to int64) ([]*Netflow, error) {
	filter := bson.M{"owner": owner, TimeCreated: bson.M{"$gte": from, "$lte": to}}
	cursor, err1 := findMany(TableNetflow, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Netflow, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Netflow)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetNetflowCount() int64 {
	num, _ := getCount(TableNetflow)
	return num
}
