package ytdlp

import (
	"net/url"
	"os"
	"strings"

	"github.com/ayamine47/sumika/lib/config"
)

func getCookiePath(u *url.URL) string {
	host := strings.ToLower(u.Hostname())

	cookiePath, exists := config.CurrentConfig.CookieList[host]
	if !exists {
		return ""
	}

	if _, err := os.Stat(cookiePath); err != nil {
		return ""
	}

	return cookiePath
}
