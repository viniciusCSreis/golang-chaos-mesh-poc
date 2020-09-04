package api

type (
	Configs struct {
		AppConfigs
	}

	AppConfigs struct {
		NatsConfig                NatsConfig  `json:"nats"`
	}

	NatsConfig struct {
		Host string `json:"host"`
		User string `json:"user"`
	}
)