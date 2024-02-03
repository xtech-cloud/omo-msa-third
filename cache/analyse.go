package cache

import (
	"encoding/base64"
	"encoding/json"
	"github.com/micro/go-micro/v2/logger"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	pb "github.com/xtech-cloud/omo-msp-third/proto/third"
	"omo.msa.third/config"
	"omo.msa.third/proxy/nosql"
	"strconv"
	"time"
)

type TerminalRecord struct {
	AppID          string `json:"AppID"`
	SerialNumber   string `json:"SerialNumber"`
	UserID         string `json:"UserID"`
	EventID        string `json:"EventID"`
	EventParameter string `json:"EventParameter"`
	Timestamp      int64  `json:"Timestamp"`
	Uuid           string `json:"Uuid"`
}

type RecordCount struct {
	Scene     string
	SN        string
	Event     string
	Content   string
	Count     uint32
	Timestamp int64
}

func CheckAnalyse() {
	time.Sleep(time.Second * 5)
	if config.Schema.Analyse.History {
		logger.Warn("check history events....")
		start := time.Now().AddDate(0, 0, -30).UnixMilli()
		end := time.Now().AddDate(0, 0, config.Schema.Analyse.Days).UnixMilli()
		cacheCtx.checkOldEvents(start, end)
	}

	loger := &CronLog{
		log: logrus.New(),
	}
	loger.log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	cli := cron.New(cron.WithChain(cron.SkipIfStillRunning(loger)))
	id, er := cli.AddFunc(config.Schema.Analyse.Timer, func() {
		from := time.Now().AddDate(0, 0, config.Schema.Analyse.Days)
		cacheCtx.checkOldEvents(from.UnixMilli(), time.Now().UnixMilli())
	})
	if er != nil {
		logger.Warn("start cron failed that err = " + er.Error())
		return
	}
	cli.Start()
	logger.Info("start cron success that id = " + strconv.Itoa(int(id)) + " and timer = " + config.Schema.Analyse.Timer)
}

func getRecordCount(sn, event, content string, list []*RecordCount) *RecordCount {
	for _, item := range list {
		if item.SN == sn && item.Event == event && item.Content == content {
			return item
		}
	}
	return nil
}

func (mine *cacheContext) checkOldEvents(start, end int64) {
	devices, _ := nosql.GetAllDevices()
	logger.Warn("start check old events.... device length = " + strconv.Itoa(len(devices)))
	all := make([]*RecordCount, 0, len(devices)*20)
	events := make([]string, 0, 10)
	for _, item := range config.Schema.Analyse.Events {
		events = append(events, item.IDs...)
	}
	for _, device := range devices {
		if device.Status > 1 && device.Status < 99 {
			for _, event := range events {
				arr, err, _ := mine.getOldEvents(device.SN, event, start, end)
				if err == nil {
					if len(arr) > 0 {
						list := mine.checkCounts(arr, device.Scene, device.SN, event)
						all = append(all, list...)
					} else {
						logger.Warn("check old events is empty of sn " + device.SN + " and event = " + event)
					}
				} else {
					logger.Warn("check old events error = " + err.Error())
				}
			}
		}
	}

	for _, item := range all {
		motion := mine.GetMotionBy(item.Scene, item.SN, item.Event, item.Content)
		if motion == nil {
			_, _ = mine.CreateMotion(item.Scene, "", item.SN, item.Event, item.Content, "", 0, item.Count)
		} else {
			_ = motion.AddCount(item.Count, "")
			if motion.Updated > item.Timestamp {

			}
		}
	}
}

func (mine *cacheContext) checkDevice(scene, sn string) *MotionInfo {
	motion := mine.GetMotionBy(scene, sn, "", "")
	if motion == nil {
		motion, _ = mine.CreateMotion(scene, "", sn, "", "", "", 0, 1)
	}
	return motion
}

func (mine *cacheContext) checkCounts(list []*TerminalRecord, scene, sn, event string) []*RecordCount {
	arr := make([]*RecordCount, 0, 20)
	for _, item := range list {
		info := getRecordCount(sn, event, item.EventParameter, arr)
		if info == nil {
			info = &RecordCount{Scene: scene, SN: sn, Event: item.EventID, Content: item.EventParameter, Timestamp: item.Timestamp}
			arr = append(arr, info)
		}
		info.Timestamp = item.Timestamp
		info.Count += 1
	}
	return arr
}

func (mine *cacheContext) getOldEvents(sn, event string, start, end int64) ([]*TerminalRecord, error, pbstatus.ResultStatus) {
	list := make([]*TerminalRecord, 0, 10)
	var reqRecords struct {
		StartTime int64  `json:"startTime"`
		SN        string `json:"serialNumber"`
		EndTime   int64  `json:"endTime"`
		EventID   string `json:"eventID"`
	}

	reqRecords.SN = sn
	reqRecords.StartTime = start
	reqRecords.EventID = event
	reqRecords.EndTime = end
	req, er := json.Marshal(reqRecords)
	if er != nil {
		return nil, er, pbstatus.ResultStatus_FormatError
	}
	result, _, er := PostHttp(config.Schema.Basic.OgmCount, req, false)
	if er != nil {
		return nil, er, pbstatus.ResultStatus_ServerError
	}

	content := result.Get("content").String()
	bytes, err := base64.StdEncoding.DecodeString(content)
	err = json.Unmarshal(bytes, &list)
	if err != nil {
		return nil, err, pbstatus.ResultStatus_FormatError
	}
	return list, nil, pbstatus.ResultStatus_Success
}

