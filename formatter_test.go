package deb822_test

import (
	"strings"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/julien-sobczak/deb822"
)

func TestFormatter(t *testing.T) {
	input := `Package: bash
Essential: yes
Status: install ok installed
Priority: required
Section: shells
Installed-Size: 6470
Maintainer: Matthias Klose <doko@debian.org>
Architecture: amd64
Multi-Arch: foreign
Source: bash (5.1-2)
Version: 5.1-2+b1
Replaces: bash-completion (<< 20060301-0), bash-doc (<= 2.05-1)
Depends: base-files (>= 2.1.12), debianutils (>= 2.15)
Pre-Depends: libc6 (>= 2.25), libtinfo6 (>= 6)
Recommends: bash-completion (>= 20060301-0)
Suggests: bash-doc
Conflicts: bash-completion (<< 20060301-0)
Conffiles:
 /etc/bash.bashrc 89269e1298235f1b12b4c16e4065ad0d
 /etc/skel/.bash_logout 22bfb8c1dd94b5f3813a2b25da67463f
 /etc/skel/.bashrc ee35a240758f374832e809ae0ea4883a
 /etc/skel/.profile f4e81ade7d6f9fb342541152d08e7a97
Description: GNU Bourne Again SHell
 Bash is an sh-compatible command language interpreter that executes
 commands read from the standard input or from a file.  Bash also
 incorporates useful features from the Korn and C shells (ksh and csh).
 .
 Bash is ultimately intended to be a conformant implementation of the
 IEEE POSIX Shell and Tools specification (IEEE Working Group 1003.2).
 .
 The Programmable Completion Code, by Ian Macdonald, is now found in
 the bash-completion package.
Homepage: http://tiswww.case.edu/php/chet/bash/bashtop.html
`
	document := parse(t, input)
	formatter := deb822.NewFormatter()
	formatter.SetFoldedFields("Description")
	formatter.SetMultilineFields("Conffiles")
	output := formatter.Format(document)

	if output != input {
		t.Errorf("Differences found:\n%s",
			diff.CharacterDiff(input, output))
		t.Log(output)
	}
}

/* Test Helper */

func parse(t *testing.T, input string) deb822.Document {
	parser, _ := deb822.NewParser(strings.NewReader(input))
	document, err := parser.Parse()
	if err != nil {
		t.Fatalf("Unable to parse DEB822 document: %v", err)
	}
	return document
}
