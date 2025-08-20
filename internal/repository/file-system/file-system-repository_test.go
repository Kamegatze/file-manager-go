package file_system_repository_test

import (
	"fmt"
	"testing"

	"github.com/Kamegatze/file-manager-go/internal/configuration"
	"github.com/Kamegatze/file-manager-go/internal/entity"
	"github.com/Kamegatze/file-manager-go/internal/repository"
	file_system_repository "github.com/Kamegatze/file-manager-go/internal/repository/file-system"
	user_repository "github.com/Kamegatze/file-manager-go/internal/repository/user"
	utility_test "github.com/Kamegatze/file-manager-go/internal/utility/test"
	pointer_utility "github.com/Kamegatze/file-manager-go/pkg/pointer-utility"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/assert"
)

var repo file_system_repository.FileSystemRepository = file_system_repository.NewFileSystemRepository()

func RunContainer(t *testing.T) configuration.DatasourceConfig {
	return utility_test.RunContainer(t, "file:../../../migrations")
}

func PreTestConfiguration(t *testing.T) *uuid.UUID {
	userRepository := repository.NewGeneralRepository[entity.User, uuid.UUID](sqlbuilder.PostgreSQL)
	userRepository.WithCallBack(user_repository.UserRowMapper).WithTableName("users")
	user, err := userRepository.Insert(&entity.User{
		LastName:  pointer_utility.NewPointer("sh"),
		FirstName: pointer_utility.NewPointer("al"),
		Username:  pointer_utility.NewPointer("kam")})

	if err != nil {
		t.Fatal(err)
	}

	return user.Id
}

func TestInsert(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	userId := PreTestConfiguration(t)

	fileSystem, err := repo.Insert(&entity.FileSystem{
		OwnerId: userId,
		IsFile:  pointer_utility.NewPointer(false),
		Name:    pointer_utility.NewPointer("root"),
		Path:    pointer_utility.NewPointer("/")})

	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, fileSystem.Id)
	assert.Equal(t, *userId, *fileSystem.OwnerId)
	assert.Nil(t, fileSystem.ParentId)
	assert.Equal(t, "rw------", *fileSystem.Rights)
	assert.Equal(t, false, *fileSystem.IsFile)
	assert.Equal(t, "root", *fileSystem.Name)
	assert.Equal(t, "/", *fileSystem.Path)
	assert.NotNil(t, fileSystem.CreatedAt)
	assert.Nil(t, fileSystem.UpdatedAt)
	assert.Equal(t, false, *fileSystem.Deleted)
}

func TestInserError(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	userId := PreTestConfiguration(t)

	abstractRepo := repository.NewGeneralRepository[entity.FileSystem, uuid.UUID](sqlbuilder.MySQL)

	repoCustom := file_system_repository.NewFileSystemRepositoryWithAbstractRepo(
		abstractRepo.WithTableName("file_system").WithCallBack(file_system_repository.FileSystemRowMapper),
	)

	fileSystem, err := repoCustom.Insert(&entity.FileSystem{
		OwnerId: userId,
		IsFile:  pointer_utility.NewPointer(false),
		Name:    pointer_utility.NewPointer("root"),
		Path:    pointer_utility.NewPointer("/")})

	assert.Nil(t, fileSystem)
	assert.NotNil(t, err)
	assert.True(t, assert.Error(t, err))
}

func TestGetAll(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	userId := PreTestConfiguration(t)

	_, err := repo.Insert(&entity.FileSystem{
		OwnerId: userId,
		IsFile:  pointer_utility.NewPointer(false),
		Name:    pointer_utility.NewPointer("root"),
		Path:    pointer_utility.NewPointer("/")})

	if err != nil {
		t.Fatal(err)
	}

	fileSystems, err := repo.GetAll()

	assert.Nil(t, err)
	assert.Len(t, fileSystems, 1)
}

func TestGetById(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	userId := PreTestConfiguration(t)

	fileSystem, err := repo.Insert(&entity.FileSystem{
		OwnerId: userId,
		IsFile:  pointer_utility.NewPointer(false),
		Name:    pointer_utility.NewPointer("root"),
		Path:    pointer_utility.NewPointer("/")})

	if err != nil {
		t.Fatal(err)
	}

	fileSystemGetById, err := repo.GetById(fileSystem.Id)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, *fileSystem.Id, *fileSystemGetById.Id)
	assert.Equal(t, *fileSystem.OwnerId, *fileSystemGetById.OwnerId)
	assert.Equal(t, *fileSystem.IsFile, *fileSystemGetById.IsFile)
	assert.Equal(t, *fileSystem.Name, *fileSystemGetById.Name)
	assert.Equal(t, *fileSystem.Path, *fileSystemGetById.Path)
	assert.Equal(t, *fileSystem.Rights, *fileSystemGetById.Rights)
	assert.Equal(t, *fileSystem.CreatedAt, *fileSystemGetById.CreatedAt)
	assert.Equal(t, *fileSystem.Deleted, *fileSystemGetById.Deleted)
	assert.Nil(t, fileSystemGetById.ParentId)
	assert.Nil(t, fileSystemGetById.UpdatedAt)
}

