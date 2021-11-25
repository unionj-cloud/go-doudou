package test

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/unionj-cloud/go-doudou/stringutils"
	"time"
)

// SetupMySQLContainer starts mysql server docker container
func SetupMySQLContainer(logger *logrus.Logger, initdb string, dbname string) (func(), string, int, error) {
	logger.Info("setup MySQL Container")
	ctx := context.Background()
	if stringutils.IsEmpty(dbname) {
		dbname = "test"
	}
	req := testcontainers.ContainerRequest{
		Image:        "mysql:latest",
		ExposedPorts: []string{"3306/tcp", "33060/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "1234",
			"MYSQL_DATABASE":      dbname,
		},
		BindMounts: map[string]string{
			initdb: "/docker-entrypoint-initdb.d",
		},
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server - GPL").WithStartupTimeout(60 * time.Second),
	}

	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		logger.Errorf("error starting mysql container: %s", err)
		panic(fmt.Sprintf("%v", err))
	}

	closeContainer := func() {
		logger.Info("terminating container")
		err := mysqlC.Terminate(ctx)
		if err != nil {
			logger.Errorf("error terminating mysql container: %s", err)
			panic(fmt.Sprintf("%v", err))
		}
	}

	host, _ := mysqlC.Host(ctx)
	p, _ := mysqlC.MappedPort(ctx, "3306/tcp")
	port := p.Int()

	return closeContainer, host, port, nil
}
