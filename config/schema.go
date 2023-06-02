package config

type ServiceConfig struct {
	TTL      int64  `json:"ttl"`
	Interval int64  `json:"interval"`
	Address  string `json:"address"`
}

type LoggerConfig struct {
	Level string `json:"level"`
	File  string `json:"file"`
	Std   bool   `json:"std"`
}

type DBConfig struct {
	Type     string `json:"type"`
	User     string `json:"user"`
	Password string `json:"password"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Name     string `json:"name"`
}

type BasicConfig struct {
	SynonymMax int    `ini:"synonym"`
	TagMax     int    `ini:"tag"`
	OgmCount   string `json:"count"`
	OgmToken   string `json:"token"`
}

type AnalyseConfig struct {
	History bool           `json:"history"`
	Timer   string         `json:"timer"`
	Days    int            `json:"days"`
	Events  []*EventConfig `json:"events"`
}

type EventConfig struct {
	Name string   `json:"name"`
	Type uint32   `json:"type"`
	IDs  []string `json:"ids"`
}

type SchemaConfig struct {
	Service  ServiceConfig `json:"service"`
	Logger   LoggerConfig  `json:"logger"`
	Database DBConfig      `json:"database"`
	Basic    BasicConfig   `json:"basic"`
	Analyse  AnalyseConfig `json:"analyse"`
}
