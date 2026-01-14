package telegram

import (
	"context"
	"fmt"
	"time"

	"github.com/ckayt/tetra/internal/config"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/rs/zerolog/log"
)

type Bot struct {
	client     *bot.Bot
	conf       *config.Config
	msgQueue   chan string
	testAction func(context.Context) string // callback for /test command
}

func New(cfg *config.Config, testAction func(context.Context) string) (*Bot, error) {
	b := &Bot{
		conf:       cfg,
		msgQueue:   make(chan string, 100), // Buffer for burst alerts
		testAction: testAction,
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(b.handler),
		bot.WithCheckInitTimeout(30 * time.Second),
	}

	// Create bot instance
	tBot, err := bot.New(cfg.TelegramToken, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}
	b.client = tBot

	// Register commands
	tBot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, b.startHandler)
	tBot.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, b.helpHandler)
	tBot.RegisterHandler(bot.HandlerTypeMessageText, "/test", bot.MatchTypeExact, b.testHandler)
	tBot.RegisterHandler(bot.HandlerTypeMessageText, "/speed", bot.MatchTypeExact, b.testHandler)

	return b, nil
}

func (b *Bot) Start(ctx context.Context) {
	// Start message sender routine
	go b.senderLoop(ctx)

	// Start polling
	log.Info().Msg("Starting Telegram bot polling...")
	b.client.Start(ctx)
}

func (b *Bot) Send(msg string) {
	select {
	case b.msgQueue <- msg:
	default:
		log.Warn().Msg("Telegram message queue full, dropping message")
	}
}

func (b *Bot) senderLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-b.msgQueue:
			b.sendMessageWithRetry(ctx, msg)
		}
	}
}

func (b *Bot) sendMessageWithRetry(ctx context.Context, text string) {
	backoff := time.Second
	maxBackoff := 30 * time.Second
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		_, err := b.client.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    b.conf.ChatID,
			Text:      text,
			ParseMode: models.ParseModeHTML,
		})
		if err == nil {
			return
		}

		log.Error().Err(err).Msgf("Failed to send telegram message (attempt %d/%d). Retrying in %v...", i+1, maxRetries, backoff)

		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff):
		}

		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}
	log.Error().Msg("Failed to send telegram message after max retries")
}

func (b *Bot) startHandler(ctx context.Context, bb *bot.Bot, update *models.Update) {
	msg := "👋 <b>Hello!</b> I am Tetra, your internet connection monitor.\n\n" +
		"I will periodically check your internet speed and notify you if it drops below the configured thresholds.\n" +
		"Use /help to see available commands."
	_, err := b.client.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      msg,
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to send start message")
	}
}

func (b *Bot) helpHandler(ctx context.Context, bb *bot.Bot, update *models.Update) {
	msg := "📋 <b>Available Commands:</b>\n" +
		"/test - Run an immediate speed test\n" +
		"/help - Show this help message\n" +
		"/start - Welcome message"
	_, err := b.client.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      msg,
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to send help message")
	}
}

func (b *Bot) testHandler(ctx context.Context, bb *bot.Bot, update *models.Update) {
	// Notify user test started
	_, err := b.client.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "🚀 <b>Starting manual speed test...</b> Please wait.",
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to send test starting message")
	}

	// Execute test (blocking call in this goroutine is fine for low volume,
	// but strictly speaking should be async if many users. Single user bot implies it's fine)
	// Actually, the prompt says "runs test ... and sends result".

	resultMsg := b.testAction(ctx)

	_, err = b.client.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      resultMsg,
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to send test result message")
	}
}

func (b *Bot) handler(ctx context.Context, bb *bot.Bot, update *models.Update) {
	// Default handler, ignore unknown messages
}
