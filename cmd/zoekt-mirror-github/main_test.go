package main

import (
	"testing"

	"github.com/google/go-github/v78/github"
)

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}

func TestHasIntersection(t *testing.T) {
	tests := []struct {
		name     string
		s1       []string
		s2       []string
		expected bool
	}{
		{
			name:     "empty slices",
			s1:       []string{},
			s2:       []string{},
			expected: false,
		},
		{
			name:     "first slice empty",
			s1:       []string{},
			s2:       []string{"a", "b"},
			expected: false,
		},
		{
			name:     "second slice empty",
			s1:       []string{"a", "b"},
			s2:       []string{},
			expected: false,
		},
		{
			name:     "no intersection",
			s1:       []string{"a", "b"},
			s2:       []string{"c", "d"},
			expected: false,
		},
		{
			name:     "single element intersection",
			s1:       []string{"a", "b"},
			s2:       []string{"b", "c"},
			expected: true,
		},
		{
			name:     "multiple intersections",
			s1:       []string{"a", "b", "c"},
			s2:       []string{"b", "c", "d"},
			expected: true,
		},
		{
			name:     "identical slices",
			s1:       []string{"a", "b", "c"},
			s2:       []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "subset",
			s1:       []string{"a"},
			s2:       []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "case sensitive - no match",
			s1:       []string{"A"},
			s2:       []string{"a"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasIntersection(tt.s1, tt.s2)
			if result != tt.expected {
				t.Errorf("hasIntersection(%v, %v) = %v, want %v", tt.s1, tt.s2, result, tt.expected)
			}
		})
	}
}

func TestFilterRepositories_Visibility(t *testing.T) {
	tests := []struct {
		name        string
		repos       []*github.Repository
		visibility  []string
		expectedLen int
		noArchived  bool
	}{
		{
			name: "no visibility filter",
			repos: []*github.Repository{
				{Visibility: stringPtr("public")},
				{Visibility: stringPtr("private")},
				{Visibility: stringPtr("internal")},
			},
			visibility:  []string{},
			expectedLen: 3,
			noArchived:  false,
		},
		{
			name: "filter public only",
			repos: []*github.Repository{
				{Visibility: stringPtr("public")},
				{Visibility: stringPtr("private")},
				{Visibility: stringPtr("internal")},
			},
			visibility:  []string{"public"},
			expectedLen: 1,
			noArchived:  false,
		},
		{
			name: "filter private only",
			repos: []*github.Repository{
				{Visibility: stringPtr("public")},
				{Visibility: stringPtr("private")},
				{Visibility: stringPtr("internal")},
			},
			visibility:  []string{"private"},
			expectedLen: 1,
			noArchived:  false,
		},
		{
			name: "filter internal only",
			repos: []*github.Repository{
				{Visibility: stringPtr("public")},
				{Visibility: stringPtr("private")},
				{Visibility: stringPtr("internal")},
			},
			visibility:  []string{"internal"},
			expectedLen: 1,
			noArchived:  false,
		},
		{
			name: "filter multiple visibilities",
			repos: []*github.Repository{
				{Visibility: stringPtr("public")},
				{Visibility: stringPtr("private")},
				{Visibility: stringPtr("internal")},
			},
			visibility:  []string{"public", "private"},
			expectedLen: 2,
			noArchived:  false,
		},
		{
			name: "no matching visibility",
			repos: []*github.Repository{
				{Visibility: stringPtr("public")},
				{Visibility: stringPtr("private")},
			},
			visibility:  []string{"internal"},
			expectedLen: 0,
			noArchived:  false,
		},
		{
			name: "all visibilities match",
			repos: []*github.Repository{
				{Visibility: stringPtr("public")},
				{Visibility: stringPtr("private")},
				{Visibility: stringPtr("internal")},
			},
			visibility:  []string{"public", "private", "internal"},
			expectedLen: 3,
			noArchived:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterRepositories(tt.repos, []string{}, []string{}, tt.noArchived, tt.visibility)
			if len(result) != tt.expectedLen {
				t.Errorf("filterRepositories() returned %d repos, want %d", len(result), tt.expectedLen)
			}
		})
	}
}

