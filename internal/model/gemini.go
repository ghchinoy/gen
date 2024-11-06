package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"cloud.google.com/go/vertexai/genai"
)

// UseGeminiModel calls Gemini's generate content method
func UseGeminiModel(ctx context.Context, modelName string, cfg Config, args []string) error {
	log.Printf("Gemini [%s]", modelName)

	var promptParts []genai.Part
	for _, arg := range args {
		if argLooksLikeGCSURL(arg) {
			part := genai.FileData{
				MIMEType: mime.TypeByExtension(filepath.Ext(arg)),
				FileURI:  arg,
			}
			promptParts = append(promptParts, part)

		} else if argLooksLikeURL(arg) {
			part, err := getPartFromURL(arg)
			if err != nil {
				return err
			}
			promptParts = append(promptParts, part)
		} else if argLooksLikeFilename(arg) {
			part, err := getPartFromFile(arg)
			if err != nil {
				return err
			}
			promptParts = append(promptParts, part)
		} else {

			promptParts = append(promptParts, genai.Text(arg))
		}
	}

	var buf bytes.Buffer
	if err := GenerateContentGemini(ctx, modelName, cfg, &buf, promptParts); err != nil {
		log.Printf("error generating content: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", buf.String())
	return nil
}

// GenerateContentGemini calls Gemini's generate content method
func GenerateContentGemini(ctx context.Context, modelName string, cfg Config, w io.Writer, parts []genai.Part) error {
	// TODO - There are differences between this function and the matching function in palm.go
	// due to when the config file contents are read.

	// TODO - Unlike matching functions in palm.go and anthropic.go, this one is public.  Should the
	// others be made public or should this one be made private.

	client, err := genai.NewClient(ctx, cfg.ProjectID, cfg.RegionID)
	if err != nil {
		return fmt.Errorf("error creating a client: %v", err)
	}
	defer client.Close()

	gemini := client.GenerativeModel(modelName)

	if cfg.ConfigFile != "" {
		modelConfig, err := os.ReadFile(cfg.ConfigFile)
		if err != nil {
			return fmt.Errorf("error reading model config file: %w", err)
		}
		var config genai.GenerationConfig
		err = json.Unmarshal(modelConfig, &config)
		if err != nil {
			return fmt.Errorf("error unmarshalling GenerationConfig from file: %w", err)
		}
		gemini.GenerationConfig = config
		if cfg.LogType != "none" {
			log.Printf("config: %v", config)
		}
	}

	resp, err := gemini.GenerateContent(ctx, parts...)
	if err != nil {
		// needs more sensible parsing of error message
		if strings.Contains(err.Error(), "lookup -aiplatform.googleapis.com:") {
			log.Print("missing REGION")
		}
		if strings.Contains(err.Error(), "RESOURCE_PROJECT_INVALID") {
			log.Print("missing PROJECT_ID")
		}
		return fmt.Errorf("error generating content: %w", err)
	}

	if cfg.OutputType == "json" {
		rb, _ := json.MarshalIndent(resp, "", "  ")
		fmt.Fprintln(w, string(rb))
	} else {
		if len(resp.Candidates) > 0 {
			var all []string
			for _, v := range resp.Candidates[0].Content.Parts {
				all = append(all, fmt.Sprintf("%s", v))
			}
			fmt.Fprintf(w, "%s", strings.Join(all, " "))
		} else {
			log.Printf("Candidate length 0")
		}
	}
	return nil
}

// thanks to eilben's https://github.com/eliben/gemini-cli/blob/main/internal/commands/prompt.go

// argLooksLikeFilename says if command-line argument looks like a filename,
// which we consider to have an alphabetical extension following a dot separator,
// but not look like a URL.
func argLooksLikeFilename(arg string) bool {
	re := regexp.MustCompile(`\.[a-zA-Z]+$`)
	return re.MatchString(arg) && strings.Index(arg, "://") < 0
}

func argLooksLikeGCSURL(arg string) bool {
	return strings.HasPrefix(arg, "gs://")
}

func argLooksLikeURL(arg string) bool {
	_, err := url.ParseRequestURI(arg)
	return err == nil
}

func getPartFromFile(path string) (genai.Part, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	mimeType := mime.TypeByExtension(filepath.Ext(path))
	log.Printf("%s bytes read: %d", mimeType, len(b))

	ext := filepath.Ext(path)
	switch strings.TrimSpace(ext) {
	case ".jpg", ".jpeg":
		return genai.ImageData("jpeg", b), nil
	case ".png":
		return genai.ImageData("png", b), nil
	case ".gif":
		return genai.ImageData("gif", b), nil
	case ".webp":
		return genai.ImageData("gif", b), nil
	case ".pdf", ".wav", ".mp3", ".mpeg", ".mov", ".mp4", ".avi", ".mpg", ".wmv", ".mpegs", ".flv":
		return genai.Blob{
			MIMEType: mimeType,
			Data:     b,
		}, nil
	default:
		return nil, fmt.Errorf("invalid file extension: %s", ext)
	}
}

func getPartFromURL(url string) (genai.Part, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image from url: %w", err)
	}
	defer resp.Body.Close()

	urlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image bytes: %w", err)
	}

	mimeType := resp.Header.Get("Content-Type")
	parts := strings.Split(mimeType, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid mime type %v", mimeType)
	}

	return genai.ImageData(parts[1], urlData), nil
}
