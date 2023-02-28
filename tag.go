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
		Presence  bool
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
	elements := strings.Split(tagString, "|")
	if len(elements) == 0 {
		return nil
	}
	tag.Required = strings.Contains(strings.ToLower(tagString), "required")
	tag.Omitempty = strings.Contains(strings.ToLower(tagString), "omitempty")
	tag.SkipPath = strings.Contains(strings.ToLower(tagString), "skippath")
	tag.Presence = strings.Contains(strings.ToLower(tagString), "presence")

	for _, checkElement := range elements {
		check := Check{}
		checkSpec := strings.Split(checkElement, ";")
		for _, checkSpecElem := range checkSpec {
			pair := strings.Split(checkSpecElem, "=")
			switch len(checkSpec) {
			case 2:
				switch strings.ToLower(strings.TrimSpace(pair[0])) {
				case "message":
					check.Message = strings.TrimSpace(pair[1])
				case "name":
					check.Name, check.Parameters = extractNameWithParams(strings.TrimSpace(pair[1]))
				}
			case 1:
				check.Name, check.Parameters = extractNameWithParams(strings.TrimSpace(checkSpecElem))
			}
		}
		switch strings.ToLower(check.Name) {
		case "omitempty", "skippath", "presence":
			continue
		}
		tag.Checks = append(tag.Checks, check)
	}
	return tag
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
