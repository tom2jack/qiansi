package req

type PageParam struct {
	LastId   int
	Page     int `json:"Page"`
	PageSize int `json:"PageSize"`
}

func (p *PageParam) Offset() int {
	return (p.Page - 1) * p.PageSize
}
