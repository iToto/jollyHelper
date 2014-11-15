// Package model provides the interface for the DB for all objects
package models

import (
	"encoding/json"
	gouuid "github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2/bson"
)

// Struct2Map returns a a map version of a structure
func Struct2Map(i interface{}) bson.M {
	v, err := json.Marshal(i)
	if err != nil {
		return nil
	}
	var m bson.M
	json.Unmarshal(v, &m)

	return m
}

// NewUid creates a new v4 UUID
func NewUid() string {
	uuid, _ := gouuid.NewV4()
	return uuid.String()
}
