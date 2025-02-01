package crncl

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

// Add things here as they come up. Things to think about:
//  - nav menu? especially if we allow generating of pages from markdown
//  - image storage location? handy for saving generated images as well as pre-populating load location
//  - social links?

// Config is the struct that models the config to be passed as a toml file to set up the blog.
type Config struct {
	BaseURL     string `toml:"base_url"`
	Port        int    `toml:"port,omitempty"`
	Title       string `toml:"title"`
	Description string `toml:"description"`
	EnableRSS   bool   `toml:"enable_rss"`
	ColorScheme Colors `toml:"color_scheme"`
}

// Colors models the color scheme for the website for both light and dark mode. These are
// string representation of Tailwind CSS colors.
type Colors struct {
	// Light mode set of colors. Used by default.
	Main      string `toml:"main"`
	Secondary string `toml:"secondary"`
	Accent    string `toml:"accent"`
	Text      string `toml:"text"`

	// Dark mode set of colors.
	DarkMain      string `toml:"dark_main"`
	DarkSecondary string `toml:"dark_secondary"`
	DarkAccent    string `toml:"dark_accent"`
	DarkText      string `toml:"dark_text"`
}

// GetConfig parses the config from the config file - "config.toml" - at the root of the project.
func GetConfig() (Config, error) {
	return GetConfigFromFile(filepath.Join(basepath, "config.toml"))
}

// GetConfigFromFile parses the config from the given filepath.
func GetConfigFromFile(filepath string) (Config, error) {
	var cfg Config
	_, err := toml.DecodeFile(filepath, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("issue finding and parsing toml file %s: %w", filepath, err)
	}

	return cfg, nil
}
