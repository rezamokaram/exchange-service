package service

import (
	// "context"
	// "errors"
	// "fmt"
	// "math/rand/v2"
	// "time"

	// "github.com/rezamokaram/exchange-service/api/pb"
	// notifDomain "github.com/rezamokaram/exchange-service/internal/notification/domain"
	"context"

	"github.com/rezamokaram/exchange-service/api/pb"
	notifPort "github.com/rezamokaram/exchange-service/internal/notification/port"

	// "github.com/rezamokaram/exchange-service/internal/user"
	"github.com/rezamokaram/exchange-service/internal/crypto/domain"
	cryptoPort "github.com/rezamokaram/exchange-service/internal/crypto/port"
	// "github.com/rezamokaram/exchange-service/pkg/jwt"
	// timeutils "github.com/rezamokaram/exchange-service/pkg/time"
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
