package randomuser

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/samber/do/v2"
	"github.com/stepanbukhtii/easy-tools/elog"
	"github.com/stepanbukhtii/go-blueprint/internal/config"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"resty.dev/v3"
)

const (
	serviceName = "random-user"

	getRandomUserEndpoint = "/api/"
)

type Client interface {
	GetRandomUser(ctx context.Context) (User, error)
}

type client struct {
	httpClient *resty.Client
}

func NewClient(injector do.Injector) (Client, error) {
	cfg := do.MustInvoke[config.Config](injector)

	return &client{
		httpClient: resty.NewWithClient(&http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}).
			SetBaseURL(cfg.RandomUser.BaseURL).
			SetHeader("Content-Type", "application/json").
			SetTimeout(30 * time.Second).
			SetResponseBodyUnlimitedReads(true).
			EnableTrace(),
	}, nil
}

func (c *client) GetRandomUser(ctx context.Context) (User, error) {
	l := &elog.RestyLogger{ServiceName: serviceName}

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetResult(UserResponse{}).
		SetError(UserResponse{}).
		Execute(http.MethodGet, getRandomUserEndpoint)
	if err != nil {
		l.Error(resp, err, "get random user")
		return User{}, err
	}

	if !resp.IsSuccess() {
		err := resp.Error().(*UserResponse)
		l.Error(resp, err, "error code received")
		return User{}, err
	}

	data, ok := resp.Result().(*UserResponse)
	if !ok || len(data.Results) == 0 {
		err := errors.New("failed to convert results")
		l.Error(resp, err, err.Error())
		return User{}, nil
	}

	l.Info(resp, "get random user")

	return data.Results[0], nil
}