func TestFilterRepositories_Archived(t *testing.T) {
	tests := []struct {
		name        string
		repos       []*github.Repository
		noArchived  bool
		expectedLen int
	}{
		{
			name: "include archived repos",
			repos: []*github.Repository{
				{Archived: boolPtr(false)},
				{Archived: boolPtr(true)},
				{Archived: boolPtr(false)},
			},
			noArchived:  false,
			expectedLen: 3,
		},
		{
			name: "exclude archived repos",
			repos: []*github.Repository{
				{Archived: boolPtr(false)},
				{Archived: boolPtr(true)},
				{Archived: boolPtr(false)},
			},
			noArchived:  true,
			expectedLen: 2,
		},
		{
			name: "all repos archived - exclude archived",
			repos: []*github.Repository{
				{Archived: boolPtr(true)},
				{Archived: boolPtr(true)},
			},
			noArchived:  true,
			expectedLen: 0,
		},
		{
			name: "no repos archived - exclude archived",
			repos: []*github.Repository{
				{Archived: boolPtr(false)},
				{Archived: boolPtr(false)},
			},
			noArchived:  true,
			expectedLen: 2,
		},
		{
			name: "nil archived field - include archived",
			repos: []*github.Repository{
				{Archived: nil},
				{Archived: boolPtr(false)},
			},
			noArchived:  false,
			expectedLen: 2,
		},
		{
			name: "nil archived field - exclude archived",
			repos: []*github.Repository{
				{Archived: nil},
				{Archived: boolPtr(false)},
			},
			noArchived:  true,
			expectedLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterRepositories(tt.repos, []string{}, []string{}, tt.noArchived, []string{})
			if len(result) != tt.expectedLen {
				t.Errorf("filterRepositories() returned %d repos, want %d", len(result), tt.expectedLen)
			}
		})
	}
}

func TestFilterRepositories_Topics(t *testing.T) {
	tests := []struct {
		name           string
		repos          []*github.Repository
		includedTopics []string
		excludedTopics []string
		expectedLen    int
	}{
		{
			name: "no topic filters",
			repos: []*github.Repository{
				{Topics: []string{"go", "tool"}},
				{Topics: []string{"python"}},
				{Topics: []string{}},
			},
			includedTopics: []string{},
			excludedTopics: []string{},
			expectedLen:    3,
		},
		{
			name: "include single topic",
			repos: []*github.Repository{
				{Topics: []string{"go", "tool"}},
				{Topics: []string{"python"}},
				{Topics: []string{"go"}},
			},
			includedTopics: []string{"go"},
			excludedTopics: []string{},
			expectedLen:    2,
		},
		{
			name: "include multiple topics",
			repos: []*github.Repository{
				{Topics: []string{"go", "tool"}},
				{Topics: []string{"python"}},
				{Topics: []string{"rust"}},
			},
			includedTopics: []string{"go", "python"},
			excludedTopics: []string{},
			expectedLen:    2,
		},
		{
			name: "exclude single topic",
			repos: []*github.Repository{
				{Topics: []string{"go", "tool"}},
				{Topics: []string{"python"}},
				{Topics: []string{"rust"}},
			},
			includedTopics: []string{},
			excludedTopics: []string{"go"},
			expectedLen:    2,
		},
		{
			name: "exclude multiple topics",
			repos: []*github.Repository{
				{Topics: []string{"go", "tool"}},
				{Topics: []string{"python"}},
				{Topics: []string{"rust"}},
			},
			includedTopics: []string{},
			excludedTopics: []string{"go", "python"},
			expectedLen:    1,
		},
		{
			name: "include and exclude topics",
			repos: []*github.Repository{
				{Topics: []string{"go", "tool"}},
				{Topics: []string{"go", "deprecated"}},
				{Topics: []string{"python"}},
			},
			includedTopics: []string{"go"},
			excludedTopics: []string{"deprecated"},
			expectedLen:    1,
		},
		{
			name: "repo with no topics",
			repos: []*github.Repository{
				{Topics: []string{}},
				{Topics: nil},
			},
			includedTopics: []string{"go"},
			excludedTopics: []string{},
			expectedLen:    0,
		},
		{
			name: "repo matches both include and exclude - exclude wins",
			repos: []*github.Repository{
				{Topics: []string{"go", "deprecated"}},
			},
			includedTopics: []string{"go"},
			excludedTopics: []string{"deprecated"},
			expectedLen:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterRepositories(tt.repos, tt.includedTopics, tt.excludedTopics, false, []string{})
			if len(result) != tt.expectedLen {
				t.Errorf("filterRepositories() returned %d repos, want %d", len(result), tt.expectedLen)
			}
		})
	}
}

