package notice

type NoticeType string

const (
	System    NoticeType = "system"
	AboutUs   NoticeType = "aboutus"
	Version   NoticeType = "version"
	Usage     NoticeType = "usage"
	Privacy   NoticeType = "privacy"
	Feedback  NoticeType = "feedback"
	Agreement NoticeType = "agreement"
)

func (nt NoticeType) String() string {
	return string(nt)
}

type Notice struct {
	NoticeType NoticeType
	Message    string
}
