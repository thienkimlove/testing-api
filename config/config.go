package config

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/spf13/viper"
	"go.tekoapis.com/kitchen/database"
	"go.tekoapis.com/kitchen/log"
	"go.tekoapis.com/testing/rpcimpl/server"
)

type Config struct {
	Log              log.Config           `json:"log"`
	Server           server.Config        `json:"server"`
	MySQL            database.MySQLConfig `json:"mysql"`
	MigrationsFolder string               `json:"migrations_folder"`
}

func Load() (*Config, error) {

	// You should set default config value here
	c := &Config{
		MySQL:            database.MySQLDefaultConfig(),
		Log:              log.DefaultConfig(),
		Server:           server.DefaultConfig(),
		MigrationsFolder: "file://migrations",
	}

	// --- hacking to load reflect structure config into env ----//
	viper.SetConfigType("json")
	configBuffer, err := json.Marshal(c)

	if err != nil {
		return nil, err
	}

	viper.ReadConfig(bytes.NewBuffer(configBuffer))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// -- end of hacking --//
	viper.AutomaticEnv()
	err = viper.Unmarshal(c)
	return c, err
}
