package config

import (
	"goravel/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("inertia", map[string]any{
		"root_view":  "resources/inertia/app.gohtml",
		"version":    config.Env("INERTIA_VERSION", ""),
		"ssr":        false,
		"flash_keys": []string{"success", "error", "warning", "info", "message"},
		"vite": map[string]any{
			"public_path": "public",
			"build_dir":   "build",
			"hot_file":    "public/hot",
			// Only public/hot enables development assets; otherwise use the manifest.
			"dev_url": "",
		},
	})
}
