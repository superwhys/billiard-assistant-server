package roomcode

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-puzzles/puzzles/predis"
	"github.com/pkg/errors"
)

const (
	roomCodePrefix = "billiard:room:code"
	lockPrefix     = "billiard:room:code:lock"
)

type RoomCodeGenerator struct {
	redisClient *predis.RedisClient
	expiration  time.Duration
}

func NewRoomCodeGenerator(redisClient *predis.RedisClient, expiration time.Duration) *RoomCodeGenerator {
	return &RoomCodeGenerator{
		redisClient: redisClient,
		expiration:  expiration,
	}
}

func (g *RoomCodeGenerator) Lock(roomId int) error {
	return g.redisClient.LockWithBlock(g.getLockKey(roomId), 10, g.expiration)
}

func (g *RoomCodeGenerator) Unlock(roomId int) error {
	return g.redisClient.UnLock(g.getLockKey(roomId))
}

func (g *RoomCodeGenerator) GenerateCode(ctx context.Context, roomId int) (int, error) {
	if err := g.Lock(roomId); err != nil {
		return 0, errors.Wrap(err, "lock room failed")
	}
	defer g.Unlock(roomId)

	existingCode, err := g.GetRoomCode(ctx, roomId)
	if err == nil {
		return existingCode, nil
	}

	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		code := rand.Intn(1000000)
		codeKey := g.getCodeKey(code)
		roomKey := g.getRoomKey(roomId)

		commands := [][]any{
			{"SETEX", codeKey, int(g.expiration.Seconds()), roomId},
			{"SETEX", roomKey, int(g.expiration.Seconds()), code},
		}

		if err := g.redisClient.DoWithTransactionPipeline(nil, commands...); err != nil {
			continue
		}

		return code, nil
	}

	return 0, errors.New("无法生成唯一房间码，请稍后重试")
}

func (g *RoomCodeGenerator) GetRoomId(ctx context.Context, code int) (int, error) {
	key := g.getCodeKey(code)

	var roomId int
	err := g.redisClient.Get(key, &roomId)
	if err != nil {
		return 0, errors.Wrap(err, "房间码无效或已过期")
	}

	return roomId, nil
}

func (g *RoomCodeGenerator) GetRoomCode(ctx context.Context, roomId int) (int, error) {
	key := g.getRoomKey(roomId)

	var code int
	err := g.redisClient.Get(key, &code)
	if err != nil {
		return 0, errors.Wrap(err, "获取房间码失败")
	}

	return code, nil
}

func (g *RoomCodeGenerator) DeleteCode(ctx context.Context, roomId int) error {
	if err := g.Lock(roomId); err != nil {
		return errors.Wrap(err, "lock room failed")
	}
	defer g.Unlock(roomId)

	code, err := g.GetRoomCode(ctx, roomId)
	if err != nil {
		return err
	}

	codeKey := g.getCodeKey(code)
	roomKey := g.getRoomKey(roomId)

	commands := [][]any{
		{"DEL", codeKey},
		{"DEL", roomKey},
	}

	return g.redisClient.DoWithTransactionPipeline([]string{codeKey, roomKey}, commands...)
}

func (g *RoomCodeGenerator) getCodeKey(code int) string {
	return fmt.Sprintf("%s:%d", roomCodePrefix, code)
}

func (g *RoomCodeGenerator) getRoomKey(roomId int) string {
	return fmt.Sprintf("%s:room:%d", roomCodePrefix, roomId)
}

func (g *RoomCodeGenerator) getLockKey(roomId int) string {
	return fmt.Sprintf("%s:%d", lockPrefix, roomId)
}
