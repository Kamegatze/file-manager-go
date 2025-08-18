package repository_test

import (
	"context"
	"database/sql"
	"file-manager/configuration"
	"file-manager/entity"
	"file-manager/repository"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var generalRepository repository.GeneralRepository[entity.User, uuid.UUID] = createGeneralRepository(sqlbuilder.PostgreSQL)

func createGeneralRepository(flavor sqlbuilder.Flavor) repository.GeneralRepository[entity.User, uuid.UUID] {
	repo := repository.NewGeneralRepository[entity.User, uuid.UUID](flavor)
	repo.WithTableName("users")
	repo.WithCallBack(UserRowMapper)
	return repo
}

func UserRowMapper(rows *sql.Rows) (entity.User, error) {
	user := entity.User{}
	if err := rows.Scan(
		&user.Id,
		&user.LastName,
		&user.FirstName,
		&user.Username); err != nil {
		return user, err
	}

	return user, nil
}

func RunContainer(t *testing.T) configuration.DatasourceConfig {
	ctx := context.Background()
	pgContainer, err := postgres.Run(
		ctx,
		"postgres:15.3-alpine",
		postgres.WithDatabase("file-manager"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)))

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	endpoint, err := pgContainer.PortEndpoint(ctx, "5432/tcp", "")

	if err != nil {
		t.Fatalf("failed get endpoint for container, pgcontainer: %s", err)
	}

	host := strings.Split(endpoint, ":")[0]
	port, err := strconv.Atoi(strings.Split(endpoint, ":")[1])
	if err != nil {
		t.Fatalf("failed get port for container, error: %s", err)
	}

	datasourceConfig := configuration.DatasourceConfig{
		Host:     host,
		Port:     int32(port),
		Username: "postgres",
		Password: "postgres",
		Driver:   "postgres",
		Database: "file-manager",
	}

	if err := datasourceConfig.Migration("file:../migrations"); err != nil {
		t.Fatalf("failed migration up with error: %s", err)
	}

	if err := configuration.RunDatabase(datasourceConfig); err != nil {
		t.Fatalf("failed start database, error: %s", err)
	}
	return datasourceConfig
}

func TestInsert(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := generalRepository.Insert(entity.User{LastName: "sh", FirstName: "al", Username: "kam"})

	if err != nil {
		t.Fatalf("error insert record, error: %s", err)
	}

	assert.NotNil(t, user)
	assert.False(t, user.Id.String() == "")
	assert.Equal(t, "sh", user.LastName)
	assert.Equal(t, "al", user.FirstName)
	assert.Equal(t, "kam", user.Username)
}

func TestInsertError(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	repo := createGeneralRepository(sqlbuilder.MySQL)

	user, err := repo.Insert(entity.User{LastName: "sh", FirstName: "al", Username: "kam"})

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.True(t, assert.Error(t, err))
}

func TestGetAll(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	_, err := generalRepository.Insert(entity.User{LastName: "sh", FirstName: "al", Username: "kam"})

	if err != nil {
		t.Fatalf("error insert record, error: %s", err)
	}

	users, err := generalRepository.GetAll()
	if err != nil {
		t.Fatalf("error get all record, error: %s", err)
	}

	assert.Len(t, users, 1)
}

func TestGetById(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := generalRepository.Insert(entity.User{LastName: "sh", FirstName: "al", Username: "kam"})

	if err != nil {
		t.Fatalf("error insert record, error: %s", err)
	}

	var id uuid.UUID = user.Id

	userByGetById, err := generalRepository.GetById(id)
	if err != nil {
		t.Fatalf("fatal get entity by id: %s, with error: %s", id.String(), err)
	}

	assert.Equal(t, id, userByGetById.Id)
	assert.Equal(t, user.LastName, userByGetById.LastName)
	assert.Equal(t, user.FirstName, userByGetById.FirstName)
	assert.Equal(t, user.Username, userByGetById.Username)
}

func TestGetByIdError(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	id := uuid.New()
	user, err := generalRepository.GetById(id)

	if err == nil {
		t.Fatalf("founding user by id: %s", id.String())
	}

	assert.Nil(t, user)
	assert.Equal(t, fmt.Errorf("not found entity by id: #%v in table: %s", id, "users"), err)
}

func TestUpdate(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := generalRepository.Insert(entity.User{LastName: "sh", FirstName: "al", Username: "kam"})

	if err != nil {
		t.Fatalf("error insert record, error: %s", err)
	}

	user.FirstName = "ivan"

	userUpdate, err := generalRepository.Update(*user)

	if err != nil {
		t.Fatalf("fatal update entity by id: %s, error: %s", user.Id.String(), err)
	}

	assert.Equal(t, user.Id, userUpdate.Id)
	assert.Equal(t, "ivan", userUpdate.FirstName)
	assert.Equal(t, user.LastName, userUpdate.LastName)
	assert.Equal(t, user.Username, userUpdate.Username)
}

func TestUpdateError(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := generalRepository.Update(entity.User{LastName: "sh", FirstName: "al", Username: "kam"})

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("error update in table %s by id: %v", "users", uuid.UUID{}), err)

}

func TestDeleteById(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := generalRepository.Insert(entity.User{LastName: "sh", FirstName: "al", Username: "kam"})

	if err != nil {
		t.Fatalf("error insert record, error: %s", err)
	}

	id, err := generalRepository.DeleteById(user.Id)

	if err != nil {
		t.Fatalf("error deleted by id: %s, error: %s", user.Id.String(), err)
	}

	entity, _ := generalRepository.GetById(*id)

	assert.Nil(t, entity)
}

func TestDeleteByIdError(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	id := uuid.New()

	returnId, err := generalRepository.DeleteById(id)

	assert.Nil(t, returnId)
	assert.NotNil(t, err)
	assert.True(t, assert.Error(t, err))
}

func TestDelete(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user, err := generalRepository.Insert(entity.User{LastName: "sh", FirstName: "al", Username: "kam"})

	if err != nil {
		t.Fatalf("error insert record, error: %s", err)
	}

	id, err := generalRepository.Delete(*user)

	if err != nil {
		t.Fatalf("error deleted by id: %s, error: %s", user.Id.String(), err)
	}

	entity, _ := generalRepository.GetById(*id)

	assert.Nil(t, entity)
}

func TestDeleteError(t *testing.T) {
	datasourceConfig := RunContainer(t)

	defer datasourceConfig.Close()

	user := entity.User{Id: uuid.New(), LastName: "sh", FirstName: "al", Username: "kam"}

	id, err := generalRepository.Delete(user)

	assert.Nil(t, id)
	assert.NotNil(t, err)
	assert.True(t, assert.Error(t, err))
}
