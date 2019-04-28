package database

// mongodb://myDBReader:D1fficultP%40ssw0rd@mongodb0.example.com:27017/admin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	client      *mongo.Client
	db          *mongo.Database
	collections *collections
}

type MongoConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	AuthDB   string `json:"auth_db"`
	DataDB   string `json:"data_db"`
}

type collections struct {
	messages,
	reactions,
	status,
	statusroles,
	statuschanges,
	voicechannels,
	memberchange *mongo.Collection
}

func (m *MongoDB) Connect(params interface{}) (err error) {
	cfg, ok := params.(*MongoConfig)
	if !ok {
		return errors.New("invalid config data type")
	}

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.AuthDB)
	if m.client, err = mongo.NewClient(options.Client().ApplyURI(uri)); err != nil {
		return
	}

	if err = m.client.Connect(ctxTimeout(10 * time.Second)); err != nil {
		return
	}

	err = m.client.Ping(ctxTimeout(2*time.Second), readpref.Primary())

	m.db = m.client.Database(cfg.AuthDB)

	m.collections = &collections{
		messages:      m.db.Collection("messages"),
		reactions:     m.db.Collection("reactions"),
		status:        m.db.Collection("status"),
		statusroles:   m.db.Collection("statusroles"),
		statuschanges: m.db.Collection("statuschanges"),
		voicechannels: m.db.Collection("voicechannels"),
		memberchange:  m.db.Collection("memberchange"),
	}

	return
}

func (m *MongoDB) Close() {
	m.client.Disconnect(ctxTimeout(5 * time.Second))
}

// --- DB FUNCTIONS --------------------

func (m *MongoDB) AddMessage(msg *Message) error {
	return m.insert(m.collections.messages, msg)
}

func (m *MongoDB) AddStatusUpdate(su *StatusUpdate) error {
	return m.insert(m.collections.statuschanges, su)
}

func (m *MongoDB) AddMembChange(mc *MemberChange) error {
	return m.insert(m.collections.memberchange, mc)
}

func (m *MongoDB) AddMembStatus(ms *MemberStatus) error {
	return m.insert(m.collections.status, ms)
}

func (m *MongoDB) AddMembStatusRoles(mrsc *MemberStatusRolesCollection) error {
	return m.insert(m.collections.statusroles, mrsc)
}

func (m *MongoDB) AddReaction(r *Reaction) error {
	return m.insert(m.collections.reactions, r)
}

func (m *MongoDB) AddVoice(v *Voice) error {
	return m.insert(m.collections.voicechannels, v)
}

// --- HELPERS -------------------------

func (m *MongoDB) insert(collection *mongo.Collection, v interface{}) error {
	_, err := collection.InsertOne(ctxTimeout(5*time.Second), v)
	return err
}

func ctxTimeout(d time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), d)
	return ctx
}
