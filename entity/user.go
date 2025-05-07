package entity

import (
	"errors"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserModel struct {
	users  map[int]*User
	nextID int
	mu     sync.RWMutex
}

func NewUserModel() *UserModel {
	return &UserModel{
		users:  make(map[int]*User),
		nextID: 1,
	}
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (m *UserModel) Create(user *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, u := range m.users {
		if u.Username == user.Username {
			return errors.New("username already exists")
		}
	}

	user.ID = m.nextID
	m.nextID++
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := user.SetPassword(user.Password); err != nil {
		return err
	}

	m.users[user.ID] = user
	return nil
}

func (m *UserModel) GetByID(id int) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *UserModel) GetByUsername(username string) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *UserModel) GetAll() []*User {
	m.mu.RLock()
	defer m.mu.RUnlock()

	users := make([]*User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users
}

func (m *UserModel) Update(user *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	existingUser, exists := m.users[user.ID]
	if !exists {
		return errors.New("user not found")
	}

	if user.Username != existingUser.Username {
		for _, u := range m.users {
			if u.Username == user.Username {
				return errors.New("username already exists")
			}
		}
	}

	existingUser.Username = user.Username
	if user.Password != "" {
		existingUser.Password = user.Password
		if err := existingUser.SetPassword(user.Password); err != nil {
			return err
		}
	}
	existingUser.Role = user.Role
	existingUser.UpdatedAt = time.Now()

	return nil
}

func (m *UserModel) Delete(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(m.users, id)
	return nil
}
