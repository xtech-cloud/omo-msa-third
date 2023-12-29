package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.third/proxy"
	"time"
)

type Honor struct {
	UID      primitive.ObjectID `bson:"_id"`
	ID       uint64             `json:"id" bson:"id"`
	Name     string             `json:"name" bson:"name"`
	Created  int64              `json:"created" bson:"created"`
	Updated  int64              `json:"updated" bson:"updated"`
	Deleted  int64              `json:"deleted" bson:"deleted"`
	Creator  string             `json:"creator" bson:"creator"`
	Operator string             `json:"operator" bson:"operator"`

	Scene    string                `json:"scene" bson:"scene"`
	Remark   string                `json:"remark" bson:"remark"`
	Type     uint8                 `json:"type" bson:"type"`
	Status   uint8                 `json:"status" bson:"status"`
	Style    string                `json:"style" bson:"style"`
	Parent   string                `json:"parent" bson:"parent"`
	Contents []*proxy.HonorContent `json:"contents" bson:"contents"`
}

func CreateHonor(info *Honor) error {
	_, err := insertOne(TableHonor, &info)
	return err
}

func GetHonorNextID() uint64 {
	num, _ := getSequenceNext(TableHonor)
	return num
}

func GetHonors() ([]*Honor, error) {
	cursor, err1 := findAllEnable(TableHonor, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Honor, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Honor)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveHonor(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db Honor uid is empty ")
	}
	_, err := removeOne(TableHonor, uid, operator)
	return err
}

func GetHonor(uid string) (*Honor, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Honor uid is empty of GetHonor")
	}

	result, err := findOne(TableHonor, uid)
	if err != nil {
		return nil, err
	}
	model := new(Honor)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetHonorsByScene(scene string) ([]*Honor, error) {
	filter := bson.M{"scene": scene, TimeDeleted: 0}
	cursor, err1 := findMany(TableHonor, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Honor, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Honor)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetHonorsByParent(parent string) ([]*Honor, error) {
	filter := bson.M{"parent": parent, TimeDeleted: 0}
	cursor, err1 := findMany(TableHonor, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Honor, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Honor)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetHonorCount() int64 {
	num, _ := getCount(TableHonor)
	return num
}

func UpdateHonorBase(uid, name, remark, style, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "style": style, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableHonor, uid, msg)
	return err
}

func UpdateHonorContents(uid, operator string, list []*proxy.HonorContent) error {
	msg := bson.M{"contents": list, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableHonor, uid, msg)
	return err
}

func UpdateHonorStatus(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableHonor, uid, msg)
	return err
}
