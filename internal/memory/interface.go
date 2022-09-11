package memory


// OSMonitor represents the OS monitor interface
type OSMonitor interface {
	GetUsedPercentage() (*float64, error)
}