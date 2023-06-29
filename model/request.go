package model

type Pagination struct {
	PageSize int `form:"pageSize"`
	Current  int `form:"current"`
}

func (p *Pagination) GetOffset() int {
	return p.PageSize * (p.Current - 1)
}
