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

type ChannelService struct{}

func switchChannel(info *cache.ChannelInfo) *pb.ChannelInfo {
	tmp := new(pb.ChannelInfo)

	return tmp
}

func (mine *ChannelService) AddOne(ctx context.Context, in *pb.ReqChannelAdd, out *pb.ReplyChannelInfo) error {
	path := "Channel.addOne"
	inLog(path, in)

	info := new(cache.ChannelInfo)
	//
	//err := cache.Context().CreateChannel(info)
	//if err != nil {
	//	out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
	//	return nil
	//}
	out.Info = switchChannel(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ChannelService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyChannelInfo) error {
	path := "Channel.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the Channel uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	//info := cache.Context().GetChannel(in.Uid)
	//if info == nil {
	//	out.Status = outError(path,"the Channel not found", pb.ResultCode_NotExisted)
	//	return nil
	//}
	//out.Info = switchChannel(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ChannelService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "Channel.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the Channel uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	//err := cache.Context().RemoveChannel(in.Uid, in.Operator)
	//if err != nil {
	//	out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
	//	return nil
	//}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *ChannelService) GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplyChannelList) error {
	path := "channel.getList"
	inLog(path, in)
	array := cache.Context().GetAllChannel()
	out.List = make([]*pb.ChannelInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchChannel(val))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ChannelService) UpdateOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyChannelInfo) error {
	path := "channel.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the Channel uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	//info := cache.Context().GetChannel(in.Uid)
	//if info == nil {
	//	out.Status = outError(path,"the Channel not found", pb.ResultCode_NotExisted)
	//	return nil
	//}
	//var err error
	//
	//if len(in.Name) > 0 || len(in.Remark) > 0 {
	//	err = info.UpdateBase(in.Name, in.Remark, in.Operator)
	//}
	//
	//if err != nil {
	//	out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
	//	return nil
	//}
	//out.Info = switchChannel(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ChannelService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyChannelList) error {
	path := "channel.getByFilter"
	inLog(path, in)
	//var array []*cache.MotionInfo
	//var max uint32 = 0
	//var pages uint32 = 0
	//
	//out.Total = max
	//out = pages
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ChannelService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "channel.getStatistic"
	inLog(path, in)
	if in.Field == "count" {

	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *ChannelService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "channel.updateByFilter"
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
	if in.Field == "count" {
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
