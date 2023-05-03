package repositories

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"

	"golang.org/x/crypto/bcrypt"

	"github.com/reemployed/reemployed/models"
)

type UserRepository interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id string) error
}

type fileUserRepository struct {
	filename string
	mutex    sync.Mutex
}

func NewFileUserRepository(filename string) UserRepository {
	return &fileUserRepository{filename: filename}
}

func (repo *fileUserRepository) GetAllUsers() ([]models.User, error) {
	data, err := ioutil.ReadFile(repo.filename)
	if err != nil {
		return nil, err
	}
	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *fileUserRepository) GetUserByID(id string) (*models.User, error) {
	data, err := ioutil.ReadFile(repo.filename)
	if err != nil {
		return nil, err
	}
	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, fmt.Errorf("User not found")
}

func (repo *fileUserRepository) GetUserByEmail(email string) (*models.User, error) {
	data, err := ioutil.ReadFile(repo.filename)
	if err != nil {
		return nil, err
	}
	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, fmt.Errorf("User not found")
}

func (repo *fileUserRepository) CreateUser(user *models.User) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// Hash the user's password before saving it to the JSON file
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	data, err := ioutil.ReadFile(repo.filename)
	if err != nil {
		return err
	}
	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return err
	}
	user.ID = generateUserID(users)
	users = append(users, *user)
	data, err = json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(repo.filename, data, 0644)
}

func (repo *fileUserRepository) UpdateUser(user *models.User) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	data, err := ioutil.ReadFile(repo.filename)
	if err != nil {
		return err
	}
	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return err
	}
	for i, u := range users {
		if u.ID == user.ID {
			users[i] = *user
			break
		}
	}
	data, err = json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(repo.filename, data, 0644)
}

func (repo *fileUserRepository) DeleteUser(id string) error {
	data, err := ioutil.ReadFile(repo.filename)
	if err != nil {
		return err
	}
	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return err
	}
	for i, u := range users {
		if u.ID == id {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}
	data, err = json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(repo.filename, data, 0644)
}

func generateUserID(users []models.User) string {
	if len(users) == 0 {
		return "1"
	}
	lastUser := users[len(users)-1]
	newID, _ := strconv.Atoi(lastUser.ID)
	newID += 1
	return fmt.Sprintf("%d", newID)
}
