package repository

import (
	"errors"
	"strconv"
	"time"
)
const (
	tokenFolderPrefix = "user_tokens:"
	accessTokenKey    = "access_token"
	refreshTokenKey   = "refresh_token"
)

func (r *Repository) SaveJWTTokenPair(userID uint, accessToken, refreshToken string) error {
	expiration := 7 * 24 * time.Hour

	userIDStr := strconv.FormatUint(uint64(userID), 10)
	tokenFolderKey := tokenFolderPrefix + userIDStr

	err := r.rd.HSet(tokenFolderKey, accessTokenKey, accessToken).Err()
	if err != nil {
		return err
	}

	err = r.rd.HSet(tokenFolderKey, refreshTokenKey, refreshToken).Err()
	if err != nil {
		// If saving Refresh Token fails, remove the saved Access Token
		r.rd.HDel(tokenFolderKey, accessTokenKey)
		return err
	}

	// Set expiration for the token folder
	err = r.rd.Expire(tokenFolderKey, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteJWTTokenPair(userID uint) error {
	userIDStr := strconv.FormatUint(uint64(userID), 10)
	tokenFolderKey := tokenFolderPrefix + userIDStr

	err := r.rd.Del(tokenFolderKey).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CheckTokenInRedis(userID uint, accessToken string) error {
	userIDStr := strconv.FormatUint(uint64(userID), 10)
	tokenFolderKey := tokenFolderPrefix + userIDStr
	accessTokenKey := "access_token"

	storedToken, err := r.rd.HGet(tokenFolderKey, accessTokenKey).Result()
	if err != nil || storedToken != accessToken {
		return errors.New("токен не действителен")
	}

	return nil
}

