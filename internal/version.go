package internal

// Version holds the application version information.
type Version struct {
	version string
	commit  string
	date    string
}

// NewVersion creates a new Version instance.
func NewVersion(version, commit, date string) Version {
	return Version{
		version: version,
		commit:  commit,
		date:    date,
	}
}

// String returns the version string.
// Either <tag> (<date>) or <commit> (<date>).
func (v Version) String() string {
	var str string
	if v.version == "dev" {
		str = v.commit
	} else {
		str = v.version
	}
	return str + " (" + v.date + ")"
}
