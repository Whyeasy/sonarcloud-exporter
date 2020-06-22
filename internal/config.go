package internal

//Config struct for SonarCloud Token and Exporter
type Config struct {
	Token         string
	ListenAddress string
	ListenPath    string
	Organization  string
}
