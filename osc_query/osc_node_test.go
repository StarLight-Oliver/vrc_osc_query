package osc_query_test

import (
	"testing"

	"github.com/StarLight-Oliver/vrc_osc_query/osc_query"
	"github.com/stretchr/testify/assert"
)

func TestNewOscNodeTree(t *testing.T) {
	node := osc_query.NewOscNodeTree("TestApp", "localhost", 8000)

	assert.Equal(t, "root node", node.Description)
	assert.Equal(t, "/", node.FullPath)
	assert.Equal(t, "TestApp", node.HostInfo.Name)
	assert.Equal(t, "localhost", node.HostInfo.Ip)
	assert.Equal(t, 8000, node.HostInfo.Port)
	assert.Equal(t, "UDP", node.HostInfo.Transport)
	assert.True(t, node.HostInfo.Extensions["ACCESS"])
}

func TestAddChild(t *testing.T) {
	root := osc_query.NewOscNodeTree("TestApp", "localhost", 8000)

	// Add a child to the root node
	child1 := root.AddChild("/child1", osc_query.OscTypeInt, "child1 desc")
	assert.Equal(t, "child1 desc", child1.Description)
	assert.Equal(t, "/child1", child1.FullPath)
	assert.Equal(t, "i", child1.Type)
	assert.Equal(t, 3, child1.Access)

	// Add a child to a non-root node
	child2 := root.AddChild("/child2Folder/child2", osc_query.OscTypeFloat, "child2 desc")
	assert.Equal(t, "child2 desc", child2.Description)
	assert.Equal(t, "/child2Folder/child2", child2.FullPath)
	assert.Equal(t, "f", child2.Type)
	assert.Equal(t, 3, child2.Access)

	// Add a child with a full path that includes multiple directories
	child3 := root.AddChild("/dir1/dir2/child3", osc_query.OscTypeBool, "child3 desc")
	assert.Equal(t, "child3 desc", child3.Description)
	assert.Equal(t, "/dir1/dir2/child3", child3.FullPath)
	assert.Equal(t, "T", child3.Type)
	assert.Equal(t, 3, child3.Access)
}
