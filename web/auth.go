package web

import (
	"crypto/rsa"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bottleneckco/statuses-backend/db"
	"github.com/bottleneckco/statuses-backend/model"
	"github.com/lestrrat/go-jwx/jwk"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type payloadOAuth struct {
	GrantType     string `json:"grant_type"`
	GoogleIDToken string `json:"google_id_token"`
	RefreshToken  string `json:"refresh_token"`
}

type responseOAuth struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

const (
	jwksURL = "https://www.googleapis.com/oauth2/v3/certs"
)

var jwtPrivateKey *rsa.PrivateKey

func init() {
	var err error
	jwtPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(os.Getenv("JWT_PRIVATE_KEY")))
	if err != nil {
		log.Panic(err)
		return
	}
}

func getGoogleJWK(token *jwt.Token) (interface{}, error) {

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

func oauth(c *gin.Context) {
	var payload payloadOAuth
	err := c.BindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": "bad payload"})
		return
	}

	switch payload.GrantType {
	case "password":
		// Verify and parse Google ID Token
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(payload.GoogleIDToken, claims, getGoogleJWK)
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
		// Create DB user
		user := model.User{
			Name:    claims["name"].(string),
			Email:   claims["email"].(string),
			Picture: claims["picture"].(string),
			Token:   payload.GoogleIDToken,
		}
		err = db.DB.FirstOrCreate(&user, model.User{Email: claims["email"].(string)}).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "error occurred"})
			log.Println(err)
			return
		}
		generateTokenPair(c, user.Email)
		break
	case "refresh_token":
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(payload.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
			return &jwtPrivateKey.PublicKey, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": "bad payload"})
			log.Println(err)
			return
		}
		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": "bad payload"})
			log.Println("token invalid")
			return
		}
		generateTokenPair(c, claims["username"].(string))
		break
	}
}

func generateTokenPair(c *gin.Context, username string) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Minute * 30).Unix(),
	})
	accessTokenString, err := accessToken.SignedString(jwtPrivateKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "error occurred"})
		log.Println(err)
		return
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"username": username,
	})

	refreshTokenString, err := refreshToken.SignedString(jwtPrivateKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "error occurred"})
		log.Println(err)
		return
	}

	user := model.User{}
	err = db.DB.Where("email = ?", username).First(&user).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "error occurred"})
		log.Println(err)
		return
	}
	err = db.DB.Model(&user).Update("token", refreshTokenString).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "error occurred"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, responseOAuth{
		TokenType:    "bearer",
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    30,
	})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return &jwtPrivateKey.PublicKey, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": "bad payload"})
			log.Println(err)
			return
		}
		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "unauthorised"})
			return
		}
		user := model.User{}
		err = db.DB.Where("email = ?", claims["username"].(string)).First(&user).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "error occurred"})
			log.Println(err)
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

// jwks render JWKS to API consumers
func jwks(c *gin.Context) {
	jwk, err := jwk.New(&jwtPrivateKey.PublicKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "error occurred"})
		log.Println(err)
		return
	}
	// Compute kid
	jwk.Set("alg", "RS256")
	jwk.Set("kid", "placeholder")
	c.JSON(http.StatusOK, gin.H{"keys": []interface{}{
		jwk,
	}})
}
