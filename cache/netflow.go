package cache

import (
	pb "github.com/xtech-cloud/omo-msp-third/proto/third"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.third/proxy"
	"omo.msa.third/proxy/nosql"
	"sort"
	"strconv"
	"time"
)

type NetflowInfo struct {
	Type uint8
	Size uint64
	baseInfo
	Scene    string
	Quote    string
	Template string
	Target   string
	Contents []*proxy.ContentInfo
}

func (mine *cacheContext) CreateNetflow(in *pb.ReqNetflowAdd) (*NetflowInfo, error) {
	db := new(nosql.Netflow)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetNetflowNextID()
	db.Created = time.Now().Unix()
	db.Creator = in.Operator
	db.Name = ""
	db.Remark = ""
	db.Owner = in.Scene
	db.Type = uint8(in.Type)
	db.Template = in.Template
	db.Target = in.Target
	db.Size = 0
	db.Contents = make([]*proxy.ContentInfo, 0, len(in.Contents))
	for _, content := range in.Contents {
		tmp := new(proxy.ContentInfo)
		tmp.Type = uint8(content.Type)
		tmp.UID = content.Uid
		tmp.Size = content.Size
		tmp.Group = content.Group
		tmp.Children = content.Children
		db.Size += content.Size
		db.Contents = append(db.Contents, tmp)
	}
	err := nosql.CreateNetflow(db)
	if err == nil {
		info := new(NetflowInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) GetNetflow(uid string) (*NetflowInfo, error) {
	db, err := nosql.GetNetflow(uid)
	if err == nil {
		info := new(NetflowInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) GetNetflowByScene(uid string) ([]*NetflowInfo, error) {
	dbs, err := nosql.GetNetflowsByOwner(uid)
	list := make([]*NetflowInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			info := new(NetflowInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list, err
}

func (mine *cacheContext) GetNetflowByLatest(uid string, num int) ([]*NetflowInfo, error) {
	dbs, err := nosql.GetNetflowsByOwner(uid)
	list := make([]*NetflowInfo, 0, int(num))
	if err == nil {
		sort.Slice(dbs, func(i, j int) bool {
			return dbs[i].Created > dbs[j].Created
		})
		var arr []*nosql.Netflow
		if len(dbs) > num {
			arr = dbs[:num-1]
		} else {
			arr = dbs
		}
		for _, db := range arr {
			info := new(NetflowInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list, err
}

func (mine *cacheContext) GetNetflowByStamp(scene, stamp string) (uint64, *pb.PairInfo) {
	utc, er := strconv.Atoi(stamp)
	if er != nil {
		return 0, nil
	}
	day := time.Unix(int64(utc), 0)
	from := time.Date(day.Year(), day.Month(), day.Day(), 1, 0, 0, 0, time.UTC)
	to := time.Date(day.Year(), day.Month(), day.Day(), 24, 0, 0, 0, time.UTC)
	list, _ := nosql.GetNetflowsByDuration(scene, from.Unix(), to.Unix())
	num := uint64(len(list))
	var size uint64 = 0
	for _, netflow := range list {
		size += netflow.Size
	}
	return num, &pb.PairInfo{Id: uint64(utc), Key: stamp, Count: size}
}

func (mine *cacheContext) GetNetflowByStamps(scene string, arr []string) (uint64, []*pb.PairInfo) {
	var from time.Time
	var to time.Time
	length := len(arr)
	pairs := make([]*pb.PairInfo, 0, length)
	var count uint64 = 0
	if length < 1 {
		return 0, nil
	} else if length == 1 {
		num, pair := mine.GetNetflowByStamp(scene, arr[0])
		count = num
		pairs = append(pairs, pair)
	} else {
		utc1, er := strconv.Atoi(arr[0])
		if er != nil {
			return 0, nil
		}
		day1 := time.Unix(int64(utc1), 0)
		utc2, er := strconv.Atoi(arr[length-1])
		if er != nil {
			return 0, nil
		}
		day2 := time.Unix(int64(utc2), 0)
		from = time.Date(day1.Year(), day1.Month(), day1.Day(), 0, 0, 0, 0, time.UTC)
		to = time.Date(day2.Year(), day2.Month(), day2.Day(), 23, 59, 59, 0, time.UTC)
		list, _ := nosql.GetNetflowsByDuration(scene, from.Unix(), to.Unix())
		if er == nil {
			for _, s := range arr {
				num, size := getNFCountInList(list, s)
				utc, _ := strconv.Atoi(s)
				pairs = append(pairs, &pb.PairInfo{Id: uint64(utc), Key: s, Count: size})
				count += num
			}
		}
	}
	return count, pairs
}

func getNFCountInList(list []*nosql.Netflow, stamp string) (uint64, uint64) {
	utc, er := strconv.Atoi(stamp)
	if er != nil {
		return 0, 0
	}
	day := time.Unix(int64(utc), 0)
	from := time.Date(day.Year(), day.Month(), day.Day(), 1, 0, 0, 0, time.UTC)
	to := time.Date(day.Year(), day.Month(), day.Day(), 24, 0, 0, 0, time.UTC)
	var count uint64 = 0
	var size uint64 = 0
	for _, item := range list {
		if item.Created > from.Unix() && item.Created < to.Unix() {
			count += 1
			size += item.Size
		}
	}
	return count, size
}

func (mine *NetflowInfo) initInfo(db *nosql.Netflow) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Name = db.Name
	mine.Type = db.Type
	mine.Created = db.Created
	mine.Creator = db.Creator
	mine.Scene = db.Owner
	mine.Target = db.Target
	mine.Size = db.Size
	mine.Template = db.Template
	mine.Quote = db.Quote
	mine.Contents = db.Contents
}
