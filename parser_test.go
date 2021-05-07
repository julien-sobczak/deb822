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
		t.Fatalf("Invalid number of paragraphs. Expected: %d, got: %d", 1, len(document.Paragraphs))
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

func TestParseRelease(t *testing.T) {
	input := `Origin: Debian
Label: Debian
Suite: testing
Codename: bullseye
Changelogs: https://metadata.ftp-master.debian.org/changelogs/@CHANGEPATH@_changelog
Date: Sat, 24 Apr 2021 14:12:29 UTC
Valid-Until: Sat, 01 May 2021 14:12:29 UTC
Acquire-By-Hash: yes
No-Support-for-Architecture-all: Packages
Architectures: all amd64 arm64 armel armhf i386 mips64el mipsel ppc64el s390x
Components: main contrib non-free
Description: Debian x.y Testing distribution - Not Released
MD5Sum:
 03f3b0f43ad6546101388b42ac724c54 466687751 main/Contents-all
 3823a8e4b3dd0fd5d4d8534ae18ccdfc    63791 main/Contents-all.diff/Index
 7d1fa31bedebd0de5fc39adf51cd3b98 30508932 main/Contents-all.gz
 7d1d9026690e4f6e208b51f4b39cb413 20366179 main/binary-all/Packages
 6b1f42a359f3b38977fb3a3441348226    63565 main/binary-all/Packages.diff/Index
 51a60d83f8edcf4191f6c5a1cb614b54  5208565 main/binary-all/Packages.gz
 0007f0860158f774977132f7c8dd3301  3919376 main/binary-all/Packages.xz
 c8bfdb5a00c15cbbf7e8e413bee667fb      101 main/binary-all/Release
 15f5a47cd121185a9af8cb79100ede86 45517628 main/binary-amd64/Packages
 ba06eabeb8d00b6fdc9b4aa75bdceff9    63565 main/binary-amd64/Packages.diff/Index
 69b3c6a4ff61ce3ddbf7ea449c5b7848 11109799 main/binary-amd64/Packages.gz
 a8d7291b6399254992e6d0fba67972ce  8194800 main/binary-amd64/Packages.xz
 74ccfb2ce9d4955e4cff92dede081800      103 main/binary-amd64/Release
 65792a893ced0c1a566d984477351ad9      104 main/source/Release
 a0ace8b8a65ecaedcdf3fbf8d8f1b044 44271077 main/source/Sources
 137aaeb7887054f89d724d9753ab5771    63565 main/source/Sources.diff/Index
 4f61b6b17165783b44da9897542d3dc7 11418555 main/source/Sources.gz
 4b7876328fab3c60819dfce675831f1d  8631172 main/source/Sources.xz
SHA256:
 a72d3f755de9b64b54dc1be0c26be5f61b4edf50d31cb767e502e7eef596943b 466687751 main/Contents-all
 c02db20480eeab9452c59982462989e2f7fc45f5462559957e656364d3679a6c    63791 main/Contents-all.diff/Index
 533e801ce8a38b84534c151295cc9f16e8b8bf457478f5cde36ae2a79e4ac1c6 30508932 main/Contents-all.gz
 da1642aa1db718eeeed1bc814578d06d2362e24b342a225a71de0cde9288c77b 124601846 main/Contents-amd64
 f60cd130789961244a7fe6b29d20c25a67c9350341b2dc2bce2c08f3c87976bc    63791 main/Contents-amd64.diff/Index
 6b2df70a836e749de809a0a8bec2e703a7118ce763168f4b3181aafb2558d979 10017612 main/Contents-amd64.gz
 6bd135ffc3b0de5087648e4373e83318af6a54957e8961e3d43056e70bbbf768 672050474 main/Contents-source
 61ee7cca92cfa4b673769b376a6ccdfbd9ad99bbbb725fcfc9a78ac9ad7e2ef0    63867 main/Contents-source.diff/Index
 98716af6e5e41bd0532da1ef7abec188bc4b3271f8519e060f735d175cc7101b 72377517 main/Contents-source.gz
 b981b54b98c670c7ce6f8427cf202434a7febe546758d8f3e9e9a1f573b38eeb 20366179 main/binary-all/Packages
 1ba02da84e9f94787a3a0a79f4b1b387e4c85e79b5d0f399063c3a1feb111cd1    63565 main/binary-all/Packages.diff/Index
 9e65018af6ece13fff2ea8315b0516857ec11aed5682bcfea31f9bb37eda06b5  5208565 main/binary-all/Packages.gz
 af3283e24c91b477dcb65c22a9423a0a468a89178bbdb8b53041ca6cefacd6a2  3919376 main/binary-all/Packages.xz
 72d0382efbdf0aaf86a59cf49e56a0e74eb1a9b6c840755c43d6346d9e6956ae      101 main/binary-all/Release
 f402576704942892a52848577f0575d78d03fd8fe245cb48268a73ad8df82f12 45517628 main/binary-amd64/Packages
 326f37f85626ffdff6dc2123a36c71bd5ef33bb7f00239cfe867435d99fa9a96    63565 main/binary-amd64/Packages.diff/Index
 609f1f4c53ca1e9dd8d72fdb9338f3163dc2d03197c348760080c931b20085d6 11109799 main/binary-amd64/Packages.gz
 78e674696a93c38f32c4cce850f55d12e782b9f1a519c47e3a0955e18014b8ff  8194800 main/binary-amd64/Packages.xz
 440e6f0db6250a47bbb7d38357cdf8a1775084c2cd05acec6a4dd66488ceee3e      103 main/binary-amd64/Release`

	parser, _ := deb822.NewParser(strings.NewReader(input))
	document, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse the document: %v", err)
	}
	if len(document.Paragraphs) != 1 {
		t.Fatalf("Invalid number of paragraphs. Expected: %d, got: %d", 1, len(document.Paragraphs))
	}

	paragraph := document.Paragraphs[0]

	expectedValue := `03f3b0f43ad6546101388b42ac724c54 466687751 main/Contents-all
3823a8e4b3dd0fd5d4d8534ae18ccdfc    63791 main/Contents-all.diff/Index
7d1fa31bedebd0de5fc39adf51cd3b98 30508932 main/Contents-all.gz
7d1d9026690e4f6e208b51f4b39cb413 20366179 main/binary-all/Packages
6b1f42a359f3b38977fb3a3441348226    63565 main/binary-all/Packages.diff/Index
51a60d83f8edcf4191f6c5a1cb614b54  5208565 main/binary-all/Packages.gz
0007f0860158f774977132f7c8dd3301  3919376 main/binary-all/Packages.xz
c8bfdb5a00c15cbbf7e8e413bee667fb      101 main/binary-all/Release
15f5a47cd121185a9af8cb79100ede86 45517628 main/binary-amd64/Packages
ba06eabeb8d00b6fdc9b4aa75bdceff9    63565 main/binary-amd64/Packages.diff/Index
69b3c6a4ff61ce3ddbf7ea449c5b7848 11109799 main/binary-amd64/Packages.gz
a8d7291b6399254992e6d0fba67972ce  8194800 main/binary-amd64/Packages.xz
74ccfb2ce9d4955e4cff92dede081800      103 main/binary-amd64/Release
65792a893ced0c1a566d984477351ad9      104 main/source/Release
a0ace8b8a65ecaedcdf3fbf8d8f1b044 44271077 main/source/Sources
137aaeb7887054f89d724d9753ab5771    63565 main/source/Sources.diff/Index
4f61b6b17165783b44da9897542d3dc7 11418555 main/source/Sources.gz
4b7876328fab3c60819dfce675831f1d  8631172 main/source/Sources.xz`

	if paragraph.Value("MD5Sum") != expectedValue {
		t.Errorf("Mismatch found in 'MD5Sum'. Differences:\n%s",
			diff.CharacterDiff(paragraph.Value("MD5Sum"), expectedValue))
	}

}
