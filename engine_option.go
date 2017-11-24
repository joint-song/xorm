package xorm

import (
	"os"
	"time"

	"github.com/go-xorm/core"
)

// Option config Engine behavior
type Option func(x *Engine)

func EnableTraceOption(enable bool) Option {
	return func(x *Engine) {
		core.EnableTraceOption(enable)(x.db)
	}
}

func ColumnMapperOption(im core.IMapper) Option {
	return func(x *Engine) {
		x.ColumnMapper = im
	}
}

func TableMapperOption(im core.IMapper) Option {
	return func(x *Engine) {
		x.TableMapper = im
	}
}

func LoggerOption(l core.ILogger) Option {
	return func(x *Engine) {
		x.logger = l
	}
}

func ShowSQLOption(show bool) Option {
	return func(x *Engine) {
		x.showSQL = show
	}
}

func ShowExecTimeOption(show bool) Option {
	return func(x *Engine) {
		x.showExecTime = show
	}
}

func DisableGlobalCacheOption(disable bool) Option {
	return func(x *Engine) {
		x.disableGlobalCache = disable
	}
}

func TZLocationOption(local *time.Location) Option {
	return func(x *Engine) {
		x.SetTZLocation(local)
	}
}

func TZDatabaseOption(local *time.Location) Option {
	return func(x *Engine) {
		x.SetTZDatabase(local)
	}
}

func MaxIdleOpenOption(conns int) Option {
	return func(x *Engine) {
		x.SetMaxOpenConns(conns)
	}
}

func MaxIdleConnsOption(conns int) Option {
	return func(x *Engine) {
		x.SetMaxIdleConns(conns)
	}
}

func MapperOption(mapper core.IMapper) Option {
	return func(x *Engine) {
		x.SetMapper(mapper)
	}
}

func DefaultCacherOption(cacher core.Cacher) Option {
	return func(x *Engine) {
		x.SetDefaultCacher(cacher)
	}
}

func ConnMaxLifetimeOption(d time.Duration) Option {
	return func(x *Engine) {
		x.SetConnMaxLifetime(d)
	}
}

func defaultOptions() []Option {
	logger := NewSimpleLogger(os.Stdout)
	logger.SetLevel(core.LOG_INFO)
	dbTZOption := func(x *Engine) {
		local := time.Local
		if x.dialect.DBType() == core.SQLITE {
			local = time.UTC
		}
		x.SetTZDatabase(local)
	}

	return []Option{
		LoggerOption(logger),
		TZLocationOption(time.Local),
		dbTZOption,
		MapperOption(new(core.SnakeMapper)),
	}
}
