// Package randword provides a client to the random word generator API.
package randword

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

// Client is an HTTP client for a random word generator API.
type Client struct {
	baseURL string
	client  http.Client
}

// NewClient constructs a new client that is configured to use the
// given endpoint specified by baseURL.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
	}
}

// HumanReadableSlug returns a randomly generated URL slug comprised of
// hyphen-separated words.
func (c *Client) HumanReadableSlug() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tasks, _ := errgroup.WithContext(ctx)

	var (
		adjective string
		noun      string
		animal    string
	)

	tasks.Go(func() error {
		var err error
		adjective, err = c.Adjective()
		return err
	})
	tasks.Go(func() error {
		var err error
		noun, err = c.Noun()
		return err
	})
	tasks.Go(func() error {
		var err error
		animal, err = c.Animal()
		return err
	})
	err := tasks.Wait()
	if err != nil {
		return "", err
	}

	return strings.Join([]string{adjective, animal, noun}, "-"), nil
}

// Noun requests a random noun.
func (c *Client) Noun() (string, error) {
	return c.getWord("/random/noun")
}

// Adjective requests a random adjective.
func (c *Client) Adjective() (string, error) {
	return c.getWord("/random/adjective")
}

// Animal requests a random animal name.
func (c *Client) Animal() (string, error) {
	return c.getWord("/random/animal")
}

func (c *Client) getWord(endpoint string) (string, error) {
	url := c.baseURL + endpoint

	rsp, err := c.client.Get(url)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	var words []string
	err = json.NewDecoder(rsp.Body).Decode(&words)
	if err != nil {
		return "", err
	}

	if len(words) == 0 {
		return "", errors.New("service returned empty list")
	}

	if len(words[0]) == 0 {
		return "", errors.New("service returned empty string")
	}

	// Collapse the output down into one word, just in case.
	var s strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(words[0]))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		s.WriteString(scanner.Text())
	}

	return s.String(), nil
}
