// Package localfiles hosts the LocalFiles type.
package localfiles

import (
	"context"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/bertjohnson/logger"
	contexttype "github.com/bertjohnson/logger/types/context"
	secretprovidertype "github.com/bertjohnson/secretprovider/types"
	json "github.com/json-iterator/go"
)

// ReadAllSecrets reads all secrets.
func (l *LocalFiles) ReadAllSecrets(ctx context.Context, secretChannel chan *secretprovidertype.Secret, errorChannel chan error) {
	// Validate parameters.
	if ctx == nil {
		errorChannel <- errors.New("context is required")

		close(secretChannel)
		close(errorChannel)

		return
	}

	// Add to context.
	if objectIDs, ok := ctx.Value(contexttype.ObjectIDs).(string); ok {
		ctx = context.WithValue(ctx, contexttype.ObjectIDs, objectIDs+"&secretproviderid="+l.ID) // nolint
	} else {
		ctx = context.WithValue(ctx, contexttype.ObjectIDs, "secretproviderid="+l.ID) // nolint
	}

	// Read directory.
	uri := l.basePath
	fis, err := ioutil.ReadDir(uri)
	if err != nil {
		errorChannel <- err

		close(secretChannel)
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
				// Read and deserialize file.
				data, err := ioutil.ReadFile(uri + pathSeparator + fileName)
				if err != nil {
					errorChannel <- err

					close(secretChannel)
					close(errorChannel)

					return
				}
				dataMap := make(map[string]interface{})
				err = json.Unmarshal(data, &dataMap)
				if err != nil {
					errorChannel <- err

					close(secretChannel)
					close(errorChannel)

					return
				}

				// Return secret.
				secret := secretprovidertype.Secret{
					Data: dataMap,
					Path: fileName[0 : len(fileName)-7],
				}

				secretChannel <- &secret
			}
		}
	}

	// Log.
	logger.Verbose(ctx, "Read all secrets.")

	close(secretChannel)
	close(errorChannel)
}
