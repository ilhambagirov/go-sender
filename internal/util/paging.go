package util

import "gorm.io/gorm"

type Paging struct {
	Limit int
	Page  int
}

func NewPaging(paging *Paging) *Paging {
	return &Paging{paging.Limit, paging.Page}
}

func (p *Paging) PaginatedResult(db *gorm.DB) *gorm.DB {
	offset := (p.Page - 1) * p.Limit
	return db.Offset(offset).Limit(p.Limit)
}
