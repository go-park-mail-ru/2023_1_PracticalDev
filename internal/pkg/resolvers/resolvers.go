package resolvers

import (
	"context"
	"errors"
	"strconv"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
)

var ErrNoAliveServices = errors.New("no alive services")

func NewResolver(client *consulapi.Client, serviceName string, logger *zap.Logger) (*manual.Resolver, error) {
	health, _, err := client.Health().Service(serviceName, "", false, nil)
	if err != nil {
		logger.Error("cant get alive services", zap.Error(err))
		return nil, err
	}
	servers := make([]resolver.Address, 0, len(health))
	for _, item := range health {
		addr := item.Service.Address +
			":" + strconv.Itoa(item.Service.Port)
		servers = append(servers, resolver.Address{Addr: addr})
	}
	if len(servers) == 0 {
		logger.Error("no alive services")
		return nil, ErrNoAliveServices
	}
	logger.Debug("Discovered alive servers", zap.Int("amount", len(servers)))

	rslver := manual.NewBuilderWithScheme("shortenerresolver")
	rslver.InitialState(resolver.State{
		Addresses: servers,
	})
	return rslver, nil
}

func NewGRPCConnWithResolver(ctx context.Context, client *consulapi.Client, serviceName string, logger *zap.Logger) (*grpc.ClientConn, error) {
	rslver, err := NewResolver(client, serviceName, logger)
	if err != nil {
		logger.Error("failed to get alive services from consul", zap.Error(err), zap.String("service", serviceName))
		return nil, err
	}
	conn, err := grpc.Dial(
		rslver.Scheme()+":///",
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithResolvers(rslver),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("failed to create grpc con", zap.Error(err), zap.String("service", serviceName))
		return nil, err
	}
	return conn, nil
}

func RunOnlineServiceDiscovery(ctx context.Context, client *consulapi.Client, nameResolver *manual.Resolver, serviceName string, logger *zap.Logger) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			health, _, err := client.Health().Service(serviceName, "", false, nil)
			if err != nil {
				logger.Error("cant get alive services")
				return
			}

			servers := make([]resolver.Address, 0, len(health))
			for _, item := range health {
				addr := item.Service.Address +
					":" + strconv.Itoa(item.Service.Port)
				servers = append(servers, resolver.Address{Addr: addr})
			}
			_ = nameResolver.CC.UpdateState(resolver.State{
				Addresses: servers,
			})
		}
	}
}
