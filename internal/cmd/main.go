package main

import (
	"context"
	"net/http"

	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
	"github.com/siemasusel/companiesapi/internal/command"
	kafkalib "github.com/siemasusel/companiesapi/internal/event/kafka"
	"github.com/siemasusel/companiesapi/internal/query"
	"github.com/siemasusel/companiesapi/internal/repository/mysql"
	httplib "github.com/siemasusel/companiesapi/internal/server/http"
	"golang.org/x/exp/slog"
)

const addr = ":8080"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db := newDatabase()
	kafkaWriter := newKafkaWriter()
	defer kafkaWriter.Close()

	companyRepo := mysql.NewCompanyRepository(db)
	eventPublisher := kafkalib.NewEventPublisher(kafkaWriter)

	commands := command.NewCommands(companyRepo, eventPublisher)
	queries := query.NewQueries(companyRepo)

	handler := httplib.NewRouter(slog.Default().Handler(), commands, queries)

	runHTTPServerOnAddr(ctx, addr, handler)
}

func runHTTPServerOnAddr(ctx context.Context, addr string, handler http.Handler) {
	server := http.Server{
		Addr:    addr,
		Handler: handler,
	}
	slog.InfoContext(ctx, "Listening for HTTP requests", "addr", addr)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.ErrorContext(ctx, "HTTP server failed")
		}
	}()

	<-ctx.Done()
	if err := server.Shutdown(context.Background()); err != nil {
		slog.ErrorContext(ctx, "HTTP server shutdown error", "err", err.Error())
	}
}

func newDatabase() *sqlx.DB {
	dataSourceName := os.Getenv("DATABASE_URL")
	if dataSourceName == "" {
		panic("missing DATABASE_URL environment variable")
	}

	db, err := sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	return db
}

func newKafkaWriter() *kafka.Writer {
	kafkaURL := os.Getenv("KAFKA_URL")
	if kafkaURL == "" {
		panic("missing KAFKA_URL environment variable")
	}

	return &kafka.Writer{
		Addr:                   kafka.TCP(kafkaURL),
		Topic:                  "company",
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
}
