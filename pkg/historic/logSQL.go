package historic

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/go-clog/clog"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"

	// Import SQLLite3
	_ "github.com/mattn/go-sqlite3"

	"github.com/rodkranz/monitor/pkg/historic/migrations"
	"github.com/rodkranz/monitor/pkg/message"
)

// SQLConfig is the configuration of SQLWriter depends
type SQLConfig struct {
	Path    string
	SQLite3 string
}

// SQLWriter configure and return the Writer function for log.
func SQLWriter(cfg SQLConfig) (w Writer, err error) {
	if err := os.MkdirAll(path.Dir(cfg.Path), os.ModePerm); err != nil {
		clog.Warn("fail to create directories: %v", err)
		return w, fmt.Errorf("fail to create directories: %v", err)
	}
	connStr := "file:" + cfg.Path + "?cache=shared&mode=rwc"

	x, err := xorm.NewEngine("sqlite3", connStr)
	if err != nil {
		clog.Warn("fail to connect to database: %v", err)
		return w, fmt.Errorf("fail to connect to database: %v", err)
	}

	x.SetMapper(core.GonicMapper{})
	x.SetLogger(xorm.NewSimpleLogger3(os.Stdout, xorm.DEFAULT_LOG_PREFIX, xorm.DEFAULT_LOG_FLAG, core.LOG_WARNING))

	if err = migrations.Migrate(x); err != nil {
		clog.Warn("migrate: %v", err)
		return w, fmt.Errorf("migrate: %v", err)
	}

	if err = x.StoreEngine("InnoDB").Sync2(new(msg)); err != nil {
		clog.Warn("sync database struct error: %v\n", err)
		return w, fmt.Errorf("sync database struct error: %v", err)
	}

	x.ShowSQL(true)

	return func(msg *message.Message) error {
		has, err := x.Where("message_id = ?", *msg.MessageId).Get(new(msg))
		if err != nil {
			return err
		}

		if has {
			clog.Info("[messageID:%s] already imported.", *msg.MessageId)
			clog.Warn("[messageID:%s] message ignored.", *msg.MessageId)
			return nil
		}

		obj := &msg{
			MessageID:   *msg.MessageId,
			Body:        *msg.Body,
			Topic:       msg.Topic(),
			CreatedUnix: time.Now().Unix(),
		}

		if _, err := x.Insert(obj); err != nil {
			clog.Warn("Error insert %#v: %s", obj, err)
			return err
		}

		clog.Warn("[messageID:%s] message registered [id:%d].", *msg.MessageId, obj.ID)
		return nil
	}, nil
}

// msg is model of information that will persist at database.
type msg struct {
	ID        int64  `xorm:"pk autoincr"`
	MessageID string `xorm:"VARCHAR(255)"`
	Topic     string `xorm:"VARCHAR(255)"`
	Body      string `xorm:"Text"`

	Created     time.Time `xorm:"-"`
	CreatedUnix int64
}

// BeforeInsert before insert any record define the current time
func (m *msg) BeforeInsert() {
	m.CreatedUnix = time.Now().Unix()
}

// AfterSet after set columns convert unix data to time.Time.
func (m *msg) AfterSet(colName string, _ xorm.Cell) {
	if colName == "created_unix" {
		m.Created = time.Unix(m.CreatedUnix, 0).Local()
	}
}
