package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/bccfilkom/drophere-go/domain"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	// A private key for context that only this package can access. This is important
	// to prevent collisions between different context uses
	userCtxKey      = &contextKey{"user"}
	errInvalidToken = errors.New("jwt: invalid token")
)

type contextKey struct {
	name string
}

// JWTAuthenticator struct
type JWTAuthenticator struct {
	key      []byte
	duration time.Duration
	algo     string
	userRepo domain.UserRepository
}

// NewJWT func
func NewJWT(secret string, duration time.Duration, algo string, userRepo domain.UserRepository) *JWTAuthenticator {
	return &JWTAuthenticator{
		key:      []byte(secret),
		duration: duration,
		algo:     algo,
		userRepo: userRepo,
	}
}

// Authenticate func
func (j *JWTAuthenticator) Authenticate(u *domain.User) (*domain.UserCredentials, error) {
	expiry := time.Now().Add(j.duration)
	token := jwt.NewWithClaims(jwt.GetSigningMethod(j.algo), jwt.MapClaims{
		"user_id": u.ID,
		"exp":     expiry.Unix(),
	})

	tokenS, err := token.SignedString(j.key)
	if err != nil {
		return nil, err
	}
	return &domain.UserCredentials{
		Token:  tokenS,
		Expiry: &expiry,
	}, nil
}

func (j *JWTAuthenticator) validateAndGetUserID(token string) (uint, error) {
	payloadI, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(j.algo) != token.Method {
			return nil, errInvalidToken
		}

		return j.key, nil
	})

	if err != nil {
		return 0, err
	}

	if !payloadI.Valid {
		return 0, errInvalidToken
	}

	claims := payloadI.Claims.(jwt.MapClaims)

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errInvalidToken
	}

	return uint(userID), nil
}

func writeGqlError(w http.ResponseWriter, msg string) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"errors": []map[string]string{
			{
				"message": msg,
			},
		},
	})
}

// Middleware func
func (j *JWTAuthenticator) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// cast inner function to HandlerFunc
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			// Allow unauthenticated users in
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			spaceIdx := strings.IndexByte(authHeader, ' ')
			if spaceIdx < 0 {
				writeGqlError(w, "Invalid Authorization header")
				return
			}
			authHeaderPrefix := authHeader[:spaceIdx]
			authToken := authHeader[spaceIdx+1:]

			if authHeaderPrefix != "bearer" && authHeaderPrefix != "Bearer" {
				writeGqlError(w, "Invalid Authorization header")
				return
			}

			userID, err := j.validateAndGetUserID(authToken)
			if err != nil {
				writeGqlError(w, "Invalid or expired token")
				return
			}

			// get the user from the database
			user, err := j.userRepo.FindByID(userID)
			if err != nil {
				writeGqlError(w, "Server Error")
				return
			}

			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, user)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// GetAuthenticatedUser finds the user from the context. REQUIRES Middleware to have run.
func (j *JWTAuthenticator) GetAuthenticatedUser(ctx context.Context) *domain.User {
	raw, _ := ctx.Value(userCtxKey).(*domain.User)
	return raw
}
