package test

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/unionj-cloud/go-doudou/pathutils"
)

func SetupMySQLContainer(logger *logrus.Logger) (func(), string, int, error) {
	logger.Info("setup MySQL Container")
	ctx := context.Background()
	mountPath := pathutils.Abs("sql")

	req := testcontainers.ContainerRequest{
		Image:        "mysql:latest",
		ExposedPorts: []string{"3306/tcp", "33060/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "1234",
			"MYSQL_DATABASE":      "test",
		},
		BindMounts: map[string]string{
			mountPath: "/docker-entrypoint-initdb.d",
		},
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server - GPL"),
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

	//host, _ := mysqlC.Host(ctx)
	p, _ := mysqlC.MappedPort(ctx, "3306/tcp")
	port := p.Int()
	host, _ := mysqlC.Name(ctx)

	return closeContainer, host, port, nil
}
