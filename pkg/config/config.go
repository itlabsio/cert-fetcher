package config

type SecretsConfigList struct {
	Secrets []SecretsConfig `yaml:"secrets" json:"secrets"`
}
type SecretsConfig struct {
	VaultPath  string   `yaml:"vaultPath,omitempty" json:"vaultPath,omitempty"`
	SecretName string   `yaml:"secretName,omitempty" json:"secretName,omitempty"`
	Namespaces []string `yaml:"namespaces,flow,omitempty" json:"namespaces,flow,omitempty"`
}

type decoder interface {
	Decode(v interface{}) error
}

func ReadConfig(confSource decoder) (SecretsConfigList, error) {
	c := SecretsConfigList{
		Secrets: []SecretsConfig{},
	}
	if err := confSource.Decode(&c); err != nil {
		return c, err
	}
	return c, nil
}