func (mine *cacheContext) getEventsBySN(sn string, start, end int64) ([]*TerminalRecord, error, pbstatus.ResultStatus) {
	list := make([]*TerminalRecord, 0, 10)
	var reqRecords struct {
		SN        string `json:"serialNumber"`
		StartTime int64  `json:"startTime"`
		EndTime   int64  `json:"endTime"`
	}

	reqRecords.SN = sn
	reqRecords.StartTime = start
	reqRecords.EndTime = end
	req, er := json.Marshal(reqRecords)
	if er != nil {
		return nil, er, pbstatus.ResultStatus_FormatError
	}
	result, _, er := PostHttp(config.Schema.Basic.OgmCount, req, false)
	if er != nil {
		return nil, er, pbstatus.ResultStatus_ServerError
	}

	content := result.Get("content").String()
	bytes, err := base64.StdEncoding.DecodeString(content)
	err = json.Unmarshal(bytes, &list)
	if err != nil {
		return nil, err, pbstatus.ResultStatus_FormatError
	}
	return list, nil, pbstatus.ResultStatus_Success
}

func (mine *cacheContext) GetWeekCount(sn, event string) uint32 {
	day := time.Now()
	weekDay := day.Weekday()
	from := time.Now().AddDate(0, 0, -int(weekDay))
	to := time.Date(day.Year(), day.Month(), day.Day(), 23, 0, 0, 0, time.UTC)
	dbs, _, _ := mine.getOldEvents(sn, event, from.UnixMilli(), to.UnixMilli())
	return uint32(len(dbs))
}

func (mine *cacheContext) GetMouthCount(sn, event string) uint32 {
	day := time.Now()
	from := time.Date(day.Year(), day.Month(), 1, 1, 0, 0, 0, time.UTC)
	to := time.Date(day.Year(), day.Month(), day.Day(), 24, 0, 0, 0, time.UTC)
	dbs, _, _ := mine.getOldEvents(sn, event, from.UnixMilli(), to.UnixMilli())
	return uint32(len(dbs))
}

func (mine *cacheContext) GetEventsCount(sn string, utc int64) uint32 {
	day := time.Unix(utc, 0)
	from := time.Date(day.Year(), day.Month(), day.Day(), 1, 0, 0, 0, time.UTC)
	to := time.Date(day.Year(), day.Month(), day.Day(), 24, 0, 0, 0, time.UTC)
	list, err, _ := mine.getEventsBySN(sn, from.UnixMilli(), to.UnixMilli())
	if err != nil {
		return 0
	}
	return uint32(len(list))
}

func (mine *cacheContext) GetEventsByList(sn string, arr []string) (uint64, []*pb.PairInfo) {
	var from time.Time
	var to time.Time
	length := len(arr)
	pairs := make([]*pb.PairInfo, 0, length)
	var count uint64 = 0
	if length < 1 {
		return 0, nil
	} else if length == 1 {
		utc, er := strconv.Atoi(arr[0])
		if er != nil {
			return 0, nil
		}
		day := time.Unix(int64(utc), 0)
		from = time.Date(day.Year(), day.Month(), day.Day(), 1, 0, 0, 0, time.UTC)
		to = time.Date(day.Year(), day.Month(), day.Day(), 24, 0, 0, 0, time.UTC)
		list, _, _ := mine.getEventsBySN(sn, from.UnixMilli(), to.UnixMilli())
		pairs = append(pairs, &pb.PairInfo{Id: uint64(utc), Key: arr[0], Count: uint64(len(list))})
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
		list, er, _ := mine.getEventsBySN(sn, from.UnixMilli(), to.UnixMilli())
		if er == nil {
			for _, s := range arr {
				num := getCountInList(list, s)
				utc, _ := strconv.Atoi(s)
				pairs = append(pairs, &pb.PairInfo{Id: uint64(utc), Key: s, Count: uint64(num)})
				count += uint64(num)
			}
		}
	}
	return count, pairs
}

func getCountInList(list []*TerminalRecord, stamp string) uint32 {
	utc, er := strconv.Atoi(stamp)
	if er != nil {
		return 0
	}
	day := time.Unix(int64(utc), 0)
	from := time.Date(day.Year(), day.Month(), day.Day(), 1, 0, 0, 0, time.UTC)
	to := time.Date(day.Year(), day.Month(), day.Day(), 24, 0, 0, 0, time.UTC)
	var count uint32 = 0
	for _, record := range list {
		if record.Timestamp > from.UnixMilli() && record.Timestamp < to.UnixMilli() {
			count += 1
		}
	}
	return count
}
