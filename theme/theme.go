package theme

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/scmn-dev/tran/constants"
)

type Theme struct {
	SelectedTreeItemColor                lipgloss.AdaptiveColor
	UnselectedTreeItemColor              lipgloss.AdaptiveColor
	ActiveBoxBorderColor                 lipgloss.AdaptiveColor
	InactiveBoxBorderColor               lipgloss.AdaptiveColor
	SpinnerColor                         lipgloss.AdaptiveColor
	StatusBarSelectedFileForegroundColor lipgloss.AdaptiveColor
	StatusBarSelectedFileBackgroundColor lipgloss.AdaptiveColor
	StatusBarBarForegroundColor          lipgloss.AdaptiveColor
	StatusBarBarBackgroundColor          lipgloss.AdaptiveColor
	StatusBarTotalFilesForegroundColor   lipgloss.AdaptiveColor
	StatusBarTotalFilesBackgroundColor   lipgloss.AdaptiveColor
	StatusBarLogoForegroundColor         lipgloss.AdaptiveColor
	StatusBarLogoBackgroundColor         lipgloss.AdaptiveColor
	ErrorColor                           lipgloss.AdaptiveColor
	DefaultTextColor                     lipgloss.AdaptiveColor
}

// appColors contains the different types of colors.
type appColors struct {
	white              string
	darkGray           string
	red                string
	black              string
}

// Colors contains the different kinds of colors and their values.
var colors = appColors{
	white:              "#FFFDF5",
	darkGray:           constants.DARK_GRAY_COLOR,
	red:                "#cc241d",
	black:              "#000000",
}

// themeMap represents the mapping of different themes.
var themeMap = map[string]Theme{
	"default": {
		SelectedTreeItemColor:                lipgloss.AdaptiveColor{Dark: constants.PRIMARY_COLOR, Light: constants.PRIMARY_COLOR},
		UnselectedTreeItemColor:              lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
		ActiveBoxBorderColor:                 lipgloss.AdaptiveColor{Dark: constants.PRIMARY_COLOR, Light: constants.PRIMARY_COLOR},
		InactiveBoxBorderColor:               lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
		SpinnerColor:                         lipgloss.AdaptiveColor{Dark: constants.PRIMARY_COLOR, Light: constants.PRIMARY_COLOR},
		StatusBarSelectedFileForegroundColor: lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarSelectedFileBackgroundColor: lipgloss.AdaptiveColor{Dark: "#4880EC", Light: "#4880EC"},
		StatusBarBarForegroundColor:          lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarBarBackgroundColor:          lipgloss.AdaptiveColor{Dark: colors.darkGray, Light: colors.darkGray},
		StatusBarTotalFilesForegroundColor:   lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarTotalFilesBackgroundColor:   lipgloss.AdaptiveColor{Dark: "#1E6AFF", Light: "#1E6AFF"},
		StatusBarLogoForegroundColor:         lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarLogoBackgroundColor:         lipgloss.AdaptiveColor{Dark: "#1747A6", Light: "#1747A6"},
		ErrorColor:                           lipgloss.AdaptiveColor{Dark: colors.red, Light: colors.red},
		DefaultTextColor:                     lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
	},
}

// GetTheme returns a theme based on the given name.
func GetTheme(theme string) Theme {
	return themeMap["default"]
}
