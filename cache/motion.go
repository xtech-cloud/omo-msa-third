package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.third/proxy/nosql"
	"time"
)

type MotionInfo struct {
	baseInfo
	Scene   string
	AppID   string
	SN      string
	UserID  string
	EventID string
	Count   uint32
}

func (mine *cacheContext) CreateMotion(scene, app, sn, eveID, operator string, count uint32) (*MotionInfo, error) {
	db := new(nosql.Motion)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetMotionNextID()
	db.CreatedTime = time.Now()
	db.Creator = operator
	db.Operator = operator
	db.UpdatedTime = time.Now()
	db.AppID = app
	db.Scene = scene
	db.SN = sn
	db.EventID = eveID
	db.Count = count

	err := nosql.CreateMotion(db)
	if err == nil {
		info := new(MotionInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) GetMotionsBySN(scene, sn string) []*MotionInfo {
	dbs, err := nosql.GetMotionsBySN(scene, sn)
	list := make([]*MotionInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			info := new(MotionInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetMotion(uid string) (*MotionInfo, error) {
	db, err := nosql.GetMotion(uid)
	if err == nil {
		info := new(MotionInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) RemoveMotion(uid, operator string) error {
	return nosql.RemoveMotion(uid, operator)
}

func (mine *cacheContext) GetMotionsByEvent(scene, id string) []*MotionInfo {
	dbs, err := nosql.GetMotionsByEventID(scene, id)
	list := make([]*MotionInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			info := new(MotionInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetMotionsBy(scene, sn, event string) []*MotionInfo {
	dbs, err := nosql.GetMotionsBy(scene, sn, event)
	list := make([]*MotionInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			info := new(MotionInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *MotionInfo) initInfo(db *nosql.Motion) {
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Scene = db.Scene
	mine.EventID = db.EventID
	mine.SN = db.SN
	mine.AppID = db.AppID
	mine.Count = db.Count
	mine.UserID = db.UserID
}

func (mine *MotionInfo) UpdateCount(offset uint32, operator string) error {
	err := nosql.UpdateMotionCount(mine.UID, operator, mine.Count+offset)
	if err == nil {
		mine.Count = mine.Count + offset
		mine.UpdateTime = time.Now()
	}
	return err
}
