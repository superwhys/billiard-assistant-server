package shared

type SaGameType int

const (
	ChineseEightBall SaGameType = iota
	Snooker
	Genting
)

func (gt SaGameType) String() string {
	switch gt {
	case ChineseEightBall:
		return "中式八球"
	case Snooker:
		return "斯诺克"
	case Genting:
		return "云顶之弈"
	default:
		return "未知"
	}
}
