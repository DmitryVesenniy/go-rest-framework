package common

import (
	"encoding/json"
	"net/http"
)

type RefreshToken struct {
	RefreshToken string `json:"refreshToken" serializer:"required"`
}

func GetRefreshToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(REFRESH_TOKEN_KEY)
	var refreshToken string
	if err != nil {
		refreshTokenSerializer := RefreshToken{}
		err := json.NewDecoder(r.Body).Decode(&refreshTokenSerializer)
		if err != nil {
			return "", err
		}
		refreshToken = refreshTokenSerializer.RefreshToken
	} else {
		refreshToken = cookie.Value
	}

	return refreshToken, nil
}
