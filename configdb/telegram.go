// configdb/telegram.go
package configdb

import (
	"database/sql"
	"errors"
)

// TelegramSettings bundles the two pieces the result-builder needs.
type TelegramSettings struct {
	UserID   string
	BotToken string
	ChatID   string
}

/*
GetTelegramSettingsBySlug returns the bot-token + chat-ID that belong to the
operator who owns <slug>.  (We first translate slug ➜ user_id, then read the
telegram_settings table.)

	settings, err := configdb.GetTelegramSettingsBySlug(slug)
*/
func GetTelegramSettingsBySlug(slug string) (*TelegramSettings, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// make sure the table exists (safe no-op if it’s already there)
	if err := ensureLinksTable(db); err != nil {
		return nil, err
	}

	var tg TelegramSettings
	// look up the owner of the slug, then fetch their Telegram row
	err = db.QueryRow(`
		SELECT ts.user_id, ts.bot_token, ts.chat_id
		  FROM user_links  ul
		  JOIN telegram_settings ts ON ts.user_id = ul.user_id
		 WHERE ul.slug = ? LIMIT 1`, slug).
		Scan(&tg.UserID, &tg.BotToken, &tg.ChatID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // slug exists but no Telegram record yet
	}
	if err != nil {
		return nil, err // genuine DB failure
	}
	return &tg, nil
}
