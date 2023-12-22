package repository

import (
	"backend/modules/auth/helpers/apperrors"
	"backend/modules/auth/models"

	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisTokenRepository struct {
	Redis *redis.Client
}

func NewRedisTokenRepository(redisClient *redis.Client) models.TokenRepository {
	return &redisTokenRepository{
		Redis: redisClient,
	}
}

func (redisRepo *redisTokenRepository) SetRefreshToken(ctx context.Context, accountIDCode string, tokenID string, tokenExpireTime time.Duration) error {
	// We store the account IDCode with the token id so we can scan
	// over the account's tokens without blocking and delete them
	// inc case of token leakage
	key := fmt.Sprintf("%s:%s", accountIDCode, tokenID)
	if err := redisRepo.Redis.Set(ctx, key, 0, tokenExpireTime).Err(); err != nil {
		log.Printf("Could not SET the Redis refresh token for account:%s with token ID:%s\n", accountIDCode, tokenID)
		return apperrors.NewInternal()
	}
	return nil
}

func (redisRepo *redisTokenRepository) DeleteRefreshToken(ctx context.Context, accountIDCode string, tokenID string) error {
	key := fmt.Sprintf("%s:%s", accountIDCode, tokenID)
	result := redisRepo.Redis.Del(ctx, key)

	if err := result.Err(); err != nil {
		log.Printf("Could not DELETE Redis refresh token to for accountIDCode:tokenID: %s:%s: %v\n", accountIDCode, tokenID, err)
		return apperrors.NewInternal()
	}

	// Val returns the count of deleted keys. If no keys are deleted then the Refresh Token is invalid.
	if result.Val() < 1 {
		log.Printf("Refresh Token to redis for accountIDCode:tokenID %s:%s does not exist\n", accountIDCode, tokenID)
		return apperrors.NewAuthorization("Invalid refresh token")
	}

	return nil
}

func (redisRepo *redisTokenRepository) DeleteAccountRefreshTokens(ctx context.Context, accountID string) error {
	pattern := fmt.Sprintf("%s*", accountID)
	log.Print(pattern)

	iterator := redisRepo.Redis.Scan(ctx, 0, pattern, 5).Iterator()
	failCount := 0

	for iterator.Next(ctx) {
		if err := redisRepo.Redis.Del(ctx, iterator.Val()).Err(); err != nil {
			log.Printf("Faield to delete refresh token when: %s\n", iterator.Val())
			failCount++
		}
	}

	// check the last value
	if err := iterator.Err(); err != nil {
		log.Printf("Failed to delete refresh token: %s\n", iterator.Val())
		failCount++
	}

	if failCount > 0 {
		return apperrors.NewInternal()
	}

	return nil
}
