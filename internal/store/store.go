package store

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Name string 		`json:"name"`
	ID string			`json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type Store interface {
	All(ctx context.Context) ([]User, error)
	Create(ctx context.Context, u User) (User, error)
	Get(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, id string, u User) (User, error)
	Delete(ctx context.Context, id string) error
}

type inMemoryStore struct {
	mu sync.Mutex
	users []User
}

var ErrNotFound = errors.New("user not found")

func NewInMemoryStore() *inMemoryStore {
	return &inMemoryStore{
		users: make([]User, 0),
	}
}

func (s *inMemoryStore) All(ctx context.Context) ([]User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <- ctx.Done():
		return nil, ctx.Err()
	default:
	}

	out := make([]User, len(s.users))
	copy(out, s.users)
	return out, nil
}

func (s *inMemoryStore) Create(ctx context.Context, u User) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		return User{}, ctx.Err()
	default:
	}

	u.ID = uuid.NewString()
	u.CreatedAt = time.Now()
	s.users = append(s.users, u)

	return u, nil
}

func (s *inMemoryStore) Get(ctx context.Context, id string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		return User{}, ctx.Err()
	default:
	}

	for _, u := range s.users {
		if u.ID == id {
			return u, nil
		}
	}

	return User{}, ErrNotFound
}

func (s *inMemoryStore) Update(ctx context.Context, id string, u User) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		return User{}, ctx.Err()
	default:
	}
	for i, existing := range s.users {
		if existing.ID == id {
			u.ID = id
			u.CreatedAt = existing.CreatedAt
			s.users[i] = u
			return u, nil
		}
	}
	return User{}, ErrNotFound
}

func (s *inMemoryStore) Delete(ctx context.Context, id string) (error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	for i, u := range s.users {
		if u.ID == id {
			s.users = append(s.users[:i], s.users[i + 1:]...)
			return nil
		}
	}

	return ErrNotFound
}