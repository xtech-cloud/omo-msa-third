package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Holiday struct {
	UID      primitive.ObjectID `bson:"_id"`
	ID       uint64             `json:"id" bson:"id"`
	Name     string             `json:"name" bson:"name"`
	Created  int64              `json:"created" bson:"created"`
	Updated  int64              `json:"updated" bson:"updated"`
	Deleted  int64              `json:"deleted" bson:"deleted"`
	Creator  string             `json:"creator" bson:"creator"`
	Operator string             `json:"operator" bson:"operator"`

	Remark string `json:"remark" bson:"remark"`
	Owner  string `json:"owner" bson:"owner"`
	Type   uint8  `json:"type" bson:"type"`
	From   int64  `json:"from" bson:"from"`
	End    int64  `json:"end" bson:"end"`
}

func CreateHoliday(info *Holiday) error {
	_, err := insertOne(TableHoliday, &info)
	return err
}

func GetHolidayNextID() uint64 {
	num, _ := getSequenceNext(TableHoliday)
	return num
}

func GetHolidays() ([]*Holiday, error) {
	cursor, err1 := findAllEnable(TableHoliday, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Holiday, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Holiday)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveHoliday(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db Holiday uid is empty ")
	}
	_, err := removeOne(TableHoliday, uid, operator)
	return err
}

func GetHoliday(uid string) (*Holiday, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Holiday uid is empty of GetHoliday")
	}

	result, err := findOne(TableHoliday, uid)
	if err != nil {
		return nil, err
	}
	model := new(Holiday)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetHolidaysByOwner(owner string) ([]*Holiday, error) {
	filter := bson.M{"owner": owner, TimeDeleted: 0}
	cursor, err1 := findMany(TableHoliday, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Holiday, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Holiday)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetHolidaysByType(owner string, tp uint32) ([]*Holiday, error) {
	filter := bson.M{"owner": owner, "type": tp, TimeDeleted: 0}
	cursor, err1 := findMany(TableHoliday, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Holiday, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Holiday)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetHolidayCount() int64 {
	num, _ := getCount(TableHoliday)
	return num
}

func UpdateHolidayBase(uid, name, remark, operator string, from, end int64) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "from": from, "end": end, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableHoliday, uid, msg)
	return err
}
