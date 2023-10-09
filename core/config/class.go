package config

type LogConfig struct {
	Debug bool `json:"debug"`
	Both  bool `json:"both"`
	File  bool `json:"file"`
}

type MongoConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Db       string `json:"db"`
}

type MysqlConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Db       string `json:"db"`
	Charset  string `json:"charset"`
}

type RedisConfig struct {
	Prefix   string `json:"prefix"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	AuthPass string `json:"auth_pass"`
	Db       int    `json:"db"`
}

type ServiceConfig struct {
	TslCrt       string                    `json:"tslCrt"`
	TslKey       string                    `json:"tslKey"`
	ServiceNodes map[int]ServiceNodeConfig `json:"services"`
}

type ServiceNodeConfig struct {
	ClientPort string `json:"clientPort"`
	UseSSL     bool   `json:"useSSL"`
}
