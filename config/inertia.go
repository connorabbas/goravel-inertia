package config

import (
	"goravel/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("inertia", map[string]any{
		"root_view": "resources/inertia/app.gohtml",
		"version":   config.Env("INERTIA_VERSION", ""),
		"ssr":       false,
		// Add custom *_toast keys here to expose them through props.flash.
		"flash_keys": []string{"success", "error", "warning", "info", "message", "success_toast", "error_toast", "warning_toast", "warn_toast", "info_toast", "neutral_toast", "message_toast"},
		"vite": map[string]any{
			"public_path": "public",
			"build_dir":   "build",
			"hot_file":    "public/hot",
			// Only public/hot enables development assets; otherwise use the manifest.
			"dev_url": "",
		},
	})
}
