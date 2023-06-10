package config

const defaultJson string = `{
	"service": {
		"address": ":7080",
		"ttl": 15,
		"interval": 10
	},
	"logger": {
		"level": "trace",
		"file": "logs/server.log",
		"std": true
	},
	"database": {
		"name": "rgsCloud",
		"ip": "192.168.1.10",
		"port": "27017",
		"user": "root",
		"password": "pass2019",
		"type": "mongodb"
	},
	"basic": {
		"synonym": 6,
		"tag": 6,
		"count":"http://api.xtech.cloud:28032/v1/xtc/analytics/generator/record",
        "token":"ogm.dev"
	},
    "analyse":{
		"history": true,
		"timer":"38 18 * * *",
		"days":-1,
		"events": [
			{
				"name": "点击量",
				"type": 1,
				"ids": [
					"/XTC/IntegerationBoard/Open",
					"/Meex/EBook/Play"
				]
			},
			{
				"name": "点赞数",
				"type": 2,
				"ids": [
					"/XTC/IntegerationBoard/Like"
				]
			},
			{
				"name": "运行时长",
				"type": 3,
				"ids": [
					"Analytics.Awake"
				]
			}
		]
    }
}
`
