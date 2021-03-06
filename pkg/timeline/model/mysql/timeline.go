package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	DBName    = "timeline"
	TableName = "timeline"
)

type TimeLine struct {
	ID         uint32    `json::"id"`
	Value      string    `json:"value"`
	Label      string    `json:"label"`
	Title      string    `json:"title"`
	EventTime  time.Time `json:"event_time"`
	CreateTime time.Time `json:"create_time"`
}

const (
	mysqlCreateDatabase = iota
	mysqlTimeLineCreateTable
	mysqlTimeLineInsert
	mysqlTimelineDelete
	mysqlTimeLineUpdateByID
	mysqlTimeLineSelectAllUnDeletedOrderByCreateTime
	mysqlTimeLineSelectAllUnDeletedOrderByEventTime
	mysqlTimeLineSelectByLabelUnDeleted
)

var (
	errInvalidUserCreateDefaultUser = errors.New("[user] invalid default user create ")

	TimeLineSQLString = []string{
		fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s`, DBName),
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s(
			id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(512) NOT NULL DEFAULT '' COMMENT'每个节点标题',
			deleted BOOLEAN NOT NULL DEFAULT FALSE COMMENT '删除位',
			value VARCHAR(2048) NOT NULL COMMENT'timeline 每个节点具体内容',
			label ENUM ('大事件', '政策新规', '反腐', '重磅发言') NOT  NULL DEFAULT '大事件' COMMENT '标签',
			event_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT'事件发生时间',
			create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '数据插入时间,为了排序'
		)ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`, DBName, TableName),
		fmt.Sprintf(`INSERT INTO %s.%s (title, deleted, value, label, event_time) VALUES (?, ?,?, ?,?)`, DBName, TableName),
		fmt.Sprintf(`UPDATE %s.%s SET deleted=? WHERE id=?`, DBName, TableName),
		fmt.Sprintf(`UPDATE %s.%s SET title = ?, value = ?, label=?, event_time=? WHERE id=?`, DBName, TableName),
		fmt.Sprintf(`SELECT id, value,label, title,event_time, create_time FROM %s.%s WHERE deleted = ? ORDER BY create_time DESC`, DBName, TableName),
		fmt.Sprintf(`SELECT id, value,label, title,event_time, create_time FROM %s.%s WHERE deleted = ? ORDER BY event_time DESC`, DBName, TableName),
		fmt.Sprintf(`SELECT id, value,label, title,event_time FROM %s.%s WHERE deleted = ? AND label=? ORDER BY event_time DESC`, DBName, TableName),
	}
)

func CreateDatabase(db *sql.DB) error {
	_, err := db.Exec(TimeLineSQLString[mysqlCreateDatabase])
	if err != nil {
		return err
	}

	return nil
}

func CreateTimeLine(db *sql.DB) error {
	_, err := db.Exec(TimeLineSQLString[mysqlTimeLineCreateTable])
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func InsertTimeLine(db *sql.DB, title, value, label string, eventTime time.Time) error {
	var deleted bool = false
	result, err := db.Exec(TimeLineSQLString[mysqlTimeLineInsert], title, deleted, value, label, eventTime)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidUserCreateDefaultUser
	}

	return nil
}

func DeleteTimeLine(db *sql.DB, id uint32) error {
	result, err := db.Exec(TimeLineSQLString[mysqlTimelineDelete], true, id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidUserCreateDefaultUser
	}

	return nil
}

func UpdateTimeLineByID(db *sql.DB, id uint32, title, value, label string, eventTime time.Time) error {
	_, err := db.Exec(TimeLineSQLString[mysqlTimeLineUpdateByID], title, value, label, eventTime, id)
	if err != nil {
		return err
	}

	return nil
}

func SelectAllUnDeletedTimeLineOrderByCreateTime(db *sql.DB) ([]*TimeLine, error) {
	var (
		TimeLines []*TimeLine

		ID         uint32
		Value      string
		Label      string
		Title      string
		EventTime  time.Time
		CreateTime time.Time
	)

	rows, err := db.Query(TimeLineSQLString[mysqlTimeLineSelectAllUnDeletedOrderByCreateTime], false)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&ID, &Value, &Label, &Title, &EventTime, &CreateTime); err != nil {

			return nil, err
		}

		TimeLine := &TimeLine{
			ID:         ID,
			Value:      Value,
			Label:      Label,
			Title:      Title,
			EventTime:  EventTime,
			CreateTime: CreateTime,
		}

		TimeLines = append(TimeLines, TimeLine)
	}

	return TimeLines, nil
}

func SelectAllUnDeletedTimeLineOrderByEventTime(db *sql.DB) ([]*TimeLine, error) {
	var (
		TimeLines []*TimeLine

		ID         uint32
		Value      string
		Label      string
		Title      string
		EventTime  time.Time
		CreateTime time.Time
	)

	rows, err := db.Query(TimeLineSQLString[mysqlTimeLineSelectAllUnDeletedOrderByEventTime], false)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&ID, &Value, &Label, &Title, &EventTime, &CreateTime); err != nil {

			return nil, err
		}

		TimeLine := &TimeLine{
			ID:         ID,
			Value:      Value,
			Label:      Label,
			Title:      Title,
			EventTime:  EventTime,
			CreateTime: CreateTime,
		}

		TimeLines = append(TimeLines, TimeLine)
	}

	return TimeLines, nil
}

func SelectByLabelUnDeletedTimeLine(db *sql.DB, label string) ([]*TimeLine, error) {
	var (
		TimeLines []*TimeLine

		ID        uint32
		Value     string
		Label     string
		Title     string
		EventTime time.Time
	)

	rows, err := db.Query(TimeLineSQLString[mysqlTimeLineSelectByLabelUnDeleted], false, label)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&ID, &Value, &Label, &Title, &EventTime); err != nil {

			return nil, err
		}

		TimeLine := &TimeLine{
			ID:        ID,
			Value:     Value,
			Label:     Label,
			Title:     Title,
			EventTime: EventTime,
		}

		TimeLines = append(TimeLines, TimeLine)
	}

	return TimeLines, nil
}
