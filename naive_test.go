package graph_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trelore/go-graph"
)

func TestMaximalIndependentSet(t *testing.T) {
	tests := []struct {
		name    string
		in      []byte
		want    []graph.NodeID
		wantErr bool
	}{
		{
			name: "empty",
			in:   []byte(``),
			want: []graph.NodeID{},
		},
		{
			name: "simple",
			in:   []byte(`DQc`),
			want: []graph.NodeID{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := graph.ImportG6(tt.in).MaximumIndependentSet()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
