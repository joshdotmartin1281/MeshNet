package app

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"MeshNet/internal/api"
	"MeshNet/internal/domain"
)

type fakeHasher struct {
	hashFunc func([]byte) string
}

func (f *fakeHasher) Hash(data []byte) string {
	if f.hashFunc != nil {
		return f.hashFunc(data)
	}

	return "test-hash"
}

func (f *fakeHasher) Name() string {
	return "test"
}

func (f *fakeHasher) Version() string {
	return "1"
}

type fakeProcessor struct {
	name    string
	version string
	process func([]byte, domain.Transform) ([]byte, error)
}

func (f *fakeProcessor) Name() string {
	return f.name
}

func (f *fakeProcessor) Version() string {
	return f.version
}

func (f *fakeProcessor) Process(
	data []byte,
	transform domain.Transform,
) ([]byte, error) {
	if f.process != nil {
		return f.process(data, transform)
	}

	return data, nil
}

type fakeRepository struct {
	saveFunc      func(context.Context, *domain.Object, *domain.Payload) error
	getFunc       func(context.Context, string) (*domain.Object, *domain.Payload, error)
	getByHashFunc func(context.Context, string) (*domain.Object, *domain.Payload, error)
	listFunc      func(context.Context) ([]*domain.Object, error)
	deleteFunc    func(context.Context, string) error

	savedObject  *domain.Object
	savedPayload *domain.Payload
}

func (f *fakeRepository) Save(
	ctx context.Context,
	obj *domain.Object,
	payload *domain.Payload,
) error {
	f.savedObject = obj
	f.savedPayload = payload

	if f.saveFunc != nil {
		return f.saveFunc(ctx, obj, payload)
	}

	return nil
}

func (f *fakeRepository) Get(
	ctx context.Context,
	id string,
) (*domain.Object, *domain.Payload, error) {
	if f.getFunc != nil {
		return f.getFunc(ctx, id)
	}

	return nil, nil, nil
}

func (f *fakeRepository) GetByHash(
	ctx context.Context,
	hash string,
) (*domain.Object, *domain.Payload, error) {
	if f.getByHashFunc != nil {
		return f.getByHashFunc(ctx, hash)
	}

	return nil, nil, nil
}

func (f *fakeRepository) List(
	ctx context.Context,
) ([]*domain.Object, error) {
	if f.listFunc != nil {
		return f.listFunc(ctx)
	}

	return nil, nil
}

func (f *fakeRepository) Delete(
	ctx context.Context,
	id string,
) error {
	if f.deleteFunc != nil {
		return f.deleteFunc(ctx, id)
	}

	return nil
}

func newTestService(repo *fakeRepository) *Service {
	uppercase := &fakeProcessor{
		name:    "uppercase",
		version: "1",
		process: func(
			data []byte,
			transform domain.Transform,
		) ([]byte, error) {
			return bytes.ToUpper(data), nil
		},
	}

	return New(
		repo,
		&fakeHasher{},
		NewProcessor(uppercase),
	)
}

func TestServicePut(t *testing.T) {
	repo := &fakeRepository{}
	service := newTestService(repo)

	data := []byte("hello world")

	resp, err := service.Put(
		context.Background(),
		api.PutRequest{
			Source: domain.SourceCLI,
			Data:   data,
		},
	)
	if err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	if resp.Object == nil {
		t.Fatal("Put() returned nil object")
	}

	obj := resp.Object

	if obj.ID == "" {
		t.Error("Put() generated empty ID")
	}

	if obj.Source != domain.SourceCLI {
		t.Errorf(
			"Source = %q, want %q",
			obj.Source,
			domain.SourceCLI,
		)
	}

	if obj.Size != int64(len(data)) {
		t.Errorf(
			"Size = %d, want %d",
			obj.Size,
			len(data),
		)
	}

	if obj.Hash == "" {
		t.Error("Put() generated empty hash")
	}

	if obj.CreatedAt.IsZero() {
		t.Error("Put() generated zero CreatedAt")
	}

	if repo.savedObject != obj {
		t.Error("repository received a different object")
	}

	if repo.savedPayload == nil {
		t.Fatal("repository received nil payload")
	}

	if repo.savedPayload.ObjectID != obj.ID {
		t.Errorf(
			"Payload.ObjectID = %q, want %q",
			repo.savedPayload.ObjectID,
			obj.ID,
		)
	}

	if !bytes.Equal(repo.savedPayload.Data, data) {
		t.Errorf(
			"Payload.Data = %v, want %v",
			repo.savedPayload.Data,
			data,
		)
	}
}

