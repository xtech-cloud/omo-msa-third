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

type HolidayService struct{}

func switchHoliday(info *cache.HolidayInfo) *pb.HolidayInfo {
	tmp := new(pb.HolidayInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Updated = info.Updated
	tmp.Created = info.Created
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator
	tmp.Type = uint32(info.Type)
	tmp.Owner = info.Owner
	tmp.From = info.From
	tmp.End = info.End
	return tmp
}

func (mine *HolidayService) AddOne(ctx context.Context, in *pb.ReqHolidayAdd, out *pb.ReplyHolidayInfo) error {
	path := "schedule.addOne"
	inLog(path, in)

	info, err := cache.Context().CreateHoliday(in)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchHoliday(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *HolidayService) UpdateBase(ctx context.Context, in *pb.ReqHolidayUpdate, out *pb.ReplyHolidayInfo) error {
	path := "schedule.updateBase"
	inLog(path, in)

	info, err := cache.Context().GetHoliday(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	err = info.UpdateBase(in)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchHoliday(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *HolidayService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyHolidayInfo) error {
	path := "schedule.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the schedule uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetHoliday(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchHoliday(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *HolidayService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "schedule.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the schedule uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveHoliday(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *HolidayService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyHolidayList) error {
	path := "schedule.getByFilter"
	inLog(path, in)
	var array []*cache.HolidayInfo
	//var max uint32 = 0
	//var pages uint32 =
	if in.Field == "" {
		array = cache.Context().GetHolidayByOwner(in.Scene)
	} else if in.Field == "type" {
		tp, err := strconv.Atoi(in.Value)
		if err != nil {
			out.Status = outError(path, err.Error(), pbstatus.ResultStatus_FormatError)
			return nil
		}
		array = cache.Context().GetThisYearHolidayByType(in.Scene, uint32(tp))
	} else if in.Field == "year" {
		year, err := strconv.Atoi(in.Value)
		if err != nil {
			out.Status = outError(path, err.Error(), pbstatus.ResultStatus_FormatError)
			return nil
		}
		array = cache.Context().GetHolidayByYear(in.Scene, year)
	}
	out.List = make([]*pb.HolidayInfo, 0, len(array))
	for _, info := range array {
		item := switchHoliday(info)
		out.List = append(out.List, item)
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *HolidayService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "schedule.getStatistic"
	inLog(path, in)
	if in.Field == "count" {

	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *HolidayService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "schedule.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the motion uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	//info, err1 := cache.Context().GetHoliday(in.Uid)
	//if err1 != nil {
	//	out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_NotExisted)
	//	return nil
	//}
	var err error
	if in.Field == "status" {

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
