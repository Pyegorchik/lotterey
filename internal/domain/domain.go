package domain

type Draw struct {
	ID             int    `json:"id"`
	Status         string `json:"status"`
	WinningNumbers string `json:"winning_numbers,omitempty"`
	CreatedAt      string `json:"created_at"`
	ClosedAt       string `json:"closed_at,omitempty"`
}

type Ticket struct {
	ID       int    `json:"id"`
	DrawID   int    `json:"draw_id"`
	Numbers  string `json:"numbers"`
	IsWinner bool   `json:"is_winner,omitempty"`
}

type TicketRequest struct {
	Numbers []int `json:"numbers"`
	DrawID  int   `json:"draw_id"`
}

type DrawResults struct {
	Draw           Draw     `json:"draw"`
	WinningNumbers []int    `json:"winning_numbers"`
	Tickets        []Ticket `json:"tickets"`
	WinnerCount    int      `json:"winner_count"`
}

// DrawStatus константы для статусов тиража
const (
	DrawStatusActive = "active"
	DrawStatusClosed = "closed"
)
