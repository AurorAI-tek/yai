package config

import "github.com/ekkinox/yai/ai"

const (
	provider_type      = "PROVIDER_TYPE"
	provider_key       = "PROVIDER_KEY"
	provider_model     = "PROVIDER_MODEL"
	provider_proxy     = "PROVIDER_PROXY"
	provider_temperature = "PROVIDER_TEMPERATURE"
	provider_max_tokens = "PROVIDER_MAX_TOKENS"
)

type AiConfig struct {
	providerType ProviderType
	key          string
	model        string
	proxy        string
	temperature  float64
	maxTokens    int
}

type ProviderType = ai.ProviderType

func (c AiConfig) GetProviderType() ProviderType {
	return c.providerType
}

func (c AiConfig) GetKey() string {
	return c.key
}

func (c AiConfig) GetModel() string {
	return c.model
}

func (c AiConfig) GetProxy() string {
	return c.proxy
}

func (c AiConfig) GetTemperature() float64 {
	return c.temperature
}

func (c AiConfig) GetMaxTokens() int {
	return c.maxTokens
}
