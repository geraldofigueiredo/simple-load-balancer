package entity

import (
	"net/http/httputil"
	"net/url"
	"reflect"
	"sync"
	"testing"
)

func TestBackend_IsAlive(t *testing.T) {
	type fields struct {
		Name         string
		URL          *url.URL
		Alive        bool
		mux          sync.Mutex
		ReverseProxy *httputil.ReverseProxy
	}
	tests := []struct {
		name      string
		fields    fields
		wantAlive bool
	}{
		{
			name: "Test Backend IsAlive",
			fields: fields{
				Name:         "Test Backend",
				URL:          &url.URL{Scheme: "http", Host: "localhost:8080"},
				Alive:        true,
				ReverseProxy: nil,
				mux:          sync.Mutex{},
			},
			wantAlive: true,
		},
		{
			name: "Test Backend IsAlive",
			fields: fields{
				Name:         "Test Backend",
				URL:          &url.URL{Scheme: "http", Host: "localhost:8080"},
				Alive:        false,
				ReverseProxy: nil,
				mux:          sync.Mutex{},
			},
			wantAlive: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Backend{
				Name:         tt.fields.Name,
				URL:          tt.fields.URL,
				Alive:        tt.fields.Alive,
				mux:          tt.fields.mux,
				ReverseProxy: tt.fields.ReverseProxy,
			}
			if gotAlive := b.IsAlive(); gotAlive != tt.wantAlive {
				t.Errorf("IsAlive() = %v, want %v", gotAlive, tt.wantAlive)
			}
		})
	}
}

func TestBackend_SetAlive(t *testing.T) {
	type fields struct {
		Name         string
		URL          *url.URL
		Alive        bool
		mux          sync.Mutex
		ReverseProxy *httputil.ReverseProxy
	}
	type args struct {
		alive bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Backend{
				Name:         tt.fields.Name,
				URL:          tt.fields.URL,
				Alive:        tt.fields.Alive,
				mux:          tt.fields.mux,
				ReverseProxy: tt.fields.ReverseProxy,
			}
			b.SetAlive(tt.args.alive)
		})
	}
}

func TestNewBackend(t *testing.T) {
	type args struct {
		name   string
		urlStr string
	}
	tests := []struct {
		name string
		args args
		want *Backend
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBackend(tt.args.name, tt.args.urlStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBackend() = %v, want %v", got, tt.want)
			}
		})
	}
}
