package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken     string `json:"-"`
	ChatID            int64
	DownloadThreshold float64
	UploadThreshold   float64
	CheckInterval     time.Duration
	DailyReportHour   int
	TimeZone          string
	LogLevel          string
}

func (c Config) String() string {
	return fmt.Sprintf("Config{ChatID:%d, Levels: DL=%.0f/UL=%.0f}", c.ChatID, c.DownloadThreshold, c.UploadThreshold)
}

func Load() (*Config, error) {
	// Load .env file, but don't fail if it doesn't exist (environment variables might be set directly)
	_ = godotenv.Load()

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_TOKEN is required")
	}

	chatIDStr := os.Getenv("CHAT_ID")
	if chatIDStr == "" {
		return nil, fmt.Errorf("CHAT_ID is required")
	}
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid CHAT_ID: %w", err)
	}

	cfg := &Config{
		TelegramToken:     token,
		ChatID:            chatID,
		DownloadThreshold: getEnvFloat("DOWNLOAD_THRESHOLD", 80.0),
		UploadThreshold:   getEnvFloat("UPLOAD_THRESHOLD", 100.0),
		CheckInterval:     getEnvDuration("CHECK_INTERVAL_MIN", 30*time.Minute),
		DailyReportHour:   getEnvInt("DAILY_REPORT_HOUR", 8),
		TimeZone:          getEnvString("TZ", "Europe/Kyiv"),
		LogLevel:          getEnvString("LOG_LEVEL", "info"),
	}

	return cfg, nil
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	// Try parsing as duration string (e.g. "1m", "30s")
	d, err := time.ParseDuration(val)
	if err == nil {
		return d
	}
	// Fallback: try parsing as simple integer (assumed minutes)
	i, err := strconv.Atoi(val)
	if err == nil {
		return time.Duration(i) * time.Minute
	}
	return defaultVal
}

func getEnvFloat(key string, defaultVal float64) float64 {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defaultVal
	}
	return f
}

func getEnvInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}

func getEnvString(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
