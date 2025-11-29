package platform

type RunnerConfig struct {
	Name        string
	Image       string
	Command     []string
	StoragePath string
	Labels      map[string]string
}
