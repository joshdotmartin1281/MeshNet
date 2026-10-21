package file

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"MeshNet/uauth"
)

type userRecord struct {
	CreatedAt time.Time         `json:"created_at"`
	Tokens    []string          `json:"tokens"`
	Grants    map[string]string `json:"grants"`
}

type state struct {
	Users map[string]*userRecord `json:"users"`
}

type Store struct {
	mu     sync.RWMutex
	path   string
	data   *state
	tokens map[string]string
}

var (
	_ auth.Provider = (*Store)(nil)
	_ auth.Manager  = (*Store)(nil)
)

func New(path string) (*Store, error) {
	s := &Store{
		path: path,
		data: &state{Users: map[string]*userRecord{}},
	}

	if path != "" {
		raw, err := os.ReadFile(path)

		switch {
		case err == nil:
			if err := json.Unmarshal(raw, s.data); err != nil {
				return nil, fmt.Errorf("parse %s: %w", path, err)
			}

			s.data.normalize()

		case errors.Is(err, os.ErrNotExist):

		default:
			return nil, err
		}
	}

	s.reindex()

	return s, nil
}

func (st *state) normalize() {
	if st.Users == nil {
		st.Users = map[string]*userRecord{}
	}

	for _, u := range st.Users {
		if u.Grants == nil {
			u.Grants = map[string]string{}
		}
	}
}

func (st *state) clone() (*state, error) {
	raw, err := json.Marshal(st)
	if err != nil {
		return nil, err
	}

	var out state

	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}

	out.normalize()

	return &out, nil
}

func (s *Store) reindex() {
	idx := make(map[string]string)

	for name, u := range s.data.Users {
		for _, h := range u.Tokens {
			idx[h] = name
		}
	}

	s.tokens = idx
}

func (s *Store) save(st *state) error {
	if s.path == "" {
		return nil
	}

	raw, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(s.path)

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, ".auth-*.tmp")
	if err != nil {
		return err
	}

	name := tmp.Name()

	if _, err := tmp.Write(raw); err != nil {
		tmp.Close()
		os.Remove(name)
		return err
	}

	if err := tmp.Close(); err != nil {
		os.Remove(name)
		return err
	}

	if err := os.Rename(name, s.path); err != nil {
		os.Remove(name)
		return err
	}

	return nil
}

func (s *Store) mutate(fn func(*state) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	next, err := s.data.clone()
	if err != nil {
		return err
	}

	if err := fn(next); err != nil {
		return err
	}

	if err := s.save(next); err != nil {
		return err
	}

	s.data = next
	s.reindex()

	return nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func newToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return "mn_" + hex.EncodeToString(b), nil
}

func (s *Store) Authenticate(_ context.Context, token string) (auth.Identity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	name, ok := s.tokens[hashToken(token)]
	if !ok {
		return auth.Identity{}, auth.ErrUnauthorized
	}

	return auth.Identity{UserID: name}, nil
}

func (s *Store) Authorize(ctx context.Context, collection string, need auth.Level) error {
	id, ok := auth.FromContext(ctx)
	if !ok {
		return auth.ErrUnauthorized
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.data.Users[id.UserID]
	if !ok {
		return auth.ErrForbidden
	}

	raw, ok := user.Grants[collection]
	if !ok {
		return auth.ErrForbidden
	}

	have, err := auth.ParseLevel(raw)
	if err != nil || have < need {
		return auth.ErrForbidden
	}

	return nil
}

func (s *Store) CreateUser(_ context.Context, name string) (auth.User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return auth.User{}, errors.New("name is required")
	}

	now := time.Now().UTC()

	err := s.mutate(func(st *state) error {
		if _, exists := st.Users[name]; exists {
			return auth.ErrUserExists
		}

		st.Users[name] = &userRecord{
			CreatedAt: now,
			Grants:    map[string]string{},
		}

		return nil
	})
	if err != nil {
		return auth.User{}, err
	}

	return auth.User{Name: name, CreatedAt: now}, nil
}

func (s *Store) ListUsers(_ context.Context) ([]auth.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]auth.User, 0, len(s.data.Users))

	for name, rec := range s.data.Users {
		u := auth.User{Name: name, CreatedAt: rec.CreatedAt}

		for collection, raw := range rec.Grants {
			level, err := auth.ParseLevel(raw)
			if err != nil {
				continue
			}

			u.Grants = append(u.Grants, auth.Grant{Collection: collection, Level: level})
		}

		sort.Slice(u.Grants, func(i, j int) bool {
			return u.Grants[i].Collection < u.Grants[j].Collection
		})

		users = append(users, u)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].Name < users[j].Name
	})

	return users, nil
}

func (s *Store) IssueToken(_ context.Context, user string) (string, error) {
	token, err := newToken()
	if err != nil {
		return "", err
	}

	err = s.mutate(func(st *state) error {
		u, ok := st.Users[user]
		if !ok {
			return auth.ErrUserNotFound
		}

		u.Tokens = append(u.Tokens, hashToken(token))

		return nil
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Store) RevokeTokens(_ context.Context, user string) error {
	return s.mutate(func(st *state) error {
		u, ok := st.Users[user]
		if !ok {
			return auth.ErrUserNotFound
		}

		u.Tokens = nil

		return nil
	})
}

func (s *Store) Grant(_ context.Context, user, collection string, level auth.Level) error {
	if strings.TrimSpace(collection) == "" {
		return errors.New("collection is required")
	}

	return s.mutate(func(st *state) error {
		u, ok := st.Users[user]
		if !ok {
			return auth.ErrUserNotFound
		}

		u.Grants[collection] = level.String()

		return nil
	})
}

func (s *Store) Ungrant(_ context.Context, user, collection string) error {
	return s.mutate(func(st *state) error {
		u, ok := st.Users[user]
		if !ok {
			return auth.ErrUserNotFound
		}

		if _, ok := u.Grants[collection]; !ok {
			return fmt.Errorf("%s has no grant on %q", user, collection)
		}

		delete(u.Grants, collection)

		return nil
	})
}
