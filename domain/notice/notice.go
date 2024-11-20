package notice

type NoticeType string

const (
	System    NoticeType = "system"
	AboutUs   NoticeType = "aboutus"
	Version   NoticeType = "version"
	Usage     NoticeType = "usage"
	Privacy   NoticeType = "privacy"
	Agreement NoticeType = "agreement"
)

func (nt NoticeType) String() string {
	return string(nt)
}

type Notice struct {
	NoticeType NoticeType
	Message    string
}
