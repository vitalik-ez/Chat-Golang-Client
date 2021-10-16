package room

type Config struct {
	ServerBasePath string
}

type Client struct {
	Config *Config
	Name   string
}

func NewClient(serverBasepath string) *Client {
	return &Client{
		Config: &Config{ServerBasePath: serverBasepath},
	}
}
