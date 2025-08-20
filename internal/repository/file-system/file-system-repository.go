package file_system_repository

import (
	"database/sql"
	"time"

	"github.com/Kamegatze/file-manager-go/internal/entity"
	"github.com/Kamegatze/file-manager-go/internal/repository"
	pointer_utility "github.com/Kamegatze/file-manager-go/pkg/pointer-utility"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
)

type FileSystemRepository interface {
	GetAll() ([]entity.FileSystem, error)
	GetAllPageable(pageable repository.Pageable) ([]entity.FileSystem, error)
	GetById(id *uuid.UUID) (*entity.FileSystem, error)
	Insert(entity *entity.FileSystem) (*entity.FileSystem, error)
	Update(entity *entity.FileSystem) (*entity.FileSystem, error)
	Delete(entity *entity.FileSystem) (*uuid.UUID, error)
	DeleteById(id *uuid.UUID) (*uuid.UUID, error)
}

type FileSystemRepositoryImpl struct {
	abstractRepo repository.Repository[entity.FileSystem, uuid.UUID]
}

func NewFileSystemRepository() FileSystemRepositoryImpl {
	var abstractRepo repository.GeneralRepository[entity.FileSystem, uuid.UUID] = repository.NewGeneralRepository[entity.FileSystem, uuid.UUID](sqlbuilder.PostgreSQL)
	abstractRepo.WithCallBack(FileSystemRowMapper).WithTableName("file_system")
	repo := FileSystemRepositoryImpl{abstractRepo: abstractRepo}
	return repo
}

func NewFileSystemRepositoryWithAbstractRepo(abstractRepo repository.Repository[entity.FileSystem, uuid.UUID]) FileSystemRepositoryImpl {
	return FileSystemRepositoryImpl{abstractRepo: abstractRepo}
}

func FileSystemRowMapper(rows *sql.Rows) (entity.FileSystem, error) {
	fileSystem := entity.FileSystem{}
	if err := rows.Scan(
		&fileSystem.Id,
		&fileSystem.OwnerId,
		&fileSystem.ParentId,
		&fileSystem.Rights,
		&fileSystem.IsFile,
		&fileSystem.Name,
		&fileSystem.Path,
		&fileSystem.CreatedAt,
		&fileSystem.UpdatedAt,
		&fileSystem.Deleted); err != nil {
		return fileSystem, err
	}
	return fileSystem, nil
}

func (repo *FileSystemRepositoryImpl) WithAbstractRepo(abstractRepo repository.Repository[entity.FileSystem, uuid.UUID]) *FileSystemRepositoryImpl {
	repo.abstractRepo = abstractRepo
	return repo
}

func (repo FileSystemRepositoryImpl) GetAll() ([]entity.FileSystem, error) {
	return repo.abstractRepo.GetAll()
}

func (repo FileSystemRepositoryImpl) GetAllPageable(pageable repository.Pageable) ([]entity.FileSystem, error) {
	return repo.abstractRepo.GetAllPageable(pageable)
}

func (repo FileSystemRepositoryImpl) GetById(id *uuid.UUID) (*entity.FileSystem, error) {
	return repo.abstractRepo.GetById(id)
}

func (repo FileSystemRepositoryImpl) Insert(entity *entity.FileSystem) (*entity.FileSystem, error) {
	return repo.abstractRepo.Insert(entity)
}

func (repo FileSystemRepositoryImpl) Update(entity *entity.FileSystem) (*entity.FileSystem, error) {
	entity.UpdatedAt = pointer_utility.NewPointer(time.Now())
	return repo.abstractRepo.Update(entity)
}

func (repo FileSystemRepositoryImpl) Delete(entity *entity.FileSystem) (*uuid.UUID, error) {
	return repo.abstractRepo.Delete(entity)
}

func (repo FileSystemRepositoryImpl) DeleteById(id *uuid.UUID) (*uuid.UUID, error) {
	return repo.abstractRepo.DeleteById(id)
}
