package outbound

type DataSourceID interface {
	String() string
	AsDataSourceID() any
}
