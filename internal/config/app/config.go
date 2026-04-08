package config

import (
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ConfPath string `json:"conf_path"`
	Port     string `json:"bot_port"`

	//env
	DbDNS         string
	BotToken      string
	GigaAuthKey   string
	SaluteAuthKey string
	//param
	Param_Get_Giga_Token time.Duration  `json:"param_get_gtoken"`
	Param_get_status     time.Duration  `json:"param_get_status"`
	PoolCert             *x509.CertPool // пул сертификатов
}

// NewConfig - создание конфигурации приложения
func NewConfig() (*Config, error) {
	cfg := Config{}
	flag.StringVar(&cfg.ConfPath, "config", "./conf.json", "путь для файла c конфигурацией")

	if cfg.ConfPath != "" {
		err := cfg.applyConfigIfEmpty()
		if err != nil {
			return nil, err
		}
	}

	pool, err := setCert()
	if err != nil {
		return nil, err
	}
	cfg.PoolCert = pool

	// обязателльные параметры
	flag.StringVar(&cfg.DbDNS, "d", getDef("DATABASE_URI", ""), "db url")
	flag.StringVar(&cfg.BotToken, "t", getDef("TELEGRAM_BOT_TOKEN", ""), "токен tgBot")
	flag.StringVar(&cfg.GigaAuthKey, "g", getDef("GIGACHAT_AUTH_KEY", ""), "gigaChat auth key")
	flag.StringVar(&cfg.SaluteAuthKey, "s", getDef("SALUTESPEECH_AUTH_KEY", ""), "saluteSpeech auth key")
	flag.Parse()

	if v := os.Getenv("param_get_gtoken"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		cfg.Param_Get_Giga_Token = time.Duration(val) * time.Minute
	}
	return &cfg, nil
}

func getDef(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func (cfg *Config) applyConfigIfEmpty() error {
	// Чтение конфигурационного файла JSON, если указан путь
	var confApp Config
	if cfg.ConfPath != "" {
		conf, err := os.ReadFile(cfg.ConfPath)
		if err != nil {
			return fmt.Errorf("ошибка чтения файла конфигурации: %w", err)
		}
		err = json.Unmarshal(conf, &confApp)
		if err != nil {
			return fmt.Errorf("ошибка разбора файла конфигурации: %w", err)
		}
	}

	if confApp.Port != "" {
		cfg.Port = confApp.Port
	}
	if confApp.Param_Get_Giga_Token != 0 {
		cfg.Param_Get_Giga_Token = confApp.Param_Get_Giga_Token * time.Second
	}
	if confApp.Param_get_status != 0 {
		cfg.Param_get_status = confApp.Param_get_status * time.Second
	}
	return nil
}

func setCert() (*x509.CertPool, error) {
	// Создаем новый пул сертификатов
	certPool := x509.NewCertPool()

	gostBytes, err := os.ReadFile("russian_trusted_root_ca_gost_2025.cer")
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл сертификата - russian_trusted_root_ca_gost_2025.cer: %v", err)
	}

	// Распарсиваем DER-байты в структуру сертификата
	gost2025, err := x509.ParseCertificate(gostBytes)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить сертификат gostCerts: %v", err)
	}
	// Добавляем распарсенный сертификат в пул
	certPool.AddCert(gost2025)

	rootCa, err := os.ReadFile("russian_trusted_root_ca.cer")
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл сертификата - russian_trusted_root_ca.cer: %v", err)
	}

	// Добавляем наш сертификат в пул
	if !certPool.AppendCertsFromPEM(rootCa) {
		return nil, fmt.Errorf("ну удалось добавить сертификат в пул - rootCa")
	}

	subCa, err := os.ReadFile("russian_trusted_sub_ca_2024.cer")
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл сертификата - russian_trusted_sub_ca_2024.cer: %v", err)
	}

	// Добавляем наш сертификат в пул
	if !certPool.AppendCertsFromPEM(subCa) {
		return nil, fmt.Errorf("ну удалось добавить сертификат в пул - subCa")
	}

	caGostByte, err := os.ReadFile("russian_trusted_root_ca_gost_2025.cer")
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл сертификата - russian_trusted_root_ca_gost_2025.cer: %v", err)
	}
	// Распарсиваем DER-байты в структуру сертификата
	cagostBytes, err := x509.ParseCertificate(caGostByte)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить сертификат gostCerts: %v", err)
	}
	// Добавляем распарсенный сертификат в пул
	certPool.AddCert(cagostBytes)

	sabCa2, err := os.ReadFile("russian_trusted_sub_ca.cer")
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл сертификата - russian_trusted_sub_ca.cer: %v", err)
	}

	// Добавляем наш сертификат в пул
	if !certPool.AppendCertsFromPEM(sabCa2) {
		return nil, fmt.Errorf("ну удалось добавить сертификат в пул - sabCa2")
	}

	return certPool, nil
}
