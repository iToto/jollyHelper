// Package oauth2_mgostore provides an interface for a Mongodb Database
package oauth2_mgostore

import (
	"github.com/RangelReale/osin"
)

// MongoClient structure that holds Mongo Client information
// UID, Secret, RedirectUri, UserData
type MongoClient struct {
	Uid         string      `bson:"uid" json:"uid"`
	Secret      string      `bson:"secret" json:"secret"`
	RedirectUri string      `bson:"redirect_uri" json:"redirect_uri"`
	UserData    interface{} `bson:"user_data" json:"user_data"`
}

// GetId getter method for the Uid property
func (d *MongoClient) GetId() string {
	return d.Uid
}

// GetSecret getter method for the Secret property
func (d *MongoClient) GetSecret() string {
	return d.Secret
}

// GetRedirectUri getter method for the RedirectUri property
func (d *MongoClient) GetRedirectUri() string {
	return d.RedirectUri
}

// GetUserData getter method for the UserData property
func (d *MongoClient) GetUserData() interface{} {
	return d.UserData
}

// CopyFrom copies a MongoClient values from one instance to another
func (d *MongoClient) CopyFrom(client osin.Client) {
	d.Uid = client.GetId()
	d.Secret = client.GetSecret()
	d.RedirectUri = client.GetRedirectUri()
	d.UserData = client.GetUserData()
}
