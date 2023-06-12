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

type MotionService struct{}

func switchMotion(info *cache.MotionInfo) *pb.MotionInfo {
	tmp := new(pb.MotionInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Name = info.Name
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Created = info.CreateTime.Unix()
	tmp.Creator = info.Creator
	tmp.Count = info.Count
	tmp.Sn = info.SN
	tmp.Content = info.Content
	tmp.Event = info.EventID
	tmp.App = info.AppID
	tmp.Scene = info.Scene
	return tmp
}

func (mine *MotionService) AddOne(ctx context.Context, in *pb.ReqMotionAdd, out *pb.ReplyMotionInfo) error {
	path := "motion.addOne"
	inLog(path, in)

	info, err := cache.Context().CreateMotion(in.Scene, in.App, in.Sn, in.Event, in.Content, in.Operator, in.Count)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchMotion(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *MotionService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyMotionInfo) error {
	path := "motion.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the motion uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetMotion(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchMotion(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *MotionService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "motion.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the motion uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveMotion(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *MotionService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyMotionList) error {
	path := "motion.getByFilter"
	inLog(path, in)
	var array []*cache.MotionInfo
	var max uint32 = 0
	var pages uint32 = 0

	if in.Field == "sn" {
		array = cache.Context().GetMotionsBySN(in.Scene, in.Value)
	} else if in.Field == "event" {
		array = cache.Context().GetMotionsByEvent(in.Scene, in.Value)
	} else if in.Field == "association" {
		if len(in.List) > 2 {
			item := cache.Context().GetMotionBy(in.Scene, in.List[0], in.List[1], in.List[2])
			array = make([]*cache.MotionInfo, 0, 1)
			array = append(array, item)
		}
	} else if in.Field == "content" {
		if len(in.List) > 1 {
			array = cache.Context().GetMotionsByContent(in.Scene, in.List[0], in.List[1])
		}
	} else {

	}
	out.List = make([]*pb.MotionInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchMotion(val))
	}
	out.Total = max
	out.Pages = pages
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *MotionService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "motion.getStatistic"
	inLog(path, in)
	if in.Field == "count" {
		if len(in.List) > 2 {
			item := cache.Context().GetMotionBy(in.Scene, in.List[0], in.List[1], in.List[2])
			out.Count = item.Count
		}
	} else if in.Field == "content" {
		for _, eve := range in.List {
			array := cache.Context().GetMotionsByEveContent(in.Scene, eve, in.Value)
			for _, info := range array {
				out.Count += info.Count
			}
		}
	} else if in.Field == "events" {
		for _, eve := range in.List {
			array := cache.Context().GetMotionsBySNEvent(in.Scene, in.Value, eve)
			for _, info := range array {
				out.Count += info.Count
			}
		}
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *MotionService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "motion.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the motion uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err1 := cache.Context().GetMotion(in.Uid)
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Field == "offset" {
		num, er := strconv.ParseInt(in.Value, 10, 32)
		if er == nil {
			err = info.AddCount(uint32(num), in.Operator)
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
