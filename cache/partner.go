package cache

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.third/proxy/nosql"
	"omo.msa.third/tool"
	"time"
)

type PartnerInfo struct {
	Status uint8
	baseInfo
	Phone    string
	Remark   string
	Secret   string
	Tags     []string
}

func (mine *cacheContext)CreatePartner(name,remark, phone, creator string, tags []string) (*PartnerInfo, error) {
	db := new(nosql.Partner)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetPartnerNextID()
	db.CreatedTime = time.Now()
	db.Creator = creator
	db.Operator = creator
	db.Name = name
	db.Remark = remark
	db.Secret = ""
	db.Phone = phone
	db.Tags = tags
	db.Status = 0
	if db.Tags == nil {
		db.Tags = make([]string, 0, 1)
	}

	err := nosql.CreatePartner(db)
	if err == nil {
		info := new(PartnerInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext)RemovePartner(uid, operator string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	err := nosql.RemovePartner(uid, operator)
	return err
}

func (mine *cacheContext)GetAllPartners() []*PartnerInfo {
	list := make([]*PartnerInfo, 0, 10)
	array,err := nosql.GetPartners()
	if err == nil{
		for _, db := range array {
			info:= new(PartnerInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext)GetPartner(uid string) *PartnerInfo {
	if len(uid) < 1 {
		return nil
	}
	db,err := nosql.GetPartner(uid)
	if err == nil{
		info:= new(PartnerInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext)GetPartnerBySecret(uid string) *PartnerInfo {
	if len(uid) < 1 {
		return nil
	}
	db,err := nosql.GetPartnerBySecret(uid)
	if err == nil{
		info:= new(PartnerInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *PartnerInfo)initInfo(db *nosql.Partner)  {
	mine.UID = db.UID.Hex()
	mine.Remark = db.Remark
	mine.ID = db.ID
	mine.Name = db.Name
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Phone = db.Phone
	mine.Status = db.Status
	mine.Secret = db.Secret
	mine.Tags = db.Tags
}

func (mine *PartnerInfo)UpdateBase(name, remark,operator string) error {
	if len(name) <1 {
		name = mine.Name
	}
	if len(remark) <1 {
		remark = mine.Remark
	}
	err := nosql.UpdatePartnerBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *PartnerInfo)UpdateTags(operator string, tags []string) error {
	err := nosql.UpdatePartnerTags(mine.UID, operator, tags)
	if err == nil {
		mine.Tags = tags
		mine.Operator = operator
	}
	return err
}

func (mine *PartnerInfo)UpdateCover(cover, operator string) error {
	err := nosql.UpdatePartnerCover(mine.UID, cover, operator)
	if err == nil {
		mine.Operator = operator
	}
	return err
}

func (mine *PartnerInfo)CreateSecret(operator string) error {
	key := uuid.NewV4().String()
	md5 := tool.StrMD5(key)
	err := nosql.UpdatePartnerSecret(mine.UID, md5, operator)
	if err == nil {
		mine.Secret = md5
	}
	return err
}