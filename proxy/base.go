package proxy

type EntityInfo struct {
	UID  string `json:"uid" bson:"uid"`
	Name string `json:"name" bson:"name"`
}

type TopicResult struct {
	Percent uint32 `json:"percent" bson:"percent"`
	Count   uint32 `json:"count" bson:"count"`
	State   uint32 `json:"state" bson:"state"`
}

type DurationInfo struct {
	Begin int64 `json:"begin" bson:"begin"`
	End   int64 `json:"end" bson:"end"`
}

type ContentInfo struct {
	UID      string   `json:"uid" bson:"uid"`
	Size     uint64   `json:"size" bson:"size"`
	Group    string   `json:"group" bson:"group"`
	Type     uint8    `json:"type" bson:"type"`
	Children []string `json:"children" bson:"children"`
}

type HonorContent struct {
	Name   string   `json:"name" bson:"name"`
	Remark string   `json:"remark" bson:"remark"`
	Quotes []string `json:"quotes" bson:"quotes"`
}
