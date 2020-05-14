package req

type PageParam struct {
	LastId   int
	Page     int `binding:"min=1" json:"Page"`
	PageSize int `binding:"min=1,max=50" json:"PageSize"`
}

func (p *PageParam) Offset() int {
	return (p.Page - 1) * p.PageSize
}
