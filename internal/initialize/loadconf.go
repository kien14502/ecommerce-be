package initialize

import (
	"log"

	"github.com/kien14502/ecommerce-be/global"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func LoadConfigInit() {
	mode := pflag.StringP("server.mode", "e", "local", "Server mode (local/development/production)")
	pflag.Parse()
	v := viper.New()
	v.SetConfigName(*mode)
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")

	// Read config
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Unmarshal to global config
	if err := v.Unmarshal(&global.Config); err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}

	log.Printf("Config loaded from: %s", v.ConfigFileUsed())
}
