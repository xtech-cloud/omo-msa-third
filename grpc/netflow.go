package grpc

import (
	"context"
	"fmt"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	pb "github.com/xtech-cloud/omo-msp-third/proto/third"
	"omo.msa.third/cache"
)

type NetflowService struct{}

func switchNetflow(info *cache.NetflowInfo) *pb.NetflowInfo {
	tmp := new(pb.NetflowInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Name = info.Name
	tmp.Created = info.Created
	tmp.Creator = info.Creator

	tmp.Scene = info.Scene
	tmp.Quote = info.Quote
	tmp.Size = info.Size
	tmp.Type = uint32(info.Type)
	tmp.Template = info.Template
	tmp.Target = info.Target
	tmp.Contents = make([]*pb.ContentInfo, 0, len(info.Contents))
	for _, content := range info.Contents {
		con := new(pb.ContentInfo)
		con.Group = content.Group
		con.Uid = content.UID
		con.Type = uint32(content.Type)
		con.Size = content.Size
		con.Children = content.Children
		tmp.Contents = append(tmp.Contents, con)
	}
	return tmp
}

func (mine *NetflowService) AddOne(ctx context.Context, in *pb.ReqNetflowAdd, out *pb.ReplyNetflowInfo) error {
	path := "netflow.addOne"
	inLog(path, in)

	var err error
	var info *cache.NetflowInfo
	info, err = cache.Context().CreateNetflow(in)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchNetflow(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *NetflowService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyNetflowInfo) error {
	path := "netflow.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the netflow uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetNetflow(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchNetflow(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *NetflowService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "netflow.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the netflow uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *NetflowService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyNetflowList) error {
	path := "netflow.getByFilter"
	inLog(path, in)
	var array []*cache.NetflowInfo
	if in.Field == "latest" {
		array, _ = cache.Context().GetNetflowByLatest(in.Scene, int(in.Number))
	} else {
		array, _ = cache.Context().GetNetflowByScene(in.Scene)
	}
	out.List = make([]*pb.NetflowInfo, 0, len(array))
	for _, info := range array {
		item := switchNetflow(info)
		out.List = append(out.List, item)
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *NetflowService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "netflow.getStatistic"
	inLog(path, in)
	if in.Field == "size" {
		array, _ := cache.Context().GetNetflowByScene(in.Scene)
		var size uint64 = 0
		for _, info := range array {
			size += info.Size
		}
		out.Count = size
	} else if in.Field == "stamps" {
		out.Count, out.List = cache.Context().GetNetflowByStamps(in.Scene, in.List)
	} else if in.Field == "today" {
		out.Count, out.List = cache.Context().GetNetflowByStamps(in.Scene, in.List)
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *NetflowService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "netflow.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the motion uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}
