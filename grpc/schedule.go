package grpc

import (
	"context"
	"errors"
	"fmt"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	pb "github.com/xtech-cloud/omo-msp-third/proto/third"
	"omo.msa.third/cache"
	"strconv"
)

type ScheduleService struct{}

func switchSchedule(info *cache.ScheduleInfo) *pb.ScheduleInfo {
	tmp := new(pb.ScheduleInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Updated = info.Updated
	tmp.Created = info.Created
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator

	tmp.Owner = info.Owner
	tmp.Status = uint32(info.Status)
	tmp.Quote = info.Quote
	tmp.Color = info.Color
	tmp.Type = uint32(info.Type)
	tmp.Ignore = uint32(info.Ignore)
	tmp.Date = &pb.DurationInfo{Begin: info.Date.Begin, End: info.Date.End}
	tmp.Time = &pb.DurationInfo{Begin: info.Time.Begin, End: info.Time.End}
	tmp.Targets = info.Targets
	tmp.Weekdays = info.Weekdays
	return tmp
}

func (mine *ScheduleService) AddOne(ctx context.Context, in *pb.ReqScheduleAdd, out *pb.ReplyScheduleInfo) error {
	path := "schedule.addOne"
	inLog(path, in)

	var err error
	var info *cache.ScheduleInfo
	if len(in.Uid) > 2 {
		info, err = cache.Context().GetSchedule(in.Uid)
		if err == nil {
			err = info.UpdateBase(in)
		}
	} else {
		info, err = cache.Context().CreateSchedule(in)
	}

	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchSchedule(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ScheduleService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyScheduleInfo) error {
	path := "schedule.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the schedule uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetSchedule(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchSchedule(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ScheduleService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "schedule.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the schedule uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveSchedule(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *ScheduleService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyScheduleList) error {
	path := "schedule.getByFilter"
	inLog(path, in)
	var array []*cache.ScheduleInfo
	//var max uint32 = 0
	//var pages uint32 = 0
	if in.Field == "" {
		array = cache.Context().GetScheduleByOwner(in.Scene)
	} else if in.Field == "now" {
		info := cache.Context().GetNowSchedule(in.Scene, in.Value)
		array = make([]*cache.ScheduleInfo, 0, 1)
		array = append(array, info)
	} else if in.Field == "today" {
		array = cache.Context().GetTodaySchedules(in.Scene)
	}
	out.List = make([]*pb.ScheduleInfo, 0, len(array))
	for _, info := range array {
		item := switchSchedule(info)
		out.List = append(out.List, item)
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ScheduleService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "schedule.getStatistic"
	inLog(path, in)
	if in.Field == "count" {

	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *ScheduleService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "schedule.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the motion uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err1 := cache.Context().GetSchedule(in.Uid)
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Field == "status" {
		num, er := strconv.ParseInt(in.Value, 10, 32)
		if er == nil {
			err = info.UpdateStatus(in.Operator, uint8(num))
		} else {
			err = er
		}
	} else {
		err = errors.New("the field not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}
