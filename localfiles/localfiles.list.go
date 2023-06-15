// Package localfiles hosts the LocalFiles type.
package localfiles

import (
	"context"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/bertjohnson/logger"
	contexttype "github.com/bertjohnson/logger/types/context"
)

// ListSecrets lists secret paths.
func (l *LocalFiles) ListSecrets(ctx context.Context, pathChannel chan string, errorChannel chan error) {
	// Validate parameters.
	if ctx == nil {
		errorChannel <- errors.New("context is required")

		close(pathChannel)
		close(errorChannel)

		return
	}

	// Add to context.
	ctx = context.WithValue(ctx, contexttype.SecretProviderID, l.ID) // nolint

	// Read directory.
	uri := l.basePath
	fis, err := ioutil.ReadDir(uri)
	if err != nil {
		errorChannel <- err

		close(pathChannel)
		close(errorChannel)

		return
	}

	// Loop through directory.
	for _, fi := range fis {
		fileName := fi.Name()
		if !fi.IsDir() {
			lastDot := strings.LastIndex(fileName, ".")
			fileExtension := strings.ToLower(fileName[lastDot+1:])
			switch fileExtension {
			case "secret":
				pathChannel <- fileName[0 : len(fileName)-7]
			}
		}
	}

	// Log.
	logger.Verbose(ctx, "Listed secrets.")

	close(pathChannel)
	close(errorChannel)
}
