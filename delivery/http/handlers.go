package http

import (
	"CoalCompany/domain"
)

type HTTPHandlers struct {
	enterprise *domain.Enterprise
}

func NewHTTPHandlers(enterprise *domain.Enterprise) *HTTPHandlers {
	return &HTTPHandlers{
		enterprise: enterprise,
	}
}
