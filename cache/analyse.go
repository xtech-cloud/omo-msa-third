package cache

import (
	"encoding/base64"
	"encoding/json"
	"github.com/micro/go-micro/v2/logger"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
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
	time.Sleep(time.Second * 10)
	if config.Schema.Analyse.History {
		logger.Warn("check history events....")
		start := time.Now().AddDate(0, 0, -20).UnixMilli()
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
		from := time.Now().AddDate(0, 0, config.Schema.Analyse.Days).UnixMilli()
		cacheCtx.checkOldEvents(from, time.Now().UnixMilli())
	})
	if er != nil {
		logger.Warn("start cron failed that err = " + er.Error())
		return
	}
	cli.Start()
	logger.Info("start cron success that id = " + strconv.Itoa(int(id)) + " and timer = " + config.Schema.Analyse.Timer)
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
					list := mine.checkCounts(arr, device.Scene, device.SN, event)
					all = append(all, list...)
				}
			}
		}
	}

	for _, item := range all {
		motion := mine.GetMotionBy(item.Scene, item.SN, item.Event, item.Content)
		if motion == nil {
			_, _ = mine.CreateMotion(item.Scene, "", item.SN, item.Event, item.Content, "", item.Count)
		} else {
			_ = motion.AddCount(item.Count, "")
			if motion.UpdateTime.UnixMilli() > item.Timestamp {

			}
		}
	}
}

func (mine *cacheContext) checkDevice(scene, sn string) *MotionInfo {
	motion := mine.GetMotionBy(scene, sn, "", "")
	if motion == nil {
		motion, _ = mine.CreateMotion(scene, "", sn, "", "", "", 0)
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
	byes, err := base64.StdEncoding.DecodeString(content)
	err = json.Unmarshal(byes, list)
	if err != nil {
		return nil, err, pbstatus.ResultStatus_FormatError
	}
	return list, nil, pbstatus.ResultStatus_Success
}

func getRecordCount(sn, event, content string, list []*RecordCount) *RecordCount {
	for _, item := range list {
		if item.SN == sn && item.Event == event && item.Content == content {
			return item
		}
	}
	return nil
}
