package dto

type Pagination struct {
	Total    int
	Current  int
	PageSize int
}

type PaginationParam struct {
	Current  int
	PageSize int
}

func (a *PaginationParam) GetCurrent() int {
	return a.Current
}

func (a *PaginationParam) GetPageSize() int {
	pageSize := a.PageSize
	if a.PageSize == 0 {
		pageSize = 15
	}

	return pageSize
}
