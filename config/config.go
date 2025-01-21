package config

type Config struct {
	Port            string `env:"SRV_PORT" envDefault:"8080"`
	GitHubAuthToken string `env:"GH_AUTH_TOKEN" envDefault:""`
}
