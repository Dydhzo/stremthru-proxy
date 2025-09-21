package shared

type StatsHandler interface {
	AddBytes(bytes int64)
	IncrementConnections()
	DecrementConnections()
}

var globalStatsHandler StatsHandler

func RegisterStatsHandler(handler StatsHandler) {
	globalStatsHandler = handler
}

func GetStatsHandler() StatsHandler {
	return globalStatsHandler
}