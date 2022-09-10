package memory

type OSMonitor interface {
	GetUsedPercentage() (*float64, error)
}