package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	httpServer "github.com/RezaMokaram/ExchangeService/api/handlers/http"
	"github.com/RezaMokaram/ExchangeService/api/pb"
	"github.com/RezaMokaram/ExchangeService/app"
	"github.com/RezaMokaram/ExchangeService/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserSignUpAndSignIn(t *testing.T) {
	cfg := config.MustReadConfig[config.ExchangeConfig]("./user_test_config.yaml")
	appContainer := app.NewMustApp(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		require.NoError(t, httpServer.Run(appContainer, cfg.Server), "Expected no error")
	}(ctx)
	time.Sleep(time.Second)

	signUpPayload := pb.UserSignUpRequest{
		FirstName: "test",
		LastName:  "test",
		Phone:     "+989191111111",
		Password:  "test",
	}

	jsonData, err := json.Marshal(&signUpPayload)
	require.NoError(t, err, "Expected no error")

	signUpURL := fmt.Sprintf("http://%s:%d/api/v1/sign-up", cfg.Server.HttpHost, cfg.Server.HttpPort)
	req, err := http.NewRequest("POST", signUpURL, bytes.NewBuffer(jsonData))
	require.NoError(t, err, "Expected no error")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	signUpResp, err := client.Do(req)
	require.NoError(t, err, "Expected no error")
	defer signUpResp.Body.Close()
	assert.Equal(t, http.StatusOK, signUpResp.Status, "Expected http status ok")

	signInPayload := pb.UserSignInRequest{
		Phone:    signUpPayload.Phone,
		Password: signUpPayload.Password,
	}
	signInURL := fmt.Sprintf("http://%s:%d/api/v1/sign-in", cfg.Server.HttpHost, cfg.Server.HttpPort)
	jsonData, err = json.Marshal(&signInPayload)
	require.NoError(t, err, "Expected no error")
	req, err = http.NewRequest("POST", signInURL, bytes.NewBuffer(jsonData))
	require.NoError(t, err, "Expected no error")
	req.Header.Set("Content-Type", "application/json")
	signInResp, err := client.Do(req)
	require.NoError(t, err, "Expected no error")
	defer signInResp.Body.Close()
	assert.Equal(t, http.StatusOK, signInResp.Status, "Expected http status ok")

	cancel()
	t.Log("E2E TestUserSignUpAndSignIn passed successfully")
}
