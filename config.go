package main

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

const configFile = "config.json"

// LoadConfig loads the dictionary from the disk
func LoadConfig() *Config {
	var config *Config
	data, err := ioutil.ReadFile(configFile)
	if err != nil || json.Unmarshal(data, &config) != nil {
		config = &Config{
			Channels: make(map[int64]*Channel),
		}
	}
	go saveTask(config)
	return config
}

func saveTask(config *Config) {
	var lastSave = time.Now().Add(-time.Minute)
	config.observer = make(chan int)
	for range config.observer {
		if sinceSave := time.Since(lastSave); sinceSave > 30*time.Second {
			lastSave = time.Now()
			if data, err := json.Marshal(config); err == nil {
				ioutil.WriteFile(configFile, data, 0644)
			}
		} else {
			go func() {
				time.Sleep(40*time.Second - sinceSave)
				config.observer <- 0
			}()
		}
	}
}

// AddChannel creates a new channel and appends it to the config
func (c *Config) AddChannel(id int64) {
	c.Channels[id] = &Channel{
		Users:   &SparseArray{},
		Portals: make(map[string][]int),
	}
	c.changed()
}

// GetChannel returns the specified channel or nil if missing
func (c *Config) GetChannel(id int64) *Channel {
	return c.Channels[id]
}

// FindChannelForUser …
func (c *Config) FindChannelForUser(userID int) *Channel {
	for _, ch := range c.Channels {
		if ch.HasUser(userID) {
			return ch
		}
	}
	return nil
}

// AddUser creates a new channel and appends it to the config
func (c *Config) AddUser(channelID int64, userID int) {
	if ch := c.Channels[channelID]; ch != nil {
		ch.Users.Insert(userID)
		c.changed()
	}
}

// HasUser returns true when the channel includes the user id
func (c *Channel) HasUser(userID int) bool {
	return c.Users.Has(userID)
}

// AddPortal adds a portal to the specified user in this channel
func (c *Channel) AddPortal(userID int, portalCode string) {
	c.Portals[portalCode] = append(c.Portals[portalCode], userID)
}

func (c *Config) changed() {
	c.observer <- 0
}
