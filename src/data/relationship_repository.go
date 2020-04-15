package data

import (
	"database/sql"
	"fmt"
	"friendMgmt/models"
	"strconv"
	"strings"
)

type IRelationshipRepository interface {
	CreateRelationship(relationship *models.Relationship) int64
	DeleteRelationships(ids []int64) bool
	GetFriendList(id int64) []string
	GetCommonFriendList(id int64, withId int64) []string
	GetValidUsersCanReceiveUpdates(senderId int64, mentionIds []int64) []string
	CheckRelationshipTwoWay(requestUserId int64, targetUserId int64, status int64) []int64
	CheckRelationshipOneWay(requestUserId int64, targetUserId int64, status int64) []int64
}

type RelationshipRepository struct {
	DB *sql.DB
}

func (repo RelationshipRepository) GetFriendList(id int64) []string {
	query := `
		select u.email
		from user u inner join 
		(select TargetUserId id from relationship
		where RequestUserId =? and status = 1
		union
		select RequestUserId id from relationship
		where TargetUserId =? and status = 1) ids
		on u.id = ids.id;
	`

	rows, err := repo.DB.Query(query, id, id)
	if err != nil {
		fmt.Println(err)
	}

	var emails []string
	for rows.Next() {
		var email string
		rows.Scan(&email)
		emails = append(emails, email)
	}

	return emails
}

func (repo RelationshipRepository) GetCommonFriendList(id int64, withId int64) []string {
	query := `
	select u.email
	from user u inner join
	(select l.id from
	(select TargetUserId id from relationship
	where RequestUserId =? and status = 1
	union
	select RequestUserId id from relationship
	where TargetUserId =? and status = 1) l
	inner join
	(select TargetUserId id from relationship
	where RequestUserId =? and status = 1
	union
	select RequestUserId id from relationship
	where TargetUserId =? and status = 1) r
	on l.id = r.id) c
	on u.id = c.id;
	`

	rows, err := repo.DB.Query(query, id, id, withId, withId)
	if err != nil {
		fmt.Println(err)
	}

	var emails []string
	for rows.Next() {
		var email string
		rows.Scan(&email)
		emails = append(emails, email)
	}

	return emails
}

func (repo RelationshipRepository) CreateRelationship(relationship *models.Relationship) int64 {
	query := `
		INSERT INTO relationship (RequestUserId, TargetUserId, Status)
		VALUES (?,?,?)
	`

	rows, err := repo.DB.Prepare(query)
	if err != nil {
		return -1
	}

	res, err := rows.Exec(relationship.RequestUserId, relationship.TargetUserId, relationship.Status)
	if err != nil {
		return -1
	}

	insertedId, err := res.LastInsertId()

	if err != nil {
		return -1
	}

	return insertedId
}

func (repo RelationshipRepository) DeleteRelationships(ids []int64) bool {

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	stmt := `DELETE FROM relationship WHERE id in (?` + strings.Repeat(",?", len(args)-1) + `)`
	// rows, err := p.DB.Exec(stmt, args...)
	rows, err := repo.DB.Prepare(stmt)

	if err != nil {
		return false
	}

	rows.Exec(args...)

	return true
}

func (repo RelationshipRepository) CheckRelationshipTwoWay(requestUserId int64, targetUserId int64, status int64) []int64 {
	query := `
	SELECT id
	FROM relationship
	where (requestuserid =? and targetuserid =? and status =?)
	OR (targetuserid =? and requestuserid =? and status =?)
	`

	rows, err := repo.DB.Query(query, requestUserId, targetUserId, status, requestUserId, targetUserId, status)
	if err != nil {
		fmt.Println(err)
	}

	var ids []int64
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}

	return ids
}

func (repo RelationshipRepository) CheckRelationshipOneWay(requestUserId int64, targetUserId int64, status int64) []int64 {
	stmt := `
		SELECT id
		FROM relationship
		where requestuserid = %s and targetuserid = %s AND status = %s
	`

	query := fmt.Sprintf(
		stmt,
		strconv.FormatInt(requestUserId, 10),
		strconv.FormatInt(targetUserId, 10),
		strconv.FormatInt(status, 10))

	fmt.Printf(query)

	rows, _ := repo.DB.Query(query)

	var ids []int64
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}

	fmt.Println(len(ids))

	return ids
}

func (repo RelationshipRepository) GetValidUsersCanReceiveUpdates(senderId int64, mentionIds []int64) []string {
	strSenderId := strconv.FormatInt(senderId, 10)

	var stmt string
	var query string

	if len(mentionIds) > 0 {
		strMentionedIds := make([]string, len(mentionIds))
		for i, id := range mentionIds {
			strMentionedIds[i] = strconv.FormatInt(id, 10)
		}

		stmt = `
			select u.email from user u
			inner join (
			select rs.id from
			(select TargetUserId id from relationship
			where RequestUserId = %s and status = 1
			union
			select RequestUserId id from relationship
			where TargetUserId = %s and status in (1,2)
			union
			select id from user
			where id in (%s)) rs
			where rs.id not in (
			select RequestUserId id
			from relationship
			where RequestUserId in (%s)
			and TargetUserId = %s
			and status = 3
			)) ids on u.id = ids.id
		`

		query = fmt.Sprintf(stmt, strSenderId, strSenderId, strings.Join(strMentionedIds, ","), strings.Join(strMentionedIds, ","), strSenderId)
	} else {
		stmt = `
			select u.email from user u
			inner join (
			select TargetUserId id from relationship
			where RequestUserId = %s and status = 1
			union
			select RequestUserId id from relationship
			where TargetUserId = %s and status in (1,2)
			) ids on u.id = ids.id
		`
		query = fmt.Sprintf(stmt, strSenderId, strSenderId)
	}

	rows, err := repo.DB.Query(query)
	if err != nil {
		fmt.Println(err)
	}

	var emails []string
	for rows.Next() {
		var email string
		rows.Scan(&email)
		emails = append(emails, email)
	}

	return emails
}
