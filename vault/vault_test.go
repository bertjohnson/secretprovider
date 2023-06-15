package vault

import (
	"bufio"
	"context"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/bertjohnson/logger"
	"github.com/bertjohnson/logger/types/env"
	secretprovidertype "github.com/bertjohnson/secretprovider/types"
	"github.com/bertjohnson/startup"
)

var (
	// Context.
	ctx context.Context

	// Vault client.
	vaultClient *Vault

	// Local Vault development server process.
	vaultCommand *exec.Cmd
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

	// Run Vault.
	vaultCommand = Start(ctx)

	// Create Vault client.
	vaultClient, err = New(ctx, &secretprovidertype.SecretProvider{
		ID:   "a70cd531-9c59-47af-b755-db5a5de6e395",
		Type: "Vault",
	})
	if err != nil {
		logger.Fatal(ctx, "Error creating Vault client: "+err.Error())
	}

	// Run tests.
	retCode := m.Run()

	// Kill local Vault development server after tests have run.
	if vaultCommand.Process != nil {
		vaultCommand.Process.Kill()
	}

	// Exit.
	os.Exit(retCode)
}

// Start starts Vault.
func Start(ctx context.Context) *exec.Cmd {
	logger.Wait(ctx)
	logger.Info(ctx, "Starting Vault container.")
	exec.Command("docker", "kill", "vault").Run()
	exec.Command("docker", "pull", "vault:latest").Run()
	vaultCommand := exec.Command("docker", "run", "-p", "8200:8200", "--name=vault", "--cap-add=IPC_LOCK", "--rm", "--stop-timeout=30", "hashicorp/vault:latest")
	stdout, err := vaultCommand.StdoutPipe()
	if err != nil {
		logger.Fatal(ctx, "Error reading standard err: "+err.Error())
	}
	if err != nil {
		logger.Fatal(ctx, "Error starting local Vault dev server: "+err.Error())
	}
	unsealKeyChannel := make(chan string)
	masterTokenChannel := make(chan string)
	go func() {
		scanner := bufio.NewScanner(stdout)
		err = vaultCommand.Start()
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "Error initializing listener") {
				logger.Fatal(ctx, line)
			}

			unsealKeyPos := strings.Index(line, "Unseal Key: ")
			if unsealKeyPos > -1 && len(line) >= (unsealKeyPos+56) {
				unsealKeyChannel <- line[unsealKeyPos+12 : unsealKeyPos+56]
			}

			rootTokenPos := strings.Index(line, "Root Token: ")
			if rootTokenPos > -1 && len(line) >= (rootTokenPos+40) {
				masterTokenChannel <- line[rootTokenPos+12 : rootTokenPos+40]
			}
		}
	}()
	// Wait for unseal key and worker token.
	unsealKey := <-unsealKeyChannel
	masterToken := <-masterTokenChannel

	// Register settings.
	if os.Getenv(env.VaultAddr) == "" {
		os.Setenv(env.VaultAddr, "http://localhost:8200")
	}
	if os.Getenv(env.VaultToken) == "" {
		os.Setenv(env.VaultToken, masterToken)
	}
	if os.Getenv(env.VaultUnsealShards) == "" {
		os.Setenv(env.VaultUnsealShards, unsealKey)
	}

	logger.Info(ctx, "Started Vault container.")

	return vaultCommand
}
