package rest

type ErrorResponse struct {
	Err string `json:"error"`
}

type SuccessResponse struct {
	Status    bool   `json:"status"`
	Operation string `json:"operation"`
}
