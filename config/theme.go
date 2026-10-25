package config

// theme holds the UI colors, as hex or ANSI codes; an empty color keeps the terminal default.
type theme struct {
	Accent  string `yaml:"accent"`
	Text    string `yaml:"text"`
	Muted   string `yaml:"muted"`
	Subtle  string `yaml:"subtle"`
	Error   string `yaml:"error"`
	Success string `yaml:"success"`
}
