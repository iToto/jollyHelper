// Package oauth2_mgostore provides an interface for a Mongodb Database
package oauth2_mgostore

import (
	"github.com/RangelReale/osin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

// collection names for the entities
const (
	CLIENT_COL    = "oauth2_clients"
	AUTHORIZE_COL = "oauth2_authorizations"
	ACCESS_COL    = "oauth2_accesses"
)

const REFRESHTOKEN = "refreshtoken"

// MongoStorage structure that holds the dbName and mongo session
type MongoStorage struct {
	dbName  string
	session *mgo.Session
}

// New creates a new Mon Session
func New(session *mgo.Session, dbName string) *MongoStorage {
	r := &MongoStorage{
		dbName:  dbName,
		session: session,
	}

	index := mgo.Index{
		Key:        []string{REFRESHTOKEN},
		Unique:     false, // refreshtoken is sometimes empty
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}
	accesses := r.session.DB(dbName).C(ACCESS_COL)
	err := accesses.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	return r
}

// Clone returns the Storage that is passed in
// (not sure what this is used for)
func (s *MongoStorage) Clone() osin.Storage {
	return s
}

// Close not yet implemented...
func (s *MongoStorage) Close() {
}

// GetClient gets the MongoClient object for a given MongoStorage
func (s *MongoStorage) GetClient(id string) (osin.Client, error) {
	log.Printf("GetClient: %s\n", id)
	session := s.session.Copy()
	defer session.Close()
	clients := session.DB(s.dbName).C(CLIENT_COL)

	var client osin.Client
	client = &MongoClient{}
	err := clients.Find(bson.M{"uid": id}).One(client)

	return client, err
}

// SetClient setter method for a MongoClient on a MongoStorage
func (s *MongoStorage) SetClient(id string, client osin.Client) error {
	log.Printf("SetClient: %s\n", id)
	session := s.session.Copy()
	defer session.Close()
	clients := session.DB(s.dbName).C(CLIENT_COL)
	_, err := clients.Upsert(bson.M{"uid": id}, client)
	return err
	//s.clients[id] = client
	//return nil
}

// SaveAuthorize updates an authorization collection
func (s *MongoStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	log.Printf("SaveAuthorize: %s\n", data.Code)
	session := s.session.Copy()
	defer session.Close()
	authorizations := session.DB(s.dbName).C(AUTHORIZE_COL)
	_, err := authorizations.UpsertId(data.Code, data)
	return err
}

// LoadAuthorize loads the authorizeData from the db
func (s *MongoStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	log.Printf("LoadAuthorize: %s\n", code)
	session := s.session.Copy()
	defer session.Close()
	authorizations := session.DB(s.dbName).C(AUTHORIZE_COL)
	authData := new(osin.AuthorizeData)

	var rstl map[string]interface{}
	err := authorizations.FindId(code).One(&rstl)
	//becessary because osin break change when convert osin.Client
	// to interface and not yet to struct
	if rstl != nil {
		var client osin.Client
		rawClient := rstl["client"].(map[string]interface{})
		client = &MongoClient{
			Uid:         rawClient["uid"].(string),
			Secret:      rawClient["secret"].(string),
			RedirectUri: rawClient["redirect_uri"].(string),
			UserData:    rawClient["user_data"],
		}

		authData.Client = client
		authData.Code = rstl["code"].(string)
		authData.ExpiresIn = int32(rstl["expiresin"].(int))
		authData.Scope = rstl["scope"].(string)
		authData.RedirectUri = rstl["redirecturi"].(string)
		authData.State = rstl["state"].(string)
		authData.CreatedAt = rstl["createdat"].(time.Time)
		authData.UserData = rstl["userdata"]
	}
	//--
	return authData, err
}

// RemoveAuthorize removes the authorizeData from a db
func (s *MongoStorage) RemoveAuthorize(code string) error {
	log.Printf("RemoveAuthorize: %s\n", code)
	session := s.session.Copy()
	defer session.Close()
	authorizations := session.DB(s.dbName).C(AUTHORIZE_COL)
	return authorizations.RemoveId(code)

	//delete(s.authorize, code)
	//return nil
}

// SaveAccess saves an access token to the access collection
func (s *MongoStorage) SaveAccess(data *osin.AccessData) error {
	log.Printf("SaveAccess: %s\n", data.AccessToken)
	session := s.session.Copy()
	defer session.Close()
	accesses := session.DB(s.dbName).C(ACCESS_COL)
	_, err := accesses.UpsertId(data.AccessToken, data)
	return err

	//s.access[data.AccessToken] = data
	//if data.RefreshToken != "" {
	//	s.refresh[data.RefreshToken] = data.AccessToken
	//}
	//return nil
}

// LoadAccess loads an access token
func (s *MongoStorage) LoadAccess(code string) (*osin.AccessData, error) {
	log.Printf("LoadAccess: %s\n", code)
	session := s.session.Copy()
	defer session.Close()
	accesses := session.DB(s.dbName).C(ACCESS_COL)
	accData := new(osin.AccessData)
	err := accesses.FindId(code).One(accData)
	return accData, err

	//if d, ok := s.access[code]; ok {
	//	return d, nil
	//}
	//return nil, errors.New("Access not found")
}

// RemoveAccess removes an access token
func (s *MongoStorage) RemoveAccess(code string) error {
	log.Printf("RemoveAccess: %s\n", code)
	session := s.session.Copy()
	defer session.Close()
	accesses := session.DB(s.dbName).C(ACCESS_COL)
	return accesses.RemoveId(code)

	//delete(s.access, code)
	//return nil
}

// LoadRefresh loads the refresh token
func (s *MongoStorage) LoadRefresh(code string) (*osin.AccessData, error) {
	log.Printf("LoadRefresh: %s\n", code)
	session := s.session.Copy()
	defer session.Close()
	accesses := session.DB(s.dbName).C(ACCESS_COL)
	accData := new(osin.AccessData)
	err := accesses.Find(bson.M{REFRESHTOKEN: code}).One(accData)
	return accData, err

	//if d, ok := s.refresh[code]; ok {
	//	return s.LoadAccess(d)
	//}
	//return nil, errors.New("Refresh not found")
}

// RemoveRefresh removes the refresh token
func (s *MongoStorage) RemoveRefresh(code string) error {
	log.Printf("RemoveRefresh: %s\n", code)
	session := s.session.Copy()
	defer session.Close()
	accesses := session.DB(s.dbName).C(ACCESS_COL)
	return accesses.Update(bson.M{REFRESHTOKEN: code}, bson.M{
		"$unset": bson.M{
			REFRESHTOKEN: 1,
		}})

	//delete(s.refresh, code)
	//return nil
}
