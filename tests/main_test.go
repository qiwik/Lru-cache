package tests

import (
	"github.com/qiwik/Lru-cache/pkg/models"
	"testing"
	"github.com/stretchr/testify"
)

func Test_Init(t *testing.T) {
	newCache := models.NewLRUCache(2)
	testify
}
