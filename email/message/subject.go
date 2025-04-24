package message

import (
	"fmt"
	"strings"
	"text/template"
)

var subjectTemplates = make(map[string]*template.Template)

var subjects = map[string]string{
	Welcome: "Important information about your {{ .AppName }} account",
}

func init() {
	for name, tmpl := range subjects {
		parsed, err := template.New(name).Parse(tmpl)
		if err != nil {
			panic(fmt.Sprintf("failed to parse email subject template '%s': %v", name, err))
		}
		subjectTemplates[name] = parsed
	}
}

func execSubjectTemplate(templateName string, fields Fields) (string, error) {
	subjectTemplate, ok := subjects[templateName]
	if !ok {
		return "", fmt.Errorf("subject template '%s' not found", templateName)
	}

	templ, err := template.New("subject").Parse(subjectTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse subject template '%s': %w", templateName, err)
	}

	var sb strings.Builder
	if err = templ.Execute(&sb, fields); err != nil {
		return "", fmt.Errorf("failed to execute subject template '%s': %w", templateName, err)
	}
	return sb.String(), nil
}
