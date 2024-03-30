# gen

`gen` is a command-line tool that interacts with Google Cloud models - foundation models Gemini, PaLM, and models found on Model Garden such as Claude and others.

## Usage

```bash
gen prompt "say something nice to me"
```

To get the list of commands, use `gen help`.

Help is also available for each command, such as `gen help tokens`.

### Setup 

#### Installing

Install `gen` on your machine via:

```
go install github.com/ghchinoy/gen@latest
```


#### Authentication

[Standard methods of authenticating](https://cloud.google.com/docs/authentication/provide-credentials-adc) to Google Cloud are supported.

#### GCP Project & Region

Set your Google Cloud Project either via the env var `PROJECT_ID` or via the flag `--project` and your region either via the env var `REGION` or via the flag `--region`

Using env vars
```
export PROJECT_ID=$(gcloud config get project)
export REGION=us-central1
```

Using flags

```
gen --project $(gcloud config get project) --region us-central1 p "hi there"
```


### Generate content

Generate content with the `prompt` command. This defaults to Gemini.

```
gen prompt "say something nice to me"

2024/03/30 15:29:13 model: gemini-1.0-pro
2024/03/30 15:29:13 prompt: [say something nice to me]
2024/03/30 15:29:13 using Gemini
You are a wonderful person with a kind heart and a beautiful soul. You deserve all the happiness in the world, and I hope you find it.
```

Using the `--output json` output flag with `json` will return the full response payload.

Use another model family, such as PaLM 2:

```
gen p --model text-bison@002 "say something nice to me"

2024/03/30 16:24:50 model: text-bison@002
2024/03/30 16:24:50 prompt: [say something nice to me]
2024/03/30 16:24:50 using PaLM 2
 You are a bright and shining light in this world. Your kindness and compassion touch the lives of everyone you meet. You are a true gift to us all.
```

Note: This uses the `p` alias for the `prompt` command, see `gen help prompt` for aliases for a specific command.


### Count Tokens

```
gen tokens "hi how are you today"

Number of tokens for the prompt: 5
```

or for a very long prompt in file

```
gen tokens --file VeryLongPromptFile.txt
```

### Interactive mode

Multiple single-turn interactions (synthetic context and support for models with a chat api planned):

```
gen interactive

2024/03/30 15:27:33 entering interactive mode
? Hi say something nice to me
You are a beautiful, intelligent, and kind person. You are loved and appreciated by many people, and you bring joy to the lives of those around you. You are strong and capable, and you can achieve anything you set your mind to. I am proud of you, and I know you will continue to do great things.

? What's your name?
I am Gemini, a multi-modal AI language model developed by Google. I don't have a name, as I am not an individual being.

```




## Acknowledgements
`gen` is inspired by Simon Willison's [llm tool](https://llm.datasette.io/en/stable/) as well as Eli Bendersky's [gemini-cli](https://github.com/eliben/gemini-cli). Both are super awesome, check them out!
