# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Changed
- Refactored the `internal/model/gemini.go` to use the `google.golang.org/genai` SDK.
- The `internal/model/client.go` now acts as a dispatcher, using the `genai` SDK for Gemini models and the `aiplatform` SDK for other models.

### Fixed
- Resolved all compilation errors that arose from the initial major refactoring of the model clients.
- Re-created missing request/response struct definitions (`AnthropicRequest`, `LlamaRequest`, etc.) in a new `internal/model/structs.go` file.
- Corrected the initialization of the `aiplatform.PredictionClient` in the `PaLMClient`, `AnthropicClient`, and `MetaClient` structs.
- Fixed an issue where `project` and `region` flags were not being correctly read from `GEN_` environment variables.
- Corrected the default model for the `prompt` command to `gemini-2.5-flash`.
- Refactored the `prompt` command to use `RunE` for proper error propagation, removing calls to `log.Fatal` and `os.Exit`.
