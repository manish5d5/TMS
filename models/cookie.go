package models

import (
	"net/http"
	"strings"
)

type CookiesModel struct {
	AccessCookie  *http.Cookie
	RefreshCookie *http.Cookie
}

// ChangeSameSiteForDevelopment adjusts SameSite for local development.
func (c *CookiesModel) ChangeSameSiteForDevelopment(r *http.Request) {
	origin := r.Header.Get("Origin")
	if strings.Contains(origin, "localhost") {
		c.AccessCookie.SameSite = http.SameSiteNoneMode
		c.RefreshCookie.SameSite = http.SameSiteNoneMode
	}
}