func TestServicePutHash(t *testing.T) {
	repo := &fakeRepository{}

	hasher := &fakeHasher{
		hashFunc: func(data []byte) string {
			if !bytes.Equal(data, []byte("hello world")) {
				t.Errorf(
					"Hasher received %q, want %q",
					data,
					"hello world",
				)
			}

			return "expected-hash"
		},
	}

	service := New(
		repo,
		hasher,
		NewProcessor(),
	)

	resp, err := service.Put(
		context.Background(),
		api.PutRequest{
			Source: domain.SourceCLI,
			Data:   []byte("hello world"),
		},
	)
	if err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	if resp.Object.Hash != "expected-hash" {
		t.Errorf(
			"Hash = %q, want %q",
			resp.Object.Hash,
			"expected-hash",
		)
	}
}

func TestServicePutEmptyData(t *testing.T) {
	repo := &fakeRepository{}
	service := newTestService(repo)

	resp, err := service.Put(
		context.Background(),
		api.PutRequest{
			Source: domain.SourceCLI,
			Data:   nil,
		},
	)

	if err == nil {
		t.Fatal("Put() error = nil, want error")
	}

	if err.Error() != "empty payload" {
		t.Errorf(
			"Put() error = %q, want %q",
			err,
			"empty payload",
		)
	}

	if resp.Object != nil {
		t.Errorf(
			"Put() returned object = %v, want nil",
			resp.Object,
		)
	}

	if repo.savedObject != nil {
		t.Error("repository was called for empty payload")
	}
}

func TestServicePutCancelledContext(t *testing.T) {
	repo := &fakeRepository{}
	service := newTestService(repo)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	resp, err := service.Put(
		ctx,
		api.PutRequest{
			Source: domain.SourceCLI,
			Data:   []byte("hello"),
		},
	)

	if !errors.Is(err, context.Canceled) {
		t.Errorf(
			"Put() error = %v, want context.Canceled",
			err,
		)
	}

	if resp.Object != nil {
		t.Errorf(
			"Put() returned object = %v, want nil",
			resp.Object,
		)
	}

	if repo.savedObject != nil {
		t.Error("repository was called with cancelled context")
	}
}

func TestServicePutRepositoryError(t *testing.T) {
	expectedErr := errors.New("repository failure")

	repo := &fakeRepository{
		saveFunc: func(
			ctx context.Context,
			obj *domain.Object,
			payload *domain.Payload,
		) error {
			return expectedErr
		},
	}

	service := newTestService(repo)

	resp, err := service.Put(
		context.Background(),
		api.PutRequest{
			Source: domain.SourceCLI,
			Data:   []byte("hello"),
		},
	)

	if !errors.Is(err, expectedErr) {
		t.Errorf(
			"Put() error = %v, want %v",
			err,
			expectedErr,
		)
	}

	if resp.Object != nil {
		t.Errorf(
			"Put() returned object = %v, want nil",
			resp.Object,
		)
	}
}

