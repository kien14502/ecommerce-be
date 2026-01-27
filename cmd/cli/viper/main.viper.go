package viper

import "github.com/spf13/viper"

func main() {
	viper := viper.New()
	viper.AddConfigPath("../../../configs")
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	port := viper.GetInt("server.port")
	println("Server will run on port:", port)
}
