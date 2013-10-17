package release

import (
	"database/sql"
	"fmt"
	"github.com/jsli/cp_release/constant"
)

type CpComponent struct {
	Id           int64
	CpId         int64
	Flag         int
	LastModifyTs int64
	RelPath      string
}

func (cc *CpComponent) save(dal *Dal, table string) (int64, error) {
	prepare := fmt.Sprintf("INSERT %s SET cp_id=?, flag=?, last_modify_ts=?, path=?", table)
	stmt, err := dal.Link.Prepare(prepare)

	if err != nil {
		return -1, err
	}
	res, err := stmt.Exec(cc.CpId, cc.Flag, cc.LastModifyTs, cc.RelPath)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	return id, err
}

func DeleteCpComponent(dal *Dal, delete_sql string) (int64, error) {
	stmt, err := dal.Link.Prepare(delete_sql)

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

func findCpComponent(dal *Dal, query string) (*CpComponent, error) {
	row := dal.Link.QueryRow(query)
	cc := CpComponent{}
	err := row.Scan(&cc.Id, &cc.CpId, &cc.Flag, &cc.LastModifyTs, &cc.RelPath)
	if err != nil {
		return nil, err
	}

	return &cc, nil
}

func findCpComponentList(dal *Dal, query string) ([]*CpComponent, error) {
	rows, err := dal.Link.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	ccs := make([]*CpComponent, 0, 100)
	for rows.Next() {
		cc := CpComponent{}
		err := rows.Scan(&cc.Id, &cc.CpId, &cc.Flag, &cc.LastModifyTs, &cc.RelPath)
		if err != nil || cc.Id < 0 {
			continue
		}
		ccs = append(ccs, &cc)
	}
	return ccs, nil
}

type Arbi struct {
	CpComponent
}

func (arbi Arbi) String() string {
	return fmt.Sprintf("Arbi(id=%d, cpid=%d, flag=%d, ts=%d, path=%s)",
		arbi.Id, arbi.CpId, arbi.Flag, arbi.LastModifyTs, arbi.RelPath)
}

func (arbi *Arbi) Save(dal *Dal) (int64, error) {
	return arbi.save(dal, constant.TABLE_ARBI)
}

func (arbi *Arbi) Delete(dal *Dal) (int64, error) {
	delete_sql := fmt.Sprintf("DELETE FROM %s where id=%d AND flag=%d", constant.TABLE_ARBI, arbi.Id)
	return DeleteCpComponent(dal, delete_sql)
}

func DeleteArbiByCpId(dal *Dal, cp_id int64) (int64, error) {
	delete_sql := fmt.Sprintf("DELETE FROM %s where cp_id=%d", constant.TABLE_ARBI, cp_id)
	return DeleteCpComponent(dal, delete_sql)
}

func FindArbiByCpId(dal *Dal, cp_id int64) (*Arbi, error) {
	query := fmt.Sprintf("SELECT * FROM %s where cp_id='%s' AND flag=%d", constant.TABLE_ARBI, cp_id, constant.AVAILABLE_FLAG)
	return FindArbi(dal, query)
}

func FindArbiByPath(dal *Dal, path string) (*Arbi, error) {
	query := fmt.Sprintf("SELECT * FROM %s where path='%s' AND flag=%d", constant.TABLE_ARBI, path, constant.AVAILABLE_FLAG)
	return FindArbi(dal, query)
}

func FindArbi(dal *Dal, query string) (*Arbi, error) {
	cc, err := findCpComponent(dal, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	arbi := &Arbi{}
	arbi.Id = cc.Id
	arbi.CpId = cc.CpId
	arbi.Flag = cc.Flag
	arbi.LastModifyTs = cc.LastModifyTs
	arbi.RelPath = cc.RelPath

	return arbi, nil
}

func FindArbiList(dal *Dal, query string) ([]*Arbi, error) {
	ccs, err := findCpComponentList(dal, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	arbis := make([]*Arbi, 0, 10)
	for _, cc := range ccs {
		arbi := &Arbi{}
		arbi.Id = cc.Id
		arbi.CpId = cc.CpId
		arbi.Flag = cc.Flag
		arbi.LastModifyTs = cc.LastModifyTs
		arbi.RelPath = cc.RelPath
		arbis = append(arbis, arbi)
	}

	return arbis, nil
}

type Grbi struct {
	CpComponent
}

func (grbi Grbi) String() string {
	return fmt.Sprintf("Grbi(id=%d, cpid=%d, flag=%d, ts=%d, path=%s)",
		grbi.Id, grbi.CpId, grbi.Flag, grbi.LastModifyTs, grbi.RelPath)
}

func (grbi *Grbi) Save(dal *Dal) (int64, error) {
	return grbi.save(dal, constant.TABLE_GRBI)
}

func (grbi *Grbi) Delete(dal *Dal) (int64, error) {
	delete_sql := fmt.Sprintf("DELETE FROM %s where id=%d", constant.TABLE_GRBI, grbi.Id)
	return DeleteCpComponent(dal, delete_sql)
}

func DeleteGrbiByCpId(dal *Dal, cp_id int64) (int64, error) {
	delete_sql := fmt.Sprintf("DELETE FROM %s where cp_id=%d", constant.TABLE_GRBI, cp_id)
	return DeleteCpComponent(dal, delete_sql)
}

func FindGrbiByCpId(dal *Dal, cp_id int64) (*Grbi, error) {
	query := fmt.Sprintf("SELECT * FROM %s where cp_id='%s' AND flag=%d", constant.TABLE_GRBI, cp_id, constant.AVAILABLE_FLAG)
	return FindGrbi(dal, query)
}

func FindGrbiByPath(dal *Dal, path string) (*Grbi, error) {
	query := fmt.Sprintf("SELECT * FROM %s where path='%s' AND flag=%d", constant.TABLE_GRBI, path, constant.AVAILABLE_FLAG)
	return FindGrbi(dal, query)
}

func FindGrbi(dal *Dal, query string) (*Grbi, error) {
	cc, err := findCpComponent(dal, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	grbi := &Grbi{}
	grbi.Id = cc.Id
	grbi.CpId = cc.CpId
	grbi.Flag = cc.Flag
	grbi.LastModifyTs = cc.LastModifyTs
	grbi.RelPath = cc.RelPath

	return grbi, nil
}

func FindGrbiList(dal *Dal, query string) ([]*Grbi, error) {
	ccs, err := findCpComponentList(dal, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	grbis := make([]*Grbi, 0, 10)
	for _, cc := range ccs {
		grbi := &Grbi{}
		grbi.Id = cc.Id
		grbi.CpId = cc.CpId
		grbi.Flag = cc.Flag
		grbi.LastModifyTs = cc.LastModifyTs
		grbi.RelPath = cc.RelPath
		grbis = append(grbis, grbi)
	}

	return grbis, nil
}


type Rfic struct {
	CpComponent
}

func (rfic Rfic) String() string {
	return fmt.Sprintf("Rfic(id=%d, cpid=%d, flag=%d, ts=%d, path=%s)",
		rfic.Id, rfic.CpId, rfic.Flag, rfic.LastModifyTs, rfic.RelPath)
}

func (rfic *Rfic) Save(dal *Dal) (int64, error) {
	return rfic.save(dal, constant.TABLE_RFIC)
}

func (rfic *Rfic) Delete(dal *Dal) (int64, error) {
	delete_sql := fmt.Sprintf("DELETE FROM %s where id=%d AND flag=%d", constant.TABLE_RFIC, rfic.Id)
	return DeleteCpComponent(dal, delete_sql)
}

func DeleteRficByCpId(dal *Dal, cp_id int64) (int64, error) {
	delete_sql := fmt.Sprintf("DELETE FROM %s where cp_id=%d", constant.TABLE_RFIC, cp_id)
	return DeleteCpComponent(dal, delete_sql)
}

func FindRficByCpId(dal *Dal, cp_id int64) (*Rfic, error) {
	query := fmt.Sprintf("SELECT * FROM %s where cp_id='%s' AND flag=%d", constant.TABLE_RFIC, cp_id, constant.AVAILABLE_FLAG)
	return FindRfic(dal, query)
}

func FindRficByPath(dal *Dal, path string) (*Rfic, error) {
	query := fmt.Sprintf("SELECT * FROM %s where path='%s' AND flag=%d", constant.TABLE_RFIC, path, constant.AVAILABLE_FLAG)
	return FindRfic(dal, query)
}

func FindRfic(dal *Dal, query string) (*Rfic, error) {
	cc, err := findCpComponent(dal, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	rfic := &Rfic{}
	rfic.Id = cc.Id
	rfic.CpId = cc.CpId
	rfic.Flag = cc.Flag
	rfic.LastModifyTs = cc.LastModifyTs
	rfic.RelPath = cc.RelPath

	return rfic, nil
}

func FindRficList(dal *Dal, query string) ([]*Rfic, error) {
	ccs, err := findCpComponentList(dal, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	rfics := make([]*Rfic, 0, 10)
	for _, cc := range ccs {
		rfic := &Rfic{}
		rfic.Id = cc.Id
		rfic.CpId = cc.CpId
		rfic.Flag = cc.Flag
		rfic.LastModifyTs = cc.LastModifyTs
		rfic.RelPath = cc.RelPath
		rfics = append(rfics, rfic)
	}

	return rfics, nil
}