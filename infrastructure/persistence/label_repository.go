package persistence

import (
	"Repository-Pattern/domain/model"
	"Repository-Pattern/domain/repositories"
	"database/sql"
	"fmt"
	"time"
)

type LabelRepo struct {
	db *sql.DB
}

func NewLabelRepository(db *sql.DB) *LabelRepo {
	return &LabelRepo{db}
}

//LabelRepo implements the repository.LabelRepository interface
var _ repositories.LabelRepository = &LabelRepo{}

func (r *LabelRepo) SaveLabel(label *model.Label) (*model.Label, error) {
	label.Prepare()
	queryInsert := fmt.Sprintf("INSERT INTO %s (name, created_at, updated_at, deleted_at) "+
		"VALUES ($1, $2, $3, $4) RETURNING id", "labels")
	prepare, err := r.db.Prepare(queryInsert)
	if err != nil {
		return label, err
	}
	err = prepare.QueryRow(&label.Name, &label.CreatedAt, &label.UpdatedAt, nil).Scan(&label.ID)
	if err != nil {
		return label, err
	}	
	return label, err
}

func (r *LabelRepo) GetLabel(id uint64) (*model.Label, error) {
	var label model.Label
	querySelect := fmt.Sprint("SELECT id, name, created_at, updated_at FROM labels WHERE deleted_at IS NULL AND id=$1")
	prepare, err := r.db.Prepare(querySelect)
	if err != nil {
		return &label, err
	}
	err = prepare.QueryRow(id).Scan(&label.ID, &label.Name, &label.CreatedAt, &label.UpdatedAt)
	if err != nil {
		fmt.Println(err)
	}
	return &label, nil
}

func (r *LabelRepo) GetAllLabel() ([]model.Label, error) {
	var labels []model.Label
	queryGetUsers := fmt.Sprintf("SELECT id, name, created_at, updated_at FROM labels WHERE deleted_at IS NULL")
	rows, err := r.db.Query(queryGetUsers)
	if err != nil {
		return labels, err
	}
	for rows.Next() {
		var label model.Label
		err := rows.Scan(&label.ID, &label.Name, &label.CreatedAt, &label.UpdatedAt)
		if err != nil {
			return labels, err
		}
		labels = append(labels, label)
	}
	defer rows.Close()
	return labels, nil
}

func (r *LabelRepo) UpdateLabel(label *model.Label) (*model.Label, error) {
	label.UpdatedAt = time.Now()
	queryUpdate := fmt.Sprint("UPDATE labels SET name=$1, updated_at=$2 WHERE id=$3")
	prepare, err := r.db.Prepare(queryUpdate)
	if err != nil {
		return label, err
	}
	_, err = prepare.Exec(label.Name, label.UpdatedAt, label.ID)
	return label, err
}

func (r *LabelRepo) DeleteLabel(id uint64) error {
	var label model.Label
	label.DeletedAt = time.Now()
	querySoftDelete := fmt.Sprint("UPDATE labels SET deleted_at=$1 WHERE id=$2")
	prepare, err := r.db.Prepare(querySoftDelete)
	if err != nil {
		return err
	}
	_, err = prepare.Exec(label.DeletedAt, id)
	if err != nil {
		return err
	}
	return nil
}
