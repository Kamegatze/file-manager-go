package repository

import (
	"database/sql"
	"file-manager/internal/configuration"
	pointer_utility "file-manager/pkg/pointer-utility"
	"fmt"
	"log"
	"reflect"

	"github.com/huandu/go-sqlbuilder"
)

type Repository[E interface{}, ID any] interface {
	GetAll() ([]E, error)
	GetAllPageable(pageable Pageable) ([]E, error)
	GetById(id *ID) (*E, error)
	Insert(entity *E) (*E, error)
	Update(entity *E) (*E, error)
	Delete(entity *E) (*ID, error)
	DeleteById(id *ID) (*ID, error)
}

type RowMapper[E interface{}, ID any] func(rows *sql.Rows) (E, error)

type GetTag func() string

type GeneralRepository[E interface{}, ID any] struct {
	tableName string
	callBack  RowMapper[E, ID]
	tagColumn GetTag
	tagId     GetTag
	flavors   sqlbuilder.Flavor
}

func NewGeneralRepository[E interface{}, ID any](flavors sqlbuilder.Flavor) GeneralRepository[E, ID] {
	repo := GeneralRepository[E, ID]{}
	repo.tagColumn = func() string { return "col" }
	repo.tagId = func() string { return "id" }
	repo.flavors = flavors
	return repo
}

func (generalRepository *GeneralRepository[E, ID]) WithTableName(tableName string) *GeneralRepository[E, ID] {
	generalRepository.tableName = tableName
	return generalRepository
}

func (generalRepository *GeneralRepository[E, ID]) WithCallBack(callBack RowMapper[E, ID]) *GeneralRepository[E, ID] {
	generalRepository.callBack = callBack
	return generalRepository
}

func (generalRepository *GeneralRepository[E, ID]) WithTagColumn(tagColumn GetTag) *GeneralRepository[E, ID] {
	generalRepository.tagColumn = tagColumn
	return generalRepository
}

func (generalRepository *GeneralRepository[E, ID]) WithTagId(tagId GetTag) *GeneralRepository[E, ID] {
	generalRepository.tagId = tagId
	return generalRepository
}

func (generalRepository GeneralRepository[E, ID]) GetAll() ([]E, error) {
	keys, _, _ := generalRepository.metaInfoGet((*E)(nil))

	db, err := configuration.GetDB()
	if err != nil {
		return nil, err
	}

	buider := sqlbuilder.Select(keys...).From(generalRepository.tableName)
	buider.SetFlavor(generalRepository.flavors)
	sql, _ := buider.Build()

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	slice := []E{}

	for rows.Next() {
		entity, err := generalRepository.callBack(rows)
		if err != nil {
			return slice, err
		}
		slice = append(slice, entity)
	}
	return slice, nil
}

func (generalRepository GeneralRepository[E, ID]) GetAllPageable(pageable Pageable) ([]E, error) {
	keys, _, _ := generalRepository.metaInfoGet((*E)(nil))

	db, err := configuration.GetDB()
	if err != nil {
		return nil, err
	}

	buider := sqlbuilder.Select(keys...).From(generalRepository.tableName)
	if len(pageable.OrderBy()) > 0 {
		buider.OrderBy(pageable.OrderBy()...)
	}
	buider.Limit(pageable.Limit())
	buider.Offset(pageable.Offset())
	buider.SetFlavor(generalRepository.flavors)
	sql, args := buider.Build()

	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	slice := []E{}

	for rows.Next() {
		entity, err := generalRepository.callBack(rows)
		if err != nil {
			return slice, err
		}
		slice = append(slice, entity)
	}
	return slice, nil
}

func (generalRepository GeneralRepository[E, ID]) GetById(id *ID) (*E, error) {
	if err := pointer_utility.PointerIsNil(id, "id is nil"); err != nil {
		return nil, err
	}

	keys, _, indexId := generalRepository.metaInfoGet((*E)(nil))

	db, err := configuration.GetDB()
	if err != nil {
		return nil, err
	}

	builder := sqlbuilder.Select(keys...).From(generalRepository.tableName)
	builder.Where(builder.Equal(keys[indexId], id))
	builder.SetFlavor(generalRepository.flavors)
	sql, args := builder.Build()

	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if rows.Next() {
		entity, err := generalRepository.callBack(rows)
		if err != nil {
			return &entity, err
		}
		return &entity, nil
	}

	return nil, fmt.Errorf("not found entity by id: #%v in table: %s", id, generalRepository.tableName)
}

