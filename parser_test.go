package deb822_test

import (
	"strings"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/julien-sobczak/deb822"
)

func TestParserSimple(t *testing.T) {
	input := `Package: foo
Architecture: any
Description: Foo package.
 .
 This description is a multiline field.

Package: bar
Architecture: amd64 sparc
Description: Bar package.
`
	parser, _ := deb822.NewParser(strings.NewReader(input))
	document, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse the document: %v", err)
	}
	if len(document.Paragraphs) != 2 {
		t.Fatalf("Invalid number of paragraphs. Expected: %d, got: %d", 2, len(document.Paragraphs))
	}
	if document.Paragraphs[0].Value("Package") != "foo" {
		t.Errorf("Wrong package name: %v", document.Paragraphs[0].Value("Package"))
	}
	if document.Paragraphs[1].Value("Package") != "bar" {
		t.Errorf("Wrong package name: %v", document.Paragraphs[1].Value("Package"))
	}
}

func TestParserAdvanced(t *testing.T) {
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
	parser, _ := deb822.NewParser(strings.NewReader(input))
	document, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse the document: %v", err)
	}
	if len(document.Paragraphs) != 1 {
		t.Fatalf("Invalid number of paragraphs. Expected: %d, got: %d", 2, len(document.Paragraphs))
	}

	paragraph := document.Paragraphs[0]

	// Simple value
	expectedDepends := `base-files (>= 2.1.12), debianutils (>= 2.15)`
	if paragraph.Value("Depends") != expectedDepends {
		t.Errorf("Mismatch found in 'Depends'. Differences:\n%s",
			diff.CharacterDiff(paragraph.Value("Depends"), expectedDepends))
	}

	// Folded value
	expectedDescription := `GNU Bourne Again SHell
Bash is an sh-compatible command language interpreter that executes
commands read from the standard input or from a file.  Bash also
incorporates useful features from the Korn and C shells (ksh and csh).

Bash is ultimately intended to be a conformant implementation of the
IEEE POSIX Shell and Tools specification (IEEE Working Group 1003.2).

The Programmable Completion Code, by Ian Macdonald, is now found in
the bash-completion package.`
	if paragraph.Value("Description") != expectedDescription {
		t.Errorf("Mismatch found in 'Description'. Differences:\n%s",
			diff.CharacterDiff(paragraph.Value("Description"), expectedDescription))
	}

	// Multiline value
	expectedConffiles := `/etc/bash.bashrc 89269e1298235f1b12b4c16e4065ad0d
/etc/skel/.bash_logout 22bfb8c1dd94b5f3813a2b25da67463f
/etc/skel/.bashrc ee35a240758f374832e809ae0ea4883a
/etc/skel/.profile f4e81ade7d6f9fb342541152d08e7a97`
	if paragraph.Value("Conffiles") != expectedConffiles {
		t.Errorf("Mismatch found in 'Conffiles'. Differences:\n%s",
			diff.CharacterDiff(paragraph.Value("Conffiles"), expectedConffiles))
	}
}
