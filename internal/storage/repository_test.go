package storage_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"MeshNet/internal/domain"
	"MeshNet/internal/storage/file"
	"MeshNet/internal/storage/memory"
)

type repositoryFactory func(t *testing.T) domain.ObjectRepository

func TestRepository(t *testing.T) {
	factories := map[string]repositoryFactory{
		"memory": func(t *testing.T) domain.ObjectRepository {
			return memory.New()
		},

		"file": func(t *testing.T) domain.ObjectRepository {
			store, err := file.New(t.TempDir())
			if err != nil {
				t.Fatalf("file.New() error = %v", err)
			}

			return store
		},
	}

	for name, factory := range factories {
		t.Run(name, func(t *testing.T) {
			t.Run("SaveAndGet", func(t *testing.T) {
				testSaveAndGet(t, factory(t))
			})

			t.Run("GetNotFound", func(t *testing.T) {
				testGetNotFound(t, factory(t))
			})

			t.Run("GetByHash", func(t *testing.T) {
				testGetByHash(t, factory(t))
			})

			t.Run("GetByHashNotFound", func(t *testing.T) {
				testGetByHashNotFound(t, factory(t))
			})

			t.Run("DuplicateHash", func(t *testing.T) {
				testDuplicateHash(t, factory(t))
			})

			t.Run("DuplicateID", func(t *testing.T) {
				testDuplicateID(t, factory(t))
			})

			t.Run("Delete", func(t *testing.T) {
				testDelete(t, factory(t))
			})

			t.Run("DeleteNotFound", func(t *testing.T) {
				testDeleteNotFound(t, factory(t))
			})

			t.Run("ListEmpty", func(t *testing.T) {
				testListEmpty(t, factory(t))
			})

			t.Run("List", func(t *testing.T) {
				testList(t, factory(t))
			})

			t.Run("PathIsolation", func(t *testing.T) {
				testPathIsolation(t, factory(t))
			})

			t.Run("ListPathIsolation", func(t *testing.T) {
				testListPathIsolation(t, factory(t))
			})

			t.Run("BinaryPayload", func(t *testing.T) {
				testBinaryPayload(t, factory(t))
			})

			t.Run("CancelledContext", func(t *testing.T) {
				testCancelledContext(t, factory(t))
			})
		})
	}
}

func testSaveAndGet(
	t *testing.T,
	repo domain.ObjectRepository,
) {
	t.Helper()

	ctx := context.Background()

	obj := testObject("object-1", "hash-1")
	payload := testPayload(
		obj.ID,
		[]byte("hello world"),
	)

	if err := repo.Save(ctx, obj, payload); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	gotObject, gotPayload, err := repo.Get(
		ctx,
		obj.Path,
		obj.ID,
	)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if gotObject.ID != obj.ID {
		t.Errorf(
			"ID = %q, want %q",
			gotObject.ID,
			obj.ID,
		)
	}

	if gotObject.Hash != obj.Hash {
		t.Errorf(
			"Hash = %q, want %q",
			gotObject.Hash,
			obj.Hash,
		)
	}

	if gotObject.Path != obj.Path {
		t.Errorf(
			"Path = %q, want %q",
			gotObject.Path,
			obj.Path,
		)
	}

	if gotObject.Size != obj.Size {
		t.Errorf(
			"Size = %d, want %d",
			gotObject.Size,
			obj.Size,
		)
	}

	if !gotObject.CreatedAt.Equal(obj.CreatedAt) {
		t.Errorf(
			"CreatedAt = %v, want %v",
			gotObject.CreatedAt,
			obj.CreatedAt,
		)
	}

	if gotPayload.ObjectID != payload.ObjectID {
		t.Errorf(
			"Payload.ObjectID = %q, want %q",
			gotPayload.ObjectID,
			payload.ObjectID,
		)
	}

	if !bytes.Equal(gotPayload.Data, payload.Data) {
		t.Errorf(
			"Payload.Data = %v, want %v",
			gotPayload.Data,
			payload.Data,
		)
	}
}

