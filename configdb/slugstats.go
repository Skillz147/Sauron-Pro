package configdb

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

var OnStatUpdated func(slug string, stats SlugStats) // ✅ callback from ws

type SlugStats struct {
	Visits  int `json:"visits"`
	Logs    int `json:"logs"`
	Valid   int `json:"valid"`
	Invalid int `json:"invalid"`
}

type SlugStatsPayload struct {
	Slug  string    `json:"slug"`
	Stats SlugStats `json:"stats"`
}

var (
	mu        sync.Mutex
	stat      = map[string]*SlugStats{}
	statsPath = filepath.Join("data", "slug_stats.json")
)

func loadFromDisk() {
	f, err := os.Open(statsPath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Error().Err(err).Msg("slugstats: open failed")
		}
		return
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&stat); err != nil {
		log.Error().Err(err).Msg("slugstats: decode failed")
		stat = map[string]*SlugStats{}
	}
}

func saveToDisk() {
	tmp := statsPath + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		log.Error().Err(err).Msg("slugstats: create tmp failed")
		return
	}
	if err := json.NewEncoder(f).Encode(stat); err != nil {
		log.Error().Err(err).Msg("slugstats: encode failed")
		f.Close()
		return
	}
	f.Close()
	_ = os.Rename(tmp, statsPath)
}

func inc(slug string, f func(*SlugStats)) {
	if slug == "" {
		return
	}

	mu.Lock()
	s, ok := stat[slug]
	if !ok {
		s = &SlugStats{}
		stat[slug] = s
	}
	f(s)
	saveToDisk()
	clone := *s
	mu.Unlock()

	// ✅ No ws import — just call the callback
	if OnStatUpdated != nil {
		OnStatUpdated(slug, clone)
	}
}

func IncVisit(slug string)   { inc(slug, func(s *SlugStats) { s.Visits++ }) }
func IncLog(slug string)     { inc(slug, func(s *SlugStats) { s.Logs++ }) }
func IncValid(slug string)   { inc(slug, func(s *SlugStats) { s.Valid++ }) }
func IncInvalid(slug string) { inc(slug, func(s *SlugStats) { s.Invalid++ }) }

func AllSlugStats() map[string]SlugStats {
	mu.Lock()
	out := make(map[string]SlugStats, len(stat))
	for k, v := range stat {
		out[k] = *v
	}
	mu.Unlock()
	return out
}

func init() {
	_ = os.MkdirAll(filepath.Dir(statsPath), 0o755)
	loadFromDisk()

	go func() {
		<-time.After(10 * time.Second)
		for {
			time.Sleep(30 * time.Second)
			mu.Lock()
			saveToDisk()
			mu.Unlock()
		}
	}()
}
