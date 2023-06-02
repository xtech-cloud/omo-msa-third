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

type PartnerService struct{}

func switchPartner(info *cache.PartnerInfo) *pb.PartnerInfo {
	tmp := new(pb.PartnerInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Created = info.CreateTime.Unix()
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator
	tmp.Tags = info.Tags
	tmp.Phone = info.Phone
	tmp.Secret = info.Secret
	return tmp
}

func (mine *PartnerService) AddOne(ctx context.Context, in *pb.ReqPartnerAdd, out *pb.ReplyPartnerInfo) error {
	path := "partner.addOne"
	inLog(path, in)

	info, err := cache.Context().CreatePartner(in.Name, in.Remark, in.Phone, in.Operator, in.Tags)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchPartner(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *PartnerService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyPartnerInfo) error {
	path := "partner.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the partner uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetPartner(in.Uid)
	if info == nil {
		out.Status = outError(path, "the partner not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchPartner(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *PartnerService) GetBySecret(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyPartnerInfo) error {
	path := "partner.getBySecret"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the partner uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetPartnerBySecret(in.Uid)
	if info == nil {
		out.Status = outError(path, "the partner not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchPartner(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *PartnerService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "partner.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the Partner uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemovePartner(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *PartnerService) GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplyPartnerList) error {
	path := "partner.getList"
	inLog(path, in)
	array := cache.Context().GetAllPartners()
	out.List = make([]*pb.PartnerInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchPartner(val))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *PartnerService) UpdateOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyPartnerInfo) error {
	path := "partner.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the Partner uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetPartner(in.Uid)
	if info == nil {
		out.Status = outError(path, "the Partner not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	//var err error

	//if len(in.Name) > 0 || len(in.Remark) > 0 {
	//	err = info.UpdateBase(in.Name, in.Remark, in.Operator)
	//}

	//if err != nil {
	//	out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
	//	return nil
	//}
	out.Info = switchPartner(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *PartnerService) CreateSecret(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyPartnerSecret) error {
	path := "partner.createSecret"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the partner uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}

	info := cache.Context().GetPartner(in.Uid)
	if info == nil {
		out.Status = outError(path, "the partner not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.CreateSecret(in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Secret = info.Secret
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *PartnerService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyPartnerList) error {
	path := "partner.getByFilter"
	inLog(path, in)
	//var array []*cache.MotionInfo
	//var max uint32 = 0
	//var pages uint32 = 0

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *PartnerService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "partner.getStatistic"
	inLog(path, in)
	if in.Field == "count" {

	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *PartnerService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "partner.updateByFilter"
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
