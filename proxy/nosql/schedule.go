package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.third/proxy"
	"time"
)

type Schedule struct {
	UID      primitive.ObjectID `bson:"_id"`
	ID       uint64             `json:"id" bson:"id"`
	Name     string             `json:"name" bson:"name"`
	Created  int64              `json:"created" bson:"created"`
	Updated  int64              `json:"updated" bson:"updated"`
	Deleted  int64              `json:"deleted" bson:"deleted"`
	Creator  string             `json:"creator" bson:"creator"`
	Operator string             `json:"operator" bson:"operator"`

	Status   uint8              `json:"status" bson:"status"`
	Type     uint8              `json:"type" bson:"type"`
	Ignore   uint8              `json:"ignore" bson:"ignore"`
	Remark   string             `json:"remark" bson:"remark"`
	Owner    string             `json:"owner" bson:"owner"`
	Quote    string             `json:"quote" bson:"quote"` //引用对象，可能是播放列表或者区域
	Color    string             `json:"color" bson:"color"`
	Date     proxy.DurationInfo `json:"date" bson:"date"`
	Time     proxy.DurationInfo `json:"time" bson:"time"`
	Weekdays []uint32           `json:"weekdays" bson:"weekdays"`
	Targets  []string           `json:"targets" bson:"targets"`
}

func CreateSchedule(info *Schedule) error {
	_, err := insertOne(TableSchedule, &info)
	return err
}

func GetScheduleNextID() uint64 {
	num, _ := getSequenceNext(TableSchedule)
	return num
}

func GetSchedules() ([]*Schedule, error) {
	cursor, err1 := findAllEnable(TableSchedule, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Schedule, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Schedule)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveSchedule(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db Schedule uid is empty ")
	}
	_, err := removeOne(TableSchedule, uid, operator)
	return err
}

func GetSchedule(uid string) (*Schedule, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Schedule uid is empty of GetSchedule")
	}

	result, err := findOne(TableSchedule, uid)
	if err != nil {
		return nil, err
	}
	model := new(Schedule)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetScheduleByTimeDate(uid string) (*Schedule, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Schedule uid is empty of GetSchedule")
	}

	result, err := findOne(TableSchedule, uid)
	if err != nil {
		return nil, err
	}
	model := new(Schedule)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetSchedulesByOwner(owner string) ([]*Schedule, error) {
	filter := bson.M{"owner": owner, TimeDeleted: 0}
	cursor, err1 := findMany(TableSchedule, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Schedule, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Schedule)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetScheduleCount() int64 {
	num, _ := getCount(TableSchedule)
	return num
}

func UpdateScheduleBase(uid, name, remark, operator, quote, color string, tp, ignore uint8, date, timee proxy.DurationInfo, list []string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "quote": quote, "date": date, "time": timee, "targets": list,
		"color": color, "ignore": ignore, "type": tp, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableSchedule, uid, msg)
	return err
}

func UpdateScheduleTags(uid, operator string, tags []string) error {
	msg := bson.M{"tags": tags, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableSchedule, uid, msg)
	return err
}

func UpdateScheduleStatus(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableSchedule, uid, msg)
	return err
}
