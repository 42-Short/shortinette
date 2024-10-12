package attributes

import (
	"fmt"
	"os"
	"strings"
)

// Lints for attributes in filePath.
//
// Args:
//   - filePath (string): Path which is to be checked.
//   - requiredAttributes (map[string]bool): Attributes which are required to be in the file. One of them not being
//     found will trigger an error.
//   - forbiddenAttributes (map[string]bool): Attributes which are forbidden in the file. One of them being found
//     will trigger an error.
//
// Errors caught are documented in err.Error() in the following format:
//	missing attributes: #![no_std], ..., ...
//	forbidden attributes used: #![allow(...)], ..., ...
func Check(filePath string, requiredAttributes map[string]bool, forbiddenAttributes map[string]bool) (err error) {
	contentAsBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(contentAsBytes), "\n")

	foundForbidden := []string{}

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if forbiddenAttributes[trimmedLine] {
			foundForbidden = append(foundForbidden, trimmedLine)
		} else if requiredAttributes[trimmedLine] {
			delete(requiredAttributes, trimmedLine)
		}
	}

	feedbackMessage := ""
	if len(requiredAttributes) != 0 {
		missingAttributesList := []string{}
		for attribute := range requiredAttributes {
			missingAttributesList = append(missingAttributesList, attribute)
		}
		feedbackMessage = fmt.Sprintf("missing attributes: %v\n", strings.Join(missingAttributesList, ","))
	}
	if len(foundForbidden) != 0 {
		feedbackMessage = fmt.Sprintf("%sforbidden attributes used: %v", feedbackMessage, strings.Join(foundForbidden, ","))
	}

	if len(feedbackMessage) != 0 {
		return fmt.Errorf(feedbackMessage)
	}
	return nil
}
