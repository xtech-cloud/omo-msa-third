package grpc

import (
	"context"
	"fmt"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	pb "github.com/xtech-cloud/omo-msp-third/proto/third"
	"omo.msa.third/cache"
	"strconv"
)

type CarouselService struct{}

func switchCarouse(info *cache.CarouselInfo) *pb.CarouselInfo {
	tmp := new(pb.CarouselInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Name = info.Name
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Created = info.CreateTime.Unix()
	tmp.Creator = info.Creator
	tmp.Owner = info.Owner
	tmp.Quotes = make([]*pb.QuoteInfo, 0, len(info.Quotes))
	for _, quote := range info.Quotes {
		tmp.Quotes = append(tmp.Quotes, &pb.QuoteInfo{Uid: quote.UID,
			Asset: quote.Asset,
			Type:  uint32(quote.Type),
			Title: quote.Title, Updated: quote.Updated})
	}
	return tmp
}

func (mine *CarouselService) AddOne(ctx context.Context, in *pb.ReqCarouselAdd, out *pb.ReplyCarouselInfo) error {
	path := "carousel.addOne"
	inLog(path, in)

	info, err := cache.Context().CreateCarousel(in.Owner, in.Operator, in.Quotes)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchCarouse(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *CarouselService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyCarouselInfo) error {
	path := "carousel.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the carouse uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetCarousel(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchCarouse(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *CarouselService) UpdateOne(ctx context.Context, in *pb.ReqCarouselUpdate, out *pb.ReplyCarouselInfo) error {
	path := "carousel.updateOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the carouse uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetCarousel(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err = info.UpdateQuotes(in.Operator, in.Quotes)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchCarouse(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *CarouselService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "carousel.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the carousel uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveCarousel(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *CarouselService) GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplyCarouselList) error {
	path := "carousel.getList"
	inLog(path, in)
	var array []*cache.CarouselInfo
	var max uint32 = 0
	var pages uint32 = 0

	out.List = make([]*pb.CarouselInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchCarouse(val))
	}
	out.Total = uint64(max)
	out.PageMax = pages
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *CarouselService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyCarouselList) error {
	path := "carousel.getByFilter"
	inLog(path, in)
	var array []*cache.CarouselInfo
	var max uint32 = 0
	var pages uint32 = 0

	out.List = make([]*pb.CarouselInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchCarouse(val))
	}
	out.Total = uint64(max)
	out.PageMax = pages
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *CarouselService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "carousel.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *CarouselService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "carousel.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the carouse uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err1 := cache.Context().GetCarousel(in.Scene)
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Field == "append" {
		tp, er := strconv.ParseUint(in.Value, 10, 32)
		if er != nil {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_FormatError)
			return nil
		}
		err = info.AppendQuote(in.List[0], in.List[1], in.List[2], uint8(tp))
	} else if in.Field == "subtract" {
		err = info.SubtractQuote(in.Value)
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}
