// regexp
package g

import (
	"regexp"
)

func IsReMatch(input string, pattern string) (ok bool, err error) {
	ok, err = regexp.Match(pattern, []byte(input))
	return
}

func ReFindAllStringSubMatch(input string, pattern string) (result [][]string) {
	re := regexp.MustCompile(pattern)
	result = re.FindAllStringSubmatch(input, -1)
	return
}

func ReReplaceAllString(content string, pattern string, replaceStr string) (result string) {
	re := regexp.MustCompile(pattern)
	result = re.ReplaceAllString(content, replaceStr)
	return
}
