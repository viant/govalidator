package govalidator

import (
	"strings"
)

type (
	//Tag represents validation tag
	Tag struct {
		Checks    []Check
		Omitempty bool
		Required  bool
		SkipPath  bool
	}

	//Check represents validation check
	Check struct {
		Name       string
		Parameters []string
		Message    string
	}
)

//ParseTag parses rule
func ParseTag(tagString string) *Tag {
	tag := &Tag{}

	elements := extractElements(tagString)
	if len(elements) == 0 {
		return tag
	}
	tag.Required = strings.Contains(strings.ToLower(tagString), "required")
	tag.Omitempty = strings.Contains(strings.ToLower(tagString), "omitempty")
	tag.SkipPath = strings.Contains(strings.ToLower(tagString), "skippath")

	for _, element := range elements {
		check := Check{}
		pair := strings.Split(element, "=")
		switch len(pair) {
		case 2:
			switch strings.ToLower(strings.TrimSpace(pair[0])) {
			case "message":
				check.Message = strings.TrimSpace(pair[1])
			case "name":
				check.Name, check.Parameters = extractNameWithParams(strings.TrimSpace(pair[1]))
			default:
				check.Name = pair[0]
				check.Parameters = pair[1:]
			}
		case 1:
			check.Name, check.Parameters = extractNameWithParams(strings.TrimSpace(element))
		}
		switch strings.ToLower(check.Name) {
		case "omitempty", "skippath", "marker":
			continue
		}
		tag.Checks = append(tag.Checks, check)
	}
	return tag
}

func extractElements(decoded string) []string {
	var result []string

	var fragment = ""
	var inCall bool
	for i, c := range decoded {
		switch decoded[i] {
		case ',', '|':
			if inCall {
				fragment += string(c)
				continue
			}
			result = append(result, fragment)
			fragment = ""
		case '(':
			fragment += string(c)
			inCall = true

		case ')':
			fragment += string(c)
			inCall = false
		default:
			fragment += string(c)
		}
	}
	if fragment != "" {
		result = append(result, fragment)
	}
	return result
}

var emptyArgs = []string{}

func extractNameWithParams(text string) (string, []string) {
	index := strings.Index(text, "(")
	if index == -1 {
		return text, emptyArgs
	}
	name := text[:index]
	argsFragment := text[index+1:]
	index = strings.LastIndex(argsFragment, ")")
	if index != -1 {
		argsFragment = argsFragment[:index]
	}
	var params []string
	for _, item := range strings.Split(argsFragment, ",") {
		params = append(params, strings.TrimSpace(item))
	}
	return name, params
}
