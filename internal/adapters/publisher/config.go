package publisher

type Config struct {
	ClusterID       string `config:"CLUSTER_ID" yaml:"cluster_id"`
	PublisherClient string `config:"PUBLISHER_CLIENT_ID" yaml:"publisher_client"`
	Subject         string `config:"SUBJECT" yaml:"subject"`
	URL             string `config:"NATS_URL" yaml:"url"`
}
