package randomuser

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stepanbukhtii/go-blueprint/internal/config"
)

func TestName(t *testing.T) {
	cfg := config.Config{RandomUser: config.RandomUser{BaseURL: "https://randomuser.me"}}

	client := NewClient(cfg)

	userData, err := client.GetRandomUser(context.Background())
	require.NoError(t, err)

	fmt.Println("userData", userData.Name.First, userData.Name.Last)
}
