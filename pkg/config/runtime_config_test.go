package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
)

func TestParseWatchNamespaces(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "single namespace",
			input: "default",
			want:  []string{"default"},
		},
		{
			name:  "two namespaces",
			input: "default,kube-system",
			want:  []string{"default", "kube-system"},
		},
		{
			name:  "three namespaces",
			input: "ns-a,ns-b,ns-c",
			want:  []string{"ns-a", "ns-b", "ns-c"},
		},
		{
			name:  "whitespace trimmed",
			input: " default , kube-system ",
			want:  []string{"default", "kube-system"},
		},
		{
			name:  "duplicates removed",
			input: "default,kube-system,default",
			want:  []string{"default", "kube-system"},
		},
		{
			name:  "empty segments ignored",
			input: "default,,kube-system,",
			want:  []string{"default", "kube-system"},
		},
		{
			name:  "all empty",
			input: ",,",
			want:  []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseWatchNamespaces(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBuildRuntimeOptions_WatchNamespace(t *testing.T) {
	scheme := k8sruntime.NewScheme()

	tests := []struct {
		name                  string
		watchNamespace        string
		wantDefaultNamespaces map[string]cache.Config
	}{
		{
			name:                  "empty namespace watches all",
			watchNamespace:        corev1.NamespaceAll,
			wantDefaultNamespaces: nil,
		},
		{
			name:           "single namespace",
			watchNamespace: "default",
			wantDefaultNamespaces: map[string]cache.Config{
				"default": {},
			},
		},
		{
			name:           "multiple namespaces",
			watchNamespace: "s3-frontend,tailscale",
			wantDefaultNamespaces: map[string]cache.Config{
				"s3-frontend": {},
				"tailscale":   {},
			},
		},
		{
			name:           "multiple namespaces with whitespace",
			watchNamespace: "s3-frontend , tailscale , kube-system",
			wantDefaultNamespaces: map[string]cache.Config{
				"s3-frontend": {},
				"tailscale":   {},
				"kube-system": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rtCfg := RuntimeConfig{
				WatchNamespace: tt.watchNamespace,
			}
			opts, err := BuildRuntimeOptions(rtCfg, scheme)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantDefaultNamespaces, opts.Cache.DefaultNamespaces)
		})
	}
}
