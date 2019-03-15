package types

// CloudType type of the cloud
type CloudType string

func (ct CloudType) String() string {
	return string(ct)
}

const (
	// AWS cloud type
	AWS = CloudType("AWS")

	// GCP cloud type
	GCP = CloudType("GCP")

	// AZURE cloud type
	AZURE = CloudType("AZURE")

	// DUMMY cloud type
	DUMMY = CloudType("DUMMY")
)

// CloudProvider interface for the functions that can be used as operations/actions on the cloud providers
type CloudProvider interface {
	GetAccountName() string
	GetInstances() ([]*Instance, error)
	StopInstances([]*Instance) []error
	TerminateInstances([]*Instance) []error
	GetAccesses() ([]*Access, error)
	GetDatabases() ([]*Database, error)
	GetDisks() ([]*Disk, error)
	DeleteDisks([]*Disk) []error
	GetImages() ([]*Image, error)
}
