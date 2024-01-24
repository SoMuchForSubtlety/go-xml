package xsdgen

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExamples(t *testing.T) {
	cases := []struct {
		name        string
		sourceFiles []string
		namespace   string
	}{
		{
			name:        "base64 binary",
			sourceFiles: []string{"testdata/base64.xsd"},
			namespace:   "http://example.org/",
		},
		{
			name:        "simple union",
			sourceFiles: []string{"testdata/simple-union.xsd"},
			namespace:   "http://example.org/",
		},
		{
			name:        "enumeration constants",
			sourceFiles: []string{"testdata/enumeration.xsd"},
			namespace:   "http://example.org/",
		},
		{
			name:        "library schema",
			namespace:   "http://dyomedea.com/ns/library",
			sourceFiles: []string{"testdata/library.xsd"},
		},
		{
			name:        "purchas order schema",
			namespace:   "http://www.example.com/PO1",
			sourceFiles: []string{"testdata/po1.xsd"},
		},
		{
			name:        "US treasure SDN",
			namespace:   "http://tempuri.org/sdnList.xsd",
			sourceFiles: []string{"testdata/sdn.xsd"},
		},
		{
			name:        "SOAP",
			namespace:   "http://schemas.xmlsoap.org/soap/encoding/",
			sourceFiles: []string{"testdata/soap11.xsd"},
		},
		{
			name:        "simple struct",
			namespace:   "http://example.org/ns",
			sourceFiles: []string{"testdata/simple-struct.xsd"},
		},
		{
			name:        "mixed data",
			namespace:   "http://example.org",
			sourceFiles: []string{"testdata/mixed-complex.xsd"},
		},
		{
			name:        "dupe names",
			sourceFiles: []string{"testdata/duplicate-names.xsd", "testdata/duplicate-names2.xsd"},
		},
		{
			name:        "inderect element ref",
			namespace:   "http://example.org/",
			sourceFiles: []string{"testdata/indirect-ref.xsd"},
		},
		{
			name:        "dupe type 2",
			namespace:   "http://example.org/",
			sourceFiles: []string{"testdata/dupe-enum.xsd"},
		},
		{
			name:        "optional time",
			sourceFiles: []string{"testdata/optional-time.xsd"},
		},
		{
			name:        "anon",
			sourceFiles: []string{"testdata/anon.xsd"},
		},
		{
			name:        "anon2",
			sourceFiles: []string{"testdata/anon2.xsd"},
		},
		{
			name:        "redundant union",
			sourceFiles: []string{"testdata/redundant-union.xsd"},
		},
		{
			name:        "extension",
			sourceFiles: []string{"testdata/extension.xsd"},
		},
		{
			name:        "optional struct",
			sourceFiles: []string{"testdata/optional-struct.xsd"},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			testGen(t, c.namespace, c.sourceFiles...)
		})
	}
}

type testLogger testing.T

func (t *testLogger) Printf(format string, v ...interface{}) {
	t.Logf(format, v...)
}

func testGen(t *testing.T, ns string, files ...string) string {
	file, err := os.CreateTemp("", "xsdgen")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	var cfg Config
	cfg.Option(DefaultOptions...)
	cfg.Option(LogOutput((*testLogger)(t)))

	args := []string{"-v", "-o", file.Name()}
	if ns != "" {
		args = append(args, "-ns", ns)
	}
	err = cfg.GenCLI(append(args, files...)...)
	if err != nil {
		t.Error(err)
	}
	data, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	ext := path.Ext(files[0])
	goldenPath := files[0][0:len(files[0])-len(ext)] + ".go.golden"
	// check if we have a golden file to compare to
	if _, err := os.Stat(goldenPath); err == nil {
		goldenData, err := os.ReadFile(goldenPath)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, string(goldenData), string(data), "the generated code does not match the expected output")
	} else {
		// create a new golden file if there is none
		err = os.WriteFile(goldenPath, data, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	return string(data)
}
