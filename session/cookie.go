package session

import (
	"net/http"
	"time"
)

func AddSessionIdCookie(sessionId string, expireDate time.Time, writer http.ResponseWriter) {
	cookieSessionId := http.Cookie{
		Name:     "SessionId",
		Value:    sessionId,
		Expires:  expireDate,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(writer, &cookieSessionId)
}

func DeleteSessionIdCookie(writer http.ResponseWriter) {
	http.SetCookie(writer,
		&http.Cookie{
			Name:    "SessionId",
			Value:   "",
			Path:    "/",
			Expires: time.Now(),
		})
}
