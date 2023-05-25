package configs

type Config struct {
	DataSourceURI                   string `mapstructure:"DATASOURCE_URI"`
	DataSourceDatabase              string `mapstructure:"DATASOURCE_DATABASE"`
	DataSourceCollection            string `mapstructure:"DATASOURCE_COLLECTION"`
	DataSourceTimeoutInMilliseconds int    `mapstructure:"DATASOURCE_TIMEOUT_IN_MILLISECONDS"`
}
