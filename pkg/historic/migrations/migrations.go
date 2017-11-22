package migrations

import (
	"fmt"

	"github.com/go-clog/clog"
	"github.com/go-xorm/xorm"
)

const minDBVersion = 1

// Migration interface for create any migration
type Migration interface {
	Description() string
	Migrate(*xorm.Engine) error
}

type migration struct {
	description string
	migrate     func(*xorm.Engine) error
}

// NewMigration Return new Migration struct.
func NewMigration(desc string, fn func(*xorm.Engine) error) Migration {
	return &migration{desc, fn}
}

// Description Return the description of Migration
func (m *migration) Description() string {
	return m.description
}

// Migrate Execute function migrate to do the migration's magic
func (m *migration) Migrate(x *xorm.Engine) error {
	return m.migrate(x)
}

// Version The version table. Should have only one row with id==1
type Version struct {
	ID      int64
	Version int64
}

var migrations = []Migration{}

// Migrate database to current version
func Migrate(x *xorm.Engine) error {
	if err := x.Sync(new(Version)); err != nil {
		return fmt.Errorf("sync: %v", err)
	}

	currentVersion := &Version{ID: 1}
	has, err := x.Get(currentVersion)
	if err != nil {
		return fmt.Errorf("get: %v", err)
	} else if !has {
		currentVersion.ID = 0
		currentVersion.Version = int64(minDBVersion + len(migrations))

		if _, err = x.InsertOne(currentVersion); err != nil {
			return fmt.Errorf("insert: %v", err)
		}
	}

	v := currentVersion.Version
	if minDBVersion > v {
		clog.Fatal(0, "Actual version is not compatible with old version")
		return nil
	}

	if int(v-minDBVersion) > len(migrations) {
		currentVersion.Version = int64(len(migrations) + minDBVersion)
		_, err = x.Id(1).Update(currentVersion)
		return err
	}
	for i, m := range migrations[v-minDBVersion:] {
		clog.Info("Migration: %s", m.Description())
		if err = m.Migrate(x); err != nil {
			return fmt.Errorf("do migrate: %v", err)
		}
		currentVersion.Version = v + int64(i) + 1
		if _, err = x.Id(1).Update(currentVersion); err != nil {
			return err
		}
	}
	return nil
}

func sessionRelease(sess *xorm.Session) {
	if !sess.IsCommitedOrRollbacked {
		sess.Rollback()
	}
	sess.Close()
}
