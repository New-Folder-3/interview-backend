package conf

import (
	"interview/cmd/flags"
	"path/filepath"
)

type Database struct {
	Type        string `json:"type" env:"TYPE"`
	Host        string `json:"host" env:"HOST"`
	Port        int    `json:"port" env:"PORT"`
	User        string `json:"user" env:"USER"`
	Password    string `json:"password" env:"PASSWORD"`
	Name        string `json:"name" env:"NAME"`
	TablePrefix string `json:"tablePrefix" env:"TABLE_PREFIX"`
	SSLMode     string `json:"sslMode" env:"SSL_MODE"`
	DSN         string `json:"dsn" env:"DSN"`
	DBFile      string `json:"db_file" env:"DB_FILE"`
}

type Schema struct {
	Port   int    `json:"port" env:"PORT"`
	Listen string `json:"listen" env:"LISTEN"`
	URL    string `json:"url" env:"URL"`
}

type Key struct {
	AliyunAPIKey string `json:"aliyunApiKey" env:"ALIYUN_API_KEY"`
}

type Default struct {
	Job        string `json:"job" env:"JOB"`
	Age        int    `json:"age" env:"AGE"`
	ModelVoice string `json:"modelVoice" env:"MODEL_VOICE"`
}

type Model struct {
	TTSVoice  string `json:"voice" env:"VOICE"`
	TTSModel  string `json:"tts_model" env:"TTS_MODEL"`
	STTModel  string `json:"stt_model" env:"STT_MODEL"`
	ChatModel string `json:"chat_model" env:"CHAT_MODEL"`
}

type Mail struct {
	Host     string `json:"host" env:"HOST"`
	Port     int    `json:"port" env:"PORT"`
	Username string `json:"username" env:"USERNAME"`
	Password string `json:"password" env:"PASSWORD"`
	From     string `json:"from" env:"FROM"`
}

type Config struct {
	Database Database `json:"database" envPrefix:"DB_"`
	Schema   Schema   `json:"schema" envPrefix:"SCHEMA"`
	API      Key      `json:"api" envPrefix:"API_"`
	Default  Default  `json:"default" envPrefix:"DEFAULT"`
	Model    Model    `json:"model" envPrefix:"MODEL"`
	Mail     Mail     `json:"mail" envPrefix:"MAIL"`
}

func DefaultConfig() *Config {
	dbPath := filepath.Join(flags.DataDir, "data.db")
	return &Config{
		Database: Database{
			Type:        "sqlite3",
			Port:        0,
			TablePrefix: "x_",
			DBFile:      dbPath,
		},
		Schema: Schema{
			Port:   2233,
			Listen: "127.0.0.1",
			URL:    "https://mydomain.com",
		},
		Default: Default{
			Job:        "任何职位",
			Age:        18,
			ModelVoice: "Chelsie",
		},
		Model: Model{
			TTSVoice:  "Chelsie",
			TTSModel:  "qwen-tts",
			STTModel:  "paraformer-v2",
			ChatModel: "qwen-vl-max-latest",
		},
	}
}
