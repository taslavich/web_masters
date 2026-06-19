package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/go-chi/chi/v5"
	"gitlab.com/taslavich/web_masters/internal/config"
	httpServer "gitlab.com/taslavich/web_masters/internal/http"
	wmApiWeb "gitlab.com/taslavich/web_masters/internal/services/wmApi/web"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig[config.WmApiConfig](ctx)
	if err != nil {
		log.Fatalf("Cannot load config: %v", err)
	}
	log.Println("Config initialized!")

	addr := net.JoinHostPort(cfg.Clickhouse.Host, cfg.Clickhouse.Port)

	clickhouseConn, err := clickhouse.Open(&clickhouse.Options{
		Addr:     []string{addr},
		Protocol: clickhouse.Native,
		TLS:      &tls.Config{},
		Auth: clickhouse.Auth{
			Username: cfg.Clickhouse.Username,
			Password: cfg.Clickhouse.Password,
			Database: cfg.Clickhouse.Database,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		MaxOpenConns: 4,
		MaxIdleConns: 4,
	})
	if err != nil {
		log.Fatalf("❌ ClickHouse Open connection failed: %v", err)
	}
	defer clickhouseConn.Close()

	if err := clickhouseConn.Ping(ctx); err != nil {
		log.Fatalf("❌ ClickHouse ping failed: %v", err)
	}
	log.Println("✅ Connected to ClickHouse")

	router := chi.NewRouter()
	router = httpServer.InitHttpRouter(router)
	wmApiWeb.InitHttpRoutes(
		ctx,
		router,
		clickhouseConn,
		cfg.SspPopAdlFeeds,
		cfg.SspPopMcFeeds,
		cfg.SspIppAdlFeeds,
		cfg.SspIppMcFeeds,
		cfg.SspBanAdlFeeds,
		cfg.SspBanMcFeeds,
		cfg.SspNatAdlFeeds,
		cfg.SspNatMcFeeds,
	)
	log.Println("HTTP routes initialized")

	httpServer.RunHttpServer(ctx, router, cfg.HttpServer.Host, cfg.HttpServer.Port)
}
