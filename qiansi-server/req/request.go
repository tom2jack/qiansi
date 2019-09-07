package req

type PageParam struct {
	LastId int
	Page int `binding:"min=1"`
	PageSize int `binding:"min=1,max=50"`
}

func (p *PageParam) Offset() int {
	return (p.Page - 1) * p.PageSize
}
