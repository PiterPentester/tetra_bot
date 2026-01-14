package stats

import (
	"testing"
	"time"
)

func TestManager_GetLast24hSummary(t *testing.T) {
	mgr := NewManager(10)
	now := time.Now()

	// Add some results
	mgr.Add(Result{
		Time:     now.Add(-1 * time.Hour),
		Download: 100,
		Upload:   50,
		Ping:     20 * time.Millisecond,
	})
	mgr.Add(Result{
		Time:     now.Add(-2 * time.Hour),
		Download: 50,
		Upload:   20,
		Ping:     40 * time.Millisecond,
	})
	mgr.Add(Result{
		Time:     now.Add(-25 * time.Hour), // Should be ignored
		Download: 200,
		Upload:   100,
	})
	mgr.Add(Result{
		Time:      now.Add(-30 * time.Minute),
		Download:  10,
		Upload:    5,
		AlertSent: true, // Should count as alert
	})

	summary := mgr.GetLast24hSummary(now, 80.0, 100.0)

	if summary.TotalTests != 3 {
		t.Errorf("Expected 3 tests, got %d", summary.TotalTests)
	}
	if summary.AlertsCount != 1 {
		t.Errorf("Expected 1 alert, got %d", summary.AlertsCount)
	}

	// Avg DL: (100 + 50 + 10) / 3 = 160 / 3 = 53.333
	expectedAvg := 53.333
	if summary.AvgDownload < expectedAvg-0.1 || summary.AvgDownload > expectedAvg+0.1 {
		t.Errorf("Expected avg download ~%f, got %f", expectedAvg, summary.AvgDownload)
	}

	// Min DL: 10
	if summary.MinDownload != 10 {
		t.Errorf("Expected min download 10, got %f", summary.MinDownload)
	}

	// Check low speed events
	// DL < 80: 50 and 10.
	// UL < 100: 50, 20, 5. (All 3 are < 100).
	// So 3 events?
	if len(summary.LowSpeedEvents) != 3 {
		t.Errorf("Expected 3 low speed events, got %d", len(summary.LowSpeedEvents))
	}
}
