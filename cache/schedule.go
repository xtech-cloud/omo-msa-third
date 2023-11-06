package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-third/proto/third"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.third/proxy"
	"omo.msa.third/proxy/nosql"
	"omo.msa.third/tool"
	"time"
)

const OneDay int64 = 24 * 3600

const (
	ScheduleOpen  = 0
	ScheduleClose = 1
)

type ScheduleInfo struct {
	Status uint8
	Type   uint8
	Ignore uint8
	baseInfo
	Owner    string
	Remark   string
	Quote    string
	Color    string
	Date     proxy.DurationInfo
	Time     proxy.DurationInfo
	Targets  []string
	Tags     []string
	Weekdays []uint32
}

func (mine *cacheContext) CreateSchedule(in *pb.ReqScheduleAdd) (*ScheduleInfo, error) {
	db := new(nosql.Schedule)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetScheduleNextID()
	db.Created = time.Now().Unix()
	db.Creator = in.Operator
	db.Operator = ""
	db.Name = in.Name
	db.Remark = in.Remark
	db.Owner = in.Owner
	db.Type = uint8(in.Type)
	db.Status = 0
	db.Quote = in.Quote
	db.Color = in.Color
	db.Ignore = uint8(in.Ignore)
	from := in.Date.Begin
	end := in.Date.End
	db.Date = proxy.DurationInfo{Begin: from, End: end}
	db.Time = proxy.DurationInfo{Begin: in.Time.Begin, End: in.Time.End}
	db.Targets = in.Targets
	db.Weekdays = in.Weekdays
	//db.Tags = tags
	//if db.Tags == nil {
	//	db.Tags = make([]string, 0, 1)
	//}

	err := nosql.CreateSchedule(db)
	if err == nil {
		info := new(ScheduleInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) RemoveSchedule(uid, operator string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	err := nosql.RemoveSchedule(uid, operator)
	return err
}

func (mine *cacheContext) GetAllSchedules() []*ScheduleInfo {
	list := make([]*ScheduleInfo, 0, 10)
	array, err := nosql.GetSchedules()
	if err == nil {
		for _, db := range array {
			info := new(ScheduleInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetSchedule(uid string) (*ScheduleInfo, error) {
	if len(uid) < 1 {
		return nil, errors.New("the uid is empty")
	}
	db, err := nosql.GetSchedule(uid)
	if err == nil {
		info := new(ScheduleInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) GetNowSchedule(owner, area string) *ScheduleInfo {
	dbs, err := nosql.GetSchedulesByOwner(owner)
	if err == nil {
		now := time.Now()
		utc := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()
		t := int64(now.Hour()*3600 + now.Minute()*60 + now.Second())
		for _, db := range dbs {
			if tool.HasItem(db.Targets, area) && utc >= db.Date.Begin && utc <= db.Date.End {
				if t >= db.Time.Begin && t <= db.Time.End {
					info := new(ScheduleInfo)
					info.initInfo(db)
					return info
				}
			}
		}
	}
	return nil
}

func (mine *cacheContext) GetTodaySchedules(owner string) []*ScheduleInfo {
	now := time.Now()
	utc := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()
	dbs, err := nosql.GetSchedulesByOwnerDate(owner, utc)
	list := make([]*ScheduleInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			info := new(ScheduleInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetScheduleByOwner(uid string) []*ScheduleInfo {
	if len(uid) < 1 {
		return nil
	}
	dbs, err := nosql.GetSchedulesByOwner(uid)
	list := make([]*ScheduleInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			info := new(ScheduleInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *ScheduleInfo) initInfo(db *nosql.Schedule) {
	mine.UID = db.UID.Hex()
	mine.Remark = db.Remark
	mine.ID = db.ID
	mine.Name = db.Name
	mine.Created = db.Created
	mine.Updated = db.Updated
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Owner = db.Owner
	mine.Status = db.Status
	mine.Type = db.Type
	mine.Quote = db.Quote
	mine.Ignore = db.Ignore
	mine.Color = db.Color
	mine.Date = db.Date
	mine.Time = db.Time
	mine.Targets = db.Targets
	mine.Weekdays = db.Weekdays
}

func (mine *ScheduleInfo) UpdateBase(in *pb.ReqScheduleAdd) error {
	if len(in.Name) < 1 {
		in.Name = mine.Name
	}
	if len(in.Remark) < 1 {
		in.Remark = mine.Remark
	}
	date := proxy.DurationInfo{Begin: in.Date.Begin, End: in.Date.End}
	tim := proxy.DurationInfo{Begin: in.Time.Begin, End: in.Time.End}
	err := nosql.UpdateScheduleBase(mine.UID, in.Name, in.Remark, in.Operator, in.Quote, in.Color, uint8(in.Type), uint8(in.Ignore), date, tim, in.Targets)
	if err == nil {
		mine.Name = in.Name
		mine.Remark = in.Remark
		mine.Operator = in.Operator
		mine.Quote = in.Quote
		mine.Color = in.Color
		mine.Type = uint8(in.Type)
		mine.Ignore = uint8(in.Ignore)
		mine.Date = date
		mine.Time = tim
		mine.Targets = in.Targets
		mine.Updated = time.Now().Unix()
	}
	return err
}

func (mine *ScheduleInfo) UpdateTags(operator string, tags []string) error {
	err := nosql.UpdateScheduleTags(mine.UID, operator, tags)
	if err == nil {
		mine.Tags = tags
		mine.Operator = operator
		mine.Updated = time.Now().Unix()
	}
	return err
}

func (mine *ScheduleInfo) UpdateStatus(operator string, st uint8) error {
	err := nosql.UpdateScheduleStatus(mine.UID, operator, st)
	if err == nil {
		mine.Status = st
		mine.Operator = operator
		mine.Updated = time.Now().Unix()
	}
	return err
}
