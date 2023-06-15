package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.third/proxy/nosql"
	"strings"
	"time"
)

type MotionInfo struct {
	Count uint32
	baseInfo
	Scene   string
	AppID   string
	SN      string
	UserID  string
	EventID string
	Content string
}

func (mine *cacheContext) CreateMotion(scene, app, sn, eveID, content, operator string, count uint32) (*MotionInfo, error) {
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
	db.Content = content

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

func (mine *cacheContext) GetMotionsBySNEvent(scene, sn, id string) []*MotionInfo {
	dbs, err := nosql.GetMotionsBySNEvent(scene, sn, id)
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

func (mine *cacheContext) GetMotionBy(scene, sn, event, content string) *MotionInfo {
	db, err := nosql.GetMotionBy(scene, sn, event, content)
	if err == nil {
		info := new(MotionInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext) GetMotionsByContent(scene, sn, content string) []*MotionInfo {
	dbs, err := nosql.GetMotionsBySN(scene, sn)
	list := make([]*MotionInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			if strings.Contains(db.Content, content) {
				info := new(MotionInfo)
				info.initInfo(db)
				list = append(list, info)
			}
		}
	}
	return list
}

func (mine *cacheContext) GetMotionsByEveContent(scene, eve, content string) []*MotionInfo {
	dbs, err := nosql.GetMotionsByEventID(scene, eve)
	list := make([]*MotionInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			if strings.Contains(db.Content, content) {
				info := new(MotionInfo)
				info.initInfo(db)
				list = append(list, info)
			}
		}
	}
	return list
}

func (mine *MotionInfo) initInfo(db *nosql.Motion) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
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
	mine.Content = db.Content
}

func (mine *MotionInfo) AddCount(offset uint32, operator string) error {
	err := nosql.UpdateMotionCount(mine.UID, operator, mine.Count+offset)
	if err == nil {
		mine.Count = mine.Count + offset
		mine.UpdateTime = time.Now()
	}
	return err
}
