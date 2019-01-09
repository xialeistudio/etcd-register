package discovery

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	service, err := NewService([]string{"http://192.168.1.144:2379"}, "/article/spider", "192.168.144:10000", nil)
	assert.NoError(t, err)

	err = service.Start()
	assert.NoError(t, err)
}
