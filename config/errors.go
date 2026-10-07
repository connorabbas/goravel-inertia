package config

import "goravel/app/facades"

func init() {
	facades.Config().Add("errors", map[string]any{
		"defaults": map[string]any{
			"4xx": map[string]any{
				"detail": "Sorry, your request could not be completed.",
				"icon":   "i-lucide-circle-alert",
			},
			"5xx": map[string]any{
				"detail": "Whoops, something went wrong on our end. Please try again.",
				"icon":   "i-lucide-server-crash",
			},
		},
		"statuses": map[string]any{
			"401": map[string]any{
				"detail": "Please sign in to continue.",
				"icon":   "i-lucide-log-in",
			},
			"403": map[string]any{
				"detail": "Sorry, you are unauthorized to access this resource/action.",
				"icon":   "i-lucide-shield-alert",
			},
			"404": map[string]any{
				"detail": "Sorry, the resource you are looking for could not be found.",
				"icon":   "i-lucide-search-x",
			},
			// Presentation only; this does not implement session-expiry handling.
			"419": map[string]any{
				"detail": "The page expired, please try again.",
				"icon":   "i-lucide-clock-alert",
			},
			"429": map[string]any{
				"detail": "You have made too many requests. Please wait and try again.",
				"icon":   "i-lucide-timer",
			},
			"500": map[string]any{
				"detail": "Whoops, something went wrong on our end. Please try again.",
				"icon":   "i-lucide-server-crash",
			},
			"503": map[string]any{
				"detail": "Sorry, we are doing some maintenance. Please check back soon.",
				"icon":   "i-lucide-construction",
			},
		},
	})
}