func TestGetByIdError(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	id := uuid.New()
	user, err := repo.GetById(&id)

	if err == nil {
		t.Fatalf("founding user by id: %s", id.String())
	}

	assert.Nil(t, user)
	assert.Equal(t, fmt.Errorf("not found entity by id: #%v in table: %s", id, "file_system"), err)
}

func TestUpdate(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	userId := PreTestConfiguration(t)

	fileSystem, err := repo.Insert(&entity.FileSystem{
		OwnerId: userId,
		IsFile:  pointer_utility.NewPointer(false),
		Name:    pointer_utility.NewPointer("root"),
		Path:    pointer_utility.NewPointer("/")})

	if err != nil {
		t.Fatal(err)
	}

	fileSystem.Name = pointer_utility.NewPointer("root_user")

	fileSystemUpdate, err := repo.Update(fileSystem)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, *fileSystem.Id, *fileSystemUpdate.Id)
	assert.Equal(t, *fileSystem.OwnerId, *fileSystemUpdate.OwnerId)
	assert.Equal(t, *fileSystem.IsFile, *fileSystemUpdate.IsFile)
	assert.Equal(t, "root_user", *fileSystemUpdate.Name)
	assert.Equal(t, *fileSystem.Path, *fileSystemUpdate.Path)
	assert.Equal(t, *fileSystem.Rights, *fileSystemUpdate.Rights)
	assert.Equal(t, *fileSystem.CreatedAt, *fileSystemUpdate.CreatedAt)
	assert.Equal(t, *fileSystem.Deleted, *fileSystemUpdate.Deleted)
	assert.Nil(t, fileSystemUpdate.ParentId)
	assert.NotNil(t, fileSystemUpdate.UpdatedAt)
}

func TestUpdateError(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	userId := PreTestConfiguration(t)

	user, err := repo.Update(&entity.FileSystem{
		OwnerId: userId,
		IsFile:  pointer_utility.NewPointer(false),
		Name:    pointer_utility.NewPointer("root"),
		Path:    pointer_utility.NewPointer("/")})

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("error update in table %s by id: %v", "file_system", nil), err)

}

func TestDeleteById(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	userId := PreTestConfiguration(t)

	fileSystem, err := repo.Insert(&entity.FileSystem{
		OwnerId: userId,
		IsFile:  pointer_utility.NewPointer(false),
		Name:    pointer_utility.NewPointer("root"),
		Path:    pointer_utility.NewPointer("/")})

	if err != nil {
		t.Fatal(err)
	}

	id, err := repo.DeleteById(fileSystem.Id)

	if err != nil {
		t.Fatal(err)
	}

	entity, _ := repo.GetById(id)

	assert.Nil(t, entity)
}

func TestDeleteByIdError(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	id := uuid.New()

	returnId, err := repo.DeleteById(&id)

	assert.Nil(t, returnId)
	assert.NotNil(t, err)
	assert.True(t, assert.Error(t, err))
}

func TestDelete(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	userId := PreTestConfiguration(t)

	fileSystem, err := repo.Insert(&entity.FileSystem{
		OwnerId: userId,
		IsFile:  pointer_utility.NewPointer(false),
		Name:    pointer_utility.NewPointer("root"),
		Path:    pointer_utility.NewPointer("/")})

	if err != nil {
		t.Fatal(err)
	}

	id, err := repo.Delete(fileSystem)

	if err != nil {
		t.Fatal(err)
	}

	entity, _ := repo.GetById(id)

	assert.Nil(t, entity)
}

func TestDeleteError(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	fileSystem := &entity.FileSystem{
		IsFile: pointer_utility.NewPointer(false),
		Name:   pointer_utility.NewPointer("root"),
		Path:   pointer_utility.NewPointer("/")}

	id, err := repo.Delete(fileSystem)

	assert.Nil(t, id)
	assert.NotNil(t, err)
	assert.True(t, assert.Error(t, err))
}
