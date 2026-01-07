package branch

import "testing"

func TestGenerate(t *testing.T) {
	tests := []struct {
		name        string
		branchType  string
		ticket      string
		description []string
		want        string
	}{
		{
			name:        "feature with ticket and description",
			branchType:  "feat",
			ticket:      "PIP-1234",
			description: []string{"implement", "new", "feature"},
			want:        "feat/pip-1234-implement-new-feature",
		},
		{
			name:        "fix with ticket only",
			branchType:  "fix",
			ticket:      "PIP-5678",
			description: []string{},
			want:        "fix/pip-5678",
		},
		{
			name:        "chore without ticket",
			branchType:  "chore",
			ticket:      "",
			description: []string{"update", "dependencies"},
			want:        "chore/update-dependencies",
		},
		{
			name:        "docs with GitHub issue ticket",
			branchType:  "docs",
			ticket:      "#123",
			description: []string{"add", "api", "documentation"},
			want:        "docs/123-add-api-documentation",
		},
		{
			name:        "empty description returns just type",
			branchType:  "feat",
			ticket:      "",
			description: []string{},
			want:        "feat",
		},
		{
			name:        "description with special characters",
			branchType:  "fix",
			ticket:      "",
			description: []string{"fix", "the", "bug!", "it's", "important"},
			want:        "fix/fix-the-bug-its-important",
		},
		{
			name:        "description with underscores",
			branchType:  "tests",
			ticket:      "",
			description: []string{"add", "unit_tests"},
			want:        "tests/add-unit-tests",
		},
		{
			name:        "description with multiple spaces",
			branchType:  "feat",
			ticket:      "",
			description: []string{"add", "new", "", "feature"},
			want:        "feat/add-new-feature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Generate(tt.branchType, tt.ticket, tt.description)
			if got != tt.want {
				t.Errorf("Generate(%q, %q, %v) = %q, want %q", tt.branchType, tt.ticket, tt.description, got, tt.want)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple lowercase",
			input: "hello world",
			want:  "hello-world",
		},
		{
			name:  "uppercase conversion",
			input: "HELLO WORLD",
			want:  "hello-world",
		},
		{
			name:  "mixed case",
			input: "Hello World",
			want:  "hello-world",
		},
		{
			name:  "with underscores",
			input: "hello_world",
			want:  "hello-world",
		},
		{
			name:  "with special characters",
			input: "hello! world@ test#",
			want:  "hello-world-test",
		},
		{
			name:  "multiple spaces",
			input: "hello    world",
			want:  "hello-world",
		},
		{
			name:  "multiple hyphens",
			input: "hello---world",
			want:  "hello-world",
		},
		{
			name:  "leading and trailing spaces",
			input: "  hello world  ",
			want:  "hello-world",
		},
		{
			name:  "leading and trailing hyphens",
			input: "-hello-world-",
			want:  "hello-world",
		},
		{
			name:  "numbers and letters",
			input: "version 2.0 release",
			want:  "version-20-release",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "only special characters",
			input: "!!!@@@###",
			want:  "",
		},
		{
			name:  "ticket with description",
			input: "PIP-1234 implement new feature",
			want:  "pip-1234-implement-new-feature",
		},
		{
			name:  "GitHub issue format",
			input: "#123 fix bug",
			want:  "123-fix-bug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slugify(tt.input)
			if got != tt.want {
				t.Errorf("slugify(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
