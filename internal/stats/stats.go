package stats

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
)

type Result struct {
	Time          time.Time
	Download      float64 // Mbps
	Upload        float64 // Mbps
	Ping          time.Duration
	BytesReceived uint64
	BytesSent     uint64
	Error         error
	AlertSent     bool
}

type Summary struct {
	TotalTests     int
	AvgDownload    float64
	MinDownload    float64
	MaxDownload    float64
	AvgUpload      float64
	MinUpload      float64
	MaxUpload      float64
	AvgPing        time.Duration
	MinPing        time.Duration
	MaxPing        time.Duration
	AlertsCount    int
	LowSpeedEvents []Result
}

type Manager struct {
	mu      sync.RWMutex
	results []Result
	maxSize int
}

func NewManager(maxSize int) *Manager {
	if maxSize <= 0 {
		maxSize = 100 // Default safe size
	}
	return &Manager{
		results: make([]Result, 0, maxSize),
		maxSize: maxSize,
	}
}

func (m *Manager) Add(r Result) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Append
	m.results = append(m.results, r)

	// Trim if needed (keep latest maxSize)
	if len(m.results) > m.maxSize {
		m.results = m.results[len(m.results)-m.maxSize:]
	}
}

func (m *Manager) GetLast24hSummary(now time.Time, dlThreshold, ulThreshold float64) Summary {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cutoff := now.Add(-24 * time.Hour)
	var filtered []Result

	for _, r := range m.results {
		if r.Time.After(cutoff) {
			filtered = append(filtered, r)
		}
	}

	if len(filtered) == 0 {
		return Summary{}
	}

	s := Summary{
		TotalTests:  len(filtered),
		MinDownload: math.MaxFloat64,
		MinUpload:   math.MaxFloat64,
		MinPing:     time.Duration(math.MaxInt64),
	}

	var sumDL, sumUL float64
	var sumPing time.Duration

	for _, r := range filtered {
		if r.Error != nil {
			// Skip failed tests for avg calculations?
			// Prompt implies stats of internet quality, failed tests might mean NO internet.
			// Let's count them in TotalTests but skip metrics if values are 0.
			continue
		}

		sumDL += r.Download
		sumUL += r.Upload
		sumPing += r.Ping

		if r.Download < s.MinDownload {
			s.MinDownload = r.Download
		}
		if r.Download > s.MaxDownload {
			s.MaxDownload = r.Download
		}

		if r.Upload < s.MinUpload {
			s.MinUpload = r.Upload
		}
		if r.Upload > s.MaxUpload {
			s.MaxUpload = r.Upload
		}

		if r.Ping < s.MinPing {
			s.MinPing = r.Ping
		}
		if r.Ping > s.MaxPing {
			s.MaxPing = r.Ping
		}

		if r.AlertSent {
			s.AlertsCount++
		}

		// Identify low speed events based on thresholds provided (or just rely on AlertSent)
		// Prompt says "brief list of low-speed events if any".
		if r.Download < dlThreshold || r.Upload < ulThreshold {
			s.LowSpeedEvents = append(s.LowSpeedEvents, r)
		}
	}

	validTests := 0
	for _, r := range filtered {
		if r.Error == nil {
			validTests++
		}
	}

	if validTests > 0 {
		s.AvgDownload = sumDL / float64(validTests)
		s.AvgUpload = sumUL / float64(validTests)
		s.AvgPing = sumPing / time.Duration(validTests)
	} else {
		// Reset mins if no valid tests
		s.MinDownload = 0
		s.MinUpload = 0
		s.MinPing = 0
	}

	return s
}

func (s Summary) String() string {
	var sb strings.Builder
	sb.WriteString("📊 <b>Daily Report</b> (Last 24h)\n")
	sb.WriteString(fmt.Sprintf("Tests run: %d\n", s.TotalTests))
	if s.TotalTests > 0 {
		sb.WriteString(fmt.Sprintf("Alerts triggered: %d\n\n", s.AlertsCount))
		sb.WriteString(fmt.Sprintf("📉 <b>Download</b>:\nAvg: %.2f | Min: %.2f | Max: %.2f Mbps\n", s.AvgDownload, s.MinDownload, s.MaxDownload))
		sb.WriteString(fmt.Sprintf("📈 <b>Upload</b>:\nAvg: %.2f | Min: %.2f | Max: %.2f Mbps\n", s.AvgUpload, s.MinUpload, s.MaxUpload))
		sb.WriteString(fmt.Sprintf("📶 <b>Ping</b>:\nAvg: %dms | Min: %dms | Max: %dms\n", s.AvgPing.Milliseconds(), s.MinPing.Milliseconds(), s.MaxPing.Milliseconds()))
	}

	if len(s.LowSpeedEvents) > 0 {
		sb.WriteString("\n⚠️ <b>Low Speed Events:</b>\n")
		// Limit to last 5 to avoid spam
		count := 0
		for i := len(s.LowSpeedEvents) - 1; i >= 0; i-- {
			if count >= 5 {
				sb.WriteString("...and more\n")
				break
			}
			e := s.LowSpeedEvents[i]
			sb.WriteString(fmt.Sprintf("- %s: ▼%.1f ▲%.1f Mbps, %dms\n", e.Time.Format("15:04"), e.Download, e.Upload, e.Ping.Milliseconds()))
			count++
		}
	}
	return sb.String()
}
