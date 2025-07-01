// handlers/slugstats.go  – FULL REWRITE
package handlers

import (
	"encoding/json"
	"net/http"

	"o365/configdb"
)

/*  /stats  – returns counters **only for the caller’s slug** */
func HandleSlugStats(w http.ResponseWriter, r *http.Request) {
	slug, ok := r.Context().Value("slug").(string)
	if !ok || slug == "" {
		http.Error(w, "slug missing", http.StatusBadRequest)
		return
	}

	statsMap := configdb.AllSlugStats() // map[string]SlugStats
	stat, exists := statsMap[slug]
	if !exists {
		stat = configdb.SlugStats{} // zero counts
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(stat)
}
