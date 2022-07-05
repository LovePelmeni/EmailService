package mongo_exceptions

import (
	"errors"
	"fmt"
)

func SavedDocumentError() error {
	return errors.New("Failed to Save Document.")
}

func DeleteDocumentError() error {
	return errors.New("Failed to delete Document")
}

func UpdateDocumentError() error {
	return errors.New("Failed to Update Document")
}

func EmtpyDocumentError() error {
	return errors.New("Emtpy Document. Operation Failed.")
}

func InvalidMongoClientError() error {
	return errors.New("Invalid MongoDB Client.")
}

func OperationFailed(Operation string, Reason ...error) error {
	return errors.New(fmt.Sprintf(
	"%s Operation Failed. Reason: %s", Operation, Reason))
}
