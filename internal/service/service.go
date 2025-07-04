package service

import (
	"context"
	"database/sql"
	"fmt"
	"lottery/internal/domain"
	"lottery/internal/repository"
)

type LotteryService interface {
	CreateDraw(ctx context.Context) (*domain.Draw, error)
	CreateTicket(ctx context.Context, req *domain.TicketRequest) (*domain.Ticket, error)
	CloseDraw(ctx context.Context, drawID int) (*domain.Draw, error)
	GetDrawResults(ctx context.Context, drawID int) (*domain.DrawResults, error)
	GetDraw(ctx context.Context, drawID int) (*domain.Draw, error)
	GetAllDraws() ([]domain.Draw, error)
}

type lotteryService struct {
	repo repository.Repository
}

func NewLotteryService(repo repository.Repository) LotteryService {
	return &lotteryService{
		repo: repo,
	}
}

func (s *lotteryService) CreateDraw(ctx context.Context) (*domain.Draw, error) {
	// Проверяем, есть ли активный тираж
	hasActive, err := s.repo.HasActiveDraw(ctx)
	if err != nil {
		return nil, err
	}
	if hasActive {
		return nil, &BusinessError{Message: "There is already an active draw"}
	}

	return s.repo.CreateDraw(ctx)
}

func (s *lotteryService) CreateTicket(ctx context.Context, req *domain.TicketRequest) (*domain.Ticket, error) {
	// Валидация чисел
	if err := s.validateTicketNumbers(req.Numbers); err != nil {
		return nil, err
	}

	// Проверка существования и статуса тиража
	draw, err := s.repo.GetDrawByID(ctx, req.DrawID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &NotFoundError{Message: "Draw not found"}
		}
		return nil, err
	}

	if draw.Status != domain.DrawStatusActive {
		return nil, &BusinessError{Message: "Draw is not active"}
	}

	// Создание билета
	ticket := &domain.Ticket{
		DrawID:  req.DrawID,
		Numbers: s.formatNumbers(req.Numbers),
	}

	return s.repo.CreateTicket(ctx, ticket)
}

func (s *lotteryService) CloseDraw(ctx context.Context, drawID int) (*domain.Draw, error) {
	// Проверяем статус тиража
	draw, err := s.repo.GetDrawByID(ctx, drawID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &NotFoundError{Message: "Draw not found"}
		}
		return nil, err
	}

	if draw.Status != domain.DrawStatusActive {
		return nil, &BusinessError{Message: "Draw is not active"}
	}

	// Генерируем выигрышные числа
	winningNumbers := s.generateWinningNumbers()
	winningNumbersStr := s.formatNumbers(winningNumbers)

	// Закрываем тираж
	err = s.repo.Close(ctx, drawID, winningNumbersStr)
	if err != nil {
		return nil, err
	}

	return s.repo.GetDrawByID(ctx, drawID)
}

func (s *lotteryService) GetDrawResults(ctx context.Context, drawID int) (*domain.DrawResults, error) {
	draw, err := s.repo.GetDrawByID(ctx, drawID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &NotFoundError{Message: "Draw not found"}
		}
		return nil, err
	}

	tickets, err := s.repo.GetTicketByDrawID(ctx, drawID)
	if err != nil {
		return nil, err
	}

	results := &domain.DrawResults{
		Draw:    *draw,
		Tickets: tickets,
	}

	// Если тираж закрыт, обрабатываем результаты
	if draw.Status == domain.DrawStatusClosed && draw.WinningNumbers != "" {
		results.WinningNumbers = s.parseNumbers(draw.WinningNumbers)

		// Определяем победителей
		winnerCount := 0
		for i := range results.Tickets {
			results.Tickets[i].IsWinner = results.Tickets[i].Numbers == draw.WinningNumbers
			if results.Tickets[i].IsWinner {
				winnerCount++
			}
		}
		results.WinnerCount = winnerCount
	}

	return results, nil
}

func (s *lotteryService) GetDraw(ctx context.Context, drawID int) (*domain.Draw, error) {
	draw, err := s.repo.GetDrawByID(ctx, drawID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &NotFoundError{Message: "Draw not found"}
		}
		return nil, err
	}
	return draw, nil
}

func (s *lotteryService) GetAllDraws() ([]domain.Draw, error) {
	return s.repo.GetAll()
}

// Вспомогательные методы сервиса
func (s *lotteryService) validateTicketNumbers(numbers []int) error {
	if len(numbers) != 5 {
		return &ValidationError{Message: "Ticket must contain exactly 5 numbers"}
	}

	seen := make(map[int]bool)
	for _, num := range numbers {
		if num < 1 || num > 36 {
			return &ValidationError{Message: "Numbers must be between 1 and 36"}
		}
		if seen[num] {
			return &ValidationError{Message: "Numbers must be unique"}
		}
		seen[num] = true
	}
	return nil
}

func (s *lotteryService) formatNumbers(numbers []int) string {
	// Сортируем числа
	sorted := make([]int, len(numbers))
	copy(sorted, numbers)

	// Простая сортировка пузырьком
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// Форматируем в строку
	result := ""
	for i, num := range sorted {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%d", num)
	}
	return result
}

func (s *lotteryService) parseNumbers(numbersStr string) []int {
	if numbersStr == "" {
		return []int{}
	}

	var numbers []int
	currentNum := ""

	for _, char := range numbersStr {
		if char == ',' {
			if currentNum != "" {
				if num := parseIntSimple(currentNum); num > 0 {
					numbers = append(numbers, num)
				}
				currentNum = ""
			}
		} else {
			currentNum += string(char)
		}
	}

	// Добавляем последнее число
	if currentNum != "" {
		if num := parseIntSimple(currentNum); num > 0 {
			numbers = append(numbers, num)
		}
	}

	return numbers
}

func parseIntSimple(s string) int {
	num := 0
	for _, char := range s {
		if char >= '0' && char <= '9' {
			num = num*10 + int(char-'0')
		} else {
			return 0
		}
	}
	return num
}

func (s *lotteryService) generateWinningNumbers() []int {
	// Используем простой псевдослучайный генератор
	seed := int64(1)
	for i := 0; i < 10; i++ {
		seed = (seed*1103515245 + 12345) % (1 << 31)
	}

	var numbers []int
	used := make(map[int]bool)

	for len(numbers) < 5 {
		seed = (seed*1103515245 + 12345) % (1 << 31)
		num := int(seed%36) + 1
		if !used[num] {
			numbers = append(numbers, num)
			used[num] = true
		}
	}

	return numbers
}
