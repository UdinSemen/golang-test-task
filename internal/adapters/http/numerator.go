package http

type AddNumRequest struct {
	Num int64 `json:"num"`
}

type NumResponse struct {
	Nums []int64 `json:"nums"`
}
