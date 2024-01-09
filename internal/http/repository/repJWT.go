package repository

import (
	"errors"
	"strconv"
	"time"
)
const (
	tokenFolderPrefix = "user_tokens:"
	accessTokenKey    = "access_token"
)

func (r *Repository) SaveJWTToken(userID uint, accessToken string) error {
	expiration := 7 * 24 * time.Hour

	userIDStr := strconv.FormatUint(uint64(userID), 10)
	tokenFolderKey := tokenFolderPrefix + userIDStr

	err := r.rd.HSet(tokenFolderKey, accessTokenKey, accessToken).Err()
	if err != nil {
		return err
	}
	
	err = r.rd.Expire(tokenFolderKey, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteJWTToken(userID uint) error {
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

