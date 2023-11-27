package grpc

import (
	"context"
	"fmt"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	pb "github.com/xtech-cloud/omo-msp-third/proto/third"
	"omo.msa.third/cache"
	"strconv"
)

type RecommendService struct{}

func switchRecommend(info *cache.RecommendInfo) *pb.RecommendInfo {
	tmp := new(pb.RecommendInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Name = info.Name
	tmp.Updated = info.Updated
	tmp.Created = info.Created
	tmp.Creator = info.Creator
	tmp.Type = uint32(info.Type)
	tmp.Targets = info.Targets
	tmp.Quote = info.Quote
	tmp.Owner = info.Owner
	return tmp
}

func (mine *RecommendService) AddOne(ctx context.Context, in *pb.ReqRecommendAdd, out *pb.ReplyRecommendInfo) error {
	path := "recommend.addOne"
	inLog(path, in)

	info, err := cache.Context().CreateRecommend(in.Owner, in.Quote, in.Operator, uint8(in.Type), in.Targets)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchRecommend(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *RecommendService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyRecommendInfo) error {
	path := "recommend.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the recommend uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetRecommend(in.Uid, uint8(in.Type))
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchRecommend(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *RecommendService) UpdateOne(ctx context.Context, in *pb.ReqRecommendUpdate, out *pb.ReplyRecommendInfo) error {
	path := "recommend.updateOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the recommend uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().GetRecommendByUID(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err = info.UpdateTargets(in.Operator, in.Targets)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchRecommend(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *RecommendService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "recommend.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the recommend uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveRecommend(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *RecommendService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyRecommendList) error {
	path := "recommend.getByFilter"
	inLog(path, in)
	var array []*cache.RecommendInfo
	var max uint32 = 0
	var pages uint32 = 0
	var err error
	if in.Field == "quote" {
		array, err = cache.Context().GetRecommendsByQuote(in.Value)
	} else if in.Field == "type" {
		tp, er := strconv.Atoi(in.Value)
		if er != nil {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
		array, err = cache.Context().GetRecommendsByType(in.Scene, uint32(tp))
	} else if in.Field == "owner_quote" {
		array, err = cache.Context().GetRecommendOwnerQuote(in.Scene, in.Value)
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.RecommendInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchRecommend(val))
	}
	out.Total = uint64(max)
	out.PageMax = pages
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *RecommendService) GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplyRecommendList) error {
	path := "recommend.getList"
	inLog(path, in)
	var array []*cache.RecommendInfo
	var max uint32 = 0
	var pages uint32 = 0
	array, err := cache.Context().GetRecommendByOwner(in.Owner)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.RecommendInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchRecommend(val))
	}
	out.Total = uint64(max)
	out.PageMax = pages
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *RecommendService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "recommend.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *RecommendService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "recommend.updateByFilter"
	inLog(path, in)
	if len(in.Scene) < 1 {
		out.Status = outError(path, "the recommend uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	tp, er := strconv.ParseUint(in.Value, 10, 32)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_FormatError)
		return nil
	}
	_, err1 := cache.Context().GetRecommend(in.Scene, uint8(tp))
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error

	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}
