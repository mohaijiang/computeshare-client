// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/mohaijiang/computeshare-client/internal/biz"
	"github.com/mohaijiang/computeshare-client/internal/conf"
	"github.com/mohaijiang/computeshare-client/internal/data"
	"github.com/mohaijiang/computeshare-client/internal/server"
	"github.com/mohaijiang/computeshare-client/internal/service"
	"github.com/mohaijiang/computeshare-client/third_party/agent"
	"github.com/mohaijiang/computeshare-client/third_party/p2p"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	greeterRepo := data.NewGreeterRepo(dataData, logger)
	greeterUsecase := biz.NewGreeterUsecase(greeterRepo, logger)
	greeterService := service.NewGreeterService(greeterUsecase)
	grpcServer := server.NewGRPCServer(confServer, greeterService, logger)
	p2pClient, err := p2p.NewP2pClient(confServer)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	client, err := service.NewDockerCli()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	httpClient, cleanup2, err := agent.NewHttpConnection(confData)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	agentService := agent.NewAgentService(httpClient, p2pClient)
	vmService := service.NewVmService(client, agentService, logger)
	shell := service.NewIpfShell(confData)
	computePowerService, err := service.NewComputePowerService(shell, client, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	vmWebsocketHandler := service.NewVmWebsocketHandler(client)
	cronJob := service.NewCronJob(vmService, logger)
	httpServer := server.NewHTTPServer(confServer, greeterService, p2pClient, vmService, computePowerService, agentService, vmWebsocketHandler, cronJob, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}
