package manager

type Config struct {
	MetricsAddr string
	ClusterName string
	AWS         struct {
		ServerRoleARN          string
		PermissionsBoundaryARN string
		RDSSecurityGroup       string
		RDSSubnetGroup         string
	}
}

type Manager struct {
	config Config
}

func New(cfg Config) (*Manager, error) {
	return &Manager{
		config: cfg,
	}, nil
}
