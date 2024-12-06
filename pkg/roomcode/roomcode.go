package roomcode

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-puzzles/puzzles/predis"
	"github.com/go-puzzles/puzzles/putils"
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

func (g *RoomCodeGenerator) GenerateCode(ctx context.Context, roomId int) (string, error) {
	if err := g.Lock(roomId); err != nil {
		return "", errors.Wrap(err, "lock room failed")
	}
	defer g.Unlock(roomId)

	existingCode, err := g.GetRoomCode(ctx, roomId)
	if err == nil {
		return existingCode, nil
	}

	roomIdMarshal, err := json.Marshal(roomId)
	if err != nil {
		return "", errors.Wrap(err, "marshal roomId")
	}

	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		code := putils.RandNumeral(6)
		codeKey := g.getCodeKey(code)
		roomKey := g.getRoomKey(roomId)

		codeMarshal, err := json.Marshal(code)
		if err != nil {
			return "", errors.Wrap(err, "marshal code")
		}

		commands := [][]any{
			{"SETEX", codeKey, int(g.expiration.Seconds()), roomIdMarshal},
			{"SETEX", roomKey, int(g.expiration.Seconds()), codeMarshal},
		}

		if err := g.redisClient.DoWithTransactionPipeline(nil, commands...); err != nil {
			continue
		}

		return code, nil
	}

	return "", errors.New("无法生成唯一房间码，请稍后重试")
}

func (g *RoomCodeGenerator) GetRoomId(ctx context.Context, code string) (int, error) {
	key := g.getCodeKey(code)

	var roomId int
	err := g.redisClient.Get(key, &roomId)
	if err != nil {
		return 0, errors.Wrap(err, "房间码无效或已过期")
	}

	return roomId, nil
}

func (g *RoomCodeGenerator) GetRoomCode(ctx context.Context, roomId int) (string, error) {
	key := g.getRoomKey(roomId)

	var code string
	err := g.redisClient.Get(key, &code)
	if err != nil {
		return "", errors.Wrap(err, "获取房间码失败")
	}

	return code, nil
}

func (g *RoomCodeGenerator) GetRoomCodeBatch(ctx context.Context, roomIds ...int) ([]string, error) {
	keys := putils.Map(roomIds, func(key int) string {
		return g.getRoomKey(key)
	})

	codeResp := make(map[string]any)
	err := g.redisClient.MGet(keys, codeResp)
	if err != nil {
		return nil, errors.Wrap(err, "multiGetCode")
	}

	codes := putils.Map(keys, func(s string) string {
		c := codeResp[s]
		if c == nil {
			return ""
		}
		return codeResp[s].(string)
	})

	return codes, nil
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

func (g *RoomCodeGenerator) getCodeKey(code string) string {
	return fmt.Sprintf("%s:%s", roomCodePrefix, code)
}

func (g *RoomCodeGenerator) getRoomKey(roomId int) string {
	return fmt.Sprintf("%s:room:%d", roomCodePrefix, roomId)
}

func (g *RoomCodeGenerator) getLockKey(roomId int) string {
	return fmt.Sprintf("%s:%d", lockPrefix, roomId)
}
