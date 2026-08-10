package app

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"MeshNet/internal/domain"
)

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

func TestServicePut(t *testing.T) {
	repo := &fakeRepository{}
	service := New(repo)

	data := []byte("hello world")

	obj, err := service.Put(
		context.Background(),
		domain.SourceCLI,
		data,
	)
	if err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	if obj == nil {
		t.Fatal("Put() returned nil object")
	}

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
	service := New(repo)

	data := []byte("hello world")

	obj, err := service.Put(
		context.Background(),
		domain.SourceCLI,
		data,
	)
	if err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	const wantHash = "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

	if obj.Hash != wantHash {
		t.Errorf(
			"Hash = %q, want %q",
			obj.Hash,
			wantHash,
		)
	}

}

func TestServicePutEmptyData(t *testing.T) {
	repo := &fakeRepository{}
	service := New(repo)

	obj, err := service.Put(
		context.Background(),
		domain.SourceCLI,
		nil,
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

	if obj != nil {
		t.Errorf(
			"Put() returned object = %v, want nil",
			obj,
		)
	}

	if repo.savedObject != nil {
		t.Error("repository was called for empty payload")
	}

}

func TestServicePutCancelledContext(t *testing.T) {
	repo := &fakeRepository{}
	service := New(repo)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	obj, err := service.Put(
		ctx,
		domain.SourceCLI,
		[]byte("hello"),
	)

	if !errors.Is(err, context.Canceled) {
		t.Errorf(
			"Put() error = %v, want context.Canceled",
			err,
		)
	}

	if obj != nil {
		t.Errorf(
			"Put() returned object = %v, want nil",
			obj,
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

	service := New(repo)

	obj, err := service.Put(
		context.Background(),
		domain.SourceCLI,
		[]byte("hello"),
	)

	if !errors.Is(err, expectedErr) {
		t.Errorf(
			"Put() error = %v, want %v",
			err,
			expectedErr,
		)
	}

	if obj != nil {
		t.Errorf(
			"Put() returned object = %v, want nil",
			obj,
		)
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

	service := New(repo)

	obj, payload, err := service.Get(
		context.Background(),
		"object-1",
	)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if obj != expectedObject {
		t.Error("Get() returned unexpected object")
	}

	if payload != expectedPayload {
		t.Error("Get() returned unexpected payload")
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

	service := New(repo)

	obj, payload, err := service.GetByHash(
		context.Background(),
		"hash-1",
	)
	if err != nil {
		t.Fatalf("GetByHash() error = %v", err)
	}

	if obj != expectedObject {
		t.Error("GetByHash() returned unexpected object")
	}

	if payload != expectedPayload {
		t.Error("GetByHash() returned unexpected payload")
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

	service := New(repo)

	objects, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(objects) != len(expected) {
		t.Fatalf(
			"List() returned %d objects, want %d",
			len(objects),
			len(expected),
		)
	}

	for i := range expected {
		if objects[i] != expected[i] {
			t.Errorf(
				"objects[%d] = %v, want %v",
				i,
				objects[i],
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

	service := New(repo)

	err := service.Delete(
		context.Background(),
		"object-1",
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

	service := New(repo)

	t.Run("Get", func(t *testing.T) {
		_, _, err := service.Get(
			context.Background(),
			"missing",
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
		_, _, err := service.GetByHash(
			context.Background(),
			"missing",
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
		_, err := service.List(context.Background())

		if !errors.Is(err, expectedErr) {
			t.Errorf(
				"List() error = %v, want %v",
				err,
				expectedErr,
			)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := service.Delete(
			context.Background(),
			"missing",
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
