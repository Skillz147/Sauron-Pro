package ws

import "sync"

// slugToUser maps slug → userID
var (
	slugToUser = make(map[string]string)
	mu         sync.RWMutex
)

// RegisterSlug sets the mapping from slug to userID
func RegisterSlug(slug string, userID string) {
	mu.Lock()
	defer mu.Unlock()
	slugToUser[slug] = userID
}

// GetUserIDForSlug returns the userID for a given slug (if exists)
func GetUserIDForSlug(slug string) (string, bool) {
	mu.RLock()
	defer mu.RUnlock()
	userID, ok := slugToUser[slug]
	return userID, ok
}

// UnregisterSlug removes the mapping for a slug (optional)
func UnregisterSlug(slug string) {
	mu.Lock()
	defer mu.Unlock()
	delete(slugToUser, slug)
}
