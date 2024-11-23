package shared

type BilliardGameType int

const (
	GameTypeUnkonwon BilliardGameType = iota
	ChineseEightBall
	Snooker
	Genting
	NineBall
)

func (gt BilliardGameType) String() string {
	switch gt {
	case ChineseEightBall:
		return "中式八球"
	case NineBall:
		return "九球追分"
	case Snooker:
		return "斯诺克"
	case Genting:
		return "云顶之弈"
	default:
		return "未知"
	}
}
