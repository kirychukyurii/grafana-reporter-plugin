package model

type Pagination struct {
	Total    int
	Current  int
	PageSize int
}

type PaginationParam struct {
	Current  int
	PageSize int
}
