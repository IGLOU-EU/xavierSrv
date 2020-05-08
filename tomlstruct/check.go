package tomlstruct

// CheckList : Define struct for check list configuration file
type CheckList struct {
	Team []struct {
		Enable bool   `toml:"enable"`
		Name   string `toml:"name"`
		App    []struct {
			Name     string `toml:"name"`
			URL      string `toml:"url"`
			Response int    `toml:"response"`
		} `toml:"app"`
	} `toml:"team"`
}
