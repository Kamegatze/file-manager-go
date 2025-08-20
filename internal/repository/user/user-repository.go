package user_repository

import (
	"database/sql"

	"time"

	"github.com/Kamegatze/file-manager-go/internal/entity"
	"github.com/Kamegatze/file-manager-go/internal/repository"
	pointer_utility "github.com/Kamegatze/file-manager-go/pkg/pointer-utility"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
)

type UserRepository interface {
	GetAll() ([]entity.User, error)
	GetAllPageable(pageable repository.Pageable) ([]entity.User, error)
	GetById(id *uuid.UUID) (*entity.User, error)
	Insert(entity *entity.User) (*entity.User, error)
	Update(entity *entity.User) (*entity.User, error)
	Delete(entity *entity.User) (*uuid.UUID, error)
	DeleteById(id *uuid.UUID) (*uuid.UUID, error)
}

type UserRepositoryImpl struct {
	abstractRepo repository.Repository[entity.User, uuid.UUID]
}

func NewUserRepository() UserRepositoryImpl {
	var abstractRepo repository.GeneralRepository[entity.User, uuid.UUID] = repository.NewGeneralRepository[entity.User, uuid.UUID](sqlbuilder.PostgreSQL)
	abstractRepo.WithCallBack(UserRowMapper).WithTableName("users")
	repo := UserRepositoryImpl{abstractRepo: abstractRepo}
	return repo
}

func NewUserWithAbstractRepo(abstractRepo repository.Repository[entity.User, uuid.UUID]) UserRepositoryImpl {
	return UserRepositoryImpl{abstractRepo: abstractRepo}
}

func UserRowMapper(rows *sql.Rows) (entity.User, error) {
	user := entity.User{}
	if err := rows.Scan(
		&user.Id,
		&user.LastName,
		&user.FirstName,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt); err != nil {
		return user, err
	}
	return user, nil
}

func (repo UserRepositoryImpl) GetAll() ([]entity.User, error) {
	return repo.abstractRepo.GetAll()
}

func (repo UserRepositoryImpl) GetAllPageable(pageable repository.Pageable) ([]entity.User, error) {
	return repo.abstractRepo.GetAllPageable(pageable)
}

func (repo UserRepositoryImpl) GetById(id *uuid.UUID) (*entity.User, error) {
	return repo.abstractRepo.GetById(id)
}

func (repo UserRepositoryImpl) Insert(entity *entity.User) (*entity.User, error) {
	return repo.abstractRepo.Insert(entity)
}

func (repo UserRepositoryImpl) Update(entity *entity.User) (*entity.User, error) {
	entity.UpdatedAt = pointer_utility.NewPointer(time.Now())
	return repo.abstractRepo.Update(entity)
}

func (repo UserRepositoryImpl) Delete(entity *entity.User) (*uuid.UUID, error) {
	return repo.abstractRepo.Delete(entity)
}

func (repo UserRepositoryImpl) DeleteById(id *uuid.UUID) (*uuid.UUID, error) {
	return repo.abstractRepo.DeleteById(id)
}
