package redisclient

// package redisclient provides  interface for the
// functions which deals with the redis db operations.

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/kumareswaramoorthi/chat/udp-server/message"
)

type RedisClient interface {
	DeleteKey(m message.Message) error
	Flush() error
	SaveToDB(m message.Message) error
	RetrieveHistory() (string, error)
}

type redisClient struct {
	redis *redis.Client
	ctx   context.Context
}

// NewRedisClient returns a redisClient object.
func NewRedisClient(ctx context.Context, redis *redis.Client) RedisClient {
	return &redisClient{
		ctx:   ctx,
		redis: redis,
	}
}

// RetrieveHistory retrieves the last 20 messages stored in redis.
func (r *redisClient) RetrieveHistory() (string, error) {
	val, err := r.redis.LRange(r.ctx, "messages", 0, -1).Result()
	if err != nil {
		return "", err
	}

	//return the message history seperated by new line.
	return strings.Join(val, "\n"), nil
}

// Flush deletes the entire data in redis.
func (r *redisClient) Flush() error {
	err := r.redis.FlushDB(r.ctx).Err()
	if err != nil {
		return err
	}
	return nil
}

// DeleteKey deletes the actual message from the "messages" list in
// redis and also deletes it from the user hashes
func (r *redisClient) DeleteKey(m message.Message) error {

	// get the actual message from the users's hash by the message-id.
	msg, err := r.redis.HGet(r.ctx, "user-name="+m.UserName, m.MessageID).Result()
	if err != nil {
		return err
	}

	// remove the message from the list
	err = r.redis.LRem(r.ctx, "messages", 0, msg).Err()
	if err != nil {
		return err
	}

	// remove the message from user's hash.
	err = r.redis.HDel(r.ctx, "user-name="+m.UserName, m.MessageID).Err()
	if err != nil {
		return err
	}

	// return nil if there are no errors.
	return nil
}

// SaveToDB saves the message to the 'messages' list in redis and also adds message to user's hash in redis.
func (r *redisClient) SaveToDB(m message.Message) error {

	// append the new message to the 'messages' list in redis.
	err := r.redis.RPush(r.ctx, "messages", fmt.Sprintf("time=%s|user-name=%s|message=%s", m.TimeStamp, m.UserName, m.Content)).Err()
	if err != nil {
		return err
	}

	// add the message to the user's hash in redis.
	err = r.redis.HSet(r.ctx, "user-name="+m.UserName, m.MessageID, fmt.Sprintf("time=%s|user-name=%s|message=%s", m.TimeStamp, m.UserName, m.Content)).Err()
	if err != nil {
		return err
	}

	// if the total message count is greater than 20 the delete the top message.
	messageCount, err := r.redis.LLen(r.ctx, "messages").Result()
	if err != nil {
		return err
	}
	if messageCount > 20 {
		errs := r.redis.LPop(r.ctx, "messages").Err()
		if errs != nil {
			return errs
		}
	}
	return nil
}
