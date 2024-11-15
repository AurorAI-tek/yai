package config

import (
	"fmt"
	"strings"

	"github.com/ekkinox/yai/ai"
	"github.com/ekkinox/yai/system"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
)

type Config struct {
	ai     AiConfig
	user   UserConfig
	system *system.Analysis
}

func (c *Config) GetAiConfig() AiConfig {
	return c.ai
}

func (c *Config) GetUserConfig() UserConfig {
	return c.user
}

func (c *Config) GetSystemConfig() *system.Analysis {
	return c.system
}

func NewConfig() (*Config, error) {
	system := system.Analyse()

	viper.SetConfigName(strings.ToLower(system.GetApplicationName()))
	viper.AddConfigPath(fmt.Sprintf("%s/.config/", system.GetHomeDirectory()))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		ai: AiConfig{
			providerType: ProviderType(viper.GetString(provider_type)),
			key:         viper.GetString(provider_key),
			model:       viper.GetString(provider_model),
			proxy:       viper.GetString(provider_proxy),
			temperature: viper.GetFloat64(provider_temperature),
			maxTokens:   viper.GetInt(provider_max_tokens),
		},
		user: UserConfig{
			defaultPromptMode: viper.GetString(user_default_prompt_mode),
			preferences:       viper.GetString(user_preferences),
		},
		system: system,
	}, nil
}

func WriteConfig(key string, providerType ProviderType, model string) (*Config, error) {
	system := system.Analyse()

	// ai defaults
	viper.Set(provider_type, string(providerType))
	viper.Set(provider_key, key)
	viper.Set(provider_model, getDefaultModel(providerType, model))
	viper.SetDefault(provider_proxy, "")
	viper.SetDefault(provider_temperature, 0.2)
	viper.SetDefault(provider_max_tokens, 1000)

	// user defaults
	viper.SetDefault(user_default_prompt_mode, "exec")
	viper.SetDefault(user_preferences, "")

	if err := viper.SafeWriteConfigAs(system.GetConfigFile()); err != nil {
		return nil, err
	}

	return NewConfig()
}

func getDefaultModel(providerType ProviderType, model string) string {
	if model != "" {
		return model
	}

	switch providerType {
	case ai.OpenAIProvider:
		return openai.GPT3Dot5Turbo
	case ai.GroqProvider:
		return "llama2-70b-4096"
	default:
		return openai.GPT3Dot5Turbo
	}
}
