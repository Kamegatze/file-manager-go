package repository_test

import (
	"errors"
	"file-manager/internal/configuration"
	"file-manager/internal/entity"
	"file-manager/internal/repository"
	user_repository "file-manager/internal/repository/user"
	utility_test "file-manager/internal/utility/test"
	pointer_utility "file-manager/pkg/pointer-utility"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/assert"
)

var generalRepository repository.Repository[entity.User, uuid.UUID] = createGeneralRepository(sqlbuilder.PostgreSQL)

func createGeneralRepository(flavor sqlbuilder.Flavor) repository.Repository[entity.User, uuid.UUID] {
	repo := repository.NewGeneralRepository[entity.User, uuid.UUID](flavor)
	repo.WithTableName("users")
	repo.WithCallBack(user_repository.UserRowMapper)
	return repo
}

func RunContainer(t *testing.T) configuration.DatasourceConfig {
	return utility_test.RunContainer(t, "file:../../migrations")
}

func TestInsert(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := generalRepository.Insert(&entity.User{
		LastName:  pointer_utility.NewPointer("sh"),
		FirstName: pointer_utility.NewPointer("al"),
		Username:  pointer_utility.NewPointer("kam")})

	if err != nil {
		t.Fatalf("error insert record, error: %s", err)
	}

	assert.NotNil(t, user)
	assert.NotNil(t, user.Id)
	assert.Equal(t, "sh", *user.LastName)
	assert.Equal(t, "al", *user.FirstName)
	assert.Equal(t, "kam", *user.Username)
}

func TestInsertError(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	repo := createGeneralRepository(sqlbuilder.MySQL)

	user, err := repo.Insert(&entity.User{
		LastName:  pointer_utility.NewPointer("sh"),
		FirstName: pointer_utility.NewPointer("al"),
		Username:  pointer_utility.NewPointer("kam")})

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.True(t, assert.Error(t, err))
}

func TestInsertErrorEntittIsNil(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := generalRepository.Insert(nil)

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.True(t, assert.Error(t, err))
	assert.Equal(t, errors.New("entity is nil"), err)
}

func TestGetAll(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	_, err := generalRepository.Insert(&entity.User{
		LastName:  pointer_utility.NewPointer("sh"),
		FirstName: pointer_utility.NewPointer("al"),
		Username:  pointer_utility.NewPointer("kam")})

	if err != nil {
		t.Fatalf("error insert record, error: %s", err)
	}

	users, err := generalRepository.GetAll()
	if err != nil {
		t.Fatalf("error get all record, error: %s", err)
	}

	assert.Len(t, users, 1)
}

