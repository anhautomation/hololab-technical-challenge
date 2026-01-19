package sqlite

import (
	"context"
	"database/sql"
	"hololab-core-system/internal/domain"
)

type MessageRepo struct{ DB *sql.DB }

func (r MessageRepo) ListByBot(ctx context.Context, botID string) ([]domain.Message, error) {
	rows, err := r.DB.QueryContext(ctx, `
SELECT role,content,created_at
FROM messages
WHERE bot_id=?
ORDER BY created_at ASC, id ASC`, botID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Message
	for rows.Next() {
		var role string
		var m domain.Message
		_ = rows.Scan(&role, &m.Content, &m.CreatedAt)
		m.BotID = botID
		m.Role = domain.Role(role)
		out = append(out, m)
	}
	return out, nil
}

func (r MessageRepo) Append(ctx context.Context, m domain.Message) error {
	_, err := r.DB.ExecContext(ctx, `
INSERT INTO messages (bot_id, role, content, created_at)
VALUES (?,?,?,?)`,
		m.BotID, string(m.Role), m.Content, m.CreatedAt)
	return err
}

func (r MessageRepo) ResetByBot(ctx context.Context, botID string) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM messages WHERE bot_id=?`, botID)
	return err
}
