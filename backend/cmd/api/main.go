package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	adaptgemini "github.com/hassad/boilerplateSaaS/backend/internal/adapter/gemini"
	"github.com/hassad/boilerplateSaaS/backend/internal/adapter/postgres"
	adaptr2 "github.com/hassad/boilerplateSaaS/backend/internal/adapter/r2"
	"github.com/hassad/boilerplateSaaS/backend/internal/adapter/resend"
	adaptstripe "github.com/hassad/boilerplateSaaS/backend/internal/adapter/stripe"
	appai "github.com/hassad/boilerplateSaaS/backend/internal/app/ai"
	appauth "github.com/hassad/boilerplateSaaS/backend/internal/app/auth"
	appbilling "github.com/hassad/boilerplateSaaS/backend/internal/app/billing"
	appblog "github.com/hassad/boilerplateSaaS/backend/internal/app/blog"
	appnotif "github.com/hassad/boilerplateSaaS/backend/internal/app/notification"
	apporg "github.com/hassad/boilerplateSaaS/backend/internal/app/org"
	appstorage "github.com/hassad/boilerplateSaaS/backend/internal/app/storage"
	appteam "github.com/hassad/boilerplateSaaS/backend/internal/app/team"
	appuser "github.com/hassad/boilerplateSaaS/backend/internal/app/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/config"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
	"github.com/hassad/boilerplateSaaS/backend/pkg/jobs"
	"github.com/hassad/boilerplateSaaS/backend/pkg/jwt"
	"github.com/hassad/boilerplateSaaS/backend/pkg/observability"
	"github.com/hassad/boilerplateSaaS/backend/pkg/ws"
)

