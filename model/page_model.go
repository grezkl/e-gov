package model

type PageForm struct {
	PageNum  int    `form:"pageNum" json:"pageNum" binding:"required"`
	PageSize int    `form:"pageSize" json:"pageSize" binding:"required"`
	Name     string `form:"name" json:"name"`
}

type PageResp struct {
	Records interface{} `form:"records" json:"records"`
	Total   int64       `form:"total" json:"total"`
}
