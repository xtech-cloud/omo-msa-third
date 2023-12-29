package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.third/proxy"
	"omo.msa.third/proxy/nosql"
	"time"
)

type HonorInfo struct {
	Type   uint8
	Status uint8
	baseInfo
	Scene    string
	Remark   string
	Parent   string
	Style    string
	Contents []*proxy.HonorContent
}

func (mine *cacheContext) CreateHonor(scene, parent, name, remark, style, operator string, tp, st uint32) (*HonorInfo, error) {
	db := new(nosql.Honor)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetHonorNextID()
	db.Created = time.Now().Unix()
	db.Creator = operator
	db.Operator = operator
	db.Updated = time.Now().Unix()
	db.Scene = scene
	db.Name = name
	db.Remark = remark
	db.Parent = parent
	db.Style = style
	db.Type = uint8(tp)
	db.Status = uint8(st)

	err := nosql.CreateHonor(db)
	if err == nil {
		info := new(HonorInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) GetHonor(uid string) (*HonorInfo, error) {
	if uid == "" {
		return nil, errors.New("the uid is empty")
	}
	db, err := nosql.GetHonor(uid)
	if err == nil {
		info := new(HonorInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) RemoveHonor(uid, operator string) error {
	return nosql.RemoveHonor(uid, operator)
}

func (mine *cacheContext) GetHonorsByScene(scene string) []*HonorInfo {

	if scene == "" {
		return nil
	}
	dbs, err := nosql.GetHonorsByScene(scene)
	list := make([]*HonorInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			if db.Parent == "" {
				info := new(HonorInfo)
				info.initInfo(db)
				list = append(list, info)
			}
		}
	}
	return list
}

func (mine *cacheContext) GetHonorsByParent(parent string) []*HonorInfo {
	dbs, err := nosql.GetHonorsByParent(parent)
	list := make([]*HonorInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			info := new(HonorInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *HonorInfo) initInfo(db *nosql.Honor) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Created = db.Created
	mine.Updated = db.Updated
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Scene = db.Scene
	mine.Remark = db.Remark
	mine.Parent = db.Parent
	mine.Style = db.Style
	mine.Type = db.Type
	mine.Status = db.Status
	mine.Contents = db.Contents
}

func (mine *HonorInfo) UpdateBase(name, remark, style, operator string) error {
	err := nosql.UpdateHonorBase(mine.UID, name, remark, style, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Style = style
		mine.Operator = operator
		mine.Updated = time.Now().Unix()
	}
	return err
}

func (mine *HonorInfo) UpdateContents(operator string, list []*proxy.HonorContent) error {
	err := nosql.UpdateHonorContents(mine.UID, operator, list)
	if err == nil {
		mine.Contents = list
		mine.Operator = operator
		mine.Updated = time.Now().Unix()
	}
	return err
}

func (mine *HonorInfo) UpdateStatus(operator string, st uint32) error {
	err := nosql.UpdateHonorStatus(mine.UID, operator, uint8(st))
	if err == nil {
		mine.Status = uint8(st)
		mine.Operator = operator
		mine.Updated = time.Now().Unix()
	}
	return err
}
