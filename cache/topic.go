package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.third/proxy"
	"omo.msa.third/proxy/nosql"
	"time"
)

type TopicInfo struct {
	baseInfo
	Scene   string
	Remark  string
	Time    uint32
	Compere string
}

func (mine *cacheContext) CreateTopic(scene, name, remark, compere, operator string, secs uint32) (*TopicInfo, error) {
	db := new(nosql.Topic)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetTopicNextID()
	db.CreatedTime = time.Now()
	db.Creator = operator
	db.Operator = operator
	db.UpdatedTime = time.Now()
	db.Scene = scene
	db.Name = name
	db.Remark = remark
	db.Compere = compere
	db.Time = secs

	err := nosql.CreateTopic(db)
	if err == nil {
		info := new(TopicInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) GetTopic(uid string) (*TopicInfo, error) {
	db, err := nosql.GetTopic(uid)
	if err == nil {
		info := new(TopicInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) GetTopicRecordsByScene(scene string) ([]*nosql.TopicRecord, error) {
	return nosql.GetTopicRecordsByScene(scene)
}

func (mine *cacheContext) GetTopicRecordsBySN(sn string) ([]*nosql.TopicRecord, error) {
	return nosql.GetTopicRecordsBySN(sn)
}

func (mine *cacheContext) RemoveTopic(uid, operator string) error {
	return nosql.RemoveTopic(uid, operator)
}

func (mine *cacheContext) GetTopicsByScene(scene string) []*TopicInfo {
	dbs, err := nosql.GetTopicsByScene(scene)
	list := make([]*TopicInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			info := new(TopicInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *TopicInfo) initInfo(db *nosql.Topic) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Scene = db.Scene
	mine.Remark = db.Remark
	mine.Compere = db.Compere
	mine.Time = db.Time
}

func (mine *TopicInfo) UpdateBase(name, remark, compere, operator string, secs uint32) error {
	err := nosql.UpdateTopicBase(mine.UID, name, remark, compere, operator, secs)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Compere = compere
		mine.Operator = operator
		mine.Time = secs
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *TopicInfo) UpdateTitle(name, operator string) error {
	err := nosql.UpdateTopicTitle(mine.UID, name, operator)
	if err == nil {
		mine.Name = name
		mine.Operator = operator
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *TopicInfo) CreateRecord(sn, operator string, st uint32, date int64, results []*proxy.TopicResult) error {
	db := new(nosql.TopicRecord)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetTopicNextID()
	db.CreatedTime = time.Now()
	db.Creator = operator
	db.Operator = operator
	db.UpdatedTime = time.Now()
	db.Scene = mine.Scene
	db.Compere = mine.Compere
	db.Name = mine.Name
	db.Topic = mine.UID
	db.SN = sn
	db.Date = date
	db.State = st
	db.Results = results

	return nosql.CreateTopicRecord(db)
}
