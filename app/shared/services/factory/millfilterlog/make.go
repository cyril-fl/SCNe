package millfilterlog

import (
	"solopg/app/shared/services/logs"
	"strings"

	plumprimitive "github.com/cyril-fl/plum/primitive"

	"github.com/spf13/pflag"
)

type Mill = mill

type mill struct {
	plumprimitive.Mill

	list  []logs.Log
	flags *pflag.FlagSet
}

type Template struct {
	LogList []logs.Log
	Flags   *pflag.FlagSet
}

func New(params Template) *mill {
	return &mill{
		list:  params.LogList,
		flags: params.Flags,
	}
}

// Getters & Setters
func (m *mill) GetList() []logs.Log {
	return m.list
}

func (m *mill) GetStringList() []string {
	strings := make([]string, len(m.list))
	for i, log := range m.list {
		strings[i] = log.String()
	}
	return strings
}

func (m *mill) GetJoinedStringList() string {
	return strings.Join(m.GetStringList(), "\n")
}

// Handlers
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
