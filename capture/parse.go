// capture/parse.go
package capture

import (
	"strconv"
	"strings"
	"time"
)

/*
	----------------------------------------------------------------
	  ParseSetCookie converts one Set-Cookie header *as received from
	  the network* into a CookieEntry ready for SaveCookie().
	  It is deliberately liberal - anything unknown is ignored rather
	  than causing an error.

------------------------------------------------------------------
*/
func ParseSetCookie(sc string, ts time.Time) CookieEntry {
	var c CookieEntry

	// defaults – safe values for Cookie-Editor
	c.Path = "/"
	c.SameSite = "no_restriction"
	c.Secure = false
	c.HTTPOnly = false
	c.Session = true
	c.Timestamp = ts
	c.Raw = sc

	segments := strings.Split(sc, ";")

	/* ---- NAME=VALUE ------------------------------------------------ */
	if len(segments) > 0 {
		pair := strings.SplitN(strings.TrimSpace(segments[0]), "=", 2)
		if len(pair) == 2 {
			c.Name = pair[0]
			c.Value = pair[1]
		}
	}

	for _, seg := range segments[1:] {
		seg = strings.TrimSpace(seg)

		switch {
		case strings.HasPrefix(strings.ToLower(seg), "domain="):
			c.Domain = strings.TrimSpace(seg[7:])

		case strings.HasPrefix(strings.ToLower(seg), "path="):
			c.Path = strings.TrimSpace(seg[5:])

		case strings.HasPrefix(strings.ToLower(seg), "expires="):
			if t, err := time.Parse(time.RFC1123, strings.TrimSpace(seg[8:])); err == nil {
				c.ExpirationDate = float64(t.Unix())
				c.Session = false
			}

		case strings.HasPrefix(strings.ToLower(seg), "max-age="):
			if secs, err := strconv.Atoi(strings.TrimSpace(seg[8:])); err == nil {
				c.ExpirationDate = float64(ts.Add(time.Duration(secs) * time.Second).Unix())
				c.Session = false
			}

		case strings.EqualFold(seg, "secure"):
			c.Secure = true

		case strings.EqualFold(seg, "httponly"):
			c.HTTPOnly = true

		case strings.HasPrefix(strings.ToLower(seg), "samesite="):
			ss := strings.ToLower(strings.TrimSpace(seg[9:]))
			switch ss {
			case "lax", "strict":
				c.SameSite = ss
			default:
				c.SameSite = "no_restriction"
			}
		}
	}

	return c
}
