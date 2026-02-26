package cookies

import (
	"net/http"
)

const (
	RefreshTokenCookieName = "refresh_token"
	RefreshTokenMaxAge     = 7 * 24 * 60 * 60 // 7days
)

func SaveRefreshToken(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   RefreshTokenMaxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
}

func RemoveRefreshToken(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
}

func GetRefreshToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(RefreshTokenCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
