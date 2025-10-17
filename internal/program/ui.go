package program

import "fmt"

func RenderCursor(isCurrentLine bool) string {
	if !isCurrentLine {
		return " "
	}

	return fmt.Sprintf("\033[%sm>\033[0m", "32")
}

func RenderCurrentLine(s string, isCurrentLine bool) string {
	if !isCurrentLine {
		return s
	}

	return fmt.Sprintf("\033[%sm%s\033[0m", "32", s)
}

func RenderElementSelected(el string) string {
	return fmt.Sprintf("\033[%sm%s\033[0m", "38;2;214;112;214", el)
}
