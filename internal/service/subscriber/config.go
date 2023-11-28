package subscriber

type Config struct {
	ClusterID        string `config:"CLUSTER_ID" yaml:"cluster_id"`
	SubscriberClient string `config:"SUBSCRIBER_CLIENT_ID" yaml:"subscriber_client"`
	Subject          string `config:"SUBJECT" yaml:"subject"`
	URL              string `config:"NATS_URL" yaml:"url"`
}
