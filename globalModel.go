package main

type Config struct {
	DbUser         string `yaml:"DbUser"`
	DbPsw          string `yaml:"DbPsw"`
	DbName         string `yaml:"DbName"`
	DbHost         string `yaml:"DbHost"`
	Port           string `yaml:"Port"`
	AllowedOrigins string `yaml:"AllowedOrigins"`
}
