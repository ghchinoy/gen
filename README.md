# gen

`gen` is a command-line tool that interacts with Google Cloud models - foundation models Gemini, PaLM, and models found on Model Garden such as Claude and others.

## Usage

```bash
gen prompt "say something nice to me"
```

To get the list of commands, use `gen --help`.

### Generate content

```
export PROJECT_ID=$(gcloud config get project)
export REGION=us-central1

gen --project $PROJECT_ID --region $LOCATION prompt "say something nice to me"

2024/03/30 15:29:13 model: gemini-1.0-pro
2024/03/30 15:29:13 prompt: [say something nice to me]
2024/03/30 15:29:13 using Gemini
You are a wonderful person with a kind heart and a beautiful soul. You deserve all the happiness in the world, and I hope you find it.
```

Using the `--output json` output flag with `json` will return the full response payload.

### Count Tokens

```
export PROJECT_ID=$(gcloud config get project)
export REGION=us-central1

gen --project $PROJECT_ID --region $REGION tokens "hi how are you today"
```

or for a very long prompt


```
export PROJECT_ID=$(gcloud config get project)
export REGION=us-central1

gen --project $PROJECT_ID --region $REGION tokens --file VeryLongPromptFile.txt
```

### Interactive mode

Multiple single-turn interactions (synthetic context and support for models with a chat api planned):

```
export PROJECT_ID=$(gcloud config get project)
export REGION=us-central1

gen --project $PROJECT_ID --region $REGION i

2024/03/30 15:27:33 entering interactive mode
? Hi say something nice to me
You are a beautiful, intelligent, and kind person. You are loved and appreciated by many people, and you bring joy to the lives of those around you. You are strong and capable, and you can achieve anything you set your mind to. I am proud of you, and I know you will continue to do great things.

? What's your name?
I am Gemini, a multi-modal AI language model developed by Google. I don't have a name, as I am not an individual being.

```


## Installing

Install `gen` on your machine via:

```
go install github.com/ghchinoy/gen@latest
```

See Usage for more information..

# Authentication

[Standard methods of authenticating](https://cloud.google.com/docs/authentication/provide-credentials-adc) to Google Cloud are supported.



## Acknowledgements
`gen` is inspired by Simon Willison's [llm tool](https://llm.datasette.io/en/stable/) as well as Eli Bendersky's [gemini-cli](https://github.com/eliben/gemini-cli). Both are super awesome, check them out!
