package style

import (
	"strings"

	"charm.land/lipgloss/v2"
	"gopkg.in/yaml.v3"
)

// Yaml renders data as YAML colored with the theme: keys in accent, values in text, syntax and empty values in muted.
func Yaml(data any) (string, error) {
	var node yaml.Node
	if err := node.Encode(data); err != nil {
		return "", err
	}

	var b strings.Builder
	writeYamlNode(&b, &node, 0)

	return strings.TrimSuffix(b.String(), "\n"), nil
}

// Helpers

func writeYamlNode(b *strings.Builder, node *yaml.Node, indent int) {
	pad := strings.Repeat(" ", indent)

	switch node.Kind {
	case yaml.DocumentNode:
		for _, child := range node.Content {
			writeYamlNode(b, child, indent)
		}
	case yaml.MappingNode:
		for i := 0; i+1 < len(node.Content); i += 2 {
			b.WriteString(pad + Accent().Render(node.Content[i].Value) + Muted().Render(":"))
			writeYamlValue(b, node.Content[i+1], indent)
		}
	case yaml.SequenceNode:
		for _, item := range node.Content {
			dash := pad + Muted().Render("-")
			if item.Kind != yaml.MappingNode {
				b.WriteString(dash)
				writeYamlValue(b, item, indent)
				continue
			}
			// The item's first key goes on the dash line: "- name: yaml".
			var mapping strings.Builder
			writeYamlNode(&mapping, item, indent+2)
			b.WriteString(dash + " " + strings.TrimPrefix(mapping.String(), pad+"  "))
		}
	}
}

// writeYamlValue writes the value that follows a key or a dash, starting on the same line.
func writeYamlValue(b *strings.Builder, node *yaml.Node, indent int) {
	switch {
	case node.Kind == yaml.ScalarNode && strings.Contains(node.Value, "\n"):
		indicator := "|-"
		if strings.HasSuffix(node.Value, "\n") {
			indicator = "|"
		}
		b.WriteString(" " + Muted().Render(indicator) + "\n")
		for _, line := range strings.Split(strings.TrimSuffix(node.Value, "\n"), "\n") {
			b.WriteString(strings.Repeat(" ", indent+2) + text().Render(line) + "\n")
		}
	case node.Kind == yaml.ScalarNode:
		b.WriteString(" " + yamlScalar(node) + "\n")
	case len(node.Content) == 0:
		empty := "{}"
		if node.Kind == yaml.SequenceNode {
			empty = "[]"
		}
		b.WriteString(" " + Muted().Render(empty) + "\n")
	default:
		b.WriteString("\n")
		writeYamlNode(b, node, indent+2)
	}
}

// yamlScalar lets yaml quote the value when needed, e.g. "#C9A26B" would otherwise read as a comment.
func yamlScalar(node *yaml.Node) string {
	if node.Value == "" {
		return Muted().Render(`""`)
	}

	out, err := yaml.Marshal(node)
	if err != nil {
		return text().Render(node.Value)
	}

	return text().Render(strings.TrimSuffix(string(out), "\n"))
}

func text() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(ColorText())
}
