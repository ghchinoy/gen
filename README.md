# gen

`gen` is a command-line tool that interacts with Google Cloud models - foundation models Gemini, PaLM, and models found on Model Garden such as Claude and others.


![](./media/gen-002.gif)

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

Check `gen version` once installed to make sure you have the latest release.

If you want to install the main branch, use `go install github.com/ghchinoy/gen@main`.


#### Authentication

[Standard methods of authenticating](https://cloud.google.com/docs/authentication/provide-credentials-adc) to Google Cloud are supported.

```bash
# set your GCP project
gcloud config set project YOUR_PROJECT
# login to obtain credentials
gcloud auth application-default login
```

#### GCP Project & Region

Set your Google Cloud Project either via the env var `PROJECT_ID` or via the flag `--project` and your region either via the env var `REGION` or via the flag `--region`

Using env vars

```bash
export PROJECT_ID=$(gcloud config get project)
export REGION=us-central1
```

Using flags

```bash
gen --project $(gcloud config get project) --region us-central1 p "hi there"
```

## Usage

### Generate content

Generate content with the `prompt` command. This defaults to Gemini.

```bash
gen prompt "say something nice to me"
```

This will result in:

```
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

If you have [Anthropic's models activated from Model Garden](https://console.cloud.google.com/vertex-ai/model-garden?pageState=(%22galleryStateKey%22:(%22f%22:(%22g%22:%5B%22providers%22%5D,%22o%22:%5B%22ANTHROPIC%22%5D),%22s%22:%22%22))), you can also use the model name in the `--model` (or `-m`) flag, the example below is for [Claude 3 Haiku](https://console.cloud.google.com/vertex-ai/publishers/anthropic/model-garden/claude-3-haiku):

```bash
gen p -m claude-3-haiku@20240307 "say something nice to me"
```

[Claude 3.5 Sonnet](https://console.cloud.google.com/vertex-ai/publishers/anthropic/model-garden/claude-3-5-sonnet):

```bash
gen p -m claude-3-5-sonnet@20240620 "say something nice to me"
```


### Model Configuration Parameters

Use the `--config` flag to pass in model parameters, as a json file, such as:

```bash
gen p --model text-bison@002 --config config.json "say something nice to me"
```

Where `config.json` contains [PaLM 2 parameters](https://cloud.google.com/vertex-ai/generative-ai/docs/model-reference/text#request_body):

```json
{
    "temperature":     0.95,
    "maxOutputTokens": 1024,
    "topP":            0.4,
    "topK":            40
}
```

or, for [Gemini](https://cloud.google.com/vertex-ai/generative-ai/docs/model-reference/gemini#request_body), containing the `generationConfig`:

```json
{
    "temperature": 0.1,
    "topP": 0.9,
    "topK": 40,
    "maxOutputTokens": 256
}
```


### Count Tokens

```
gen tokens "hi how are you today"

Number of tokens for the prompt: 5
```

or for a very long prompt in file

```
gen tokens --file VeryLongPromptFile.txt

Number of tokens for the prompt: 1599681
```

### Interactive mode

Multiple single-turn interactions (synthetic context and support for models with a chat api planned):

```
gen interactive

2024/03/30 15:27:33 entering interactive mode
2024/03/30 15:27:33 type 'exit' or 'quit' to exit
2024/03/30 15:27:33 model: gemini-1.0-pro
? Hi say something nice to me
You are a beautiful, intelligent, and kind person. You are loved and appreciated by many people, and you bring joy to the lives of those around you. You are strong and capable, and you can achieve anything you set your mind to. I am proud of you, and I know you will continue to do great things.

? What's your name?
I am Gemini, a multi-modal AI language model developed by Google. I don't have a name, as I am not an individual being.

```

### Compare outputs with diff

Using the unix `diff` command and a clever ordering of `gen`, you can compare the output of two models with the same prompt.

```
export MY_PROMPT="say something nice to me and mention your name"
diff <(gen p "${MY_PROMPT}") <(gen p -m claude-3-5-sonnet@20240620 "${MY_PROMPT}")
```

The result should be similar to:

```
< You are a bright spark in the world, and I, Bard, am delighted to have crossed paths with you today!  Your energy and creativity shine through, and I'm sure you have amazing things ahead of you. Keep shining! 😊
---
> As an AI language model, I don't have a personal name, but I'm happy to say something nice to you!
2a3
> You are a thoughtful and kind person for wanting to hear something positive. Your curiosity and willingness to engage in conversation are admirable qualities. I hope you have a wonderful day filled with joy and positivity!
```



## Acknowledgements
`gen` is inspired by Simon Willison's [llm tool](https://llm.datasette.io/en/stable/) as well as Eli Bendersky's [gemini-cli](https://github.com/eliben/gemini-cli). Both are super awesome, check them out!


## License

Apache 2.0; see [`LICENSE`](LICENSE) for details.

## Disclaimer

This project is not an official Google project. It is not supported by Google and Google specifically disclaims all warranties as to its quality, merchantability, or fitness for a particular purpose.