# gen

`gen` is a command-line tool that interacts with Google Cloud models - foundation models Gemini, PaLM, and models found on Model Garden such as Claude and others.

## Usage

```bash
gen prompt "say something nice to me"
```

To get the list of commands, use `gen` by itself.

### Generate content

```
export PROJECT_ID=$(gcloud config get project)
export LOCATION=us-central1

gen --project $PROJECT_ID --region $LOCATION prompt "say something nice to me"
```

### Count Tokens

```
export PROJECT_ID=$(gcloud config get project)
export LOCATION=us-central1

gen --project $PROJECT_ID --region $LOCATION tokens "hi how are you today"
```

or for a very long prompt


```
export PROJECT_ID=$(gcloud config get project)
export LOCATION=us-central1

gen --project $PROJECT_ID --region $LOCATION tokens --file VeryLongPromptFile.txt
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
