package mysql

import (
	"database/sql"
	"strconv"

	"github.com/admpub/nging/application/library/dbmanager/driver"
	"github.com/webx-top/db/lib/factory"
)

func (m *mySQL) getDatabases() ([]string, error) {
	sqlStr := `SELECT SCHEMA_NAME FROM information_schema.SCHEMATA`
	if m.getBVersion() < 5 {
		sqlStr = `SHOW DATABASES`
	}
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return nil, err
	}
	ret := []string{}
	for rows.Next() {
		var v sql.NullString
		err := rows.Scan(&v)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v.String)
	}
	return ret, nil
}

func (m *mySQL) getTables() ([]string, error) {
	sqlStr := `SHOW TABLES`
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return nil, err
	}
	ret := []string{}
	for rows.Next() {
		var v sql.NullString
		err := rows.Scan(&v)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v.String)
	}
	return ret, nil
}

func (m *mySQL) newParam() *factory.Param {
	return factory.NewParam(m.db)
}

func (m *mySQL) getBVersion() int {
	part := driver.RegexpNotNumber.Split(m.getVersion(), 2)
	i, err := strconv.Atoi(part[0])
	if err != nil {
		m.Echo().Logger().Error(err.Error())
		return 0
	}
	return i
}

func (m *mySQL) getVersion() string {
	if len(m.version) > 0 {
		return m.version
	}
	row, err := m.newParam().SetCollection(`SELECT version()`).QueryRow()
	if err != nil {
		return err.Error()
	}
	var v sql.NullString
	err = row.Scan(&v)
	if err != nil {
		return err.Error()
	}
	m.version = v.String
	return v.String
}

func (m *mySQL) baseInfo() error {
	dbList, err := m.getDatabases()
	if err != nil {
		return err
	}
	m.Set(`dbList`, dbList)
	if len(m.DbAuth.Db) > 0 {
		tableList, err := m.getTables()
		if err != nil {
			return err
		}
		m.Set(`tableList`, tableList)
	}

	m.Set(`dbVersion`, m.getVersion())
	return nil
}
