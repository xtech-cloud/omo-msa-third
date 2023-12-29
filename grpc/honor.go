package grpc

import (
	"context"
	"errors"
	"fmt"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	pb "github.com/xtech-cloud/omo-msp-third/proto/third"
	"omo.msa.third/cache"
	"omo.msa.third/proxy"
)

type HonorService struct{}

func switchHonor(info *cache.HonorInfo) *pb.HonorInfo {
	tmp := new(pb.HonorInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Name = info.Name
	tmp.Updated = info.Updated
	tmp.Created = info.Created
	tmp.Creator = info.Creator
	tmp.Remark = info.Remark
	tmp.Scene = info.Scene
	tmp.Parent = info.Parent
	tmp.Style = info.Style
	tmp.Type = uint32(info.Type)
	tmp.Status = uint32(info.Status)
	tmp.Contents = make([]*pb.HonorContent, 0, len(info.Contents))
	for _, content := range info.Contents {
		tmp.Contents = append(tmp.Contents, &pb.HonorContent{
			Name:   content.Name,
			Remark: content.Remark,
			Quotes: content.Quotes,
		})
	}
	return tmp
}

func (mine *HonorService) AddOne(ctx context.Context, in *pb.ReqHonorAdd, out *pb.ReplyHonorInfo) error {
	path := "honor.addOne"
	inLog(path, in)

	info, err := cache.Context().CreateHonor(in.Scene, in.Parent, in.Name, in.Remark, in.Style, in.Operator, in.Type, in.Status)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchHonor(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *HonorService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyHonorInfo) error {
	path := "honor.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the honor uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetHonor(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchHonor(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *HonorService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "honor.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the honor uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveHonor(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *HonorService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyHonorList) error {
	path := "honor.getByFilter"
	inLog(path, in)
	var array []*cache.HonorInfo
	var max uint32 = 0
	var pages uint32 = 0

	if in.Field == "scene" || in.Field == "" {
		array = cache.Context().GetHonorsByScene(in.Scene)
	} else if in.Field == "parent" {
		array = cache.Context().GetHonorsByParent(in.Value)
	}
	out.List = make([]*pb.HonorInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchHonor(val))
	}
	out.Total = max
	out.Pages = pages
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *HonorService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "honor.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *HonorService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "honor.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the honor uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	_, err1 := cache.Context().GetHonor(in.Uid)
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Field == "title" {

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

func (mine *HonorService) UpdateBase(ctx context.Context, in *pb.ReqHonorUpdate, out *pb.ReplyHonorInfo) error {
	path := "honor.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the honor uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err1 := cache.Context().GetHonor(in.Uid)
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateBase(in.Name, in.Remark, in.Style, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *HonorService) UpdateStatus(ctx context.Context, in *pb.RequestFlag, out *pb.ReplyInfo) error {
	path := "honor.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the honor uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetHonor(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err = info.UpdateStatus(in.Operator, in.Flag)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *HonorService) UpdateContents(ctx context.Context, in *pb.ReqHonorContents, out *pb.ReplyInfo) error {
	path := "honor.updateContents"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the honor uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetHonor(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	arr := make([]*proxy.HonorContent, 0, len(in.List))
	for _, content := range in.List {
		arr = append(arr, &proxy.HonorContent{
			Name:   content.Name,
			Remark: content.Remark,
			Quotes: content.Quotes,
		})
	}
	err = info.UpdateContents(in.Operator, arr)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}
