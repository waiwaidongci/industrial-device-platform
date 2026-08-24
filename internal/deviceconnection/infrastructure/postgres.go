package infrastructure

import (
	"context"
	"errors"
	"github.com/example/industrial-device-platform/internal/deviceconnection/domain"
	"time"
)

type Store interface {
	Exec(context.Context, string, ...any) error
	Query(context.Context, string, ...any) ([][]any, error)
}
type PostgreSQLRepository struct{ db Store }

func NewPostgreSQLRepository(db Store) *PostgreSQLRepository { return &PostgreSQLRepository{db: db} }
func (r *PostgreSQLRepository) Save(ctx context.Context, c *domain.Connection) error {
	return errors.New("postgres adapter requires configured driver")
}
func (r *PostgreSQLRepository) Get(ctx context.Context, id string) (*domain.Connection, error) {
	return nil, errors.New("postgres adapter requires configured driver")
}
func (r *PostgreSQLRepository) ListOffline(ctx context.Context, t time.Time) ([]*domain.Connection, error) {
	return nil, errors.New("postgres adapter requires configured driver")
}
