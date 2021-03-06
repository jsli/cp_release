package release

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/howeyc/fsnotify"
	"github.com/jsli/cp_release/constant"
	"github.com/jsli/cp_release/policy"
	"github.com/jsli/gtbox/pathutil"
	"time"
)

type CpRelease struct {
	Id            int64
	Mode          string
	Sim           string
	Version       string
	VersionScalar int64
	Flag          int
	LastModifyTs  int64
	RelPath       string
	Prefix        string
}

func (cp CpRelease) String() string {
	return fmt.Sprintf("CpRelease(id=%d, mode=%s, sim=%s, version=%s, version_scalar=%d, flag=%d, last_modify_ts=%d, rel_path=%s, prefix=%s)",
		cp.Id, cp.Mode, cp.Sim, cp.Version, cp.VersionScalar, cp.Flag, cp.LastModifyTs, cp.RelPath, cp.Prefix)
}

func (cp *CpRelease) Save(dal *Dal) (int64, error) {
	insert_sql := fmt.Sprintf("INSERT %s SET mode=?, sim=?, version=?, version_scalar=?, flag=?, last_modify_ts=?, path=?, prefix=?",
		constant.TABLE_CP)
	stmt, err := dal.DB.Prepare(insert_sql)

	if err != nil {
		return -1, err
	}
	res, err := stmt.Exec(cp.Mode, cp.Sim, cp.Version, cp.VersionScalar, cp.Flag, cp.LastModifyTs, cp.RelPath, cp.Prefix)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	return id, err
}

func (cp *CpRelease) Update(dal *Dal) (int64, error) {
	update_sql := fmt.Sprintf("UPDATE %s SET flag=?, last_modify_ts=? where id =%d", constant.TABLE_CP, cp.Id)
	stmt, err := dal.DB.Prepare(update_sql)

	if err != nil {
		return -1, err
	}
	res, err := stmt.Exec(cp.Flag, cp.LastModifyTs)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	return id, err
}

func (cp *CpRelease) Delete(dal *Dal) (int64, error) {
	return DeleteCpByPath(dal, cp.RelPath)
}

func (cp *CpRelease) LoadSelfFromFileEvent(event *fsnotify.FileEvent) error {
	path := event.Name
	parent_path := pathutil.ParentPath(path)
	cp.Mode = constant.PATH_TO_MODE[parent_path[:len(parent_path)-1]]

	base_name := pathutil.BaseName(path)
	version := policy.ExtractVersion(base_name)
	if version == "" {
		return errors.New(fmt.Sprintf("Illegal version : %s", base_name))
	}
	cp.Version = version
	cp.VersionScalar = policy.QuantitateVersion(version)

	cp.Sim = constant.MODE_TO_SIM[cp.Mode]
	cp.LastModifyTs = time.Now().Unix()
	cp.Flag = constant.AVAILABLE_FLAG
	cp.RelPath = path[constant.PATH_PREFIX_LEN:]
	return nil
}

func DeleteCpByPath(dal *Dal, path string) (int64, error) {
	delete_sql := fmt.Sprintf("DELETE FROM %s where path='%s'", constant.TABLE_CP, path)
	return DeleteCp(dal, delete_sql)
}

func DeleteCp(dal *Dal, delete_sql string) (int64, error) {
	stmt, err := dal.DB.Prepare(delete_sql)

	if err != nil {
		return -1, err
	}
	res, err := stmt.Exec()
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	return id, err
}

func FindCpReleaseByPath(dal *Dal, path string) (*CpRelease, error) {
	query := fmt.Sprintf("SELECT * FROM %s where path='%s' AND flag=%d", constant.TABLE_CP, path, constant.AVAILABLE_FLAG)
	return FindCpRelease(dal, query)
}

func FindCpReleaseById(dal *Dal, id string) (*CpRelease, error) {
	query := fmt.Sprintf("SELECT * FROM %s where id=%s AND flag=%d", constant.TABLE_CP, id, constant.AVAILABLE_FLAG)
	return FindCpRelease(dal, query)
}

func FindCpRelease(dal *Dal, query string) (*CpRelease, error) {
	row := dal.DB.QueryRow(query)
	cp := CpRelease{}
	err := row.Scan(&cp.Id, &cp.Mode, &cp.Sim, &cp.Version, &cp.VersionScalar, &cp.Flag,
		&cp.LastModifyTs, &cp.RelPath, &cp.Prefix)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &cp, nil
}

func FindCpReleaseList(dal *Dal, query string) ([]*CpRelease, error) {
	rows, err := dal.DB.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	cps := make([]*CpRelease, 0, 100)
	for rows.Next() {
		cp := CpRelease{}
		err := rows.Scan(&cp.Id, &cp.Mode, &cp.Sim, &cp.Version, &cp.VersionScalar, &cp.Flag,
			&cp.LastModifyTs, &cp.RelPath, &cp.Prefix)
		if err != nil || cp.Id < 0 {
			continue
		}
		cps = append(cps, &cp)
	}
	return cps, nil
}
