package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/andreyvit/diff"
	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
	yaml "gopkg.in/yaml.v2"
)

type Checker struct {
	Interval   int64           `json:"interval" yaml:"interval"`
	StoreHash  bool            `json:"store_hash" yaml:"store_hash"`
	Debug      bool            `json:"debug" yaml:"debug"`
	Sites      map[string]Site `json:"sites" yaml:"sites"`
	ConfigPath string          `json:"-" yaml:"-"`
	Notifiers  []Notifier      `json:"-" yaml:"-"`
}

func NewChecker(config string, notifiers []Notifier) (Checker, error) {
	c := Checker{ConfigPath: config}

	data, err := ioutil.ReadFile(config)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal([]byte(data), &c)
	c.Notifiers = notifiers
	return c, err
}

func (c Checker) Run(singleRun bool) {
	log.Println("Started...")
	for {
		for i, s := range c.Sites {
			if c.Debug {
				log.Printf("Scanning '%s'\n", i)
			}
			err := s.Check(c.Notifiers)
			if err != nil {
				log.Println(err)
			}
			if c.Debug {
				log.Printf("%s => %s\n", c.Sites[i].Hash, s.Hash)
			}
			c.Sites[i] = s
		}

		if c.StoreHash {
			c.UpdateConfig()
		}

		if singleRun {
			break
		}

		time.Sleep(time.Duration(c.Interval) * time.Second)
	}
}

func (c Checker) UpdateConfig() error {
	out, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.ConfigPath, out, 0644)
	return err
}

type Site struct {
	URL       string `json:"url" yaml:"url"`
	Template  string `json:"template" yaml:"template"`
	Recipient string `json:"recipient" yaml:"recipient"`
	Selector  string `json:"selector" yaml:"selector"`
	Diff      bool   `json:"diff" yaml:"diff"`
	Hash      string `json:"hash" yaml:"hash"`
	Data      string `json:"data" yaml:"data"`
}

func (s *Site) Check(n []Notifier) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from error '%s' while checking %s\n", r, s.URL)
		}
	}()

	oldHash := s.Hash
	oldData := s.Data

	response, err := http.Get(s.URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if s.Selector == "" {
		data := string(content)
		if s.Diff {
			s.Data = data
		}
		s.Hash = GetMD5Hash(data)
	} else {
		sel, err := cascadia.Compile(s.Selector)
		if err != nil {
			return err
		}
		r := bytes.NewReader(content)
		node, err := html.Parse(r)
		if err != nil {
			return err
		}
		firstMatch := sel.MatchFirst(node)
		data := firstMatch.Data
		if s.Diff {
			s.Data = data
		}
		s.Hash = GetMD5Hash(data)
	}
	if oldHash != "" && oldHash != s.Hash {
		diff := diff.LineDiff(oldData, s.Data)
		for _, notifier := range n {
			err = notifier.Notify(s.Recipient, s.Template, diff)
			if err != nil {
				log.Printf("Error while notifing: %s\n", err.Error())
			}
		}
	}
	return nil
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
