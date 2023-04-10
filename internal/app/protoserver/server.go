package protoserver

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/CyrilSbrodov/metricService.git/cmd/config"
	"github.com/CyrilSbrodov/metricService.git/cmd/loggers"
	pb "github.com/CyrilSbrodov/metricService.git/internal/app/proto"
	"github.com/CyrilSbrodov/metricService.git/internal/storage"
)

type StoreServer struct {
	Storage storage.Storage
	logger  loggers.Logger
	cfg     config.ServerConfig
	pb.UnimplementedStorageServer
}

func NewStorageServer(cfg config.ServerConfig, logger loggers.Logger) *StoreServer {
	var store storage.Storage
	return &StoreServer{
		Storage: store,
		logger:  logger,
		cfg:     cfg,
	}
}

func (s *StoreServer) Run() {
	listen, err := net.Listen("tcp", s.cfg.GRPCAddr)
	if err != nil {
		s.logger.LogErr(err, " ")
		os.Exit(1)
	}
	srv := grpc.NewServer()
	pb.RegisterStorageServer(srv, s)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		if err = srv.Serve(listen); err != nil {
			s.logger.LogErr(err, "")
		}
	}()
	s.logger.LogInfo("server is listen:", s.cfg.GRPCAddr, "start server")

	//gracefullshutdown
	<-done

	s.logger.LogInfo("", "", "server stopped")

	srv.GracefulStop()

	s.logger.LogInfo("", "", "server stopped")
}

func (s *StoreServer) CollectMetric(ctx context.Context, in *pb.AddMetricRequest) (*pb.AddMetricResponse, error) {
	var response pb.AddMetricResponse
	var m storage.Metrics
	metric, err := s.Storage.GetMetric(*m.ToMetric(in.Metrics))
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
	err := s.Storage.CollectMetrics(metrics)
	if err != nil {
		s.logger.LogErr(err, "")
		return nil, err
	}
	return &response, nil
}

func (s *StoreServer) GetAll(ctx context.Context, in *pb.GetAllMetricsRequest) (*pb.GetAllMetricsResponse, error) {
	var response pb.GetAllMetricsResponse
	metrics, err := s.Storage.GetAll()
	if err != nil {
		s.logger.LogErr(err, "")
		return nil, err
	}
	response.AnswerToWeb.Msg = metrics
	return &response, nil
}
