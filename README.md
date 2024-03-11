# Chat CLI

A cli tool for quickly querying common llm providers. Currently supported:

- [x] Claude/Anthropic
- [ ] GPT/OpenAI
- [ ] local

The goal is to support both synchronous and streaming calls to common LLM apis.
Has multiple modes
## Basic Mode
Uses basic `fmt.Print` functionality for writing to console, `readlines` to read

Pasting text with newlines will not work in this mode

## TUI Mode
Default mode, press enter to add newlines.
Use CTRL+S to submit message