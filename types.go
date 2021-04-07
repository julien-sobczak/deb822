package deb822

// Implementation inspired on https://gitlab.yuribugelli.it/debian/go-debian/

// Document regroups paragraphs like the /var/lib/dpkg/status file.
type Document struct {
	Paragraphs []Paragraph
}

// A Paragraph is a block of RFC2822-like key value pairs. This struct contains
// two methods to fetch values, a Map called Values, and a Slice called
// Order, which maintains the ordering as defined in the RFC2822-like block
type Paragraph struct {
	Values map[string]string
	Order  []string
}

func (p Paragraph) Value(field string) string {
	if value, ok := p.Values[field]; ok {
		return value
	}
	return ""
}
