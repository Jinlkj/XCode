package config

type Config struct {
	MySQL struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
	}
	Redis struct {
		Host     string
		Port     int
		Password string
		DB       int
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
	MySQL: struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
	}{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "",
		DBName:   "code_search_user",
	},
	Redis: struct {
		Host     string
		Port     int
		Password string
		DB       int
	}{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
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
