package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ckayt/tetra/internal/config"
	"github.com/ckayt/tetra/internal/speed"
	"github.com/ckayt/tetra/internal/stats"
	"github.com/ckayt/tetra/internal/telegram"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Setup logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// Set Log Level
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err == nil {
		zerolog.SetGlobalLevel(level)
	}

	log.Info().Str("config", cfg.String()).Msg("Starting Tetra")

	// Init components
	statsMgr := stats.NewManager(100) // Keep ~100 results (approx 2 days at 30min interval)
	speedRunner := speed.NewRunner()

	// Define test action wrapper with mutex to avoid concurrent speed tests
	var testMu sync.Mutex
	runTest := func(ctx context.Context, manual bool) string {
		testMu.Lock()
		defer testMu.Unlock()

		start := time.Now()
		log.Info().Bool("manual", manual).Msg("Running speed test...")

		res := speedRunner.Run(ctx)
		duration := time.Since(start)

		log.Info().
			Float64("download", res.Download).
			Float64("upload", res.Upload).
			Dur("ping", res.Ping).
			Err(res.Error).
			Dur("duration", duration).
			Msg("Speed test completed")

		msg := formatResult(res)

		// Check thresholds if not error
		alertTriggered := false
		if res.Error == nil && !manual {
			if res.Download < cfg.DownloadThreshold || res.Upload < cfg.UploadThreshold {
				alertTriggered = true
				res.AlertSent = true
			}
		}

		statsMgr.Add(res)

		if alertTriggered {
			return fmt.Sprintf("🚨 <b>Internet Quality Alert!</b>\n%s", msg)
		}
		if manual {
			return fmt.Sprintf("✅ <b>Manual Test Result:</b>\n%s", msg)
		}
		return ""
	}

	// Define stats action
	getStats := func(ctx context.Context) string {
		summary := statsMgr.GetLast24hSummary(time.Now(), cfg.DownloadThreshold, cfg.UploadThreshold)
		return summary.String()
	}

	// Init Telegram Bot with retry
	var bot *telegram.Bot
	for {
		bot, err = telegram.New(cfg, func(ctx context.Context) string {
			return runTest(ctx, true)
		}, getStats)
		if err == nil {
			break
		}
		log.Error().Err(err).Msg("Failed to init Telegram bot, retrying in 5s...")
		time.Sleep(5 * time.Second)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start Bot in background
	go bot.Start(ctx)

	// Start Ticker
	ticker := time.NewTicker(cfg.CheckInterval)
	defer ticker.Stop()

	// Daily Report Scheduler
	go dailyReportLoop(ctx, cfg, statsMgr, bot)

	// Run initial test immediately in background (after a short delay to let things settle)
	go func() {
		time.Sleep(5 * time.Second)
		log.Info().Msg("Taking initial speed test...")
		alertMsg := runTest(ctx, false)
		if alertMsg != "" {
			bot.Send(alertMsg)
		}
	}()

	// Start Health Check Server
	go func() {
		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})
		http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
			// Could check if bot is connected or config is loaded
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ready"))
		})

		log.Info().Msg("Starting health check server on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Error().Err(err).Msg("Health check server failed")
		}
	}()

	// Handle Signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msg("Tetra is running. Press Ctrl+C to stop.")

	for {
		select {
		case <-sigChan:
			log.Info().Msg("Shutting down...")
			cancel()
			// Give some time for cleanup if needed
			time.Sleep(1 * time.Second)
			return
		case <-ticker.C:
			alertMsg := runTest(ctx, false)
			if alertMsg != "" {
				bot.Send(alertMsg)
			}
		}
	}
}

func dailyReportLoop(ctx context.Context, cfg *config.Config, statsMgr *stats.Manager, bot *telegram.Bot) {
	loc, err := time.LoadLocation(cfg.TimeZone)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load timezone, using UTC")
		loc = time.UTC
	}

	for {
		now := time.Now().In(loc)
		nextReport := time.Date(now.Year(), now.Month(), now.Day(), cfg.DailyReportHour, 0, 0, 0, loc)

		if nextReport.Before(now) {
			nextReport = nextReport.Add(24 * time.Hour)
		}

		wait := nextReport.Sub(now)
		log.Info().Time("next_report", nextReport).Dur("wait", wait).Msg("Scheduled daily report")

		select {
		case <-ctx.Done():
			return
		case <-time.After(wait):
			// Generate report
			log.Info().Msg("Generating daily report...")
			summary := statsMgr.GetLast24hSummary(time.Now(), cfg.DownloadThreshold, cfg.UploadThreshold)
			bot.Send(summary.String())

			// Wait a bit to avoid double send due to slight time discrepancies (unlikely with time.After but good practice)
			time.Sleep(1 * time.Minute)
		}
	}
}

func formatResult(r stats.Result) string {
	if r.Error != nil {
		return fmt.Sprintf("⚠️ <b>Test Failed:</b> %v", r.Error)
	}
	return fmt.Sprintf(
		"⬇️ <b>Download:</b> %.2f Mbps\n"+
			"⬆️ <b>Upload:</b> %.2f Mbps\n"+
			"📶 <b>Ping:</b> %d ms",
		r.Download, r.Upload, r.Ping.Milliseconds(),
	)
}
