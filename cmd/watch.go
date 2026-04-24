package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/stackwatch/internal/drift"
)

var watchCmd = &cobra.Command{
	Use:   "watch [stack-name]",
	Short: "Watch a CloudFormation stack for drift",
	Args:  cobra.ExactArgs(1),
	RunE:  runWatch,
}

func init() {
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	stackName := args[0]

	detector, err := drift.NewCloudFormationDetector(region, profile)
	if err != nil {
		return fmt.Errorf("failed to initialize detector: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("\nShutting down...")
		cancel()
	}()

	fmt.Printf("Watching stack %q in region %q (interval: %ds)\n", stackName, region, interval)

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	if err := detector.Check(ctx, stackName); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := detector.Check(ctx, stackName); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
