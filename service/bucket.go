package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

type BucketService[T any] struct {
	Rds       redis.Cmdable
	NewBucket func(owner string) string
	NewValue  func() T
}

func (s *BucketService[T]) List(ctx context.Context, owner string) (map[string]T, error) {
	bucket := s.NewBucket(owner)
	resp := s.Rds.HGetAll(ctx, bucket)
	if err := resp.Err(); err != nil {
		return nil, fmt.Errorf("redis get bucket all error. %w", err)
	}
	output := make(map[string]T, len(resp.Val()))
	for key, data := range resp.Val() {
		var value = s.NewValue()
		if err := sonic.Unmarshal([]byte(data), &value); err != nil {
			return nil, fmt.Errorf("unmarshal redis bucket value. %w", err)
		}
		output[key] = value
	}
	return output, nil
}

func (s *BucketService[T]) Store(ctx context.Context, owner, key string, value T) error {
	data, err := sonic.ConfigFastest.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal value to json error. %w", err)
	}
	bucket := s.NewBucket(owner)
	resp := s.Rds.HSet(ctx, bucket, key, string(data))
	if err := resp.Err(); err != nil {
		return fmt.Errorf("redis store value error. %w", err)
	}
	return nil
}

func (s *BucketService[T]) Load(ctx context.Context, owner, key string) (*T, error) {
	bucket := s.NewBucket(owner)
	resp := s.Rds.HGet(ctx, bucket, key)
	if err := resp.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("redis bucket load error. %w", resp.Err())
	}
	if resp.Val() != "" {
		var value = s.NewValue()
		err := sonic.ConfigFastest.Unmarshal([]byte(resp.Val()), &value)
		if err != nil {
			return nil, fmt.Errorf("unmarshal value error. %w", err)
		}
		return &value, nil
	}
	return nil, nil
}

func (s *BucketService[T]) Exists(ctx context.Context, owner, key string) (bool, error) {
	bucket := s.NewBucket(owner)
	resp := s.Rds.HExists(ctx, bucket, key)
	if resp.Err() != nil {
		return false, fmt.Errorf("redis bucket exists error. %w", resp.Err())
	}
	return resp.Val(), nil
}

func (s *BucketService[T]) Remove(ctx context.Context, owner, key string) error {
	bucket := s.NewBucket(owner)
	resp := s.Rds.HDel(ctx, bucket, key)
	if resp.Err() != nil {
		return fmt.Errorf("redis bucket remove error. %w", resp.Err())
	}
	return nil
}
