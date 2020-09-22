package persistence

import (
	"Repository-Pattern/domain/model"
	"Repository-Pattern/domain/repositories"
	"database/sql"
	"fmt"
)

type PostLabelRepo struct {
	db *sql.DB
}

func NewPostLabelRepository(db *sql.DB) *PostLabelRepo {
	return &PostLabelRepo{db}
}

//PostLabelRepo implements the repository.PostLabelRepository interface
var _ repositories.PostLabelRepository = &PostLabelRepo{}

func (r *PostLabelRepo) SavePostLabel(postLabel *model.PostLabel) (*model.PostLabel, error) {
	queryInsert := fmt.Sprintf("INSERT INTO %s (post_id, label_id) "+
		"VALUES ($1, $2)", "post_labels")
	prepare, err := r.db.Prepare(queryInsert)
	if err != nil {
		return postLabel, err
	}
	_, err = prepare.Exec(&postLabel.PostID, &postLabel.LabelID)
	if err != nil {
		return postLabel, err
	}
	return postLabel, err
}

func (r *PostLabelRepo) GetPostLabel(id uint64) (*model.PostLabel, error) {
	var postLabel model.PostLabel
	querySelect := fmt.Sprint("SELECT * FROM post_labels WHERE post_id=$1")
	prepare, err := r.db.Prepare(querySelect)
	if err != nil {
		return &postLabel, err
	}
	err = prepare.QueryRow(id).Scan(&postLabel.PostID, &postLabel.LabelID)
	if err != nil {
		fmt.Println(err)
	}
	return &postLabel, nil
}

func (r *PostLabelRepo) GetAllPostLabel() ([]model.PostLabel, error) {
	var postLabels []model.PostLabel
	queryGetUsers := fmt.Sprintf("SELECT * FROM post_labels")
	prepare, err := r.db.Prepare(queryGetUsers)
	if err != nil {
		return postLabels, err
	}
	rows, err := prepare.Query()
	if err != nil {
		return postLabels, err
	}
	for rows.Next() {
		var postLabel model.PostLabel
		err := rows.Scan(&postLabel.PostID, &postLabel.LabelID)
		if err != nil {
			return postLabels, err
		}
		postLabels = append(postLabels, postLabel)
	}
	defer rows.Close()
	return postLabels, nil
}

func (r *PostLabelRepo) UpdatePostLabel(postLabel *model.PostLabel) (*model.PostLabel, error) {
	queryUpdate := fmt.Sprint("UPDATE post_labels SET post_id=$1, label_id=$2 WHERE post_id=$3")
	prepare, err := r.db.Prepare(queryUpdate)
	if err != nil {
		return postLabel, err
	}
	_, err = prepare.Exec(postLabel.PostID, postLabel.LabelID, postLabel.PostID)
	return postLabel, err
}

func (r *PostLabelRepo) DeletePostLabel(id uint64) error {
	querySoftDelete := fmt.Sprint("DELETE FROM post_labels WHERE post_id=$1")
	prepare, err := r.db.Prepare(querySoftDelete)
	if err != nil {
		return err
	}
	_, err = prepare.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
