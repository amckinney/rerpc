package rerpc

import "context"

// Route is an RPC route.
type Route struct {
	Type           StreamType
	Path           string
	Implementation func(context.Context, StreamFunc)
}

// NewMux returns a new *Mux that implements the HandlerV2 interface.
func NewMux(routes ...Route) *Mux {
	routeMap := make(map[string]Route, len(routes))
	for _, route := range routes {
		routeMap[route.Path] = route
	}
	return &Mux{
		routes: routeMap,
	}
}

// Mux implements the HandlerV2 interface.
type Mux struct {
	routes map[string]Route
}

// Handle handles the given Stream according to the Route registered
// for the given path.
func (m *Mux) Handle(ctx context.Context, streamFunc StreamFunc, path string) error {
	route, ok := m.routes[path]
	if !ok {
		// TODO(alex): We might be able to un-export this constructor
		// since it's largely used as an implementation detail.
		return NewUnknownPathError(path)
	}
	route.Implementation(ctx, streamFunc)
	return nil
}

// StreamType returns the StreamType associated with the given path.
func (m *Mux) StreamType(path string) (StreamType, error) {
	route, ok := m.routes[path]
	if !ok {
		// TODO(alex): We might be able to un-export this constructor
		// since it's largely used as an implementation detail.
		return 0, NewUnknownPathError(path)
	}
	return route.Type, nil
}
