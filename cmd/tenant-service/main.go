package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/vyolayer/vyolayer/internal/tenant/config"
	"github.com/vyolayer/vyolayer/internal/tenant/infra"
	"github.com/vyolayer/vyolayer/internal/tenant/middleware"
	tenantrepo "github.com/vyolayer/vyolayer/internal/tenant/repo"
	"github.com/vyolayer/vyolayer/internal/tenant/server"
	"github.com/vyolayer/vyolayer/internal/tenant/usecase"
	"github.com/vyolayer/vyolayer/pkg/cache"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"github.com/vyolayer/vyolayer/pkg/postgres"
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"

	tenantGrpc "github.com/vyolayer/vyolayer/internal/tenant/delivery/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("iam service error: %v", err)
	}
}

func run() error {
	appLogger := logger.NewAppLogger("TENANT SERVICE")

	// Load Environment Variables
	if err := godotenv.Load(); err != nil {
		appLogger.Error("Note: No .env file found; relying on system environment variables", err.Error())
	}

	// Load Configuration
	cfg := config.Load()

	// Setup Database Connection
	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to retrieve underlying sql database object: %w", err)
	}
	defer sqlDB.Close()

	//
	memoryCache := cache.NewMemoryCache(2*time.Minute, 5*time.Minute)

	// ── Repositories ─────────────────────────────────────────────────────────
	orgRepo := tenantrepo.NewOrganizationRepo(
		db,
		appLogger.WithContext("Org Repo"),
	)
	orgMemberRepo := tenantrepo.NewOrganizationMemberRepo(
		db,
		appLogger.WithContext("Org Member Repo"),
	)
	orgRoleRepo := tenantrepo.NewOrganizationRoleRepository(
		db,
		appLogger.WithContext("Org Role Repo"),
	)
	orgPermRepo := tenantrepo.NewOrganizationPermissionRepository(
		db,
		appLogger.WithContext("Org Permission Repo"),
	)
	orgMemberRoleRepo := tenantrepo.NewMemberOrganizationRoleRepository(
		db,
		appLogger.WithContext("Member Org Role Repo"),
	)
	orgInvitationRepo := tenantrepo.NewOrganizationMemberInvitationRepo(
		db,
		appLogger.WithContext("Org Member Invitation Repo"),
	)
	tenantInfraRepo := tenantrepo.NewTenantInfraRepo(
		db,
		appLogger.WithContext("Tenant Infra Repo"),
	)
	projectRepo := tenantrepo.NewProjectRepository(
		db,
		appLogger.WithContext("Project Repo"),
	)
	projectMemberRepo := tenantrepo.NewProjectMemberRepository(
		db,
		appLogger.WithContext("Project Member Repo"),
	)

	// Infra
	dbProvisioner := infra.NewPostgresProvisioner(db)
	dbMigrator := infra.NewMigrator(db)

	tenantProvisioner := infra.NewTenantProvisioner(
		dbProvisioner,
		dbMigrator,
	)

	// ── Middleware ─────────────────────────────────────────────────────────
	middlewarePBAC := tenantrepo.NewOptimizedPermissionChecker(
		db,
		appLogger.WithContext("Optimized PBAC"),
	)
	cachePBAC := middleware.NewCachedPermissionChecker(
		middlewarePBAC,
		memoryCache,
		2*time.Minute,
	)

	// ── Use Cases ────────────────────────────────────────────────────────────
	orgUC := usecase.NewOrganizationUseCase(
		appLogger.WithContext("Org UseCase"),
		tenantProvisioner,
		orgRepo,
		orgMemberRepo,
		orgMemberRoleRepo,
		orgRoleRepo,
		orgPermRepo,
		projectRepo,
		projectMemberRepo,
		tenantInfraRepo,
	)
	orgMemberUC := usecase.NewOrganizationMemberUseCase(
		appLogger.WithContext("Org Member UseCase"),
		orgMemberRepo,
	)
	orgInvitationUC := usecase.NewOrganizationMemberInvitationUseCase(
		appLogger.WithContext("Org Member Invitation UseCase"),
		orgRepo,
		orgMemberRepo,
		orgRoleRepo,
		orgMemberRoleRepo,
		orgInvitationRepo,
	)
	projectUC := usecase.NewProjectUseCase(
		appLogger.WithContext("Project UseCase"),
		projectRepo,
		tenantInfraRepo,
	)
	projectMemberUC := usecase.NewProjectMemberUseCase(
		appLogger.WithContext("Project Member UseCase"),
		projectMemberRepo,
		projectRepo,
	)

	// ── gRPC handlers ─────────────────────────────────────────────────────────
	OrgH := tenantGrpc.NewOrganizationHandler(
		appLogger.WithContext("Org Handler"),
		orgUC,
		projectUC,
		projectMemberUC,
	)

	OrgMemberH := tenantGrpc.NewOrganizationMemberHandler(
		appLogger.WithContext("Org Member Handler"),
		orgMemberUC,
	)

	OrgInvitationH := tenantGrpc.NewOrganizationInvitationHandler(
		appLogger.WithContext("Org Invitation Handler"),
		orgInvitationUC,
	)
	ProjectH := tenantGrpc.NewProjectHandler(
		appLogger.WithContext("Project Handler"),
		orgUC,
		orgMemberUC,
		projectUC,
		projectMemberUC,
	)

	// ── Server ────────────────────────────────────────────────────────────────
	grpcSrv := server.NewGRPCServer(cfg.GRPC.GRPCPort, appLogger, cachePBAC)
	tenantV1.RegisterOrganizationServiceServer(grpcSrv.Engine, OrgH)
	tenantV1.RegisterOrganizationMemberServiceServer(grpcSrv.Engine, OrgMemberH)
	tenantV1.RegisterOrganizationInvitationServiceServer(grpcSrv.Engine, OrgInvitationH)
	tenantV1.RegisterProjectServiceServer(grpcSrv.Engine, ProjectH)

	return grpcSrv.Start()
}
