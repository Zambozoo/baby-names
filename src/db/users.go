package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/boltdb/bolt"
)

var userBucket = []byte("users")

type (
	UserDB struct {
		db *bolt.DB
	}
	User struct {
		Username        string
		PartnerUsername string
		LikedNames      map[string]struct{}
		DislikedNames   map[string]struct{}
	}
)

func NewUserDB(path string) (*UserDB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(userBucket)
		return err
	}); err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &UserDB{db: db}, nil
}

func (udb *UserDB) GetUser(username string) (*User, error) {
	var user *User
	if err := udb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(userBucket)
		userBytes := b.Get([]byte(username))
		if userBytes == nil {
			return nil
		}
		user = &User{}
		return json.Unmarshal(userBytes, user)
	}); err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (udb *UserDB) DeleteUsers(username string) (*User, error) {
	var user User
	if err := udb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(userBucket)
		userBytes := b.Get([]byte(username))
		if userBytes == nil {
			return fmt.Errorf("No such user")
		}
		if err := json.Unmarshal(userBytes, &user); err != nil {
			return err
		}

		return errors.Join(
			b.Delete([]byte(username)),
			b.Delete([]byte(user.PartnerUsername)),
		)
	}); err != nil {
		return nil, fmt.Errorf("failed to delete users: %w", err)
	}

	return &user, nil
}

func (udb *UserDB) UpdateUser(user *User) error {
	if err := udb.db.Update(func(tx *bolt.Tx) error {
		userBytes, err := json.Marshal(user)
		if err != nil {
			return err
		}
		b := tx.Bucket(userBucket)
		return b.Put([]byte(user.Username), userBytes)
	}); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (udb *UserDB) CreateUsers(username, partnerUsername string) (*User, *User, error) {
	user := &User{
		Username:        username,
		PartnerUsername: partnerUsername,
		LikedNames:      map[string]struct{}{},
		DislikedNames:   map[string]struct{}{},
	}
	partnerUser := &User{
		Username:        partnerUsername,
		PartnerUsername: username,
		LikedNames:      map[string]struct{}{},
		DislikedNames:   map[string]struct{}{},
	}

	if err := udb.db.Update(func(tx *bolt.Tx) error {
		userBytes, err := json.Marshal(user)
		if err != nil {
			return err
		}
		partnerUserBytes, err := json.Marshal(partnerUser)
		if err != nil {
			return err
		}

		b := tx.Bucket(userBucket)
		return errors.Join(
			b.Put([]byte(user.Username), userBytes),
			b.Put([]byte(partnerUser.Username), partnerUserBytes),
		)
	}); err != nil {
		return nil, nil, fmt.Errorf("failed to create users [%s, %s]: %w", username, partnerUsername, err)
	}

	return user, partnerUser, nil
}

func (u *User) LikeName(name string) {
	u.LikedNames[name] = struct{}{}
	delete(u.DislikedNames, name)
}

func (u *User) DislikeName(name string) {
	u.DislikedNames[name] = struct{}{}
	delete(u.LikedNames, name)
}
func (u *User) Matched(partnerUser *User) []string {
	matched := make(map[string]struct{}, len(u.LikedNames))
	for key := range u.LikedNames {
		if _, ok := partnerUser.LikedNames[key]; ok {
			matched[key] = struct{}{}
		}
	}

	return slices.Sorted(maps.Keys(matched))
}
