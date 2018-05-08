package config

var GlobalConfig = Config{}

type Config struct {
	PORT          int64    `toml:"port"`
	ALLOW_ORIGIN  []string `toml:"allow_origin"`
	SQLITE_FILE   string   `toml:"sqlite_path"`
	JWT_KEY       string   `toml:"jwt_key"`
	MAX_QUERY_DAY int      `toml:"max_query_day"`
	STATIC_DIR    string   `toml:"static_dir"`
}
