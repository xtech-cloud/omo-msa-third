package cache

import (
	pb "github.com/xtech-cloud/omo-msp-third/proto/third"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.third/proxy"
	"omo.msa.third/proxy/nosql"
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
