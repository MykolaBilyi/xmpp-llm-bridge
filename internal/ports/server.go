package ports

import "context"

type Server interface {
	Serve(context.Context) error
	Shutdown(context.Context) error
}
