package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	ENV_PREFIX       = "CHAT_CLI"
	LLM_API_HOST     = "llm-api-host"
	LLM_API_MODEL    = "llm-api-model"
	LLM_API_KEY      = "llm-api-key"
	LLM_API_KEY_FILE = "llm-api-key-file"
	DEBUG            = "debug"
)

func InitDefault() {
	viper.SetEnvPrefix(ENV_PREFIX)
	viper.SetEnvKeyReplacer(strings.NewReplacer("_", "-"))
	viper.AutomaticEnv()
	prepareEnv()
}
func populateDefaults() {
	viper.Set(LLM_API_HOST, "https://api.anthropic.com/v1/messages")
	viper.Set(LLM_API_KEY_FILE, "~/.config/.chat_cli/key")
	viper.Set(LLM_API_MODEL, "claude-3-opus-20240229")
}

func prepareEnv() {
	// Currently this snippet does not work with viper lib
	// https://github.com/spf13/viper/issues/188
	populateDefaults()
	underScoreReplacer := strings.NewReplacer("_", "-")
	viper.SetEnvPrefix(ENV_PREFIX)
	viper.SetEnvKeyReplacer(underScoreReplacer)
	viper.AutomaticEnv()
	// unfortunate workaround
	for _, envEntry := range os.Environ() {
		envVar := strings.Split(envEntry, "=")
		if strings.HasPrefix(envVar[0], ENV_PREFIX) {
			stripPrefix := fmt.Sprintf("%s_", ENV_PREFIX)
			newKey := strings.TrimPrefix(envVar[0], stripPrefix)
			newKey = underScoreReplacer.Replace(newKey)
			viper.Set(newKey, envVar[1])
		}
	}
}
