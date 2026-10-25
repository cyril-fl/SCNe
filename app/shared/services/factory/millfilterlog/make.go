package millfilterlog

import (
	"solopg/app/shared/services/logs"
	"strings"

	plumprimitive "github.com/cyril-fl/plum/primitive"

	"github.com/spf13/pflag"
)

// Mill filters logs from the --filter and --tail flags.
type Mill = mill

type mill struct {
	plumprimitive.Mill

	list  []logs.Log
	flags *pflag.FlagSet
}

// Template holds the logs and the flags to read.
type Template struct {
	LogList []logs.Log
	Flags   *pflag.FlagSet
}

// New returns a Mill over params.LogList.
func New(params Template) *mill {
	return &mill{
		list:  params.LogList,
		flags: params.Flags,
	}
}

// Getters & Setters

// GetList returns the current logs.
func (m *mill) GetList() []logs.Log {
	return m.list
}

// GetStringList returns the logs formatted one per entry.
func (m *mill) GetStringList() []string {
	strings := make([]string, len(m.list))
	for i, log := range m.list {
		strings[i] = log.String()
	}
	return strings
}

// GetJoinedStringList returns the logs formatted one per line.
func (m *mill) GetJoinedStringList() string {
	return strings.Join(m.GetStringList(), "\n")
}

// Handlers

// Filter keeps the logs of the --filter kind; an empty or unknown kind keeps all.
func (m *mill) Filter() *mill {
	if m.HasErr() {
		return m
	}

	filter, err := m.flags.GetString("filter")
	if err != nil {
		return m
	}

	if filter == "" {
		return m
	}

	if !logs.AssertKind(logs.Kind(filter)) {
		return m
	}

	filtered := make([]logs.Log, 0, len(m.list))
	for _, log := range m.list {
		if log.Type != logs.Kind(filter) {
			continue
		}

		filtered = append(filtered, log)
	}

	m.list = filtered

	return m
}

// Tail keeps the last --tail logs; a value out of range keeps all.
func (m *mill) Tail() *mill {
	if m.HasErr() {
		return m
	}

	tail, err := m.flags.GetInt("tail")
	if err != nil {
		return m
	}

	if tail <= 0 || tail > len(m.list) {
		return m
	}

	if tail < len(m.list) {
		m.list = m.list[len(m.list)-tail:]
	}

	return m
}

// Helpers
