package app

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stepanbukhtii/easy-tools/nats"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/stepanbukhtii/easy-tools/elog"
	"github.com/stepanbukhtii/easy-tools/kafka"
	"github.com/stepanbukhtii/easy-tools/rabbitmq"
	"github.com/stepanbukhtii/go-blueprint/internal/config"
	"github.com/stepanbukhtii/go-blueprint/internal/repository"
)

const shutdownTimeout = 30 * time.Second

type App struct {
	Config            config.Config
	Database          *sql.DB
	PgxPool           *pgxpool.Pool
	RedisClient       *redis.Client
	RabbitMQPublisher *rabbitmq.Publisher
	KafkaProducer     *kafka.Producer
	NatsPublisher     *nats.Publisher
	Repository        repository.Repository
	Services          *Services
	closer            []func(ctx context.Context)
}

func New(ctx context.Context) (*App, error) {
	var app App
	var err error

	app.Config, err = config.Read()
	if err != nil {
		return nil, fmt.Errorf("config read: %w", err)
	}

	if err := app.initOpenTelemetry(ctx); err != nil {
		return nil, fmt.Errorf("open telemetry init: %w", err)
	}

	if err := app.initSlog(ctx); err != nil {
		return nil, fmt.Errorf("slog init: %w", err)
	}

	if err := app.initDatabase(ctx); err != nil {
		return nil, fmt.Errorf("database init: %w", err)
	}

	if err := app.initRedis(ctx); err != nil {
		return nil, fmt.Errorf("redis init: %w", err)
	}

	app.RabbitMQPublisher, err = rabbitmq.NewPublisher(app.Config.RabbitMQ.ConnectionURI())
	if err != nil {
		return nil, err
	}

	app.KafkaProducer, err = kafka.NewProducer(app.Config.Kafka.Brokers...)
	if err != nil {
		return nil, err
	}

	app.NatsPublisher, err = nats.NewPublisher(app.Config.NATS.ConnectionURI())
	if err != nil {
		return nil, err
	}

	app.Repository = repository.NewRepository(app.Config.Service.Name, app.Database, app.RedisClient)

	if err := app.initServices(); err != nil {
		return nil, fmt.Errorf("init services: %w", err)
	}

	return &app, nil
}

func (a *App) Close() {
	if err := a.Database.Close(); err != nil {
		slog.With(elog.Err(err)).Error("database close")
	}

	a.PgxPool.Close()

	if err := a.RedisClient.Close(); err != nil {
		slog.With(elog.Err(err)).Error("redis close")
	}

	if err := a.RabbitMQPublisher.Close(); err != nil {
		slog.With(elog.Err(err)).Error("rabbitMQ connection close")
	}

	a.KafkaProducer.Close()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := range a.closer {
		go a.closer[i](shutdownCtx)
	}
}

func (a *App) initOpenTelemetry(ctx context.Context) error {
	otel.SetTextMapPropagator(propagation.TraceContext{})

	if a.Config.OpenTelemetry.Disabled {
		otel.SetTracerProvider(trace.NewTracerProvider(trace.WithSampler(trace.AlwaysSample())))
		return nil
	}

	res, err := resource.New(ctx)
	if err != nil {
		return fmt.Errorf("create resource: %w", err)
	}

	// tracer
	traceExporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return fmt.Errorf("create trace exporter: %w", err)
	}

	traceProvider := trace.NewTracerProvider(trace.WithBatcher(traceExporter), trace.WithResource(res))
	otel.SetTracerProvider(traceProvider)

	a.closer = append(a.closer, func(ctx context.Context) {
		if err := traceProvider.Shutdown(ctx); err != nil {
			slog.With(elog.Err(err)).Error("open telemetry trace provider close")
		}
	})

	// metrics
	metricExporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return fmt.Errorf("create metric exporter: %w", err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
	)
	otel.SetMeterProvider(meterProvider)

	if err = runtime.Start(); err != nil {
		return fmt.Errorf("start runtime: %w", err)
	}

	a.closer = append(a.closer, func(ctx context.Context) {
		if err := meterProvider.Shutdown(ctx); err != nil {
			slog.With(elog.Err(err)).Error("open telemetry metrics provider close")
		}
	})

	return nil
}

func (a *App) initSlog(ctx context.Context) error {
	if a.Config.OpenTelemetry.Disabled {
		slog.SetDefault(slog.New(elog.NewSlogHandler(a.Config.Log, a.Config.Service)))
		return nil
	}

	res, err := resource.New(ctx)
	if err != nil {
		return fmt.Errorf("create resource: %w", err)
	}

	logExporter, err := otlploggrpc.New(ctx)
	if err != nil {
		return fmt.Errorf("create exporter: %w", err)
	}
	loggerProvider := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
	)
	global.SetLoggerProvider(loggerProvider)

	a.closer = append(a.closer, func(ctx context.Context) {
		if err := loggerProvider.Shutdown(ctx); err != nil {
			slog.With(elog.Err(err)).Error("open telemetry logger provider close")
		}
	})

	stdoutHandler := elog.NewSlogHandler(a.Config.Log, a.Config.Service)
	otelHandler := otelslog.NewHandler(a.Config.Service.Name)

	slog.SetDefault(slog.New(elog.NewMultiHandler(stdoutHandler, otelHandler)))

	return nil
}

func (a *App) initDatabase(ctx context.Context) error {
	poolConfig, err := pgxpool.ParseConfig(a.Config.Database.ConnectionURI())
	if err != nil {
		return err
	}

	poolConfig.ConnConfig.Tracer = otelpgx.NewTracer()

	if a.Config.Database.MaxOpenConnections != nil {
		poolConfig.MaxConns = int32(*a.Config.Database.MaxOpenConnections)
	}
	if a.Config.Database.MaxIdleConnections != nil {
		poolConfig.MinConns = int32(*a.Config.Database.MaxIdleConnections)
	}
	poolConfig.MaxConnLifetime = time.Minute
	poolConfig.MaxConnIdleTime = time.Minute

	a.PgxPool, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("database connect: %w", err)
	}

	if err = a.PgxPool.Ping(ctx); err != nil {
		return fmt.Errorf("database ping: %w", err)
	}

	a.Database = stdlib.OpenDBFromPool(a.PgxPool)

	return nil
}

func (a *App) initRedis(ctx context.Context) error {
	var tlsConfig *tls.Config
	if !a.Config.Redis.TLSDisabled {
		tlsConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	redisOptions := &redis.Options{
		Addr:      a.Config.Redis.Addresses[0],
		DB:        a.Config.Redis.DB,
		Password:  a.Config.Redis.Password,
		TLSConfig: tlsConfig,
	}
	a.RedisClient = redis.NewClient(redisOptions)

	if err := a.RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}

	if err := redisotel.InstrumentTracing(a.RedisClient); err != nil {
		return fmt.Errorf("open telemetry tracing: %w", err)
	}

	return nil
}
