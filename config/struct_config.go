package config

// Config : Define struct for master configuration file
type Config struct {
	Overall struct {
		LoopWait      int    `toml:"loopWait"`
		LogFile       string `toml:"logFile"`
		LogVerbose    bool   `toml:"logVerbose"`
		CheckListFile string `toml:"checkListFile"`
		ErrorListFile string `toml:"errorListFile"`
	} `toml:"overall"`
	HTTP struct {
		MaxWait int `toml:"maxWait"`
	} `toml:"http"`
	Reporting struct {
		Local bool   `toml:"local"`
		Dir   string `toml:"dir"`
	} `toml:"reporting"`
}
