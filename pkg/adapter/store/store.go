package store

import "github.com/google/wire"

// ProviderSet is store provider.
var ProviderSet = wire.NewSet(NewReportStore)
