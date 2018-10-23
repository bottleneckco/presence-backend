package web

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/bottleneckco/statuses-backend/db"
	"github.com/bottleneckco/statuses-backend/model"
	"github.com/lestrrat/go-jwx/jwk"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type payloadLogin struct {
	Profile struct {
		GoogleID   string `json:"googleId"`
		ImageURL   string `json:"imageUrl"`
		Email      string `json:"email"`
		Name       string `json:"name"`
		GivenName  string `json:"givenName"`
		FamilyName string `json:"familyName"`
	} `json:"profileObj"`
	Token struct {
		AccessToken   string `json:"access_token"`
		ExpiresAt     int64  `json:"expires_at"`
		ExpiresIn     int64  `json:"expires_in"`
		FirstIssuedAt int64  `json:"first_issued_at"`
		IDToken       string `json:"id_token"`
		IDPID         string `json:"idpId"`
		LoginHint     string `json:"login_hint"`
		Scope         string `json:"scope"`
		TokenType     string `json:"token_type"`
	} `json:"tokenObj"`
}

const (
	jwksURL = "https://www.googleapis.com/oauth2/v3/certs"
)

func getKey(token *jwt.Token) (interface{}, error) {

	// TODO: cache response so we don't have to make a request every time
	// we want to verify a JWT
	set, err := jwk.FetchHTTP(jwksURL)
	if err != nil {
		return nil, err
	}

	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("expecting JWT header to have string kid")
	}

	if key := set.LookupKeyID(keyID); len(key) == 1 {
		return key[0].Materialize()
	}

	return nil, errors.New("unable to find key")
}

// Login handle login requests
func Login(c *gin.Context) {
	var payload payloadLogin
	err := c.BindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": "bad payload"})
		return
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(payload.Token.IDToken, claims, getKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": "bad payload"})
		log.Println(err)
		return
	}
	if claims["aud"] != os.Getenv("GOOGLE_CLIENT_ID") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": "bad payload"})
		log.Println("aud claim mismatch")
		return
	}
	if !token.Valid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": "bad payload"})
		log.Println("token invalid")
		return
	}
	// Google user verified!
	user := model.User{
		Name:    payload.Profile.Name,
		Email:   payload.Profile.Email,
		Picture: payload.Profile.ImageURL,
		Token:   payload.Token.IDToken,
	}
	db.DB.FirstOrCreate(&user, model.User{Email: payload.Profile.Email})
}
