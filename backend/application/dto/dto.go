package dto

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func Success() *Resp {
	return &Resp{
		Code: 0,
		Msg:  "success",
	}
}

type PageParam struct {
	Page     int64 `form:"page" json:"page"`
	PageSize int64 `form:"pageSize" json:"pageSize"`
}

type IPageParam interface {
	UnWrap() (int64, int64)
}

func (p *PageParam) UnWrap() (int64, int64) {
	if p.Page < 0 {
		p.Page = 0
	}
	if p.PageSize <= 0 || p.PageSize > 100 {
		p.PageSize = 10
	}

	return p.Page, p.PageSize
}
