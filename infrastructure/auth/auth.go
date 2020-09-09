package auth

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"time"
)

type AuthInterface interface {
	CreateAuth(uuid.UUID, *TokenDetails) error
	FetchAuth(string) (uuid.UUID, error)
	DeleteRefresh(string) error
	DeleteTokens(*AccessDetails) error
}

type ClientData struct {
	client *redis.Client
}

var _ AuthInterface = &ClientData{}

func NewAuth(client *redis.Client) *ClientData {
	return &ClientData{client: client}
}

type AccessDetails struct {
	TokenUuid string
	UserUUID  uuid.UUID
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUuid    string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

//Save token metadata to Redis
func (tk *ClientData) CreateAuth(userid uuid.UUID, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	atCreated, err := tk.client.Set(td.TokenUuid, userid.String(), at.Sub(now)).Result()
	if err != nil {
		return err
	}
	rtCreated, err := tk.client.Set(td.RefreshUuid, userid.String(), rt.Sub(now)).Result()
	if err != nil {
		return err
	}
	if atCreated == "0" || rtCreated == "0" {
		return errors.New("no record inserted")
	}
	return nil
}

//Check the metadata saved
func (tk *ClientData) FetchAuth(tokenUuid string) (uuid.UUID, error) {
	userid, err := tk.client.Get(tokenUuid).Result()
	userUUID, err := uuid.Parse(userid)
	if err != nil {
		return userUUID, err
	}
	return userUUID, nil
}

//Once a user row in the token table
func (tk *ClientData) DeleteTokens(authD *AccessDetails) error {
	//get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%d", authD.TokenUuid, authD.UserUUID)
	//delete access token
	deletedAt, err := tk.client.Del(authD.TokenUuid).Result()
	if err != nil {
		return err
	}
	//delete refresh token
	deletedRt, err := tk.client.Del(refreshUuid).Result()
	if err != nil {
		return err
	}
	//When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
	}
	return nil
}

func (tk *ClientData) DeleteRefresh(refreshUuid string) error {
	//delete refresh token
	deleted, err := tk.client.Del(refreshUuid).Result()
	if err != nil || deleted == 0 {
		return err
	}
	return nil
}
