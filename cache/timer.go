package cache

import "github.com/sirupsen/logrus"

type CronLog struct {
	log *logrus.Logger
}

func (l *CronLog) Info(msg string, keysAndValues ...interface{}) {
	l.log.WithFields(logrus.Fields{
		"data": keysAndValues,
	}).Info(msg)
}

func (l *CronLog) Error(err error, msg string, keysAndValues ...interface{}) {
	l.log.WithFields(logrus.Fields{
		"msg":  msg,
		"data": keysAndValues,
	}).Warn(msg)
}
