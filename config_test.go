package crncl

import (
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func TestGetConfigFromFile(t *testing.T) {
	testCases := []struct {
		name           string
		filepath       string
		expectedConfig Config
		expectedErr    bool
	}{
		{
			name:     "config found and parsed",
			filepath: filepath.Join(basepath, "test", "golden", "config.toml"),
			expectedConfig: Config{
				BaseURL:     "localhost",
				Port:        3030,
				Title:       "Devin's Website",
				Description: "Just a testing",
				EnableRSS:   false,
				ColorScheme: Colors{
					Main:          "slate-900",
					Secondary:     "gray-700",
					Accent:        "teal-700",
					Text:          "teal-400",
					DarkMain:      "slate-900",
					DarkSecondary: "gray-700",
					DarkAccent:    "teal-700",
					DarkText:      "teal-400",
				},
			},
		},
		{
			name:        "config file not found",
			filepath:    "",
			expectedErr: true,
		},
		{
			name:        "config file not toml",
			filepath:    filepath.Join(basepath, "test", "golden", "badconfig.json"),
			expectedErr: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			is := is.New(t)
			actual, err := GetConfigFromFile(test.filepath)
			if test.expectedErr {
				is.True(err != nil) // there should be an err
			} else {
				is.NoErr(err)                         // parse error
				is.Equal(actual, test.expectedConfig) // parsed config must match
			}

		})
	}
}
