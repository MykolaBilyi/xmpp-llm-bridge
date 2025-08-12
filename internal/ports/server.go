package ports

import "context"

type Server interface {
	Serve() error
	Shutdown(context.Context) error
}
