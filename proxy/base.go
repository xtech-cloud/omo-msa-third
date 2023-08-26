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