func testGetNotFound(
	t *testing.T,
	repo domain.ObjectRepository,
) {
	t.Helper()

	_, _, err := repo.Get(
		context.Background(),
		"test",
		"does-not-exist",
	)

	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf(
			"Get() error = %v, want domain.ErrNotFound",
			err,
		)
	}
}

func testGetByHash(
	t *testing.T,
	repo domain.ObjectRepository,
) {
	t.Helper()

	ctx := context.Background()

	obj := testObject("object-1", "hash-1")
	payload := testPayload(
		obj.ID,
		[]byte("hello"),
	)

	if err := repo.Save(ctx, obj, payload); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	gotObject, gotPayload, err := repo.GetByHash(
		ctx,
		obj.Path,
		obj.Hash,
	)
	if err != nil {
		t.Fatalf("GetByHash() error = %v", err)
	}

	if gotObject.ID != obj.ID {
		t.Errorf(
			"ID = %q, want %q",
			gotObject.ID,
			obj.ID,
		)
	}

	if gotObject.Hash != obj.Hash {
		t.Errorf(
			"Hash = %q, want %q",
			gotObject.Hash,
			obj.Hash,
		)
	}

	if gotObject.Path != obj.Path {
		t.Errorf(
			"Path = %q, want %q",
			gotObject.Path,
			obj.Path,
		)
	}

	if gotPayload.ObjectID != obj.ID {
		t.Errorf(
			"Payload.ObjectID = %q, want %q",
			gotPayload.ObjectID,
			obj.ID,
		)
	}
}

func testGetByHashNotFound(
	t *testing.T,
	repo domain.ObjectRepository,
) {
	t.Helper()

	_, _, err := repo.GetByHash(
		context.Background(),
		"test",
		"does-not-exist",
	)

	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf(
			"GetByHash() error = %v, want domain.ErrNotFound",
			err,
		)
	}
}

func testDuplicateHash(
	t *testing.T,
	repo domain.ObjectRepository,
) {
	t.Helper()

	ctx := context.Background()

	obj1 := testObject("object-1", "same-hash")
	obj2 := testObject("object-2", "same-hash")

	if err := repo.Save(
		ctx,
		obj1,
		testPayload(obj1.ID, []byte("one")),
	); err != nil {
		t.Fatalf(
			"first Save() error = %v",
			err,
		)
	}

	err := repo.Save(
		ctx,
		obj2,
		testPayload(obj2.ID, []byte("two")),
	)

	if !errors.Is(err, domain.ErrDuplicate) {
		t.Fatalf(
			"second Save() error = %v, want domain.ErrDuplicate",
			err,
		)
	}
}

func testDuplicateID(
	t *testing.T,
	repo domain.ObjectRepository,
) {
	t.Helper()

	ctx := context.Background()

	obj1 := testObject("same-id", "hash-1")
	obj2 := testObject("same-id", "hash-2")

	if err := repo.Save(
		ctx,
		obj1,
		testPayload(obj1.ID, []byte("one")),
	); err != nil {
		t.Fatalf(
			"first Save() error = %v",
			err,
		)
	}

	err := repo.Save(
		ctx,
		obj2,
		testPayload(obj2.ID, []byte("two")),
	)

	if !errors.Is(err, domain.ErrDuplicate) {
		t.Fatalf(
			"second Save() error = %v, want domain.ErrDuplicate",
			err,
		)
	}
}

func testDelete(
	t *testing.T,
	repo domain.ObjectRepository,
) {
	t.Helper()

	ctx := context.Background()

	obj := testObject("object-1", "hash-1")
	payload := testPayload(
		obj.ID,
		[]byte("hello"),
	)

	if err := repo.Save(ctx, obj, payload); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if err := repo.Delete(
		ctx,
		obj.Path,
		obj.ID,
	); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, _, err := repo.Get(
		ctx,
		obj.Path,
		obj.ID,
	)

	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf(
			"Get() after Delete() error = %v, want domain.ErrNotFound",
			err,
		)
	}

	_, _, err = repo.GetByHash(
		ctx,
		obj.Path,
		obj.Hash,
	)

	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf(
			"GetByHash() after Delete() error = %v, want domain.ErrNotFound",
			err,
		)
	}
}

