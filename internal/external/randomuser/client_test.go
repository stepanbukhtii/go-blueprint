package randomuser

import (
	"context"
	"fmt"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stepanbukhtii/go-blueprint/internal/config"
	"github.com/stretchr/testify/require"
)

func TestName(t *testing.T) {
	injector := do.New()

	do.ProvideValue(injector, config.Config{RandomUser: config.RandomUser{BaseURL: "https://randomuser.me"}})

	client, err := NewClient(injector)
	require.NoError(t, err)

	userData, err := client.GetRandomUser(context.Background())
	require.NoError(t, err)

	fmt.Println("userData", userData.Name.First, userData.Name.Last)
}
