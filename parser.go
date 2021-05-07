package deb822

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Parser struct {
	reader *bufio.Reader
}

func NewParser(reader io.Reader) (*Parser, error) {
	bufioReader := bufio.NewReader(reader)
	ret := Parser{
		reader: bufioReader,
	}
	return &ret, nil
}

func (p *Parser) Parse() (Document, error) {
	ret := Document{
		[]Paragraph{},
	}
	for {
		paragraph, err := p.next()
		if err == io.EOF {
			return ret, nil
		} else if err != nil {
			return ret, err
		}
		if len(paragraph.Values) == 0 {
			// Ignore
			continue
		}
		p.postProcess(paragraph)
		ret.Paragraphs = append(ret.Paragraphs, *paragraph)
	}
}

func (p *Parser) next() (*Paragraph, error) {
	paragraph := Paragraph{
		Order:  []string{},
		Values: map[string]string{},
	}
	var lastKey string
	var lastLineIsFirst bool

	for {
		line, err := p.reader.ReadString('\n')
		if err == io.EOF {
			/* Let's return the parsed paragraph if we have it */
			if len(paragraph.Order) > 0 {
				return &paragraph, nil
			}
			/* Else, let's go ahead and drop the EOF out raw */
			return nil, err
		} else if err != nil {
			return nil, err
		}

		if line == "\n" {
			/* Lines are ended by a blank line; so we're able to go ahead
			* and return this guy as-is. All set. Done. Finished. */
			return &paragraph, nil
		}

		/* Right, so we have a line in one of the following formats:
		 *
		 * "Key: Value"
		 * " Foobar"
		 *
		 * Foobar is seen as a continuation of the last line, and the
		 * Key line is a Key/Value mapping.
		 */

		if strings.HasPrefix(line, " ") {
			/* This is a continuation line; so we're going to go ahead and
			 * clean it up, and throw it into the list. We're going to remove
			 * the space, and if it's a line that only has a dot on it, we'll
			 * remove that too (since " .\n" is actually "\n"). We only
			 * trim off space on the right hand, because indentation under
			 * the single space is up to the data format. Not us. */

			/* TrimFunc(line[1:], unicode.IsSpace) is identical to calling
			 * TrimSpace. */
			line = strings.TrimRightFunc(line[1:], unicode.IsSpace) + "\n"

			if line == ".\n" {
				line = "\n"
			}

			if lastLineIsFirst && strings.TrimSpace(paragraph.Values[lastKey]) != "" { // Folded field?
				paragraph.Values[lastKey] += "\n"
			}

			lastLineIsFirst = false
			paragraph.Values[lastKey] += line
			continue
		}

		/* So, if we're here, we've got a key line. Let's go ahead and split
		 * this on the first key, and set that guy */
		els := strings.SplitN(line, ":", 2)
		if len(els) != 2 {
			return nil, fmt.Errorf("bad line: '%s' has no ':'", line)
		}

		/* We'll go ahead and take off any leading spaces */
		lastKey = strings.TrimSpace(els[0])
		value := strings.TrimSpace(els[1])
		lastLineIsFirst = true

		paragraph.Order = append(paragraph.Order, lastKey)
		paragraph.Values[lastKey] = value
	}
}

func (p *Parser) postProcess(paragraph *Paragraph) {
	for field, value := range paragraph.Values {
		paragraph.Values[field] = strings.TrimSpace(value)
	}
}
