package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.third/proxy/nosql"
	"time"
)

const (
	SourceActivity SourceType = 0
	SourceEntity   SourceType = 1
	SourceArticle  SourceType = 2
	SourcePhoto    SourceType = 3
	SourceCourse   SourceType = 4
)

type SourceType uint8

type RecommendInfo struct {
	Type uint8
	baseInfo
	Owner   string
	Targets []string
}

func (mine *cacheContext) CreateRecommend(owner, operator string, tp uint8, list []string) (*RecommendInfo, error) {
	db := new(nosql.Recommend)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetRecommendNextID()
	db.CreatedTime = time.Now()
	db.Name = ""
	db.Type = tp
	db.Creator = operator
	db.Owner = owner
	db.Targets = list
	if db.Targets == nil {
		db.Targets = make([]string, 0, 1)
	}
	err := nosql.CreateRecommend(db)
	if err != nil {
		return nil, err
	}
	info := new(RecommendInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetRecommend(owner string, tp uint8) (*RecommendInfo, error) {
	db, err := nosql.GetRecommendBy(owner, tp)
	if err != nil {
		return nil, err
	}
	info := new(RecommendInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetRecommendBy(owner string) ([]*RecommendInfo, error) {
	if owner == "" {
		return nil, errors.New("the owner is empty")
	}
	dbs, err := nosql.GetRecommendByOwner(owner)
	if err != nil {
		return nil, err
	}
	list := make([]*RecommendInfo, 0, len(dbs))
	for _, db := range dbs {
		info := new(RecommendInfo)
		info.initInfo(db)
		list = append(list, info)
	}
	return list, nil
}

func (mine *cacheContext) RemoveRecommend(uid, operator string) error {
	return nosql.RemoveRecommend(uid, operator)
}

func (mine *RecommendInfo) initInfo(db *nosql.Recommend) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Name = db.Name
	mine.Created = db.Created
	mine.Updated = db.Updated
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Owner = db.Owner
	mine.Type = db.Type
	mine.Targets = db.Targets
}

func (mine *RecommendInfo) UpdateTargets(operator string, list []string) error {
	//if list == nil || len(list) < 1 {
	//	return errors.New("the target list is nil")
	//}
	err := nosql.UpdateRecommendTargets(mine.UID, operator, list)
	if err == nil {
		mine.Targets = list
	}
	return err
}
