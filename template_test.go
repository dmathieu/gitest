package gitest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTemplate(t *testing.T) {
	_, err := newTemplate("basic")
	assert.Nil(t, err)
}

func TestNewMissingTemplate(t *testing.T) {
	_, err := newTemplate("unknown")
	assert.NotNil(t, err)
}
