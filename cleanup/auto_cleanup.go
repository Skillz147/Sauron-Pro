package cleanup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"o365/configdb"
	"o365/handlers"
	"o365/utils"

	"github.com/rs/zerolog/log"
)

// AutoCleanupConfig defines automatic cleanup behavior
type AutoCleanupConfig struct {
	Enabled       bool     `json:"enabled"`
	IntervalHours int      `json:"interval_hours"` // How often to run cleanup
	RetentionDays int      `json:"retention_days"` // Keep data newer than N days
	Operations    []string `json:"operations"`     // Which operations to run
}

// AutoCleanupScheduler manages automated cleanup operations
type AutoCleanupScheduler struct {
	config   AutoCleanupConfig
	stopChan chan bool
	running  bool
}

var scheduler *AutoCleanupScheduler

// DefaultAutoCleanupConfig returns the default configuration
func DefaultAutoCleanupConfig() AutoCleanupConfig {
	return AutoCleanupConfig{
		Enabled:       true,
		IntervalHours: 24,               // Run daily
		RetentionDays: 7,                // Keep 7 days of logs
		Operations:    []string{"logs"}, // Only auto-cleanup logs
	}
}

// InitAutoCleanup starts the automatic cleanup scheduler
func InitAutoCleanup() {
	config := DefaultAutoCleanupConfig()
	scheduler = &AutoCleanupScheduler{
		config:   config,
		stopChan: make(chan bool),
		running:  false,
	}

	if config.Enabled {
		scheduler.Start()
	}

	log.Info().
		Bool("enabled", config.Enabled).
		Int("interval_hours", config.IntervalHours).
		Int("retention_days", config.RetentionDays).
		Strs("operations", config.Operations).
		Msg("🧹 Auto-cleanup scheduler initialized")
}

// Start begins the automatic cleanup schedule
func (acs *AutoCleanupScheduler) Start() {
	if acs.running {
		return
	}

	acs.running = true
	go acs.scheduleLoop()

	utils.SystemLogger.Info().
		Int("interval_hours", acs.config.IntervalHours).
		Msg("🔄 Auto-cleanup scheduler started")
}

// Stop halts the automatic cleanup schedule
func (acs *AutoCleanupScheduler) Stop() {
	if !acs.running {
		return
	}

	acs.stopChan <- true
	acs.running = false

	utils.SystemLogger.Info().Msg("⏹️ Auto-cleanup scheduler stopped")
}

// scheduleLoop runs the cleanup operations on schedule
func (acs *AutoCleanupScheduler) scheduleLoop() {
	ticker := time.NewTicker(time.Duration(acs.config.IntervalHours) * time.Hour)
	defer ticker.Stop()

	// Run initial cleanup on startup (after 1 minute)
	initialTimer := time.NewTimer(1 * time.Minute)
	defer initialTimer.Stop()

	for {
		select {
		case <-initialTimer.C:
			acs.performCleanup("startup")
			initialTimer.Stop() // Disable initial timer after first run

		case <-ticker.C:
			acs.performCleanup("scheduled")

		case <-acs.stopChan:
			utils.SystemLogger.Info().Msg("🛑 Auto-cleanup scheduler stopped")
			return
		}
	}
}

