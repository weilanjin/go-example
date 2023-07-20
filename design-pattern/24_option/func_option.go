package option

type Options struct {
	Host     string
	Port     uint16
	Username string
	Password string
	Database string
}

type Option func(*Options)

func WithHost(host string) Option {
	return func(o *Options) {
		o.Host = host
	}
}

func WithPort(port uint16) Option {
	return func(o *Options) {
		o.Port = port
	}
}

func NewClient(username, password, database string, opts ...Option) *Options {
	o := &Options{
		Host:     "127.0.0.1",
		Port:     3306,
		Username: username,
		Password: password,
		Database: database,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}