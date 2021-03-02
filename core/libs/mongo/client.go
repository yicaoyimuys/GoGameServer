package mongo

import (
	"GoGameServer/core/libs/dict"
	"time"

	"gopkg.in/mgo.v2"
)

type Client struct {
	session *mgo.Session
	db      string
}

func NewClient(mongoConfig map[string]interface{}) (*Client, error) {
	addr := dict.GetString(mongoConfig, "host") + ":" + dict.GetString(mongoConfig, "port")
	user := dict.GetString(mongoConfig, "user")
	pwd := dict.GetString(mongoConfig, "password")
	db := dict.GetString(mongoConfig, "db")

	dialInfo := &mgo.DialInfo{
		Addrs:    []string{addr},
		Username: user,
		Password: pwd,
		Timeout:  5 * time.Second,
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}

	//返回数据
	client := &Client{
		session: session,
		db:      db,
	}
	return client, nil
}

func (this *Client) connect(collection string) (*mgo.Session, *mgo.Collection) {
	s := this.session.Copy()
	c := s.DB(this.db).C(collection)
	return s, c
}

func (this *Client) Insert(collection string, docs ...interface{}) error {
	ms, c := this.connect(collection)
	defer ms.Close()
	return c.Insert(docs...)
}

func (this *Client) FindOne(collection string, query, selector, result interface{}) error {
	ms, c := this.connect(collection)
	defer ms.Close()
	return c.Find(query).Select(selector).One(result)
}

func (this *Client) FindAll(collection string, query, selector, result interface{}) error {
	ms, c := this.connect(collection)
	defer ms.Close()
	return c.Find(query).Select(selector).All(result)
}

func (this *Client) Update(collection string, query, update interface{}) error {
	ms, c := this.connect(collection)
	defer ms.Close()
	return c.Update(query, update)
}

func (this *Client) UpdateAll(collection string, query, update interface{}) error {
	ms, c := this.connect(collection)
	defer ms.Close()
	_, err := c.UpdateAll(query, update)
	return err
}

func (this *Client) Remove(collection string, query interface{}) error {
	ms, c := this.connect(collection)
	defer ms.Close()
	return c.Remove(query)
}

func (this *Client) RemoveAll(collection string, query interface{}) error {
	ms, c := this.connect(collection)
	defer ms.Close()
	_, err := c.RemoveAll(query)
	return err
}
