package cache

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.third/proxy/nosql"
	"time"
)

const (
	SourceActivity SourceType = 0
	SourceEntity   SourceType = 1
	SourceArticle  SourceType = 2
	SourcePhoto    SourceType = 3
)

type SourceType uint8

type RecommendInfo struct {
	Type uint8
	baseInfo
	Owner   string //所属组织或者场景
	Quote   string //引用实体
	Targets []string
}

func (mine *cacheContext) CreateRecommend(owner, quote, operator string, tp uint8, list []string) (*RecommendInfo, error) {
	db := new(nosql.Recommend)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetRecommendNextID()
	db.CreatedTime = time.Now()
	db.Name = ""
	db.Type = tp
	db.Creator = operator
	db.Owner = owner
	db.Targets = list
	db.Quote = quote
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

func (mine *cacheContext) GetRecommendByUID(uid string) (*RecommendInfo, error) {
	db, err := nosql.GetRecommend(uid)
	if err != nil {
		return nil, err
	}
	info := new(RecommendInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetRecommendOwnerQuote(owner, quote string) ([]*RecommendInfo, error) {
	dbs, err := nosql.GetRecommendByOwnerQuote(owner, quote)
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

func (mine *cacheContext) GetRecommendOwnerTarget(owner, target string) ([]*RecommendInfo, error) {
	dbs, err := nosql.GetRecommendByOwnerTarget(owner, target)
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

func (mine *cacheContext) GetRecommendByOwner(owner string) ([]*RecommendInfo, error) {
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

func (mine *cacheContext) GetRecommendsByQuote(quote string) ([]*RecommendInfo, error) {
	if quote == "" {
		return nil, errors.New("the quote is empty")
	}
	dbs, err := nosql.GetRecommendsByQuote(quote)
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

func (mine *cacheContext) GetRecommendsByType(scene string, tp uint32) ([]*RecommendInfo, error) {
	if scene == "" {
		scene = DefaultScene
	}
	if tp > 100 {
		//按人气动态获取推荐列表： 1=书籍； 2=地点
		dbs, err := nosql.GetMotionsByTop(DefaultScene, 1, 50)
		if err != nil {
			return nil, err
		}
		list := make([]*RecommendInfo, 0, 1)
		info := new(RecommendInfo)
		info.UID = fmt.Sprintf("recommend_%s-%d", DefaultScene, 1)
		info.Type = uint8(tp)
		info.Owner = DefaultScene
		info.Targets = make([]string, len(dbs))
		for _, db := range dbs {
			info.Targets = append(info.Targets, db.Content)
		}
		info.Creator = DefaultScene
		return list, nil
	} else {
		//手动配置的推荐
		dbs, err := nosql.GetRecommendsByType(scene, tp)
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
	mine.Quote = db.Quote
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
