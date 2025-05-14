package store

import "sync"

type User struct {
	Name string `json:"name"`
	ID int 		`json:"id"`
}

type Store struct {
	mu sync.Mutex
	users []User
	next int
}

func New() *Store {
	return &Store{next: 1}
}

func (s *Store) All() []User {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]User, len(s.users))

	copy(out, s.users)
	return out
}

func (s *Store) Add(name string) User{
	s.mu.Lock()
	defer s.mu.Unlock()

	u := User{Name: name, ID: s.next}

	s.next++
	s.users = append(s.users, u)
	return u
}

func (s *Store) Get(id int) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, u := range s.users {
		if u.ID == id {
			return u, true
		}
	}

	return User{}, false
}