func TestGetAllPageable(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	for i := 0; i < 15; i++ {
		_, err := generalRepository.Insert(&entity.User{
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

	usersFirst, err := generalRepository.GetAllPageable(pageableFirst)

	if err != nil {
		t.Fatalf("error get all record via pageable with error: %s", err)
	}

	pageableLast, err := repository.NewPageableImpl(2, 10, []string{"last_name"})

	if err != nil {
		t.Fatal(err)
	}

	usersLast, err := generalRepository.GetAllPageable(pageableLast)

	if err != nil {
		t.Fatalf("error get all record via pageable with error: %s", err)
	}

	assert.Len(t, usersFirst, 10)
	assert.Len(t, usersLast, 5)
}

func TestGetById(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := generalRepository.Insert(&entity.User{
		LastName:  pointer_utility.NewPointer("sh"),
		FirstName: pointer_utility.NewPointer("al"),
		Username:  pointer_utility.NewPointer("kam")})

	if err != nil {
		t.Fatalf("error insert record, error: %s", err)
	}

	var id *uuid.UUID = user.Id

	userByGetById, err := generalRepository.GetById(id)
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
	user, err := generalRepository.GetById(&id)

	if err == nil {
		t.Fatalf("founding user by id: %s", id.String())
	}

	assert.Nil(t, user)
	assert.Equal(t, fmt.Errorf("not found entity by id: #%v in table: %s", id, "users"), err)
}

func TestGetByIdErrorIdIsNil(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := generalRepository.GetById(nil)

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.True(t, assert.Error(t, err))
	assert.Equal(t, errors.New("id is nil"), err)
}

func TestUpdate(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := generalRepository.Insert(&entity.User{
		LastName:  pointer_utility.NewPointer("sh"),
		FirstName: pointer_utility.NewPointer("al"),
		Username:  pointer_utility.NewPointer("kam")})

	if err != nil {
		t.Fatalf("error insert record, error: %s", err)
	}

	user.FirstName = pointer_utility.NewPointer("ivan")

	userUpdate, err := generalRepository.Update(user)

	if err != nil {
		t.Fatalf("fatal update entity by id: %s, error: %s", user.Id.String(), err)
	}

	assert.Equal(t, *user.Id, *userUpdate.Id)
	assert.Equal(t, "ivan", *userUpdate.FirstName)
	assert.Equal(t, *user.LastName, *userUpdate.LastName)
	assert.Equal(t, *user.Username, *userUpdate.Username)
}

func TestUpdateError(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := generalRepository.Update(&entity.User{
		LastName:  pointer_utility.NewPointer("sh"),
		FirstName: pointer_utility.NewPointer("al"),
		Username:  pointer_utility.NewPointer("kam")})

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("error update in table %s by id: %v", "users", nil), err)

}

func TestUpdateErrorEntityIsNil(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := generalRepository.Update(nil)

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.True(t, assert.Error(t, err))
	assert.Equal(t, errors.New("entity is nil"), err)
}

func TestDeleteById(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := generalRepository.Insert(&entity.User{
		LastName:  pointer_utility.NewPointer("sh"),
		FirstName: pointer_utility.NewPointer("al"),
		Username:  pointer_utility.NewPointer("kam")})

	if err != nil {
		t.Fatalf("error insert record, error: %s", err)
	}

	id, err := generalRepository.DeleteById(user.Id)

	if err != nil {
		t.Fatalf("error deleted by id: %s, error: %s", user.Id.String(), err)
	}

	entity, _ := generalRepository.GetById(id)

	assert.Nil(t, entity)
}

func TestDeleteByIdError(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	id := uuid.New()

	returnId, err := generalRepository.DeleteById(&id)

	assert.Nil(t, returnId)
	assert.NotNil(t, err)
	assert.True(t, assert.Error(t, err))
}

func TestDeleteByIdErrorIdIsNil(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	id, err := generalRepository.DeleteById(nil)

	assert.Nil(t, id)
	assert.NotNil(t, err)
	assert.True(t, assert.Error(t, err))
	assert.Equal(t, errors.New("id is nil"), err)
}

func TestDelete(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := generalRepository.Insert(&entity.User{
		LastName:  pointer_utility.NewPointer("sh"),
		FirstName: pointer_utility.NewPointer("al"),
		Username:  pointer_utility.NewPointer("kam")})

	if err != nil {
		t.Fatalf("error insert record, error: %s", err)
	}

	id, err := generalRepository.Delete(user)

	if err != nil {
		t.Fatalf("error deleted by id: %s, error: %s", user.Id.String(), err)
	}

	entity, _ := generalRepository.GetById(id)

	assert.Nil(t, entity)
}

func TestDeleteError(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user := entity.User{
		LastName:  pointer_utility.NewPointer("sh"),
		FirstName: pointer_utility.NewPointer("al"),
		Username:  pointer_utility.NewPointer("kam")}

	id, err := generalRepository.Delete(&user)

	assert.Nil(t, id)
	assert.NotNil(t, err)
	assert.True(t, assert.Error(t, err))
}

func TestDeleteErrorEntityIsNil(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	id, err := generalRepository.Delete(nil)

	assert.Nil(t, id)
	assert.NotNil(t, err)
	assert.True(t, assert.Error(t, err))
	assert.Equal(t, errors.New("entity is nil"), err)
}
