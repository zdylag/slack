package slack

import (
	"regexp"
)

func replaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
	result := ""
	lastIndex := 0

	for _, indices := range re.FindAllSubmatchIndex([]byte(str), -1) {
		groups := []string{}
		for i := 0; i < len(indices); i += 2 {
			if indices[i] == -1 {
				groups = append(groups, "")
				continue
			}
			groups = append(groups, str[indices[i]:indices[i+1]])
		}

		result += str[lastIndex:indices[0]] + repl(groups)
		lastIndex = indices[1]
	}

	return result + str[lastIndex:]
}
