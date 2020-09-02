package req

const (
	UID         = "UID"
	LOGIN_KEY   = "LOGIN-KEY"
	LOGIN_TOKEN = "LOGIN-TOKEN"
)

type PageParam struct {
	LastId   int
	Page     int `json:"Page"`
	PageSize int `json:"PageSize"`
}

func (p *PageParam) Offset() int {
	return (p.Page - 1) * p.PageSize
}
