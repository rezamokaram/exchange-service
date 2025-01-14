package service

import (
	// "context"
	// "errors"
	// "fmt"
	// "math/rand/v2"
	// "time"

	// "github.com/RezaMokaram/ExchangeService/api/pb"
	// notifDomain "github.com/RezaMokaram/ExchangeService/internal/notification/domain"
	notifPort "github.com/RezaMokaram/ExchangeService/internal/notification/port"
	// "github.com/RezaMokaram/ExchangeService/internal/user"
	// "github.com/RezaMokaram/ExchangeService/internal/user/domain"
	cryptoPort "github.com/RezaMokaram/ExchangeService/internal/crypto/port"
	// "github.com/RezaMokaram/ExchangeService/pkg/jwt"
	// timeutils "github.com/RezaMokaram/ExchangeService/pkg/time"
	// jwt2 "github.com/golang-jwt/jwt/v5"
)

type CryptoService struct {
	svc           cryptoPort.Service
	notifSvc      notifPort.Service
	authSecret    string
	expMin        uint
	refreshExpMin uint
}

func NewCryptoService(
	svc cryptoPort.Service,
	authSecret string,
	expMin uint,
	refreshExpMin uint,
	notifSvc notifPort.Service,
) *CryptoService {
	return &CryptoService{
		svc:           svc,
		authSecret:    authSecret,
		expMin:        expMin,
		refreshExpMin: refreshExpMin,
		notifSvc:      notifSvc,
	}
}
