package cache

import (
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.third/proxy/nosql"
	"omo.msa.third/tool"
	"sort"
	"strings"
	"time"
)

type MotionInfo struct {
	Count uint32
	Type  uint32
	baseInfo
	Scene   string
	AppID   string
	SN      string
	UserID  string
	EventID string
	meta    string
	content string
	bundle  string
}

func (mine *cacheContext) CreateMotion(scene, app, sn, eveID, content, operator string, tp, count uint32) (*MotionInfo, error) {
	db := new(nosql.Motion)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetMotionNextID()
	db.CreatedTime = time.Now()
	db.Creator = operator
	db.Operator = operator
	db.UpdatedTime = time.Now()
	db.AppID = app
	db.Scene = scene
	db.SN = sn
	db.EventID = eveID
	db.Count = count
	db.Type = tp
	if db.Count < 1 {
		db.Count = 1
	}
	db.Content = content

	err := nosql.CreateMotion(db)
	if err == nil {
		info := new(MotionInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) GetMotionsBySN(scene, sn string) []*MotionInfo {
	dbs, err := nosql.GetMotionsBySN(scene, sn)
	list := make([]*MotionInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			info := new(MotionInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetMotion(uid string) (*MotionInfo, error) {
	db, err := nosql.GetMotion(uid)
	if err == nil {
		info := new(MotionInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) RemoveMotion(uid, operator string) error {
	return nosql.RemoveMotion(uid, operator)
}

func (mine *cacheContext) GetMotionsByEvent(scene, id string) []*MotionInfo {
	dbs, err := nosql.GetMotionsByEventID(scene, id)
	list := make([]*MotionInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			info := new(MotionInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetMotionsBySNEvent(scene, sn, event string) []*MotionInfo {
	dbs, err := nosql.GetMotionsBySNEvent(scene, sn, event)
	list := make([]*MotionInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			info := new(MotionInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetMotionBy(scene, sn, event, content string) *MotionInfo {
	db, err := nosql.GetMotionBy(scene, sn, event, content)
	if err == nil {
		info := new(MotionInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext) GetMotionContentCount(content string) uint64 {
	if len(content) < 2 {
		return 0
	}
	var num uint64 = 0
	dbs, err := nosql.GetMotionsByContent2(content)
	if err != nil {
		return num
	}
	for _, db := range dbs {
		num += uint64(db.Count)
	}
	return num
}

func (mine *cacheContext) GetMotionsByRegex(scene, sn, event, content string) []*MotionInfo {
	if len(sn) > 0 {
		if len(event) > 2 {
			var info = mine.GetMotionBy(scene, sn, event, content)
			arr := make([]*MotionInfo, 0, 1)
			arr = append(arr, info)
			return arr
		} else {
			return mine.GetMotionsByContent(scene, sn, content)
		}
	}
	dbs, err := nosql.GetMotionsByRegex(scene, content)
	if err != nil {
		return nil
	}
	list := make([]*MotionInfo, 0, len(dbs))
	for _, db := range dbs {
		info := new(MotionInfo)
		info.initInfo(db)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetMotionsByContent(scene, sn, content string) []*MotionInfo {
	dbs, err := nosql.GetMotionsBySN(scene, sn)
	list := make([]*MotionInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			if strings.Contains(db.Content, content) {
				info := new(MotionInfo)
				info.initInfo(db)
				list = append(list, info)
			}
		}
	}
	return list
}

func (mine *cacheContext) GetRanksByBundle(scene string, num uint32, events []string) []*MotionInfo {
	dbs, err := nosql.GetMotionsByOwner(scene)
	list := make([]*MotionInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			if tool.HasItem(events, db.EventID) && len(db.Content) > 2 {
				info := new(MotionInfo)
				info.initInfo(db)
				bundle := getMotionInfo(info.bundle, list)
				if bundle != nil {
					bundle.Count += info.Count
				} else {
					list = append(list, info)
				}
			}
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Count > list[j].Count
	})
	if len(list) > int(num) {
		return list[:num]
	} else {
		return list
	}
}

func (mine *cacheContext) GetRanksByContent(scene string, num uint32, events []string) []*MotionInfo {
	dbs, err := nosql.GetMotionsByOwner(scene)
	list := make([]*MotionInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			if tool.HasItem(events, db.EventID) && len(db.Content) > 2 {
				info := new(MotionInfo)
				info.initInfo(db)
				list = append(list, info)
			}
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Count > list[j].Count
	})
	if len(list) > int(num) {
		return list[:num]
	} else {
		return list
	}
}

func (mine *cacheContext) GetMotionsByEveContent(scene, event, content string) []*MotionInfo {
	dbs, err := nosql.GetMotionsByEventID(scene, event)
	list := make([]*MotionInfo, 0, len(dbs))
	if err == nil {
		for _, db := range dbs {
			if strings.Contains(db.Content, content) {
				info := new(MotionInfo)
				info.initInfo(db)
				list = append(list, info)
			}
		}
	}
	return list
}

func (mine *cacheContext) GetMotionsTopContent(tp, top uint32) []string {
	dbs, _ := nosql.GetMotionsByTop(DefaultScene, tp, top)
	list := make([]string, 0, len(dbs))
	for _, db := range dbs {
		list = append(list, db.Content)
	}
	return list
}

func getMotionInfo(bundle string, list []*MotionInfo) *MotionInfo {
	for _, info := range list {
		if info.bundle == bundle {
			return info
		}
	}
	return nil
}

func (mine *MotionInfo) initInfo(db *nosql.Motion) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Created = db.Created
	mine.Updated = db.Updated
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Scene = db.Scene
	mine.EventID = db.EventID
	mine.SN = db.SN
	mine.AppID = db.AppID
	mine.Count = db.Count
	mine.UserID = db.UserID
	mine.Type = db.Type
	mine.meta = db.Content

	if len(mine.meta) > 2 && gjson.Valid(mine.meta) {
		result := gjson.Parse(mine.meta)
		uri := result.Get("uri").String()
		if strings.Contains(uri, "/") {
			arr := strings.Split(uri, "/")
			mine.bundle = arr[0]
			mine.content = arr[1]
		}
	}
}

func (mine *MotionInfo) AddCount(offset uint32, operator string) error {
	err := nosql.UpdateMotionCount(mine.UID, operator, mine.Count+offset)
	if err == nil {
		mine.Count = mine.Count + offset
		mine.Updated = time.Now().Unix()
	}
	return err
}

func (mine *MotionInfo) Bundle() string {
	return mine.bundle
}

func (mine *MotionInfo) Content() string {
	return mine.content
}

func (mine *MotionInfo) Meta() string {
	return mine.meta
}
