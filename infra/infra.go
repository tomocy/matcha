package infra

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
)

func createWorkspace() error {
	name := configFilename()
	if _, err := os.Stat(name); err == nil {
		return nil
	}

	dir := workspaceName()
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	f, err := os.Create(name)
	if err != nil {
		return err
	}

	return f.Close()
}

func loadConfig() (*config, error) {
	srcName := configFilename()
	src, err := os.Open(srcName)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	var loaded *config
	if err := readJSON(src, &loaded); err != nil {
		return nil, err
	}

	return loaded, nil
}

func saveConfig(config *config) error {
	destName := configFilename()
	dest, err := os.OpenFile(destName, os.O_WRONLY, 0700)
	if err != nil {
		return err
	}
	defer dest.Close()

	return writeJSON(dest, config)
}

type config struct {
	Reddit redditConfig
}

type redditConfig struct {
	AccessToken *oauth2.Token
}

func configFilename() string {
	return filepath.Join(workspaceName(), "config.json")
}

func workspaceName() string {
	return filepath.Join(os.Getenv("HOME"), ".matcha")
}

func readJSON(src io.Reader, dest interface{}) error {
	return json.NewDecoder(src).Decode(dest)
}

func writeJSON(dest io.Writer, src interface{}) error {
	return json.NewEncoder(dest).Encode(src)
}
