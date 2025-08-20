package repository_test

import (
	"fmt"
	"testing"

	"github.com/Kamegatze/file-manager-go/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestNewPageableImpl(t *testing.T) {
	pageable, err := repository.NewPageableImpl(2, 10, []string{})

	assert.Nil(t, err)
	assert.NotNil(t, pageable)
	assert.Equal(t, 10, pageable.Offset())
	assert.Equal(t, 10, pageable.Limit())
	assert.Equal(t, []string{}, pageable.OrderBy())
}

func TestNewPageablePageError(t *testing.T) {
	_, err := repository.NewPageableImpl(0, 10, []string{})

	assert.NotNil(t, err)
	assert.True(t, assert.Error(t, err))
	assert.Equal(t, fmt.Errorf("field 'page' must be more 0, but have: %d", 0), err)
}

func TestNewPageableLimitError(t *testing.T) {
	_, err := repository.NewPageableImpl(1, 0, []string{})

	assert.NotNil(t, err)
	assert.True(t, assert.Error(t, err))
	assert.Equal(t, fmt.Errorf("field 'limit' must be more 0, but have: %d", 0), err)
}
