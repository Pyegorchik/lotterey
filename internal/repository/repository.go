package repository

import (
	"context"
	"database/sql"
	"lottery/internal/domain"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

type Repository interface {
	BeginTx(ctx context.Context) (Transaction, error)

	CreateDraw(ctx context.Context) (*domain.Draw, error)
	GetDrawByID(ctx context.Context, id int) (*domain.Draw, error)
	GetAll() ([]domain.Draw, error)
	Update(ctx context.Context, draw *domain.Draw) error
	HasActiveDraw(ctx context.Context) (bool, error)
	Close(ctx context.Context, id int, winningNumbers string) error

	CreateTicket(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error)
	GetTicketByDrawID(ctx context.Context, drawID int) ([]domain.Ticket, error)
}

type Transaction interface {
	Commit() error
	Rollback() error
}

func (r *repository) BeginTx(ctx context.Context) (Transaction, error) {
	t, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *repository) CreateDraw(ctx context.Context) (*domain.Draw, error) {
	result, err := r.db.ExecContext(ctx, "INSERT INTO draws (status) VALUES (?)", domain.DrawStatusActive)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return r.GetDrawByID(ctx, int(id))
}

func (r *repository) GetDrawByID(ctx context.Context, id int) (*domain.Draw, error) {
	var draw domain.Draw
	err := r.db.QueryRowContext(ctx,
		`SELECT id, status, COALESCE(winning_numbers, ''), created_at, COALESCE(closed_at, '')
		FROM draws WHERE id = ?`, id).
		Scan(&draw.ID, &draw.Status, &draw.WinningNumbers, &draw.CreatedAt, &draw.ClosedAt)

	if err != nil {
		return nil, err
	}
	return &draw, nil
}

func (r *repository) GetAll() ([]domain.Draw, error) {
	rows, err := r.db.Query(`
		SELECT id, status, COALESCE(winning_numbers, ''), created_at, COALESCE(closed_at, '')
		FROM draws ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var draws []domain.Draw
	for rows.Next() {
		var draw domain.Draw
		err := rows.Scan(&draw.ID, &draw.Status, &draw.WinningNumbers, &draw.CreatedAt, &draw.ClosedAt)
		if err != nil {
			continue
		}
		draws = append(draws, draw)
	}
	return draws, nil
}

func (r *repository) Update(ctx context.Context, draw *domain.Draw) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE draws SET status = ?, winning_numbers = ?, closed_at = ? WHERE id = ?`,
		draw.Status, draw.WinningNumbers, draw.ClosedAt, draw.ID)
	return err
}

func (r *repository) HasActiveDraw(ctx context.Context) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM draws WHERE status = ?", domain.DrawStatusActive).Scan(&count)
	return count > 0, err
}

func (r *repository) Close(ctx context.Context, id int, winningNumbers string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE draws SET status = ?, winning_numbers = ?, closed_at = CURRENT_TIMESTAMP WHERE id = ?`,
		domain.DrawStatusClosed, winningNumbers, id)
	return err
}

func (r *repository) CreateTicket(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error) {
	result, err := r.db.ExecContext(ctx, "INSERT INTO tickets (draw_id, numbers) VALUES (?, ?)",
		ticket.DrawID, ticket.Numbers)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	ticket.ID = int(id)
	return ticket, nil
}

func (r *repository) GetTicketByDrawID(ctx context.Context, drawID int) ([]domain.Ticket, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, draw_id, numbers FROM tickets WHERE draw_id = ?", drawID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []domain.Ticket
	for rows.Next() {
		var ticket domain.Ticket
		err := rows.Scan(&ticket.ID, &ticket.DrawID, &ticket.Numbers)
		if err != nil {
			continue
		}
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}