func testDeleteNotFound(
	t *testing.T,
	repo domain.ObjectRepository,
) {
	t.Helper()

	err := repo.Delete(
		context.Background(),
		"test",
		"does-not-exist",
	)

	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf(
			"Delete() error = %v, want domain.ErrNotFound",
			err,
		)
	}
}

func testListEmpty(
	t *testing.T,
	repo domain.ObjectRepository,
) {
	t.Helper()

	objects, err := repo.List(
		context.Background(),
		"test",
	)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(objects) != 0 {
		t.Fatalf(
			"List() returned %d objects, want 0",
			len(objects),
		)
	}
}

func testList(
	t *testing.T,
	repo domain.ObjectRepository,
) {
	t.Helper()

	ctx := context.Background()

	objects := []*domain.Object{
		testObject("object-1", "hash-1"),
		testObject("object-2", "hash-2"),
		testObject("object-3", "hash-3"),
	}

	for _, obj := range objects {
		payload := testPayload(
			obj.ID,
			[]byte(obj.ID),
		)

		if err := repo.Save(ctx, obj, payload); err != nil {
			t.Fatalf(
				"Save(%q) error = %v",
				obj.ID,
				err,
			)
		}
	}

	got, err := repo.List(ctx, "test")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(got) != len(objects) {
		t.Fatalf(
			"List() returned %d objects, want %d",
			len(got),
			len(objects),
		)
	}

	found := make(map[string]bool)

	for _, obj := range got {
		found[obj.ID] = true
	}

	for _, obj := range objects {
		if !found[obj.ID] {
			t.Errorf(
				"List() missing object %q",
				obj.ID,
			)
		}
	}
}

func testPathIsolation(
	t *testing.T,
	repo domain.ObjectRepository,
) {
	t.Helper()

	ctx := context.Background()

	obj1 := testObject("same-id", "hash-a")
	obj1.Path = "path-a"

	obj2 := testObject("same-id", "hash-b")
	obj2.Path = "path-b"

	payload1 := testPayload(
		obj1.ID,
		[]byte("path a"),
	)

	payload2 := testPayload(
		obj2.ID,
		[]byte("path b"),
	)

	if err := repo.Save(ctx, obj1, payload1); err != nil {
		t.Fatalf(
			"Save(path-a) error = %v",
			err,
		)
	}

	if err := repo.Save(ctx, obj2, payload2); err != nil {
		t.Fatalf(
			"Save(path-b) error = %v",
			err,
		)
	}

	gotObject, gotPayload, err := repo.Get(
		ctx,
		"path-a",
		"same-id",
	)
	if err != nil {
		t.Fatalf(
			"Get(path-a) error = %v",
			err,
		)
	}

	if gotObject.Path != "path-a" {
		t.Errorf(
			"Get(path-a) returned Path = %q, want %q",
			gotObject.Path,
			"path-a",
		)
	}

	if !bytes.Equal(
		gotPayload.Data,
		[]byte("path a"),
	) {
		t.Errorf(
			"Get(path-a) payload = %q, want %q",
			gotPayload.Data,
			[]byte("path a"),
		)
	}

	gotObject, gotPayload, err = repo.Get(
		ctx,
		"path-b",
		"same-id",
	)
	if err != nil {
		t.Fatalf(
			"Get(path-b) error = %v",
			err,
		)
	}

	if gotObject.Path != "path-b" {
		t.Errorf(
			"Get(path-b) returned Path = %q, want %q",
			gotObject.Path,
			"path-b",
		)
	}

	if !bytes.Equal(
		gotPayload.Data,
		[]byte("path b"),
	) {
		t.Errorf(
			"Get(path-b) payload = %q, want %q",
			gotPayload.Data,
			[]byte("path b"),
		)
	}
}

