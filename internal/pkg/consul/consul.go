package consul

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/config"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

func NewConsulClient() (*consulapi.Client, error) {
	cfg := consulapi.DefaultConfig()
	cfg.Address = viper.GetString(config.ConsulConfig.Addr)
	return consulapi.NewClient(cfg)
}

func RegisterService(client *consulapi.Client) error {
	return client.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      viper.GetString(config.GrpcConfig.ServiceName) + "_" + config.GetConsulAddr(),
		Name:    viper.GetString(config.GrpcConfig.ServiceName),
		Port:    viper.GetInt(config.GrpcConfig.Port),
		Address: config.GetConsulAddr(),
	})
}
