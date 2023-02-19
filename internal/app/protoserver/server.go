package protoserver

import (
	"context"
	"crypto/rsa"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/cmd/loggers"
	pb "github.com/CyrilSbrodov/metricService.git/internal/app/proto"
	"github.com/CyrilSbrodov/metricService.git/internal/crypto"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
	"github.com/CyrilSbrodov/metricService.git/internal/storage/repositories"
	"github.com/CyrilSbrodov/metricService.git/pkg/client/postgresql"
)

type StoreServer struct {
	storage storage.Storage
	logger  loggers.Logger
	pb.UnimplementedStorageServer
}

func newStorageServer(s storage.Storage, logger loggers.Logger) *StoreServer {
	return &StoreServer{
		storage: s,
		logger:  logger,
	}
}

type ServerApp struct {
	cfg      config.ServerConfig
	logger   *loggers.Logger
	Cryptoer crypto.Cryptoer
	private  *rsa.PrivateKey
	listen   *net.Listener
	server   *grpc.Server
}

func NewServerApp() *ServerApp {
	logger := loggers.NewLogger()
	cfg := config.ServerConfigInit()
	listen, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		logger.LogErr(err, " ")
		os.Exit(1)
	}
	s := grpc.NewServer()
	if cfg.CryptoPROKey != "" {
		c := crypto.NewCrypto()
		err := c.AddCryptoKey("public.pem", cfg.CryptoPROKey, "cert.pem", cfg.CryptoPROKeyPath)
		if err != nil {
			logger.LogErr(err, "filed to create file")
			os.Exit(1)
		}
		p, err := c.LoadPrivatePEMKey(cfg.CryptoPROKey)
		if err != nil {
			logger.LogErr(err, "filed to load file")
			os.Exit(1)
		}
		return &ServerApp{
			listen:   &listen,
			server:   s,
			cfg:      *cfg,
			logger:   logger,
			Cryptoer: c,
			private:  p,
		}
	}
	return &ServerApp{
		listen:   &listen,
		server:   s,
		cfg:      *cfg,
		logger:   logger,
		Cryptoer: nil,
		private:  nil,
	}
}

func (a *ServerApp) Run() {
	var err error
	//определение БД
	var store storage.Storage
	storage := newStorageServer(store, *a.logger)

	if len(a.cfg.DatabaseDSN) != 0 {
		client, err := postgresql.NewClient(context.Background(), 5, &a.cfg, a.logger)
		if err != nil {
			a.logger.LogErr(err, "")
			os.Exit(1)
		}
		storage.storage, err = repositories.NewPGSStore(client, &a.cfg, a.logger)
		if err != nil {
			a.logger.LogErr(err, "")
			os.Exit(1)
		}
	} else {
		storage.storage, err = repositories.NewRepository(&a.cfg, a.logger)
		if err != nil {
			a.logger.LogErr(err, "")
			os.Exit(1)
		}
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		if err = a.server.Serve(*a.listen); err != nil {
			a.logger.LogErr(err, "")
		}
	}()
	a.logger.LogInfo("server is listen:", a.cfg.GRPCAddr, "start server")

	//gracefullshutdown
	<-done

	a.logger.LogInfo("", "", "server stopped")

	a.server.GracefulStop()

	a.logger.LogInfo("", "", "server stopped")
}

func (s *StoreServer) CollectMetric(ctx context.Context, in *pb.AddMetricRequest) (*pb.AddMetricResponse, error) {
	var response pb.AddMetricResponse
	var m storage.Metrics
	metric, err := s.storage.GetMetric(*m.ToMetric(in.Metrics))
	if err != nil {
		s.logger.LogErr(err, "")
		return nil, err
	}
	response.Metrics = metric.ToProto()
	return &response, nil
}

func (s *StoreServer) CollectMetrics(ctx context.Context, in *pb.AddMetricsRequest) (*pb.AddMetricsResponse, error) {
	var response pb.AddMetricsResponse
	var metrics []storage.Metrics
	var m storage.Metrics
	for _, metric := range in.Metrics {
		metrics = append(metrics, *m.ToMetric(metric))
	}
	err := s.storage.CollectMetrics(metrics)
	if err != nil {
		s.logger.LogErr(err, "")
		return nil, err
	}
	return &response, nil
}

func (s *StoreServer) GetAll(ctx context.Context, in *pb.GetAllMetricsRequest) (*pb.GetAllMetricsResponse, error) {
	var response pb.GetAllMetricsResponse
	metrics, err := s.storage.GetAll()
	if err != nil {
		s.logger.LogErr(err, "")
		return nil, err
	}
	response.AnswerToWeb.Msg = metrics
	return &response, status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}

//func (s *StoreServer) CollectOrChangeGauge(ctx context.Context, in *CollectGaugeRequest) (*CollectGaugeResponse, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method CollectOrChangeGauge not implemented")
//}
//func (s *StoreServer) CollectOrIncreaseCounter(ctx context.Context, in *CollectCounterRequest) (*CollectCounterResponse, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method CollectOrIncreaseCounter not implemented")
//}
//func (s *StoreServer) GetGauge(ctx context.Context, in *GetGaugeRequest) (*GetGaugeResponse, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method GetGauge not implemented")
//}
//func (s *StoreServer) GetCounter(ctx context.Context, in *GetCounterRequest) (*GetCounterResponse, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method GetCounter not implemented")
//}
//func (s *StoreServer) PingClient(ctx context.Context, in *PingClientRequest) (*PingClientResponse, error) {
//	return nil, status.Errorf(codes.Unimplemented, "method PingClient not implemented")
//}
