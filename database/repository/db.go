// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.completeTaskStmt, err = db.PrepareContext(ctx, completeTask); err != nil {
		return nil, fmt.Errorf("error preparing query CompleteTask: %w", err)
	}
	if q.createTaskStmt, err = db.PrepareContext(ctx, createTask); err != nil {
		return nil, fmt.Errorf("error preparing query CreateTask: %w", err)
	}
	if q.createUserStmt, err = db.PrepareContext(ctx, createUser); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUser: %w", err)
	}
	if q.deleteTaskStmt, err = db.PrepareContext(ctx, deleteTask); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteTask: %w", err)
	}
	if q.getAllTasksStmt, err = db.PrepareContext(ctx, getAllTasks); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllTasks: %w", err)
	}
	if q.getAllUsersStmt, err = db.PrepareContext(ctx, getAllUsers); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllUsers: %w", err)
	}
	if q.getTaskWithPidStmt, err = db.PrepareContext(ctx, getTaskWithPid); err != nil {
		return nil, fmt.Errorf("error preparing query GetTaskWithPid: %w", err)
	}
	if q.getTaskWithTechIdStmt, err = db.PrepareContext(ctx, getTaskWithTechId); err != nil {
		return nil, fmt.Errorf("error preparing query GetTaskWithTechId: %w", err)
	}
	if q.getUserWithPidStmt, err = db.PrepareContext(ctx, getUserWithPid); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserWithPid: %w", err)
	}
	if q.getUserWithUsernameStmt, err = db.PrepareContext(ctx, getUserWithUsername); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserWithUsername: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.completeTaskStmt != nil {
		if cerr := q.completeTaskStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing completeTaskStmt: %w", cerr)
		}
	}
	if q.createTaskStmt != nil {
		if cerr := q.createTaskStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createTaskStmt: %w", cerr)
		}
	}
	if q.createUserStmt != nil {
		if cerr := q.createUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserStmt: %w", cerr)
		}
	}
	if q.deleteTaskStmt != nil {
		if cerr := q.deleteTaskStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteTaskStmt: %w", cerr)
		}
	}
	if q.getAllTasksStmt != nil {
		if cerr := q.getAllTasksStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllTasksStmt: %w", cerr)
		}
	}
	if q.getAllUsersStmt != nil {
		if cerr := q.getAllUsersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllUsersStmt: %w", cerr)
		}
	}
	if q.getTaskWithPidStmt != nil {
		if cerr := q.getTaskWithPidStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTaskWithPidStmt: %w", cerr)
		}
	}
	if q.getTaskWithTechIdStmt != nil {
		if cerr := q.getTaskWithTechIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTaskWithTechIdStmt: %w", cerr)
		}
	}
	if q.getUserWithPidStmt != nil {
		if cerr := q.getUserWithPidStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserWithPidStmt: %w", cerr)
		}
	}
	if q.getUserWithUsernameStmt != nil {
		if cerr := q.getUserWithUsernameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserWithUsernameStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                      DBTX
	tx                      *sql.Tx
	completeTaskStmt        *sql.Stmt
	createTaskStmt          *sql.Stmt
	createUserStmt          *sql.Stmt
	deleteTaskStmt          *sql.Stmt
	getAllTasksStmt         *sql.Stmt
	getAllUsersStmt         *sql.Stmt
	getTaskWithPidStmt      *sql.Stmt
	getTaskWithTechIdStmt   *sql.Stmt
	getUserWithPidStmt      *sql.Stmt
	getUserWithUsernameStmt *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                      tx,
		tx:                      tx,
		completeTaskStmt:        q.completeTaskStmt,
		createTaskStmt:          q.createTaskStmt,
		createUserStmt:          q.createUserStmt,
		deleteTaskStmt:          q.deleteTaskStmt,
		getAllTasksStmt:         q.getAllTasksStmt,
		getAllUsersStmt:         q.getAllUsersStmt,
		getTaskWithPidStmt:      q.getTaskWithPidStmt,
		getTaskWithTechIdStmt:   q.getTaskWithTechIdStmt,
		getUserWithPidStmt:      q.getUserWithPidStmt,
		getUserWithUsernameStmt: q.getUserWithUsernameStmt,
	}
}
