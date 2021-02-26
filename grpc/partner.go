package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-third/proto/third"
	"omo.msa.third/cache"
)

type PartnerService struct {}

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

func (mine *PartnerService)AddOne(ctx context.Context, in *pb.ReqPartnerAdd, out *pb.ReplyPartnerInfo) error {
	path := "partner.addOne"
	inLog(path, in)

	info,err := cache.Context().CreatePartner(in.Name, in.Remark, in.Phone, in.Operator, in.Tags)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
		return nil
	}
	out.Info = switchPartner(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *PartnerService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyPartnerInfo) error {
	path := "partner.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the partner uid is empty", pb.ResultCode_Empty)
		return nil
	}
	info := cache.Context().GetPartner(in.Uid)
	if info == nil {
		out.Status = outError(path,"the partner not found", pb.ResultCode_NotExisted)
		return nil
	}
	out.Info = switchPartner(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *PartnerService)GetBySecret(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyPartnerInfo) error {
	path := "partner.getBySecret"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the partner uid is empty", pb.ResultCode_Empty)
		return nil
	}
	info := cache.Context().GetPartnerBySecret(in.Uid)
	if info == nil {
		out.Status = outError(path,"the partner not found", pb.ResultCode_NotExisted)
		return nil
	}
	out.Info = switchPartner(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *PartnerService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "partner.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the Partner uid is empty", pb.ResultCode_Empty)
		return nil
	}
	err := cache.Context().RemovePartner(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *PartnerService)GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplyPartnerList) error {
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

func (mine *PartnerService)UpdateOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyPartnerInfo) error {
	path := "partner.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the Partner uid is empty", pb.ResultCode_Empty)
		return nil
	}
	info := cache.Context().GetPartner(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Partner not found", pb.ResultCode_NotExisted)
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


func (mine *PartnerService)CreateSecret(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyPartnerSecret) error {
	path := "partner.createSecret"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the partner uid is empty", pb.ResultCode_Empty)
		return nil
	}

	info := cache.Context().GetPartner(in.Uid)
	if info == nil {
		out.Status = outError(path,"the partner not found", pb.ResultCode_NotExisted)
		return nil
	}
	err := info.CreateSecret(in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultCode_DBException)
		return nil
	}
	out.Secret = info.Secret
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}
