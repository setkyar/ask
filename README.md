# ask

**Ask** will make it easer to chat with OpenAI [GPT Model](https://platform.openai.com/docs/models) or [Anthropic](https://docs.anthropic.com/en/api/messages) via your terminal.

## Installation

1. Go to [release page](https://github.com/setkyar/ask/releases).
2. Download the binary file and made it executable
3. Move binary file `ask` to `/usr/local/bin/ask`
4. Start accessing `ask` via cli.
5. Get your OpenAI API key via [OpenAI](https://platform.openai.com/account/api-keys) and submit it to `ask`
6. Get your Anthropic API key via [Anthropic console](https://console.anthropic.com/settings/keys) and submit it to `ask`


Your API key and configuration will be store at `$HOME/.ask_ai_settings.yaml`.

### Usage 

You can run `$ ask` via terminal. It will ask you to fill up the API keys and ask you to choose the default model. 
After that, you can start chatting with your selected default AI model. If you want to run specific provider, you can run with the following command

```
$ ask -p claude # chat with claude
$ ask -p openai # chat with openai
```

If you don't remember which commands you can use. Just run 

```
$ ask --help
Ask is a CLI application that allows you to interact with AI providers like Claude and OpenAI.
It provides various commands to set up and use these AI services.

Usage:
  ask [flags]

Flags:
  -h, --help              help for ask
  -p, --provider string   Specify AI provider (claude or openai)
      --update-config     Update the configuration settings
```