package requests

type CreateOrder struct {
	Number       string  `json:"number"       binding:"required"`
	Total        float64 `json:"total"        binding:"required,gt=0"`
	CustomerName string  `json:"customerName" binding:"required"`
}
