package app

import (
	"log"

	"github.com/UdinSemen/golang-test-task/internal/config"
	numeratorUC "github.com/UdinSemen/golang-test-task/internal/domain/usecase/numerator"
	numeratorHTTP "github.com/UdinSemen/golang-test-task/internal/infrastructure/api/http/numerator"
	numeratorRepo "github.com/UdinSemen/golang-test-task/internal/infrastructure/repository/numerator"
	"github.com/UdinSemen/golang-test-task/pkg/closer"
	"github.com/UdinSemen/golang-test-task/pkg/connectors/postgres"
)

type serviceProvider struct {
	cfg *config.App

	postgresConnector *postgres.Connector

	numeratorRouter     *numeratorHTTP.Router
	numeratorUsecase    *numeratorUC.Usecase
	numeratorRepository *numeratorRepo.Repository
}

func newServiceProvider(cfg *config.App) *serviceProvider {
	return &serviceProvider{
		cfg: cfg,
	}
}

func (s *serviceProvider) PostgresConnector() *postgres.Connector {
	if s.postgresConnector == nil {
		con, err := postgres.NewConnector(&s.cfg.Connectors.MainStorage.Postgres)
		if err != nil {
			log.Fatal("create postgres connector: ", err)
		}
		closer.Add(con.Close)

		s.postgresConnector = con
	}

	return s.postgresConnector
}

func (s *serviceProvider) NumeratorRepository() *numeratorRepo.Repository {
	if s.numeratorRepository == nil {
		s.numeratorRepository = numeratorRepo.NewRepository(s.PostgresConnector().GetConnector())
	}

	return s.numeratorRepository
}

func (s *serviceProvider) NumeratorUsecase() *numeratorUC.Usecase {
	if s.numeratorUsecase == nil {
		s.numeratorUsecase = numeratorUC.NewUsecase(
			s.NumeratorRepository(),
		)
	}

	return s.numeratorUsecase
}

func (s *serviceProvider) NumeratorRouter() *numeratorHTTP.Router {
	if s.numeratorRouter == nil {
		s.numeratorRouter = numeratorHTTP.NewRouter(s.NumeratorUsecase())
	}

	return s.numeratorRouter
}
