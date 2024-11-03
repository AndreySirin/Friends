package storage_test

import (
	"Friends/logg"
	"Friends/storage"
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"log"
	"os"
	"testing"
)

const (
	username = "test"
	password = "test"
	database = "test"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	lg := logg.New()
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(database),
		postgres.WithUsername(username),
		postgres.WithPassword(password),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	//ip, err := postgresContainer.ContainerIP(ctx)
	//require.NoError(t, err)
	//fmt.Println(ip)

	network, err := postgresContainer.Inspect(ctx)
	require.NoError(t, err)
	port := network.NetworkSettings.Ports["5432/tcp"][0].HostPort
	fmt.Println(port)

	st, err := storage.New(lg, username, password, "0.0.0.0:"+port, database)
	require.NoError(t, err)

	err = st.DummyMigration(ctx)
	require.NoError(t, err)

	dir, err := os.Getwd()
	require.NoError(t, err)

	imgPath := dir + "/img/r.jpeg"
	imgData, err := os.ReadFile(imgPath)
	require.NoError(t, err)

	req := storage.ProductFriend{
		ID:        1,
		Name:      "ivan",
		Hobby:     "sport",
		Price:     1,
		ImageData: imgData,
	}

	err = st.AddProductFriend(ctx, req)
	require.NoError(t, err)
}
