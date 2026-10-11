package sharedtui

// keys

type EventKey string

const (
	// Navigation keys
	KEY_ENTER     EventKey = "enter"
	KEY_ESC       EventKey = "esc"
	KEY_BACKSPACE EventKey = "backspace"
	KEY_TAB       EventKey = "tab"
	KEY_SHIFT_TAB EventKey = "shift+tab"

	// Arrow keys
	KEY_UP      EventKey = "up"
	KEY_DOWN    EventKey = "down"
	KEY_LEFT    EventKey = "left"
	KEY_RIGHT   EventKey = "right"
	KEY_PG_UP   EventKey = "pgup"
	KEY_PG_DOWN EventKey = "pgdown"

	// Shortcut
	SHORT_CTRL_C    EventKey = "ctrl+c"
	SHORT_CTRL_Q    EventKey = "ctrl+q"
	SHORT_CTRL_N    EventKey = "ctrl+n"
	SHORT_SHIFT_TAB EventKey = "shift+tab"

	// Command keys
	CMD_SHIFT_ENTER EventKey = "shift+enter"
	CMD_ALT_ENTER   EventKey = "alt+enter"
	CMD_CTRL_S      EventKey = "ctrl+s"

	// Other keys
	KEY_EMPTY EventKey = ""
)
