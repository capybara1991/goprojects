package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"auth-service/internal/repository"
)

type AuthService struct {
	jwtSecret  string
	jwtTTL     time.Duration
	tokenBytes int
	refreshTTL time.Duration
	bcryptCost int
	webhookURL string
	store      *repository.Store
}

func NewAuthService(secret string, jwtMin, tokBytes, refDays, cost int, webhook string, store *repository.Store) *AuthService {
	return &AuthService{
		jwtSecret:  secret,
		jwtTTL:     time.Minute * time.Duration(jwtMin),
		tokenBytes: tokBytes,
		refreshTTL: time.Hour * 24 * time.Duration(refDays),
		bcryptCost: cost,
		webhookURL: webhook,
		store:      store,
	}
}
func (a *AuthService) Issue(userID uuid.UUID, ua, ip string) (access, refresh string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(a.jwtTTL).Unix(),
	})
	access, err = token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return
	}
	raw := make([]byte, a.tokenBytes)
	_, err = rand.Read(raw)
	if err != nil {
		return
	}
	refresh = base64.StdEncoding.EncodeToString(raw)

	hash, err := bcrypt.GenerateFromPassword(raw, a.bcryptCost)
	if err != nil {
		return
	}
	rt := &repository.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Hash:      string(hash),
		UserAgent: ua,
		IP:        ip,
		IssuedAt:  time.Now(),
	}
	err = a.store.Save(rt)
	return
}

func (a *AuthService) Refresh(oldAccess, rawRefresh, ua, ip string) (newAccess, newRefresh string, err error) {
	token, _, err := new(jwt.Parser).ParseUnverified(oldAccess, jwt.MapClaims{})
	if err != nil {
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return
	}
	rt, err := a.store.FindValid(userID)
	if err != nil {
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(rt.Hash), []byte(rawRefresh)); err != nil {
		a.store.RevokeAll(userID)
		return "", "", errors.New("invalid refresh token")
	}
	if rt.UserAgent != ua {
		a.store.RevokeAll(userID)
		return "", "", errors.New("user-agent changed, session revoked")
	}

	if rt.IP != ip && a.webhookURL != "" {
		go http.Post(a.webhookURL, "application/json", nil)
	}
	_ = a.store.MarkUsed(rt.ID)

	return a.Issue(userID, ua, ip)
}

func (a *AuthService) Logout(userID uuid.UUID) error {
	return a.store.RevokeAll(userID)
}