func (generalRepository GeneralRepository[E, ID]) Insert(entity *E) (*E, error) {

	if err := pointer_utility.PointerIsNil(entity, "entity is nil"); err != nil {
		return nil, err
	}

	keys, values, indexId := generalRepository.metaInfoGet(entity)

	keysWithoutId := []string{}
	valuesWithoutId := []any{}

	keysWithoutId = append(keysWithoutId, keys[:indexId]...)
	keysWithoutId = append(keysWithoutId, keys[indexId+1:]...)

	valuesWithoutId = append(valuesWithoutId, values[:indexId]...)
	valuesWithoutId = append(valuesWithoutId, values[indexId+1:]...)

	keysWithoutNil := make([]string, 0, len(keysWithoutId))
	valuesWithoutNil := make([]any, 0, len(keysWithoutId))

	for i, value := range valuesWithoutId {
		if reflect.ValueOf(value).Kind() != reflect.Ptr || !reflect.ValueOf(value).IsNil() {
			keysWithoutNil = append(keysWithoutNil, keysWithoutId[i])
			valuesWithoutNil = append(valuesWithoutNil, value)
		}
	}

	builder := sqlbuilder.InsertInto(generalRepository.tableName)
	builder.Cols(keysWithoutNil...)
	builder.Values(valuesWithoutNil...)
	builder.Returning(keys...)
	builder.SetFlavor(generalRepository.flavors)
	sql, args := builder.Build()

	db, err := configuration.GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(sql, args...)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if rows.Next() {
		entity, err := generalRepository.callBack(rows)
		if err != nil {
			return &entity, err
		}
		return &entity, nil
	}
	return nil, fmt.Errorf("error insert into in table %s", generalRepository.tableName)
}

func (generalRepository GeneralRepository[E, ID]) Update(entity *E) (*E, error) {

	if err := pointer_utility.PointerIsNil(entity, "entity is nil"); err != nil {
		return nil, err
	}

	keys, values, indexId := generalRepository.metaInfoGet(entity)

	builder := sqlbuilder.Update(generalRepository.tableName)
	sets := []string{}
	for i := 0; i < len(keys); i++ {
		if i != indexId {
			sets = append(sets, builder.Assign(keys[i], values[i]))
		}
	}
	builder.Set(sets...).Where(builder.Equal(keys[indexId], values[indexId])).Returning(keys...)
	builder.SetFlavor(generalRepository.flavors)
	sql, args := builder.Build()

	db, err := configuration.GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if rows.Next() {
		entity, err := generalRepository.callBack(rows)
		if err != nil {
			return &entity, err
		}
		return &entity, nil
	}
	return nil, fmt.Errorf("error update in table %s by id: %v", generalRepository.tableName, values[indexId])
}

func (generalRepository GeneralRepository[E, ID]) Delete(entity *E) (*ID, error) {

	if err := pointer_utility.PointerIsNil(entity, "entity is nil"); err != nil {
		return nil, err
	}

	keys, values, indexId := generalRepository.metaInfoGet(entity)

	if reflect.ValueOf(values[indexId]).IsNil() {
		return nil, fmt.Errorf("value from field %s no must be nil", keys[indexId])
	}

	var id ID
	if reflect.ValueOf(values[indexId]).Kind() == reflect.Ptr {
		id = *values[indexId].(*ID)
	} else {
		id = values[indexId].(ID)
	}
	return generalRepository.DeleteById(&id)
}

func (generalRepository GeneralRepository[E, ID]) DeleteById(id *ID) (*ID, error) {

	if err := pointer_utility.PointerIsNil(id, "id is nil"); err != nil {
		return nil, err
	}

	keys, _, indexId := generalRepository.metaInfoGet((*E)(nil))

	builder := sqlbuilder.DeleteFrom(generalRepository.tableName)
	builder.Where(builder.Equal(keys[indexId], id)).Returning(keys[indexId])
	builder.SetFlavor(generalRepository.flavors)
	sql, args := builder.Build()

	db, err := configuration.GetDB()
	if err != nil {
		return nil, err
	}

	var deleteId ID

	if err := db.QueryRow(sql, args...).Scan(&deleteId); err != nil {
		return nil, err
	}

	return &deleteId, nil
}

func (generalRepository GeneralRepository[E, ID]) metaInfoGet(entity *E) (keys []string, values []any, indexId int) {

	typeOf := reflect.TypeOf(entity).Elem()

	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		nameField := field.Tag.Get(generalRepository.tagColumn())
		id := field.Tag.Get(generalRepository.tagId())
		if id != "" {
			indexId = i
		}
		keys = append(keys, nameField)
		if entity != nil {
			values = append(values, reflect.ValueOf(*entity).Field(i).Interface())
		}
	}
	return keys, values, indexId
}
