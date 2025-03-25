package config

type Config struct {
	Elasticsearch struct {
		Addresses []string
	}
	Services struct {
		Auth struct {
			Port int
		}
		Search struct {
			Port int
		}
		Gateway struct {
			Port int
		}
	}
}

var DefaultConfig = Config{
	Elasticsearch: struct {
		Addresses []string
	}{
		Addresses: []string{"http://localhost:9200/"},
	},
	Services: struct {
		Auth struct {
			Port int
		}
		Search struct {
			Port int
		}
		Gateway struct {
			Port int
		}
	}{
		Auth: struct {
			Port int
		}{
			Port: 50051,
		},
		Search: struct {
			Port int
		}{
			Port: 50052,
		},
		Gateway: struct {
			Port int
		}{
			Port: 8000,
		},
	},
}
