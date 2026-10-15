package pg

import (
	"github.com/prometheus/client_golang/prometheus"
)

// PoolCollector - prometheus.Collector для *Pool, общий для всех сервисов
// web-rabis, использующих pg.Pool. Снимает pgxpool.Stat() в момент Collect
// (как и положено стейтлесс-коллектору), поэтому значения всегда свежие на
// момент скрейпа и не требуют отдельного фонового обновления.
//
// Использование в сервисе:
//
//	prometheus.MustRegister(pg.NewPoolCollector(pool, prometheus.Labels{"service": "ebook"}))
//
// constLabels обязателен хотя бы с лейблом "service" - сам пул не знает,
// какому сервису принадлежит, а в общем хранилище метрик (VictoriaMetrics)
// без этого лейбла пулы разных сервисов неразличимы.
type PoolCollector struct {
	pool *Pool

	maxConns           *prometheus.Desc
	totalConns         *prometheus.Desc
	acquiredConns      *prometheus.Desc
	idleConns          *prometheus.Desc
	constructingConns  *prometheus.Desc
	acquireCount       *prometheus.Desc
	acquireDuration    *prometheus.Desc
	canceledAcquires   *prometheus.Desc
	emptyAcquireCount  *prometheus.Desc
	newConnsCount      *prometheus.Desc
	maxLifetimeDestroy *prometheus.Desc
	maxIdleDestroy     *prometheus.Desc
}

// NewPoolCollector создаёт коллектор для pool. pool может ещё не быть
// подключён (Connect вызывается позже) - Collect в этом случае просто не
// отдаёт метрики за этот скрейп, без паники.
func NewPoolCollector(pool *Pool, constLabels prometheus.Labels) *PoolCollector {
	desc := func(name, help string) *prometheus.Desc {
		return prometheus.NewDesc("db_pool_"+name, help, nil, constLabels)
	}

	return &PoolCollector{
		pool: pool,

		maxConns:          desc("max_conns", "Maximum size of the pool."),
		totalConns:        desc("total_conns", "Total number of resources currently in the pool (constructing+acquired+idle)."),
		acquiredConns:     desc("acquired_conns", "Number of currently acquired (in-use) connections in the pool."),
		idleConns:         desc("idle_conns", "Number of currently idle connections in the pool."),
		constructingConns: desc("constructing_conns", "Number of connections with construction in progress."),

		acquireCount:       desc("acquire_count_total", "Cumulative count of successful acquires from the pool."),
		acquireDuration:    desc("acquire_duration_seconds_total", "Total duration of all successful acquires from the pool, in seconds."),
		canceledAcquires:   desc("canceled_acquire_count_total", "Cumulative count of acquires from the pool that were canceled by a context."),
		emptyAcquireCount:  desc("empty_acquire_count_total", "Cumulative count of successful acquires that waited for a resource to be released or constructed because the pool was empty."),
		newConnsCount:      desc("new_conns_count_total", "Cumulative count of new connections opened."),
		maxLifetimeDestroy: desc("max_lifetime_destroy_count_total", "Cumulative count of connections destroyed because they exceeded MaxConnLifetime."),
		maxIdleDestroy:     desc("max_idle_destroy_count_total", "Cumulative count of connections destroyed because they exceeded MaxConnIdleTime."),
	}
}

func (c *PoolCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.maxConns
	ch <- c.totalConns
	ch <- c.acquiredConns
	ch <- c.idleConns
	ch <- c.constructingConns
	ch <- c.acquireCount
	ch <- c.acquireDuration
	ch <- c.canceledAcquires
	ch <- c.emptyAcquireCount
	ch <- c.newConnsCount
	ch <- c.maxLifetimeDestroy
	ch <- c.maxIdleDestroy
}

func (c *PoolCollector) Collect(ch chan<- prometheus.Metric) {
	if c.pool == nil || c.pool.pool == nil {
		return
	}
	stat := c.pool.pool.Stat()

	ch <- prometheus.MustNewConstMetric(c.maxConns, prometheus.GaugeValue, float64(stat.MaxConns()))
	ch <- prometheus.MustNewConstMetric(c.totalConns, prometheus.GaugeValue, float64(stat.TotalConns()))
	ch <- prometheus.MustNewConstMetric(c.acquiredConns, prometheus.GaugeValue, float64(stat.AcquiredConns()))
	ch <- prometheus.MustNewConstMetric(c.idleConns, prometheus.GaugeValue, float64(stat.IdleConns()))
	ch <- prometheus.MustNewConstMetric(c.constructingConns, prometheus.GaugeValue, float64(stat.ConstructingConns()))

	ch <- prometheus.MustNewConstMetric(c.acquireCount, prometheus.CounterValue, float64(stat.AcquireCount()))
	ch <- prometheus.MustNewConstMetric(c.acquireDuration, prometheus.CounterValue, stat.AcquireDuration().Seconds())
	ch <- prometheus.MustNewConstMetric(c.canceledAcquires, prometheus.CounterValue, float64(stat.CanceledAcquireCount()))
	ch <- prometheus.MustNewConstMetric(c.emptyAcquireCount, prometheus.CounterValue, float64(stat.EmptyAcquireCount()))
	ch <- prometheus.MustNewConstMetric(c.newConnsCount, prometheus.CounterValue, float64(stat.NewConnsCount()))
	ch <- prometheus.MustNewConstMetric(c.maxLifetimeDestroy, prometheus.CounterValue, float64(stat.MaxLifetimeDestroyCount()))
	ch <- prometheus.MustNewConstMetric(c.maxIdleDestroy, prometheus.CounterValue, float64(stat.MaxIdleDestroyCount()))
}