func testListPathIsolation(
	t *testing.T,
	repo domain.ObjectRepository,
) {
	t.Helper()

	ctx := context.Background()

	obj1 := testObject("object-1", "hash-1")
	obj1.Path = "path-a"

	obj2 := testObject("object-2", "hash-2")
	obj2.Path = "path-a"

	obj3 := testObject("object-3", "hash-3")
	obj3.Path = "path-b"

	objects := []*domain.Object{
		obj1,
		obj2,
		obj3,
	}

	for _, obj := range objects {
		if err := repo.Save(
			ctx,
			obj,
			testPayload(
				obj.ID,
				[]byte(obj.ID),
			),
		); err != nil {
			t.Fatalf(
				"Save(%q) error = %v",
				obj.ID,
				err,
			)
		}
	}

	got, err := repo.List(ctx, "path-a")
	if err != nil {
		t.Fatalf(
			"List(path-a) error = %v",
			err,
		)
	}

	if len(got) != 2 {
		t.Fatalf(
			"List(path-a) returned %d objects, want 2",
			len(got),
		)
	}

	found := make(map[string]bool)

	for _, obj := range got {
		if obj.Path != "path-a" {
			t.Errorf(
				"List(path-a) returned object with Path = %q",
				obj.Path,
			)
		}

		found[obj.ID] = true
	}

	if !found["object-1"] {
		t.Error("List(path-a) missing object-1")
	}

	if !found["object-2"] {
		t.Error("List(path-a) missing object-2")
	}

	if found["object-3"] {
		t.Error("List(path-a) incorrectly returned object-3")
	}
}

func testBinaryPayload(
	t *testing.T,
	repo domain.ObjectRepository,
) {
	t.Helper()

	ctx := context.Background()

	data := []byte{
		0x00,
		0x01,
		0x02,
		0x7f,
		0x80,
		0xfe,
		0xff,
	}

	obj := testObject(
		"binary-object",
		"binary-hash",
	)

	payload := testPayload(
		obj.ID,
		data,
	)

	if err := repo.Save(ctx, obj, payload); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	_, got, err := repo.Get(
		ctx,
		obj.Path,
		obj.ID,
	)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if !bytes.Equal(got.Data, data) {
		t.Fatalf(
			"payload data changed: got %v, want %v",
			got.Data,
			data,
		)
	}
}

func testCancelledContext(
	t *testing.T,
	repo domain.ObjectRepository,
) {
	t.Helper()

	ctx, cancel := context.WithCancel(
		context.Background(),
	)
	cancel()

	obj := testObject("object-1", "hash-1")
	payload := testPayload(
		obj.ID,
		[]byte("hello"),
	)

	if err := repo.Save(
		ctx,
		obj,
		payload,
	); !errors.Is(err, context.Canceled) {
		t.Errorf(
			"Save() error = %v, want context.Canceled",
			err,
		)
	}

	if _, _, err := repo.Get(
		ctx,
		obj.Path,
		obj.ID,
	); !errors.Is(err, context.Canceled) {
		t.Errorf(
			"Get() error = %v, want context.Canceled",
			err,
		)
	}

	if _, _, err := repo.GetByHash(
		ctx,
		obj.Path,
		obj.Hash,
	); !errors.Is(err, context.Canceled) {
		t.Errorf(
			"GetByHash() error = %v, want context.Canceled",
			err,
		)
	}

	if _, err := repo.List(
		ctx,
		obj.Path,
	); !errors.Is(err, context.Canceled) {
		t.Errorf(
			"List() error = %v, want context.Canceled",
			err,
		)
	}

	if err := repo.Delete(
		ctx,
		obj.Path,
		obj.ID,
	); !errors.Is(err, context.Canceled) {
		t.Errorf(
			"Delete() error = %v, want context.Canceled",
			err,
		)
	}
}

func testObject(
	id string,
	hash string,
) *domain.Object {
	return &domain.Object{
		ID:   id,
		Hash: hash,
		Path: "test",
		Size: 11,
		CreatedAt: time.Date(
			2026,
			time.January,
			1,
			12,
			0,
			0,
			0,
			time.UTC,
		),
		Transforms: []domain.Transform{
			{
				Name:    "test-transform",
				Version: "1.0",
				Params: map[string]string{
					"foo": "bar",
				},
			},
		},
	}
}

func testPayload(
	objectID string,
	data []byte,
) *domain.Payload {
	return &domain.Payload{
		ObjectID: objectID,
		Data:     data,
	}
}
