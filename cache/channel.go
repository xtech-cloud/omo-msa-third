package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.third/proxy/nosql"
	"time"
)

const (
	OwnerTypePerson = 1
	OwnerTypeUnit = 0
)

type ChannelInfo struct {
	UID string
}

func (mine *cacheContext)GetChannel(owner string) (*ChannelInfo,error) {
	db,err := nosql.GetChannelByOwner(owner)
	if err != nil {
		return nil, err
	}
	info := new(ChannelInfo)
	info.initInfo(db)
	return info,nil
}

func (mine *cacheContext)GetAllChannel() []*ChannelInfo {
	list := make([]*ChannelInfo, 0, 10)
	return list
}

func (mine *ChannelInfo) initInfo(db *nosql.Channel) {

}

func (mine *ChannelInfo)createChannel(uid string) (*nosql.Channel,error) {
	db := new(nosql.Channel)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetChannelNextID()
	db.CreatedTime = time.Now()
	db.Owner = uid
	err := nosql.CreateChannel(db)
	if err != nil {
		return nil, err
	}else{
		return db, nil
	}
}

func (mine *ChannelInfo)HadBag(uid string) bool {

	return false
}

func (mine *ChannelInfo) AppendAsset(uid string) error {
	if len(uid) < 1 {
		return errors.New("the asset uid is empty")
	}
	if mine.HadBag(uid) {
		return nil
	}
	er := nosql.AppendChannelBag(mine.UID, uid)
	return er
}

func (mine *ChannelInfo) SubtractAsset(uid string) error {
	if !mine.HadBag(uid) {
		return nil
	}
	er := nosql.SubtractChannelBag(mine.UID, uid)
	return er
}

func (mine *ChannelInfo)UpdateBags(list []string, operator string) error {
	er := nosql.UpdateChannelBags(mine.UID, operator, list)
	return er
}