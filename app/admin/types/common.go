package types

type IDReq struct {
	ID uint `json:"id" form:"id" param:"id" uri:"id" query:"id"`
}
