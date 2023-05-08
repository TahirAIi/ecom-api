package main

import "net/http"

func (app application) IsAdmin() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := GetBearerToken(r)
			if len(token) < 0 || !app.IsAdminUser(token) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (app application) IsAdminUser(token string) bool {
	jwtClaims, err := Validate(token)
	if err != nil {
		return false
	}

	user, err := app.models.User.GetByEmail(jwtClaims.Email)
	if err != nil {
		return false
	}

	return user.IsAdmin

}
