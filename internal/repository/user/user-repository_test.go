package user_repository_test

import (
	"fmt"
	"testing"

	"github.com/Kamegatze/file-manager-go/internal/configuration"
	"github.com/Kamegatze/file-manager-go/internal/entity"
	"github.com/Kamegatze/file-manager-go/internal/repository"
	user_repository "github.com/Kamegatze/file-manager-go/internal/repository/user"
	utility_test "github.com/Kamegatze/file-manager-go/internal/utility/test"
	pointer_utility "github.com/Kamegatze/file-manager-go/pkg/pointer-utility"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/assert"
)

var repo user_repository.UserRepository = user_repository.NewUserRepository()

func RunContainer(t *testing.T) configuration.DatasourceConfig {
	return utility_test.RunContainer(t, "file:../../../migrations")
}

func TestInsert(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := repo.Insert(&entity.User{
		LastName:  pointer_utility.NewPointer("sh"),
		FirstName: pointer_utility.NewPointer("al"),
		Username:  pointer_utility.NewPointer("kam"),
	})

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "sh", *user.LastName)
	assert.Equal(t, "al", *user.FirstName)
	assert.Equal(t, "kam", *user.Username)
	assert.NotNil(t, user.CreatedAt)
	assert.Nil(t, user.UpdatedAt)
}

func TestInsertError(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	abstractRepo := repository.NewGeneralRepository[entity.User, uuid.UUID](sqlbuilder.MySQL)

	repoCustom := user_repository.NewUserWithAbstractRepo(
		abstractRepo.WithTableName("users").WithCallBack(user_repository.UserRowMapper))

	user, err := repoCustom.Insert(&entity.User{
		LastName:  pointer_utility.NewPointer("sh"),
		FirstName: pointer_utility.NewPointer("al"),
		Username:  pointer_utility.NewPointer("kam")})

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.True(t, assert.Error(t, err))
}

func TestGetAll(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	_, err := repo.Insert(&entity.User{
		LastName:  pointer_utility.NewPointer("sh"),
		FirstName: pointer_utility.NewPointer("al"),
		Username:  pointer_utility.NewPointer("kam")})

	if err != nil {
		t.Fatalf("error insert record, error: %s", err)
	}

	users, err := repo.GetAll()
	if err != nil {
		t.Fatalf("error get all record, error: %s", err)
	}

	assert.Len(t, users, 1)
}

func TestGetAllPageable(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	for i := 0; i < 15; i++ {
		_, err := repo.Insert(&entity.User{
			LastName:  pointer_utility.NewPointer("sh"),
			FirstName: pointer_utility.NewPointer("al"),
			Username:  pointer_utility.NewPointer("kam")})

		if err != nil {
			t.Fatalf("error insert record, error: %s", err)
		}
	}
	pageableFirst, err := repository.NewPageableImpl(1, 10, []string{"last_name"})

	if err != nil {
		t.Fatal(err)
	}

	usersFirst, err := repo.GetAllPageable(pageableFirst)

	if err != nil {
		t.Fatalf("error get all record via pageable with error: %s", err)
	}

	pageableLast, err := repository.NewPageableImpl(2, 10, []string{"last_name"})

	if err != nil {
		t.Fatal(err)
	}

	usersLast, err := repo.GetAllPageable(pageableLast)

	if err != nil {
		t.Fatalf("error get all record via pageable with error: %s", err)
	}

	assert.Len(t, usersFirst, 10)
	assert.Len(t, usersLast, 5)
}

func TestGetById(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := repo.Insert(&entity.User{
		LastName:  pointer_utility.NewPointer("sh"),
		FirstName: pointer_utility.NewPointer("al"),
		Username:  pointer_utility.NewPointer("kam")})

	if err != nil {
		t.Fatalf("error insert record, error: %s", err)
	}

	var id *uuid.UUID = user.Id

	userByGetById, err := repo.GetById(id)
	if err != nil {
		t.Fatalf("fatal get entity by id: %s, with error: %s", id.String(), err)
	}

	assert.Equal(t, *id, *userByGetById.Id)
	assert.Equal(t, *user.LastName, *userByGetById.LastName)
	assert.Equal(t, *user.FirstName, *userByGetById.FirstName)
	assert.Equal(t, *user.Username, *userByGetById.Username)
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
	assert.Equal(t, fmt.Errorf("not found entity by id: #%v in table: %s", id, "users"), err)
}

func TestUpdate(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := repo.Insert(&entity.User{
		LastName:  pointer_utility.NewPointer("sh"),
		FirstName: pointer_utility.NewPointer("al"),
		Username:  pointer_utility.NewPointer("kam")})

	if err != nil {
		t.Fatalf("error insert record, error: %s", err)
	}

	user.FirstName = pointer_utility.NewPointer("ivan")

	userUpdate, err := repo.Update(user)

	if err != nil {
		t.Fatalf("fatal update entity by id: %s, error: %s", user.Id.String(), err)
	}

	assert.Equal(t, *user.Id, *userUpdate.Id)
	assert.Equal(t, "ivan", *userUpdate.FirstName)
	assert.Equal(t, *user.LastName, *userUpdate.LastName)
	assert.Equal(t, *user.Username, *userUpdate.Username)
	assert.NotNil(t, userUpdate.UpdatedAt)
}

func TestUpdateError(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := repo.Update(&entity.User{
		LastName:  pointer_utility.NewPointer("sh"),
		FirstName: pointer_utility.NewPointer("al"),
		Username:  pointer_utility.NewPointer("kam")})

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("error update in table %s by id: %v", "users", nil), err)

}

func TestDeleteById(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := repo.Insert(&entity.User{
		LastName:  pointer_utility.NewPointer("sh"),
		FirstName: pointer_utility.NewPointer("al"),
		Username:  pointer_utility.NewPointer("kam")})

	if err != nil {
		t.Fatalf("error insert record, error: %s", err)
	}

	id, err := repo.DeleteById(user.Id)

	if err != nil {
		t.Fatalf("error deleted by id: %s, error: %s", user.Id.String(), err)
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

	user, err := repo.Insert(&entity.User{
		LastName:  pointer_utility.NewPointer("sh"),
		FirstName: pointer_utility.NewPointer("al"),
		Username:  pointer_utility.NewPointer("kam")})

	if err != nil {
		t.Fatalf("error insert record, error: %s", err)
	}

	id, err := repo.Delete(user)

	if err != nil {
		t.Fatalf("error deleted by id: %s, error: %s", user.Id.String(), err)
	}

	entity, _ := repo.GetById(id)

	assert.Nil(t, entity)
}

func TestDeleteError(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user := entity.User{
		LastName:  pointer_utility.NewPointer("sh"),
		FirstName: pointer_utility.NewPointer("al"),
		Username:  pointer_utility.NewPointer("kam")}

	id, err := repo.Delete(&user)

	assert.Nil(t, id)
	assert.NotNil(t, err)
	assert.True(t, assert.Error(t, err))
}
