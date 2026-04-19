package localization

import (
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

// Localizer struct holds the bundle and supports translation
type Localizer struct {
	bundle *i18n.Bundle
	logger *zap.Logger
}

// NewLocalizer initializes the i18n bundle and loads translation files
func NewLocalizer(logger *zap.Logger) *Localizer {
	// Create the i18n bundle for English (fallback language)
	bundle := i18n.NewBundle(language.English)

	// Register TOML as the format parser
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Load translation files (in TOML format)
	// Use filepath.Glob to match all files matching the pattern
	files, err := filepath.Glob("locales/active.*.toml")
	if err != nil {
		logger.Error("Error matching translation files", zap.Error(err))
		return nil
	}

	// Load all matching translation files
	for _, file := range files {
		_, err := bundle.LoadMessageFile(file)
		if err != nil {
			logger.Error("Could not load translation file", zap.String("file", file), zap.Error(err))
		}
	}

	// Return a new Localizer with the loaded bundle
	return &Localizer{bundle: bundle, logger: logger}
}

// T translates a given key using the detected language
func (l *Localizer) T(lang, key string, data map[string]any) string {
	// Create a new Localizer instance for the given language
	loc := i18n.NewLocalizer(l.bundle, lang)

	// Try to localize the message
	translated, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: data,
	})

	if err != nil {
		l.logger.Warn("Missing translation", zap.String("key", key), zap.String("lang", lang))
		return key // Return the key instead of crashing
	}

	return translated
}

// GetBundle returns the unexported bundle
func (l *Localizer) GetBundle() *i18n.Bundle {
	return l.bundle
}