func main() {
	// Structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg := config.Load()

	// Tracing. No-op (collector-free) unless OTEL_EXPORTER_OTLP_ENDPOINT is set,
	// so local runs need no collector. shutdownTracing flushes pending spans on
	// graceful shutdown.
	shutdownTracing, err := observability.SetupTracing(context.Background(), observability.TracingConfig{
		Endpoint:    cfg.OTELExporterEndpoint,
		ServiceName: cfg.OTELServiceName,
	})
	if err != nil {
		logger.Error("failed to set up tracing", slog.String("error", err.Error()))
	}

	// Database
	db := postgres.NewDB(cfg.DatabaseURL)
	defer db.Close()

	// Prometheus metrics (dedicated registry). The HTTP middleware records
	// request count/latency; a gauge reports DB pool in-use connections.
	var metrics *observability.Metrics
	if cfg.MetricsEnabled {
		metrics = observability.NewMetrics(db)
	}

	// Repositories
	userRepo := postgres.NewUserRepository(db)
	passwordResetRepo := postgres.NewPasswordResetRepository(db)
	emailVerificationRepo := postgres.NewEmailVerificationRepository(db)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(db)

	// Organizations (the tenant). orgScope is the org-scoped unit-of-work that runs
	// tenant repository work under the restricted RLS role + active-org GUC, so
	// PostgreSQL row-level security enforces tenant isolation as the last line of
	// defense. txManager is the privileged (RLS-bypassing) unit-of-work used by
	// system flows like signup and team-create.
	orgRepo := postgres.NewOrganizationRepository(db)
	orgMemberRepo := postgres.NewOrganizationMemberRepository(db)
	orgScope := postgres.NewOrgScope(db)
	txManager := postgres.NewTxManager(db)
	orgSvc := apporg.NewService(orgRepo, orgMemberRepo)

	// JWT (short-lived access tokens with configurable TTL/iss/aud)
	jwtMaker := jwt.NewMakerWithOptions(cfg.JWTSecret, cfg.AccessTokenTTL, cfg.JWTIssuer, cfg.JWTAudience)

	// External services
	var emailSvc service.EmailService
	if cfg.ResendKey != "" {
		emailSvc = resend.NewEmailServiceWithTimeout(cfg.ResendKey, cfg.ExternalCallTimeout)
	}

	// Billing repositories + service (optional — only if Stripe key set).
	// Subscriptions are tenant-scoped: the request path uses orgScope (RLS), the
	// Stripe webhook (system path) uses the raw repository + resolves org from the
	// customer's user, so subscriptions are stamped with a tenant on creation.
	var billingSvc *appbilling.Service
	if cfg.StripeKey != "" {
		subscriptionRepo := postgres.NewSubscriptionRepository(db)
		planRepo := postgres.NewPlanRepository(db)
		invoiceRepo := postgres.NewInvoiceRepository(db)
		processedEventRepo := postgres.NewProcessedEventRepository(db)
		paymentSvc := adaptstripe.NewPaymentServiceWithTimeout(cfg.StripeKey, cfg.StripeWebhookSecret, cfg.ExternalCallTimeout)
		billingSvc = appbilling.NewService(appbilling.Deps{
			Users:           userRepo,
			Orgs:            orgRepo,
			Subscriptions:   subscriptionRepo,
			SubscriptionTx:  orgScope,
			Plans:           planRepo,
			Invoices:        invoiceRepo,
			ProcessedEvents: processedEventRepo,
			Payment:         paymentSvc,
			FrontendURL:     cfg.FrontendURL,
		})
	}

	// Storage (optional — only if R2 keys set). File metadata is org-scoped via
	// orgScope so RLS isolates each tenant's files.
	var storageSvc *appstorage.Service
	if cfg.R2AccessKey != "" {
		r2Client := adaptr2.NewClient(cfg.R2AccountID, cfg.R2AccessKey, cfg.R2SecretKey)
		r2Storage := adaptr2.NewStorageServiceWithTimeout(r2Client, cfg.R2BucketName, cfg.R2PublicURL, cfg.ExternalCallTimeout)
		storageSvc = appstorage.NewService(r2Storage, orgScope)
	}

	// AI Chat (optional — only if Gemini key set). Conversations are org-scoped.
	var aiSvc *appai.Service
	var demoAI service.AIService
	if cfg.GeminiKey != "" {
		geminiClient, err := adaptgemini.NewClient(context.Background(), cfg.GeminiKey)
		if err != nil {
			logger.Error("failed to create Gemini client", slog.String("error", err.Error()))
		} else {
			geminiAI := adaptgemini.NewAIServiceWithTimeout(geminiClient, cfg.ExternalCallTimeout)
			demoAI = geminiAI // expose for the public demo endpoint
			aiSvc = appai.NewService(orgScope, geminiAI)
		}
	}

	// Notifications (org-scoped).
	notifSvc := appnotif.NewService(orgScope)

	// WebSocket hub (real-time communication)
	wsHub := ws.NewHub()
	go wsHub.Run()

	// Wire WebSocket broadcaster into notification service
	notifSvc.SetBroadcaster(wsHub)

	// Blog
	blogRepo := postgres.NewBlogRepository(db)
	blogSvc := appblog.NewService(blogRepo)

	// Teams (optional)
	teamRepo := postgres.NewTeamRepository(db)
	memberRepo := postgres.NewTeamMemberRepository(db)
	teamSvc := appteam.NewService(teamRepo, memberRepo, txManager)

	// App services. Auth owns tenant signup: Register creates the user + their
	// personal org + owner membership atomically via the (privileged) txManager,
	// and stamps the active org into the access token.
	authSvc := appauth.NewService(appauth.Deps{
		Users:           userRepo,
		Orgs:            orgRepo,
		Resets:          passwordResetRepo,
		Verifications:   emailVerificationRepo,
		RefreshTokens:   refreshTokenRepo,
		Tx:              txManager,
		Email:           emailSvc,
		JWTMaker:        jwtMaker,
		FrontendURL:     cfg.FrontendURL,
		RefreshTokenTTL: cfg.RefreshTokenTTL,
	})
	userSvc := appuser.NewService(userRepo)

	// Router. The org resolver turns an authenticated user into an authorized
	// active organization for each request (verifying membership of any explicit
	// org claim, else the user's default org).
	router := handler.NewRouter(authSvc, userSvc, billingSvc, storageSvc, aiSvc, notifSvc, blogSvc, teamSvc, wsHub, cfg.JWTSecret, cfg.FrontendURL, db, logger, demoAI, orgSvc.ResolveActiveOrg, metrics)

	// HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 120 * time.Second, // Increased for SSE streaming endpoints
		IdleTimeout:  60 * time.Second,
	}

	// Background job scheduler. Each job invocation is bounded by cfg.JobTimeout
	// so a stuck cleanup (slow query / hung call) can never run unbounded.
	scheduler := jobs.NewSchedulerWithTimeout(logger, cfg.JobTimeout)
	scheduler.Register(jobs.Job{
		Name:     "clean-expired-password-resets",
		Interval: 1 * time.Hour,
		Fn: func(ctx context.Context) error {
			return passwordResetRepo.DeleteExpired(ctx)
		},
	})
	scheduler.Register(jobs.Job{
		Name:     "clean-expired-email-verifications",
		Interval: 1 * time.Hour,
		Fn: func(ctx context.Context) error {
			return emailVerificationRepo.DeleteExpired(ctx)
		},
	})
	scheduler.Register(jobs.Job{
		Name:     "clean-expired-refresh-tokens",
		Interval: 1 * time.Hour,
		Fn: func(ctx context.Context) error {
			return refreshTokenRepo.DeleteExpired(ctx)
		},
	})
	scheduler.Register(jobs.Job{
		Name:     "log-system-stats",
		Interval: 5 * time.Minute,
		Fn: func(_ context.Context) error {
			stats := db.Stats()
			logger.Info("system stats",
				slog.Int("goroutines", runtime.NumGoroutine()),
				slog.Int("db_open_connections", stats.OpenConnections),
				slog.Int("db_in_use", stats.InUse),
				slog.Int("db_idle", stats.Idle),
			)
			return nil
		},
	})
	scheduler.Start(context.Background())

	// Start server in goroutine
	go func() {
		logger.Info("API server starting", slog.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Info("shutdown signal received", slog.String("signal", sig.String()))

	scheduler.Stop()
	wsHub.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", slog.String("error", err.Error()))
	}

	// Flush any pending spans before exit (no-op when tracing is disabled).
	if shutdownTracing != nil {
		if err := shutdownTracing(ctx); err != nil {
			logger.Error("failed to flush tracer", slog.String("error", err.Error()))
		}
	}

	db.Close()
	logger.Info("server stopped gracefully")
}
