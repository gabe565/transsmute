package cmd

import "runtime/debug"

func buildVersion(version string) string {
	var commit string
	var modified bool
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				commit = setting.Value
			case "vcs.modified":
				if setting.Value == "true" {
					modified = true
				}
			}
		}
	}

	if commit != "" {
		version += " ("
		if modified {
			version += "*"
		}
		if len(commit) > 8 {
			commit = commit[:8]
		}
		version += commit + ")"
	}
	return version
}
