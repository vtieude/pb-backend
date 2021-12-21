package entities

const ConfigKey = "pbconfig"

type PbConfig struct {
	DbUser  string `yaml:"DbUser"`
	DbPsw   string `yaml:"DbPsw"`
	DbName  string `yaml:"DbName"`
	DbHost  string `yaml:"DbHost"`
	AppPort string `yaml:"AppPort"`
}
