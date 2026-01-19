package sqlite

import (
	"context"
	"database/sql"
	"hololab-core-system/internal/domain"
)

type BotRepo struct{ DB *sql.DB }

func (r BotRepo) List(ctx context.Context) ([]domain.Bot, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id,name,job,bio,style,knowledge,created_at FROM bots ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.Bot, 0)

	for rows.Next() {
		var b domain.Bot
		if err := rows.Scan(&b.ID, &b.Name, &b.Job, &b.Bio, &b.Style, &b.Knowledge, &b.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r BotRepo) Get(ctx context.Context, id string) (domain.Bot, error) {
	var b domain.Bot
	err := r.DB.QueryRowContext(ctx, `SELECT id,name,job,bio,style,knowledge,created_at FROM bots WHERE id=?`, id).
		Scan(&b.ID, &b.Name, &b.Job, &b.Bio, &b.Style, &b.Knowledge, &b.CreatedAt)
	if err == sql.ErrNoRows {
		return domain.Bot{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Bot{}, err
	}
	return b, nil
}

func (r BotRepo) Create(ctx context.Context, b domain.Bot) error {
	_, err := r.DB.ExecContext(ctx, `
INSERT INTO bots (id,name,job,bio,style,knowledge,created_at)
VALUES (?,?,?,?,?,?,?)`,
		b.ID, b.Name, b.Job, b.Bio, b.Style, b.Knowledge, b.CreatedAt)
	return err
}

func (r BotRepo) Delete(ctx context.Context, id string) (int64, error) {
	_, _ = r.DB.ExecContext(ctx, `DELETE FROM messages WHERE bot_id=?`, id)
	res, err := r.DB.ExecContext(ctx, `DELETE FROM bots WHERE id=?`, id)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}
