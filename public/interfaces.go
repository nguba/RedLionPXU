package public

// ProfileReader defines the interface for reading brewing profiles from a data source.
type ProfileReader interface {
	ReadProfiles() ([]Profile, error)
}
