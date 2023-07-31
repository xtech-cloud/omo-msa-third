package grpc

import (
	"context"
	"errors"
	"fmt"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	pb "github.com/xtech-cloud/omo-msp-third/proto/third"
	"omo.msa.third/cache"
	"omo.msa.third/proxy"
	"omo.msa.third/proxy/nosql"
)

type TopicService struct{}

func switchTopic(info *cache.TopicInfo) *pb.TopicInfo {
	tmp := new(pb.TopicInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Name = info.Name
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Created = info.CreateTime.Unix()
	tmp.Creator = info.Creator
	tmp.Remark = info.Remark
	tmp.Time = info.Time
	tmp.Compere = info.Compere
	tmp.Scene = info.Scene
	return tmp
}

func switchTopicRecord(info *nosql.TopicRecord) *pb.TopicRecord {
	tmp := new(pb.TopicRecord)
	tmp.Uid = info.UID.Hex()
	tmp.Id = info.ID
	tmp.Created = info.CreatedTime.Unix()
	tmp.Operator = info.Creator
	tmp.Topic = info.Topic
	tmp.Date = info.Date
	tmp.Compere = info.Compere
	tmp.Scene = info.Scene
	tmp.Sn = info.SN
	tmp.State = info.State
	tmp.Results = make([]*pb.TopicResult, 0, 3)
	for _, result := range info.Results {
		tmp.Results = append(tmp.Results, &pb.TopicResult{
			State:   result.State,
			Percent: result.Percent,
			Count:   result.Count,
		})
	}
	return tmp
}

func (mine *TopicService) AddOne(ctx context.Context, in *pb.ReqTopicAdd, out *pb.ReplyTopicInfo) error {
	path := "topic.addOne"
	inLog(path, in)

	info, err := cache.Context().CreateTopic(in.Scene, in.Title, in.Remark, in.Compere, in.Operator, in.Time)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchTopic(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TopicService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyTopicInfo) error {
	path := "topic.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the topic uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetTopic(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchTopic(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TopicService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "topic.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the topic uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveTopic(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *TopicService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyTopicList) error {
	path := "topic.getByFilter"
	inLog(path, in)
	var array []*cache.TopicInfo
	var max uint32 = 0
	var pages uint32 = 0

	if in.Field == "scene" || in.Field == "" {
		array = cache.Context().GetTopicsByScene(in.Scene)
	} else {

	}
	out.List = make([]*pb.TopicInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchTopic(val))
	}
	out.Total = max
	out.Pages = pages
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *TopicService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "topic.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *TopicService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "topic.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the topic uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err1 := cache.Context().GetTopic(in.Uid)
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Field == "title" {
		err = info.UpdateTitle(in.Value, in.Operator)
	} else if in.Field == "base" {
		if len(in.List) > 2 {
			err = info.UpdateBase(in.List[0], in.List[1], in.List[2], in.Operator, 0)
		} else {
			err = errors.New("the list size is error that field base")
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

func (mine *TopicService) AddRecord(ctx context.Context, in *pb.ReqTopicRecord, out *pb.ReplyInfo) error {
	path := "topic.addRecord"
	inLog(path, in)
	if len(in.Topic) < 1 {
		out.Status = outError(path, "the topic uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err1 := cache.Context().GetTopic(in.Topic)
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	results := make([]*proxy.TopicResult, 0, 3)
	for _, result := range in.Results {
		results = append(results, &proxy.TopicResult{State: result.State, Percent: result.Percent, Count: result.Count})
	}
	err := info.CreateRecord(in.Sn, in.Operator, in.State, in.Date, results)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *TopicService) GetRecords(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyTopicRecords) error {
	path := "topic.getRecords"
	inLog(path, in)
	var array []*nosql.TopicRecord
	var max uint32 = 0
	var pages uint32 = 0
	var err error
	if in.Field == "scene" || in.Field == "" {
		array, err = cache.Context().GetTopicRecordsByScene(in.Scene)
	} else if in.Field == "sn" {
		array, err = cache.Context().GetTopicRecordsBySN(in.Value)
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.List = make([]*pb.TopicRecord, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchTopicRecord(val))
	}
	out.Total = max
	out.Pages = pages
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}
