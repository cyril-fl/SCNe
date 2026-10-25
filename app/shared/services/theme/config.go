package theme

// - Config - //

// Config holds the UI colors, as hex or ANSI codes.
type Config struct {
	Accent  string `yaml:"accent" json:"accent"`
	Text    string `yaml:"text" json:"text"`
	Muted   string `yaml:"muted" json:"muted"`
	Subtle  string `yaml:"subtle" json:"subtle"`
	Error   string `yaml:"error" json:"error"`
	Success string `yaml:"success" json:"success"`
}

// Default is the theme used until SetCache is called, and for every color the config leaves empty.
var Default = Config{
	Accent:  "#C9A26B",
	Text:    "#E3DED6",
	Muted:   "#8C877F",
	Subtle:  "#4B4843",
	Error:   "#D0614C",
	Success: "#8FB573",
}

var cache = Default

// Get returns the current theme.
func Get() Config {
	return cache
}

// SetCache sets the current theme.
func SetCache(cfg Config) {
	cache = cfg
}

// Merge returns cfg where every empty color is taken from base.
func (cfg Config) Merge(base Config) Config {
	return Config{
		Accent:  orDefault(cfg.Accent, base.Accent),
		Text:    orDefault(cfg.Text, base.Text),
		Muted:   orDefault(cfg.Muted, base.Muted),
		Subtle:  orDefault(cfg.Subtle, base.Subtle),
		Error:   orDefault(cfg.Error, base.Error),
		Success: orDefault(cfg.Success, base.Success),
	}
}

// Helpers

func orDefault(value string, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}
