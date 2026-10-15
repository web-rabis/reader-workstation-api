package pg

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrInvalidConfig возвращается, если Config заполнен не полностью.
var ErrInvalidConfig = errors.New("invalid configuration structure")

const defaultConnectTimeout = 100 * time.Second

// Config описывает подключение к Postgres. Специфики конкретного
// сервиса (какие репозитории он открывает) здесь нет и быть не должно —
// это переиспользуемая часть, общая для всех сервисов web-rabis.
type Config struct {
	URL            string
	DataBaseName   string
	MaxConns       int32         // 0 -> используется значение по умолчанию (100)
	ConnectTimeout time.Duration // 0 -> используется значение по умолчанию (100s)
}

// Pool — обёртка над pgxpool.Pool с единообразным жизненным циклом
// (Connect/Ping/Close), которым могут пользоваться все сервисы.
type Pool struct {
	connURL string
	dbName  string

	pool           *pgxpool.Pool
	config         *pgxpool.Config
	connectTimeout time.Duration
}

// New валидирует конфигурацию и готовит pgxpool.Config. Реального
// подключения ещё не происходит — для этого нужно вызвать Connect.
func New(conf Config) (*Pool, error) {
	if conf.URL == "" || conf.DataBaseName == "" {
		return nil, ErrInvalidConfig
	}

	pgxCfg, err := pgxpool.ParseConfig(conf.URL)
	if err != nil {
		return nil, err
	}

	maxConns := conf.MaxConns
	if maxConns == 0 {
		maxConns = 100
	}
	pgxCfg.MaxConns = maxConns

	connectTimeout := conf.ConnectTimeout
	if connectTimeout == 0 {
		connectTimeout = defaultConnectTimeout
	}

	return &Pool{
		connURL:        conf.URL,
		dbName:         conf.DataBaseName,
		config:         pgxCfg,
		connectTimeout: connectTimeout,
	}, nil
}

func (p *Pool) Name() string { return "PgSql" }

func (p *Pool) Connect(ctx context.Context) error {
	ctxWT, cancel := context.WithTimeout(ctx, p.connectTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctxWT, p.config)
	if err != nil {
		return err
	}
	p.pool = pool
	return nil
}

func (p *Pool) Ping(ctx context.Context) error {
	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	return conn.Ping(ctx)
}

func (p *Pool) Close() error {
	if p.pool != nil {
		p.pool.Close()
	}
	return nil
}

// DB возвращает нижележащий *pgxpool.Pool — им пользуются конкретные
// репозитории сервиса (через db/base) и его собственный DataStore,
// который добавляет доменные методы вроде Ebook()/Dictionary().
func (p *Pool) DB() *pgxpool.Pool { return p.pool }
