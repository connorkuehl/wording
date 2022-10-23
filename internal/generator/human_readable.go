package generator

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type RandomWordClient struct {
	baseURL  string
	client   http.Client
	fallback interface{ NewToken() string }
}

func NewRandomWordClient(baseURL string, fallback interface{ NewToken() string }) *RandomWordClient {
	return &RandomWordClient{
		baseURL:  baseURL,
		fallback: fallback,
	}
}

func (c *RandomWordClient) NewToken() string {
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
		log.WithError(err).Warn("word generator returned error, falling back")
		return c.fallback.NewToken()
	}

	return strings.Join([]string{adjective, animal, noun}, "-")
}

func (c *RandomWordClient) Noun() (string, error) {
	return c.getWord("/random/noun")
}

func (c *RandomWordClient) Adjective() (string, error) {
	return c.getWord("/random/adjective")
}

func (c *RandomWordClient) Animal() (string, error) {
	return c.getWord("/random/animal")
}

func (c *RandomWordClient) getWord(endpoint string) (string, error) {
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
