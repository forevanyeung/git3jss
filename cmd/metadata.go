package main

import (
	"fmt"
	"regexp"
	"strings"
)

// download
func WriteMetadata(jssScript JssScript) string {
	var b strings.Builder
	b.WriteString("\n\n### git3jss\n#\n")
	fmt.Fprintf(&b, "#%16s:  %s\n", "Name", jssScript.Name)
	for i, line := range strings.Split(jssScript.Info, "\n") {
		if i == 0 {
			fmt.Fprintf(&b, "#%16s:  %s\n", "Info", line)
		} else {
			fmt.Fprintf(&b, "#                   %0s\n", line)
		}
	}
	for i, line := range strings.Split(jssScript.Notes, "\n") {
		if i == 0 {
			fmt.Fprintf(&b, "#%16s:  %s\n", "Notes", line)
		} else {
			fmt.Fprintf(&b, "#                   %0s\n", line)
		}
	}
	fmt.Fprintf(&b, "#%16s:  %s\n", "Priority", jssScript.Priority)
	fmt.Fprintf(&b, "#%16s:  %s\n", "Requirements", jssScript.OSRequirements)
	fmt.Fprintf(&b, "#%16s:  %s\n", "Category", jssScript.CategoryName)
	if jssScript.Parameter4 != "" {
		fmt.Fprintf(&b, "#%16s:  %s\n", "Parameter4", jssScript.Parameter4)
	}
	if jssScript.Parameter5 != "" {
		fmt.Fprintf(&b, "#%16s:  %s\n", "Parameter5", jssScript.Parameter5)
	}
	if jssScript.Parameter6 != "" {
		fmt.Fprintf(&b, "#%16s:  %s\n", "Parameter6", jssScript.Parameter6)
	}
	if jssScript.Parameter7 != "" {
		fmt.Fprintf(&b, "#%16s:  %s\n", "Parameter7", jssScript.Parameter7)
	}
	if jssScript.Parameter8 != "" {
		fmt.Fprintf(&b, "#%16s:  %s\n", "Parameter8", jssScript.Parameter8)
	}
	if jssScript.Parameter9 != "" {
		fmt.Fprintf(&b, "#%16s:  %s\n", "Parameter9", jssScript.Parameter9)
	}
	if jssScript.Parameter10 != "" {
		fmt.Fprintf(&b, "#%16s:  %s\n", "Parameter10", jssScript.Parameter10)
	}
	if jssScript.Parameter11 != "" {
		fmt.Fprintf(&b, "#%16s:  %s\n", "Parameter11", jssScript.Parameter11)
	}
	b.WriteString("#\n###\n\n")

	return b.String()
}

// sync
func ReadMetadata(content string) JssScript {
	linebreak := regexp.MustCompile(`\r?\n`)
	lines := linebreak.Split(content, -1)

	// determine script language
	// interpreter := lines[0]

	metaFound := false
	var jssScript JssScript
	prevField := ""
	for _, l := range lines {
		if metaFound {
			// continue processing as long as there is still comments
			if !strings.HasPrefix(l, "#") {
				break
			}

			// stop processing at end of metadata
			if strings.HasPrefix(l, "###") {
				break
			}

			field, value, ok := strings.Cut(l, ":")
			if ok {
				commentIndex := strings.LastIndex(field, "#") + 1 //FIXME: # Name: script #
				field = strings.TrimSpace(field[commentIndex:])
				prevField = field
				value = strings.TrimSpace(value)

				switch field {
				case "Name":
					jssScript.Name = value
				case "Info":
					jssScript.Info = value
				case "Notes":
					jssScript.Notes = value
				case "Priority":
					jssScript.Priority = value
				case "Requirements":
					jssScript.OSRequirements = value
				case "Category":
					jssScript.CategoryName = value
				case "Parameter4":
					jssScript.Parameter4 = value
				case "Parameter5":
					jssScript.Parameter5 = value
				case "Parameter6":
					jssScript.Parameter6 = value
				case "Parameter7":
					jssScript.Parameter7 = value
				case "Parameter8":
					jssScript.Parameter8 = value
				case "Parameter9":
					jssScript.Parameter9 = value
				case "Parameter10":
					jssScript.Parameter10 = value
				case "Parameter11":
					jssScript.Parameter11 = value
				}
			} else {
				// allow multiline for Info and Notes only
				if prevField == "Info" && field != "" {
					commentIndex := strings.LastIndex(field, "#") + 1
					field = strings.TrimSpace(field[commentIndex:])
					jssScript.Info = jssScript.Info + "\n" + field
				}
				if prevField == "Notes" && field != "" {
					commentIndex := strings.LastIndex(field, "#") + 1
					field = strings.TrimSpace(field[commentIndex:])
					jssScript.Notes = jssScript.Notes + "\n" + field
				}
			}
		}

		// find start of git3jss metadata
		if strings.HasPrefix(l, "### git3jss") {
			metaFound = true
		}
	}

	return jssScript
}
