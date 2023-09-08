package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-third/proto/third"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.third/proxy/nosql"
	"time"
)

type HolidayInfo struct {
	Type uint8
	From int64
	End  int64
	baseInfo
	Owner  string
	Remark string
}

func (mine *cacheContext) CreateHoliday(in *pb.ReqHolidayAdd) (*HolidayInfo, error) {
	db := new(nosql.Holiday)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetHolidayNextID()
	db.Created = time.Now().Unix()
	db.Creator = in.Operator
	db.Operator = ""
	db.Name = in.Name
	db.Remark = in.Remark
	db.Owner = in.Owner
	db.From = in.From
	db.End = in.End
	db.Type = uint8(in.Type)

	err := nosql.CreateHoliday(db)
	if err == nil {
		info := new(HolidayInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) RemoveHoliday(uid, operator string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	err := nosql.RemoveHoliday(uid, operator)
	return err
}

func (mine *cacheContext) GetAllHolidays() []*HolidayInfo {
	list := make([]*HolidayInfo, 0, 10)
	array, err := nosql.GetHolidays()
	if err == nil {
		for _, db := range array {
			info := new(HolidayInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetHoliday(uid string) (*HolidayInfo, error) {
	if len(uid) < 1 {
		return nil, errors.New("the uid is empty")
	}
	db, err := nosql.GetHoliday(uid)
	if err == nil {
		info := new(HolidayInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) GetOfficialHoliday(owner string, from int64) *HolidayInfo {
	if len(owner) < 1 {
		return nil
	}
	db, err := nosql.GetHolidayByFrom(owner, from)
	if err != nil {
		return nil
	}
	info := new(HolidayInfo)
	info.initInfo(db)
	return info
}

func (mine *cacheContext) GetHolidayByOwner(uid string) []*HolidayInfo {
	if len(uid) < 1 {
		return nil
	}
	dbs, err := nosql.GetHolidaysByOwner(uid)
	list := make([]*HolidayInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			info := new(HolidayInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetThisYearHolidayByType(owner string, tp uint32) []*HolidayInfo {
	if len(owner) < 1 {
		owner = "system"
	}
	dbs, err := nosql.GetHolidaysByType(owner, tp)
	list := make([]*HolidayInfo, 0, len(dbs))
	if err == nil {
		now := time.Now()
		from := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(now.Year()+1, time.January, 1, 0, 0, 0, 0, time.UTC)
		for _, db := range dbs {
			if db.From > from.Unix() && db.From < end.Unix() {
				info := new(HolidayInfo)
				info.initInfo(db)
				list = append(list, info)
			}
		}
	}
	return list
}

func (mine *cacheContext) GetHolidayByYear(owner string, year int) []*HolidayInfo {
	if len(owner) < 1 {
		return nil
	}
	dbs, err := nosql.GetHolidaysByOwner(owner)
	list := make([]*HolidayInfo, 0, len(dbs))
	if err == nil {
		from := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC)
		for _, db := range dbs {
			if db.From > from.Unix() && db.From < end.Unix() {
				info := new(HolidayInfo)
				info.initInfo(db)
				list = append(list, info)
			}
		}
	}
	return list
}

func (mine *HolidayInfo) initInfo(db *nosql.Holiday) {
	mine.UID = db.UID.Hex()
	mine.Remark = db.Remark
	mine.ID = db.ID
	mine.Name = db.Name
	mine.Created = db.Created
	mine.Updated = db.Updated
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Owner = db.Owner
	mine.From = db.From
	mine.End = db.End
	mine.Type = db.Type
}

func (mine *HolidayInfo) UpdateBase(in *pb.ReqHolidayUpdate) error {
	if len(in.Name) < 1 {
		in.Name = mine.Name
	}
	if len(in.Remark) < 1 {
		in.Remark = mine.Remark
	}

	err := nosql.UpdateHolidayBase(mine.UID, in.Name, in.Remark, in.Operator, in.From, in.End)
	if err == nil {
		mine.Name = in.Name
		mine.Remark = in.Remark
		mine.Operator = in.Operator
		mine.From = in.From
		mine.End = in.End
		mine.Updated = time.Now().Unix()
	}
	return err
}
