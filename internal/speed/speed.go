package speed

import (
	"context"
	"fmt"
	"time"

	"github.com/ckayt/tetra/internal/stats"
	"github.com/rs/zerolog/log"
	"github.com/showwin/speedtest-go/speedtest"
)

type Runner struct{}

func NewRunner() *Runner {
	return &Runner{}
}

// Run executes the speedtest with retries.
// Returns a stats.Result.
func (r *Runner) Run(ctx context.Context) stats.Result {
	var result stats.Result
	var err error

	// Retry up to 3 times
	for i := 0; i < 3; i++ {
		if ctx.Err() != nil {
			result.Error = ctx.Err()
			return result
		}

		if i > 0 {
			log.Info().Msgf("Retrying speedtest (attempt %d/3)...", i+1)
			time.Sleep(5 * time.Second) // Wait a bit before retry
		}

		result, err = r.executeCheck(ctx)
		if err == nil {
			return result
		}
		log.Warn().Err(err).Msg("Speedtest failed")
	}

	result.Error = err
	result.Time = time.Now()
	return result
}

func (r *Runner) executeCheck(ctx context.Context) (stats.Result, error) {
	res := stats.Result{
		Time: time.Now(),
	}

	client := speedtest.New()

	// Fetch user info
	_, err := client.FetchUserInfoContext(ctx)
	if err != nil {
		return res, fmt.Errorf("failed to fetch user info: %w", err)
	}

	// Fetch servers
	serverList, err := client.FetchServerListContext(ctx)
	if err != nil {
		return res, fmt.Errorf("failed to fetch server list: %w", err)
	}

	// Find closest server
	targets, err := serverList.FindServer([]int{})
	if err != nil || len(targets) == 0 {
		return res, fmt.Errorf("failed to find server: %w", err)
	}

	server := targets[0] // Pick the best one

	// Ping
	err = server.PingTest(nil)
	if err != nil {
		return res, fmt.Errorf("ping test failed: %w", err)
	}
	res.Ping = server.Latency

	// Download
	err = server.DownloadTest()
	if err != nil {
		return res, fmt.Errorf("download test failed: %w", err)
	}
	res.Download = server.DLSpeed.Mbps()

	// Upload
	err = server.UploadTest()
	if err != nil {
		return res, fmt.Errorf("upload test failed: %w", err)
	}
	res.Upload = server.ULSpeed.Mbps()

	// Store byte counts if available (speedtest-go usually exposes them via server.Context but mostly we utilize DLSpeed/ULSpeed)
	// We won't worry about byte counts for this specific request as it's not explicitly asked for in the report,
	// but the struct has them. We'll leave them 0 for now unless we dig deep into internal counters.

	return res, nil
}
