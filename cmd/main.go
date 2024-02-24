package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/donaldknoller/chat-cli/internal/anthropic"
	anthropicClient "github.com/donaldknoller/chat-cli/internal/anthropic"
	"github.com/donaldknoller/chat-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "A simple CLI application",
	//Run:   run,
	Run: runStream,
}

func init() {
	rootCmd.Flags().String("key", "", "API Key for LLM")
	rootCmd.Flags().String("host", "", "API Host for LLM")
	viper.BindPFlag(config.LLM_API_KEY, rootCmd.Flags().Lookup("key"))
	viper.BindPFlag(config.LLM_API_HOST, rootCmd.Flags().Lookup("host"))
	viper.AutomaticEnv()
	config.InitDefault()
}

func run(cmd *cobra.Command, args []string) {

	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nExit signal received. Exiting...")
		rl.Close()
		os.Exit(0)
	}()
	client := anthropicClient.NewClient()
	useStream := true
	request := anthropic.Request{
		Model:     viper.GetString(config.LLM_API_MODEL),
		MaxTokens: 1024,
		Messages:  []anthropic.Message{},
		Stream:    &useStream,
	}
	for {
		line, readErr := rl.Readline()
		if readErr != nil { // EOF or Ctrl+D will exit the loop
			break
		}

		if strings.TrimSpace(line) == "" {
			continue
		}
		request.InsertMessage(line)
		response, clientErr := client.PostData(request)
		if clientErr != nil {
			fmt.Printf("exit due to %v", clientErr)
			break
		}
		fmt.Println("\n" + response)
		request.Merge(response)
	}
}

func runStream(cmd *cobra.Command, args []string) {

	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nExit signal received. Exiting...")
		rl.Close()
		os.Exit(0)
	}()
	client := anthropicClient.NewClient()
	useStream := true
	request := anthropic.Request{
		Model:     viper.GetString(config.LLM_API_MODEL),
		MaxTokens: 1024,
		Messages:  []anthropic.Message{},
		Stream:    &useStream,
	}
	for {
		line, readErr := rl.Readline()
		if readErr != nil { // EOF or Ctrl+D will exit the loop
			break
		}

		if strings.TrimSpace(line) == "" {
			continue
		}
		request.InsertMessage(line)
		responseChan := make(chan anthropicClient.StreamResponse)
		go client.PostDataStream(request, responseChan)
		response := ""
		for data := range responseChan {
			if data.Err != nil {
				fmt.Printf("exit due to %v", data.Err)
				return
			}
			response = response + data.Chunk
			fmt.Printf(data.Chunk)
		}
		fmt.Println()
		request.Merge(response)
	}
}
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
