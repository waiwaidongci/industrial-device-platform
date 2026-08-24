package infrastructure

import (
	"context"
	"errors"
	"github.com/example/industrial-device-platform/internal/devicecatalog/domain"
)

type DB interface {
	Exec(context.Context, string, ...any) error
	Query(context.Context, string, ...any) ([][]any, error)
}
type PostgreSQLRepository struct{ db DB }

func NewPostgreSQLRepository(db DB) *PostgreSQLRepository { return &PostgreSQLRepository{db: db} }
func (r *PostgreSQLRepository) Create(ctx context.Context, d *domain.Device) error {
	return errors.New("postgres adapter not configured")
}
func (r *PostgreSQLRepository) Get(ctx context.Context, id string) (*domain.Device, error) {
	return nil, errors.New("postgres adapter not configured")
}
func (r *PostgreSQLRepository) List(ctx context.Context, t, g string) ([]*domain.Device, error) {
	return nil, errors.New("postgres adapter not configured")
}
func (r *PostgreSQLRepository) Update(ctx context.Context, d *domain.Device) error {
	return errors.New("postgres adapter not configured")
}
func (r *PostgreSQLRepository) Delete(ctx context.Context, id string) error {
	return errors.New("postgres adapter not configured")
}
