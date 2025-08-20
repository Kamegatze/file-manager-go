package repository

import "fmt"

type Pageable interface {
	Offset() int
	Limit() int
	OrderBy() []string
}

type PageableImpl struct {
	page    int
	limit   int
	orderBy []string
}

func NewPageableImpl(page int, limit int, orderBy []string) (PageableImpl, error) {
	pageable := PageableImpl{}
	if err := pageable.WithLimit(limit); err != nil {
		return pageable, err
	}
	if err := pageable.WithPage(page); err != nil {
		return pageable, err
	}
	pageable.WithOrderBy(orderBy)
	return pageable, nil
}

func (pageable *PageableImpl) WithPage(page int) error {
	if page < 1 {
		return fmt.Errorf("field 'page' must be more 0, but have: %d", page)
	}
	pageable.page = page
	return nil
}

func (pageable *PageableImpl) WithLimit(limit int) error {
	if limit < 1 {
		return fmt.Errorf("field 'limit' must be more 0, but have: %d", limit)
	}
	pageable.limit = limit
	return nil
}

func (pageable *PageableImpl) WithOrderBy(orderBy []string) {
	pageable.orderBy = orderBy
}

func (pageable PageableImpl) Offset() int {
	return (pageable.page - 1) * pageable.limit
}

func (pageable PageableImpl) Limit() int {
	return pageable.limit
}

func (pageable PageableImpl) OrderBy() []string {
	return pageable.orderBy
}
