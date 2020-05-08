package tomlstruct

// ErrorList : Define struct for error list configuration file
type ErrorList struct {
	Team map[string]struct {
		Report []struct {
			Process     string   `toml:"process"`
			Host        string   `toml:"host,omitempty"`
			From        string   `toml:"from,omitempty"`
			User        string   `toml:"user,omitempty"`
			Passwd      string   `toml:"passwd,omitempty"`
			Recipients  []string `toml:"recipients,omitempty"`
			Subject     string   `toml:"subject,omitempty"`
			Body        string   `toml:"body"`
			Methods     string   `toml:"methods,omitempty"`
			URL         string   `toml:"url,omitempty"`
			ContentType string   `toml:"Content-Type,omitempty"`
		} `toml:"report"`
	} `toml:"team"`
}