func TestServicePutTransform(t *testing.T) {
	repo := &fakeRepository{}
	service := newTestService(repo)

	data := []byte("hello world")

	resp, err := service.Put(
		context.Background(),
		api.PutRequest{
			Source: domain.SourceCLI,
			Data:   data,
			Transforms: []domain.Transform{
				{
					Name:    "uppercase",
					Version: "1",
					Params:  map[string]string{},
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	expected := []byte("HELLO WORLD")

	if !bytes.Equal(repo.savedPayload.Data, expected) {
		t.Errorf(
			"stored payload = %q, want %q",
			repo.savedPayload.Data,
			expected,
		)
	}

	if resp.Object.Hash != "test-hash" {
		t.Errorf(
			"Hash = %q, want test-hash",
			resp.Object.Hash,
		)
	}
}

func TestServicePutStoresTransforms(t *testing.T) {
	repo := &fakeRepository{}
	service := newTestService(repo)

	transforms := []domain.Transform{
		{
			Name:    "uppercase",
			Version: "1",
			Params:  map[string]string{},
		},
	}

	_, err := service.Put(
		context.Background(),
		api.PutRequest{
			Source:     domain.SourceCLI,
			Data:       []byte("hello world"),
			Transforms: transforms,
		},
	)
	if err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	if repo.savedObject == nil {
		t.Fatal("repository received nil object")
	}

	if len(repo.savedObject.Transforms) != 1 {
		t.Fatalf(
			"Transforms length = %d, want 1",
			len(repo.savedObject.Transforms),
		)
	}

	got := repo.savedObject.Transforms[0]

	if got.Name != "uppercase" {
		t.Errorf(
			"Transform.Name = %q, want %q",
			got.Name,
			"uppercase",
		)
	}

	if got.Version != "1" {
		t.Errorf(
			"Transform.Version = %q, want %q",
			got.Version,
			"1",
		)
	}

	if got.Params == nil {
		t.Fatal("Transform.Params = nil, want non-nil map")
	}
}

func TestServiceGet(t *testing.T) {
	expectedObject := &domain.Object{
		ID:   "object-1",
		Hash: "hash-1",
	}

	expectedPayload := &domain.Payload{
		ObjectID: "object-1",
		Data:     []byte("hello"),
	}

	repo := &fakeRepository{
		getFunc: func(
			ctx context.Context,
			id string,
		) (*domain.Object, *domain.Payload, error) {
			if id != "object-1" {
				t.Errorf(
					"Get() id = %q, want object-1",
					id,
				)
			}

			return expectedObject, expectedPayload, nil
		},
	}

	service := newTestService(repo)

	resp, err := service.Get(
		context.Background(),
		api.GetRequest{
			ID: "object-1",
		},
	)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if resp.Object != expectedObject {
		t.Error("Get() returned unexpected object")
	}

	if resp.Payload != expectedPayload {
		t.Error("Get() returned unexpected payload")
	}
}

func TestServiceGetTransform(t *testing.T) {
	expectedObject := &domain.Object{
		ID:   "object-1",
		Hash: "hash-1",
	}

	expectedPayload := &domain.Payload{
		ObjectID: "object-1",
		Data:     []byte("hello"),
	}

	repo := &fakeRepository{
		getFunc: func(
			ctx context.Context,
			id string,
		) (*domain.Object, *domain.Payload, error) {
			return expectedObject, expectedPayload, nil
		},
	}

	service := newTestService(repo)

	resp, err := service.Get(
		context.Background(),
		api.GetRequest{
			ID: "object-1",
			Transforms: []domain.Transform{
				{
					Name:    "uppercase",
					Version: "1",
					Params:  map[string]string{},
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	expected := []byte("HELLO")

	if !bytes.Equal(resp.Payload.Data, expected) {
		t.Errorf(
			"Get() payload = %q, want %q",
			resp.Payload.Data,
			expected,
		)
	}
}

func TestServiceGetByHash(t *testing.T) {
	expectedObject := &domain.Object{
		ID:   "object-1",
		Hash: "hash-1",
	}

	expectedPayload := &domain.Payload{
		ObjectID: "object-1",
		Data:     []byte("hello"),
	}

	repo := &fakeRepository{
		getByHashFunc: func(
			ctx context.Context,
			hash string,
		) (*domain.Object, *domain.Payload, error) {
			if hash != "hash-1" {
				t.Errorf(
					"GetByHash() hash = %q, want hash-1",
					hash,
				)
			}

			return expectedObject, expectedPayload, nil
		},
	}

	service := newTestService(repo)

	resp, err := service.GetByHash(
		context.Background(),
		api.GetByHashRequest{
			Hash: "hash-1",
		},
	)
	if err != nil {
		t.Fatalf("GetByHash() error = %v", err)
	}

	if resp.Object != expectedObject {
		t.Error("GetByHash() returned unexpected object")
	}

	if resp.Payload != expectedPayload {
		t.Error("GetByHash() returned unexpected payload")
	}
}

func TestServiceGetByHashTransform(t *testing.T) {
	expectedObject := &domain.Object{
		ID:   "object-1",
		Hash: "hash-1",
	}

	expectedPayload := &domain.Payload{
		ObjectID: "object-1",
		Data:     []byte("hello"),
	}

	repo := &fakeRepository{
		getByHashFunc: func(
			ctx context.Context,
			hash string,
		) (*domain.Object, *domain.Payload, error) {
			return expectedObject, expectedPayload, nil
		},
	}

	service := newTestService(repo)

	resp, err := service.GetByHash(
		context.Background(),
		api.GetByHashRequest{
			Hash: "hash-1",
			Transforms: []domain.Transform{
				{
					Name:    "uppercase",
					Version: "1",
					Params:  map[string]string{},
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("GetByHash() error = %v", err)
	}

	expected := []byte("HELLO")

	if !bytes.Equal(resp.Payload.Data, expected) {
		t.Errorf(
			"GetByHash() payload = %q, want %q",
			resp.Payload.Data,
			expected,
		)
	}
}

func TestServiceList(t *testing.T) {
	expected := []*domain.Object{
		{
			ID: "object-1",
		},
		{
			ID: "object-2",
		},
	}

	repo := &fakeRepository{
		listFunc: func(
			ctx context.Context,
		) ([]*domain.Object, error) {
			return expected, nil
		},
	}

	service := newTestService(repo)

	resp, err := service.List(
		context.Background(),
		api.ListRequest{},
	)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(resp.Objects) != len(expected) {
		t.Fatalf(
			"List() returned %d objects, want %d",
			len(resp.Objects),
			len(expected),
		)
	}

	for i := range expected {
		if resp.Objects[i] != expected[i] {
			t.Errorf(
				"objects[%d] = %v, want %v",
				i,
				resp.Objects[i],
				expected[i],
			)
		}
	}
}

func TestServiceDelete(t *testing.T) {
	var deletedID string

	repo := &fakeRepository{
		deleteFunc: func(
			ctx context.Context,
			id string,
		) error {
			deletedID = id
			return nil
		},
	}

	service := newTestService(repo)

	_, err := service.Delete(
		context.Background(),
		api.DeleteRequest{
			ID: "object-1",
		},
	)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	if deletedID != "object-1" {
		t.Errorf(
			"Delete() id = %q, want object-1",
			deletedID,
		)
	}
}

func TestServiceErrorPropagation(t *testing.T) {
	expectedErr := domain.ErrNotFound

	repo := &fakeRepository{
		getFunc: func(
			ctx context.Context,
			id string,
		) (*domain.Object, *domain.Payload, error) {
			return nil, nil, expectedErr
		},
		getByHashFunc: func(
			ctx context.Context,
			hash string,
		) (*domain.Object, *domain.Payload, error) {
			return nil, nil, expectedErr
		},
		listFunc: func(
			ctx context.Context,
		) ([]*domain.Object, error) {
			return nil, expectedErr
		},
		deleteFunc: func(
			ctx context.Context,
			id string,
		) error {
			return expectedErr
		},
	}

	service := newTestService(repo)

	t.Run("Get", func(t *testing.T) {
		_, err := service.Get(
			context.Background(),
			api.GetRequest{
				ID: "missing",
			},
		)

		if !errors.Is(err, expectedErr) {
			t.Errorf(
				"Get() error = %v, want %v",
				err,
				expectedErr,
			)
		}
	})

	t.Run("GetByHash", func(t *testing.T) {
		_, err := service.GetByHash(
			context.Background(),
			api.GetByHashRequest{
				Hash: "missing",
			},
		)

		if !errors.Is(err, expectedErr) {
			t.Errorf(
				"GetByHash() error = %v, want %v",
				err,
				expectedErr,
			)
		}
	})

	t.Run("List", func(t *testing.T) {
		_, err := service.List(
			context.Background(),
			api.ListRequest{},
		)

		if !errors.Is(err, expectedErr) {
			t.Errorf(
				"List() error = %v, want %v",
				err,
				expectedErr,
			)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		_, err := service.Delete(
			context.Background(),
			api.DeleteRequest{
				ID: "missing",
			},
		)

		if !errors.Is(err, expectedErr) {
			t.Errorf(
				"Delete() error = %v, want %v",
				err,
				expectedErr,
			)
		}
	})
}
