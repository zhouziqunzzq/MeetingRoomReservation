package config

var GlobalConfig = Config{}

type Config struct {
	// HTTP server settings
	PORT         int64    `toml:"port"`
	ALLOW_ORIGIN []string `toml:"allow_origin"`
	JWT_KEY      string   `toml:"jwt_key"`
	STATIC_DIR   string   `toml:"static_dir"`

	// Database settings
	SQLITE_FILE string `toml:"sqlite_path"`

	// Other settings
	MAX_QUERY_DAY        int    `toml:"max_query_day"`
	LOCK_ADVANCED_MINUTE int64  `toml:"lock_advanced_minute"`
	LOCK_PORT            string `toml:"lock_port"`
	LOCK_CMD             string `toml:"lock_cmd"`
}
