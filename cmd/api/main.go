// Package main Product API
//
// @title           Product API
// @version         1.0
// @description     REST API untuk CRUD produk menggunakan net/http dan ServeMux
// @host            localhost:8081
// @BasePath        /
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	_ "product-api/docs"
	"product-api/internal/database"
	"product-api/internal/event"
	"product-api/internal/grpcserver"
	"product-api/internal/handler"
	"product-api/internal/observability"
	"product-api/internal/repository"
	"product-api/internal/router"
	"product-api/internal/service"
	"product-api/internal/util"
	pb "product-api/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	wd, _ := os.Getwd()
	slog.Debug(wd)
	_ = godotenv.Load()
	dsn := "root:rootpassword@tcp(localhost:3306)/products_db"
	ctx := context.Background()

	// Mysql
	db, err := database.ConnectDB(dsn)
	if err != nil {
		log.Fatal(fmt.Errorf("gagal membuat koneksi MySQL: %w", err))
	}
	defer db.Close()

	// Redis
	redisClient, err := database.ConnectRedis("localhost:6379")
	if err != nil {
		log.Fatal(fmt.Errorf("gagal membuat koneksi Redis: %w", err))
	}
	defer redisClient.Close()

	// RabbitMQ
	rabbitConn, err := database.ConnectRabbitMQ("amqp://app_user:app_password@localhost:5672/")
	if err != nil {
		log.Fatal(fmt.Errorf("gagal membuat koneksi RabbitMQ: %w", err))
	}
	defer rabbitConn.Close()

	rabbitCh, err := rabbitConn.Channel()
	if err != nil {
		log.Fatal(fmt.Errorf("gagal membuka channel RabbitMQ: %w", err))
	}
	defer rabbitCh.Close()

	if err := event.SetupProductBroker(rabbitCh); err != nil {
		log.Fatal(fmt.Errorf("gagal setup exchange/queue RabbitMQ: %w", err))
	}

	publisher := event.NewProductPublisher(rabbitCh)

	consumerCh, err := rabbitConn.Channel()
	if err != nil {
		log.Fatal(fmt.Errorf("gagal membuka channel consumer RabbitMQ: %w", err))
	}
	defer consumerCh.Close()

	util.SafeGo(ctx, "rabbit-consumer", func() {
		event.StartProductConsumer(consumerCh)
	})

	// Logger
	logger, shutdownLogger, err := observability.InitLogger("product-api")
	if err != nil {
		log.Fatal(fmt.Errorf("gagal setup logger: %w", err))
	}
	defer shutdownLogger(context.Background())

	slog.SetDefault(logger)

	repo := repository.NewProductRepositoryMySQL(db)
	svc := service.NewProductService(repo, redisClient, logger, publisher)
	h := handler.NewProductHandler(svc)
	mux := router.NewRouter(h)

	grpcSrv := grpcserver.NewProductGRPCServer(svc)
	grpcServer := grpc.NewServer()
	pb.RegisterProductGRPCServiceServer(grpcServer, grpcSrv)

	util.SafeGo(ctx, "grpc-server", func() {
		listener, err := net.Listen("tcp", ":9090")
		if err != nil {
			log.Fatal(fmt.Errorf("gagal listen gRPC : %w", err))
		}
		slog.Info("gRPC Server jalan di port 9090")
		if err := grpcServer.Serve(listener); err != nil {
			slog.Error("gRPC server gagal berjalan", "error", err)
		}
	})

	util.SafeGo(ctx, "grpc-gateway", func() {
		gateway("localhost:9090", ":8082")
	})

	slog.Info("REST server jalan di port 8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		slog.Error("REST server gagal berjalan", "error", err)
	}
}

func gateway(grpcAddr string, gatewayAddr string) {
	ctx := context.Background()
	gwMux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := pb.RegisterProductGRPCServiceHandlerFromEndpoint(ctx, gwMux, grpcAddr, opts)
	if err != nil {
		log.Fatal(fmt.Errorf("gagal daftar gateway: %w", err))
	}
	slog.Info("gateway REST (dari grpc) jalan di port 8082")
	log.Fatal(http.ListenAndServe(gatewayAddr, gwMux))
}
