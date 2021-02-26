package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Partner struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Status      uint8 `json:"status" bson:"status"`
	Cover       string             `json:"cover" bson:"cover"`
	Remark      string             `json:"remark" bson:"remark"`
	Phone       string 				`json:"phone" bson:"phone"`
	Secret      string 				`json:"secret" bson:"secret"`
	Tags        []string 			`json:"tags" bsonL:"tags"`
}

func CreatePartner(info *Partner) error {
	_, err := insertOne(TablePartner, &info)
	return err
}

func GetPartnerNextID() uint64 {
	num, _ := getSequenceNext(TablePartner)
	return num
}

func GetPartners() ([]*Partner, error) {
	cursor, err1 := findAll(TablePartner, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Partner, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Partner)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemovePartner(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db Partner uid is empty ")
	}
	_, err := removeOne(TablePartner, uid, operator)
	return err
}

func GetPartner(uid string) (*Partner, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Partner uid is empty of GetPartner")
	}

	result, err := findOne(TablePartner, uid)
	if err != nil {
		return nil, err
	}
	model := new(Partner)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetPartnerBySecret(secret string) (*Partner, error) {
	if len(secret) < 2 {
		return nil, errors.New("db partner uid is empty of GetPartnerBySecret")
	}
	msg := bson.M{"secret": secret, "deleteAt": new(time.Time)}
	result, err := findOneBy(TablePartner, msg)
	if err != nil {
		return nil, err
	}
	model := new(Partner)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetPartnerCount() int64 {
	num, _ := getCount(TablePartner)
	return num
}

func UpdatePartnerBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TablePartner, uid, msg)
	return err
}

func UpdatePartnerCover(uid, cover, operator string) error {
	msg := bson.M{"cover": cover, "operator":operator, "updatedAt": time.Now()}
	_, err := updateOne(TablePartner, uid, msg)
	return err
}

func UpdatePartnerSecret(uid, secret, operator string) error {
	msg := bson.M{"secret": secret, "operator":operator, "updatedAt": time.Now()}
	_, err := updateOne(TablePartner, uid, msg)
	return err
}

func UpdatePartnerTags(uid, operator string, tags []string) error {
	msg := bson.M{"tags": tags, "operator":operator, "updatedAt": time.Now()}
	_, err := updateOne(TablePartner, uid, msg)
	return err
}

