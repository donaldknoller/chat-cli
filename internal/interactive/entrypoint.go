package interactive

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/donaldknoller/chat-cli/internal/common_llm"
	"github.com/spf13/cobra"
	"log"
)

var (
	t *tea.Program
)

func RunTea(cmd *cobra.Command, args []string) {
	t = tea.NewProgram(InitialModel(), tea.WithAltScreen())
	if _, err := t.Run(); err != nil {
		log.Fatal(err)
	}
}

//
//func RunStream(cmd *cobra.Command, args []string) {
//	rl, err := readline.New("> ")
//	if err != nil {
//		panic(err)
//	}
//	defer rl.Close()
//
//	c := make(chan os.Signal, 1)
//	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
//
//	go func() {
//		<-c
//		fmt.Println("\nExit signal received. Exiting...")
//		rl.Close()
//		os.Exit(0)
//	}()
//	client := anthropic.NewClient()
//	useStream := true
//	request := anthropic.Request{
//		Model:     viper.GetString(config.LLM_API_MODEL),
//		MaxTokens: 1024,
//		Messages:  []common_llm.ChatMessage{},
//		Stream:    &useStream,
//	}
//	for {
//		line, readErr := rl.Readline()
//		if readErr != nil { // EOF or Ctrl+D will exit the loop
//			break
//		}
//
//		if strings.TrimSpace(line) == "" {
//			continue
//		}
//		request.AddMessage(line, common_llm.User)
//		responseChan := make(chan anthropic.StreamResponse)
//		go client.PostDataStream(request, responseChan)
//		response := ""
//		for data := range responseChan {
//			if data.Err != nil {
//				fmt.Printf("exit due to %v", data.Err)
//				return
//			}
//			response = response + data.Chunk
//			fmt.Printf(data.Chunk)
//		}
//		fmt.Println()
//		request.AddMessage(response, 0)
//	}
//}

func streamMessage(session common_llm.Session) func() tea.Msg {
	return func() tea.Msg {
		responseChan := make(chan common_llm.StreamResponse)
		go session.PostDataStream(responseChan)
		for data := range responseChan {
			t.Send(data)
		}
		return nil
	}
}
