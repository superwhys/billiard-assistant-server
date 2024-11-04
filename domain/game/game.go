package game

type Game struct {
	GameId     int
	GameType   SaGameType
	GameConfig *Config
}

type Config struct {
	MaxPlayers int
	Desc       string
}

func (c *Config) GetMaxPlayers() int {
	return c.MaxPlayers
}

type SaGameType int

const (
	ChineseEightBall SaGameType = iota
	Snooker
	YunDing
)

func (gt SaGameType) String() string {
	switch gt {
	case ChineseEightBall:
		return "中式八球"
	case Snooker:
		return "斯诺克"
	case YunDing:
		return "云顶之弈"
	default:
		return "未知"
	}
}
