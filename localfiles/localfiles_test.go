package localfiles

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/bertjohnson/logger"
	secretprovidertype "github.com/bertjohnson/secretprovider/types"
	"github.com/bertjohnson/startup"
	utilio "github.com/bertjohnson/util/io"
)

var (
	// Context.
	ctx context.Context

	// Local files client.
	localFilesClient *LocalFiles
)

// TestMain runs tests.
func TestMain(m *testing.M) {
	// Declare that the configuration is ready.
	err := startup.Ready()
	if err != nil {
		log.Fatalln("Error loading configuration values: " + err.Error())
	}

	// Wait for logger.
	ctx = context.Background()
	logger.Wait(ctx)

	// Create Local Files client.
	currentDirectory, err := os.Getwd()
	if err != nil {
		logger.Fatal(ctx, err.Error())
	}
	localFilesClient, err = New(ctx, &secretprovidertype.SecretProvider{
		ID:   "ea47bc05-9549-45df-b332-eb46d2319239",
		Type: "LocalFiles",
		URI:  currentDirectory + pathSeparator + "test",
	})
	if err != nil {
		logger.Fatal(ctx, "Error creating Vault client: "+err.Error())
	}

	// Run tests.
	retCode := m.Run()

	// Clean up.
	utilio.RemoveRecursive(currentDirectory + "/test")

	// Exit.
	os.Exit(retCode)
}