// performCleanup executes the actual cleanup operations
func (acs *AutoCleanupScheduler) performCleanup(trigger string) {
	utils.SystemLogger.Info().
		Str("trigger", trigger).
		Int("retention_days", acs.config.RetentionDays).
		Strs("operations", acs.config.Operations).
		Msg("🧹 Starting automatic cleanup")

	startTime := time.Now()

	// Create cleanup request
	cleanupReq := handlers.CleanupRequest{
		AdminKey:      configdb.GetAdminKey(),
		Operations:    acs.config.Operations,
		RetentionDays: acs.config.RetentionDays,
		DryRun:        false,
	}

	// Simulate HTTP request to reuse existing cleanup logic
	reqBody, err := json.Marshal(cleanupReq)
	if err != nil {
		utils.SystemLogger.Error().Err(err).Msg("❌ Auto-cleanup: Failed to marshal request")
		return
	}

	// Create a mock HTTP request
	req, err := http.NewRequest("POST", "/admin/cleanup", strings.NewReader(string(reqBody)))
	if err != nil {
		utils.SystemLogger.Error().Err(err).Msg("❌ Auto-cleanup: Failed to create request")
		return
	}

	// Use a response recorder to capture the response
	recorder := &responseRecorder{
		statusCode: 200,
		body:       make([]byte, 0),
	}

	// Call the cleanup handler
	handlers.HandleAdminCleanup(recorder, req)

	// Parse response
	var cleanupResp handlers.CleanupResponse
	if err := json.Unmarshal(recorder.body, &cleanupResp); err != nil {
		utils.SystemLogger.Error().Err(err).Msg("❌ Auto-cleanup: Failed to parse response")
		return
	}

	duration := time.Since(startTime)

	// Log results
	if cleanupResp.Success {
		utils.SystemLogger.Info().
			Dur("duration", duration).
			Int64("size_freed", cleanupResp.TotalSize).
			Int("operations", len(cleanupResp.Operations)).
			Str("trigger", trigger).
			Msg("✅ Auto-cleanup completed successfully")
	} else {
		utils.SystemLogger.Error().
			Dur("duration", duration).
			Str("trigger", trigger).
			Str("message", cleanupResp.Message).
			Msg("❌ Auto-cleanup failed")
	}

	// Log detailed operation results
	for operation, result := range cleanupResp.Operations {
		if result.Success {
			utils.SystemLogger.Debug().
				Str("operation", operation).
				Int("items_removed", result.ItemsRemoved).
				Int64("size_freed", result.SizeFreed).
				Str("details", result.Details).
				Msg("🧹 Auto-cleanup operation completed")
		} else {
			utils.SystemLogger.Warn().
				Str("operation", operation).
				Str("error", result.Error).
				Msg("⚠️ Auto-cleanup operation failed")
		}
	}
}

// responseRecorder is a simple HTTP response recorder for internal use
type responseRecorder struct {
	statusCode int
	headers    http.Header
	body       []byte
}

func (rr *responseRecorder) Header() http.Header {
	if rr.headers == nil {
		rr.headers = make(http.Header)
	}
	return rr.headers
}

func (rr *responseRecorder) Write(data []byte) (int, error) {
	rr.body = append(rr.body, data...)
	return len(data), nil
}

func (rr *responseRecorder) WriteHeader(statusCode int) {
	rr.statusCode = statusCode
}

// GetSchedulerStatus returns the current status of the auto-cleanup scheduler
func GetSchedulerStatus() map[string]interface{} {
	if scheduler == nil {
		return map[string]interface{}{
			"initialized": false,
		}
	}

	return map[string]interface{}{
		"initialized":    true,
		"enabled":        scheduler.config.Enabled,
		"running":        scheduler.running,
		"interval_hours": scheduler.config.IntervalHours,
		"retention_days": scheduler.config.RetentionDays,
		"operations":     scheduler.config.Operations,
	}
}

// UpdateSchedulerConfig updates the scheduler configuration
func UpdateSchedulerConfig(config AutoCleanupConfig) error {
	if scheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}

	// Stop current scheduler if running
	if scheduler.running {
		scheduler.Stop()
	}

	// Update configuration
	scheduler.config = config

	// Restart if enabled
	if config.Enabled {
		scheduler.Start()
	}

	utils.SystemLogger.Info().
		Bool("enabled", config.Enabled).
		Int("interval_hours", config.IntervalHours).
		Int("retention_days", config.RetentionDays).
		Strs("operations", config.Operations).
		Msg("🔧 Auto-cleanup configuration updated")

	return nil
}

// HandleAutoCleanupConfig provides API endpoint for managing auto-cleanup
func HandleAutoCleanupConfig(w http.ResponseWriter, r *http.Request) {
	// Check admin authentication
	adminKey := r.Header.Get("X-Admin-Key")
	if adminKey != configdb.GetAdminKey() {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Return current configuration and status
		status := GetSchedulerStatus()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)

	case http.MethodPost:
		// Update configuration
		var config AutoCleanupConfig
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := UpdateSchedulerConfig(config); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Auto-cleanup configuration updated",
			"config":  config,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
