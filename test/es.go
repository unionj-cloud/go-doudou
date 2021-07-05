package test

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/unionj-cloud/go-doudou/logutils"
)

func PrepareTestEnvironment() (func(), string, int) {
	logger := logutils.NewLogger()
	var terminateContainer func() // variable to store function to terminate container
	var host string
	var port int
	var err error
	terminateContainer, host, port, err = SetupEs6Container(logger)
	if err != nil {
		logger.Panicln("failed to setup Elasticsearch container")
	}
	return terminateContainer, host, port
}

func SetupEs6Container(logger *logrus.Logger) (func(), string, int, error) {
	logger.Info("setup Elasticsearch v6 Container")
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "elasticsearch:6.8.12",
		ExposedPorts: []string{"9200/tcp", "9300/tcp"},
		Env: map[string]string{
			"discovery.type": "single-node",
		},
		WaitingFor: wait.ForLog("started"),
	}

	esC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		logger.Errorf("error starting Elasticsearch container: %s", err)
		panic(fmt.Sprintf("%v", err))
	}

	closeContainer := func() {
		logger.Info("terminating container")
		err := esC.Terminate(ctx)
		if err != nil {
			logger.Errorf("error terminating Elasticsearch container: %s", err)
			panic(fmt.Sprintf("%v", err))
		}
	}

	host, _ := esC.Host(ctx)
	p, _ := esC.MappedPort(ctx, "9200/tcp")
	port := p.Int()

	return closeContainer, host, port, nil
}
