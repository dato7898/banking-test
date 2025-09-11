package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionRepository interface {
	Save(opaqueToken, jwt string, userID int) error
	Get(opaqueToken string) (string, int, error)
	Refresh(opaqueToken string) error
	Delete(token string) error
}

type sessionRepository struct {
	client     *redis.Client
	expiration time.Duration
}

func NewSessionRepository(addr string, password string, db int, expiration time.Duration) SessionRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &sessionRepository{
		client:     rdb,
		expiration: expiration,
	}
}

func (r *sessionRepository) Save(opaqueToken, jwt string, userID int) error {
	ctx := context.Background()
	pipe := r.client.Pipeline()

	pipe.HSet(ctx, tokenKey(opaqueToken),
		"user_id", strconv.Itoa(userID),
		"jwt", jwt,
	)
	pipe.Expire(ctx, tokenKey(opaqueToken), r.expiration)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *sessionRepository) Get(opaqueToken string) (string, int, error) {
	ctx := context.Background()
	vals, err := r.client.HGetAll(ctx, tokenKey(opaqueToken)).Result()
	if err != nil || len(vals) == 0 {
		return "", 0, errors.New("session not found")
	}

	userID, err := strconv.Atoi(vals["user_id"])
	if err != nil {
		return "", 0, err
	}

	return vals["jwt"], userID, nil
}

func (r *sessionRepository) Refresh(opaqueToken string) error {
	ctx := context.Background()
	return r.client.Expire(ctx, tokenKey(opaqueToken), r.expiration).Err()
}

func (r *sessionRepository) Delete(opaqueToken string) error {
	ctx := context.Background()
	return r.client.Del(ctx, tokenKey(opaqueToken)).Err()
}

func tokenKey(token string) string {
	return fmt.Sprintf("session:%s", token)
}
