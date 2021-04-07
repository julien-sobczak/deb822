package deb822

import (
	"strings"
)

type Formatter struct {
	foldedFields    []string
	multilineFields []string
}

func NewFormatter() *Formatter {
	return &Formatter{}
}

func (f *Formatter) SetFoldedFields(fields ...string) {
	f.foldedFields = append(f.foldedFields, fields...)
}

func (f *Formatter) SetMultilineFields(fields ...string) {
	f.multilineFields = append(f.multilineFields, fields...)
}

func (f *Formatter) IsFoldedField(field string) bool {
	for _, f := range f.foldedFields {
		if f == field {
			return true
		}
	}
	return false
}

func (f *Formatter) IsMultilineField(field string) bool {
	for _, f := range f.multilineFields {
		if f == field {
			return true
		}
	}
	return false
}

func (f *Formatter) Format(document Document) string {
	var sb strings.Builder
	for i, paragraph := range document.Paragraphs {
		if i > 0 {
			sb.WriteString("\n")
		}
		for _, field := range paragraph.Order {
			value := f.formatValue(field, paragraph.Value(field))
			sb.WriteString(field)
			sb.WriteString(":")
			if !strings.HasPrefix(value, "\n") {
				sb.WriteString(" ")
			}
			sb.WriteString(value)
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func (f *Formatter) formatValue(field, value string) string {
	if f.IsFoldedField(field) {
		ret := value
		ret = strings.ReplaceAll(ret, "\n\n", "\n.\n")
		ret = strings.ReplaceAll(ret, "\n", "\n ")
		return ret
	} else if f.IsMultilineField(field) {
		var sb strings.Builder
		for _, str := range strings.Split(value, "\n") {
			sb.WriteString("\n ")
			sb.WriteString(str)
		}
		return sb.String()
	} else { // simple field
		return value
	}
}
