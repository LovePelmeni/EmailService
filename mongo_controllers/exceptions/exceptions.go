package mongo_exceptions 

import (
	"errors"
)
 
func SavedDocumentError() error {
	return errors.New("Failed to Save Document.")
}

func DeleteDocumentError() error {
	return errors.New("Failed to delete Document")
}