package service

import (
	// "context"
	// "errors"
	// "fmt"
	// "math/rand/v2"
	// "time"

	// "github.com/RezaMokaram/ExchangeService/api/pb"
	// notifDomain "github.com/RezaMokaram/ExchangeService/internal/notification/domain"
	"context"

	"github.com/RezaMokaram/ExchangeService/api/pb"
	notifPort "github.com/RezaMokaram/ExchangeService/internal/notification/port"
	// "github.com/RezaMokaram/ExchangeService/internal/user"
	"github.com/RezaMokaram/ExchangeService/internal/crypto/domain"
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

func (s *CryptoService) CreateCrypto(ctx context.Context, req *pb.CreateCryptoRequest) (*pb.CreateCryptoResponse, error) {
	cryptoID, err := s.svc.CreateCrypto(ctx, domain.Crypto{
		Name:         req.GetName(),
		Symbol:       req.GetSymbol(),
		CurrentPrice: req.GetCurrentPrice(),
		BuyFee:       req.GetCurrentPrice() + (req.GetCurrentPrice() * 10 / 100),
		SellFee:      req.GetCurrentPrice() - (req.GetCurrentPrice() * 10 / 100),
	})

	if err != nil {
		return nil, err
	}

	return &pb.CreateCryptoResponse{
		Id: uint64(cryptoID),
	}, nil
}
