package plugin

import "github.com/google/wire"

var wireBasicSet = wire.NewSet()

func Initialize() (*Server, error) {
	wire.Build(wireExtsSet)
	return &Server{}, nil
}
