package dse

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	apmModel "github.com/jiqiang/tst/server/apm/model"
	uiModel "github.com/jiqiang/tst/server/ui/model"
)

const (
	host         = "10.10.0.158"
	port         = 9042
	keyspace     = "iepmaster"
	timeout      = 15 * time.Second
	protoVersion = 4
)

// Cluster represents a new cassandra cluster.
type Cluster struct {
	session *gocql.Session
}

// Init initializes a cassandra db cluster and session.
func (c *Cluster) Init() error {
	cluster := gocql.NewCluster(host)
	cluster.Port = port
	cluster.Keyspace = keyspace
	cluster.Timeout = timeout
	cluster.ProtoVersion = protoVersion

	var err error
	c.session, err = cluster.CreateSession()
	if err != nil {
		c.session.Close()
		return err
	}

	return nil
}

// InsertAssets inserts assets data into event_current table.
func (c Cluster) InsertAssets(enterpriseSourceKey string, siteName string, assets []apmModel.Asset) error {
	batch := gocql.NewBatch(gocql.LoggedBatch)
	sql := `INSERT INTO event_current(enterprise_uid, event_bucket, resrc_uid, event_ts) VALUES (?, ?, ?, ?)`
	for _, asset := range assets {
		resrcUID := fmt.Sprintf("%s.%s", enterpriseSourceKey, asset.SourceKey)
		batch.Query(sql, enterpriseSourceKey, siteName, resrcUID, time.Now().UTC())
	}
	err := c.session.ExecuteBatch(batch)
	if err != nil {
		return err
	}
	return nil
}

// GetAssets gets all asset data.
func (c Cluster) GetAssets() ([]uiModel.Asset, error) {
	sql := `SELECT enterprise_uid, event_bucket, event_ts FROM event_current`
	iter := c.session.Query(sql).Iter()
	assets := []uiModel.Asset{}
	var enterpriseUID, eventBucket string
	var eventTs time.Time
	for iter.Scan(&enterpriseUID, &eventBucket, &eventTs) {
		assets = append(assets, uiModel.Asset{
			EnterpriseSourceKey:        enterpriseUID,
			SiteName:                   eventBucket,
			TimeElapsedSinceLastUpdate: time.Now().Sub(eventTs).Minutes(),
		})
	}
	return assets, nil
}
