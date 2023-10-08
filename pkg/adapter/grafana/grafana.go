package grafana

import "github.com/google/wire"

// ProviderSet is Grafana client provider.
var ProviderSet = wire.NewSet(NewDashboardClient)
