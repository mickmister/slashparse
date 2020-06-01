package slashparse

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

type newSlashCommandTests struct {
	testName      string
	args          string
	want          SlashCommand
	configPath    string
	expectedError error
}

func TestNewSlashCommand(t *testing.T) {
	tests := []newSlashCommandTests{
		{
			testName:   "simple test",
			args:       "/print",
			configPath: "./examples/helloWorld/simple.yaml",
			want: SlashCommand{
				Name:        "Print",
				Description: "Echos back what you type.",
				Arguments: []Argument{
					{
						Name:        "text",
						ArgType:     "quoted text",
						Description: "text you want to print",
						ErrorMsg:    "foo is not a valid value for text. Expected format is quoted text.",
						Position:    1,
					},
				},
				Values: map[string]string{},
			},
		},
		{
			testName:      "invalid command test",
			args:          "/pssrint",
			configPath:    "./examples/helloWorld/simple.yaml",
			expectedError: errors.New("pssrint is not a valid command"),
		},
		{
			testName:   "quoted text paramater value test",
			args:       `/print "foo bar"`,
			configPath: "./examples/helloWorld/simple.yaml",
			want: SlashCommand{
				Name:        "Print",
				Description: "Echos back what you type.",
				Arguments: []Argument{
					{
						Name:        "text",
						ArgType:     "quoted text",
						Description: "text you want to print",
						ErrorMsg:    "foo is not a valid value for text. Expected format is quoted text.",
						Position:    1,
					},
				},
				Values: map[string]string{
					"text": "foo bar",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			slashDef, _ := ioutil.ReadFile(test.configPath)

			newSlash, err := NewSlashCommand(test.args, slashDef)
			if err != nil {
				assert.Equal(t, test.expectedError, err)
			}
			assert.Equal(t, test.want, newSlash)
		})
	}
}

func TestGetSlashHelp(t *testing.T) {
	testYamlPath := "./examples/helloWorld/simple.yaml"

	args := "/print"

	slashDef, _ := ioutil.ReadFile(testYamlPath)
	newSlash, _ := NewSlashCommand(args, slashDef)

	got := newSlash.GetSlashHelp()

	want := `## Print Help
* Echos back what you type. *

### Arguments

* text: text you want to print
`
	assert.Equal(t, want, got)
}

type getCommandStringTests struct {
	testName    string
	args        string
	want        string
	expectError bool
}

func TestGetCommandString(t *testing.T) {
	tests := []getCommandStringTests{
		{
			testName: "valid print example",
			args:     "/print",
			want:     "Print",
		},
		{
			testName:    "invalid print example",
			args:        "",
			want:        "",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			testYamlPath := "./examples/helloWorld/simple.yaml"

			slashDef, _ := ioutil.ReadFile(testYamlPath)
			newSlash, _ := NewSlashCommand(test.args, slashDef)
			got, err := newSlash.GetCommandString(test.args)
			if err != nil {
				assert.Equal(t, test.expectError, true)
			} else {
				assert.Equal(t, test.want, got)
			}
		})
	}
}

func TestGetPositionalArgs(t *testing.T) {
	got := GetPositionalArgs("foo \"man chu\"  \\choo wow")
	want := []string{"foo", "man chu", "\\choo", "wow"}
	assert.Equal(t, want, got)
}
