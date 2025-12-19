package models

type LoginResponse struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	UserID        int64  `json:"user_id"`
	AccessCookie  string `json:"access_cookie"`
	RefreshCookie string `json:"refresh_cookie"`
}