func TestFilterRepositories_Combined(t *testing.T) {
	tests := []struct {
		name        string
		repos       []*github.Repository
		include     []string
		exclude     []string
		noArchived  bool
		visibility  []string
		expectedLen int
	}{
		{
			name: "archived and visibility filter",
			repos: []*github.Repository{
				{Archived: boolPtr(false), Visibility: stringPtr("public")},
				{Archived: boolPtr(true), Visibility: stringPtr("public")},
				{Archived: boolPtr(false), Visibility: stringPtr("private")},
			},
			include:     []string{},
			exclude:     []string{},
			noArchived:  true,
			visibility:  []string{"public"},
			expectedLen: 1,
		},
		{
			name: "archived, visibility, and topics filter",
			repos: []*github.Repository{
				{Archived: boolPtr(false), Visibility: stringPtr("public"), Topics: []string{"go"}},
				{Archived: boolPtr(true), Visibility: stringPtr("public"), Topics: []string{"go"}},
				{Archived: boolPtr(false), Visibility: stringPtr("private"), Topics: []string{"go"}},
				{Archived: boolPtr(false), Visibility: stringPtr("public"), Topics: []string{"python"}},
			},
			include:     []string{"go"},
			exclude:     []string{},
			noArchived:  true,
			visibility:  []string{"public"},
			expectedLen: 1,
		},
		{
			name: "all filters combined",
			repos: []*github.Repository{
				{Archived: boolPtr(false), Visibility: stringPtr("public"), Topics: []string{"go", "tool"}},
				{Archived: boolPtr(false), Visibility: stringPtr("public"), Topics: []string{"go", "deprecated"}},
				{Archived: boolPtr(true), Visibility: stringPtr("public"), Topics: []string{"go", "tool"}},
				{Archived: boolPtr(false), Visibility: stringPtr("private"), Topics: []string{"go", "tool"}},
			},
			include:     []string{"go"},
			exclude:     []string{"deprecated"},
			noArchived:  true,
			visibility:  []string{"public"},
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterRepositories(tt.repos, tt.include, tt.exclude, tt.noArchived, tt.visibility)
			if len(result) != tt.expectedLen {
				t.Errorf("filterRepositories() returned %d repos, want %d", len(result), tt.expectedLen)
			}
		})
	}
}

func TestItoa(t *testing.T) {
	tests := []struct {
		name     string
		input    *int
		expected string
	}{
		{
			name:     "nil pointer",
			input:    nil,
			expected: "",
		},
		{
			name:     "zero",
			input:    intPtr(0),
			expected: "0",
		},
		{
			name:     "positive number",
			input:    intPtr(42),
			expected: "42",
		},
		{
			name:     "negative number",
			input:    intPtr(-123),
			expected: "-123",
		},
		{
			name:     "large number",
			input:    intPtr(999999),
			expected: "999999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := itoa(tt.input)
			if result != tt.expected {
				t.Errorf("itoa(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMarshalBool(t *testing.T) {
	tests := []struct {
		name     string
		input    bool
		expected string
	}{
		{
			name:     "true",
			input:    true,
			expected: "1",
		},
		{
			name:     "false",
			input:    false,
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := marshalBool(tt.input)
			if result != tt.expected {
				t.Errorf("marshalBool(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
