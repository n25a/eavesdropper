package database

type Conf struct {
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	DatabaseName string `yaml:"database_name"`
}
