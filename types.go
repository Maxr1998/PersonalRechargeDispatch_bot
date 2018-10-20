package main

import (
	"encoding/json"
	"golang.org/x/tools/container/intsets"
	"log"
)

// Config is a data structure to store the config
type Config struct {
	Channels map[int64]*Channel `json:"channels"`
	observer chan int
}

// Channel is a data structure that maps a dispatch channel to portals & users
type Channel struct {
	Users   *SparseArray     `json:"users"`
	Portals map[string][]int `json:"portals"`
}

// SparseArray is a wrapper for intsets.Sparse
type SparseArray struct {
	intsets.Sparse
}

// MarshalJSON encodes a SparseArray into json
func (s *SparseArray) MarshalJSON() ([]byte, error) {
	value := make([]int, 0)
	value = s.AppendTo(value)
	return json.Marshal(value)
}

// UnmarshalJSON decodes json into a SparseArray
func (s *SparseArray) UnmarshalJSON(data []byte) error {
	var values []int
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	for _, v := range values {
		s.Insert(v)
	}
	return nil
}

type appError struct {
	Message string
	Error   error
}

func (a *appError) log() {
	if a.Error == nil {
		log.Println(a.Message)
	}
	log.Println(a.Message + ": " + a.Error.Error())
}
