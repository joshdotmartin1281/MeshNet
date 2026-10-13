package tags

import (
	"context"
	"errors"
	"net/url"

	"MeshNet/query"
)

type SearchByTag struct{}

func NewSearchByTag() SearchByTag {
	return SearchByTag{}
}

func (SearchByTag) Name() string {
	return "tags"
}

func (SearchByTag) Init(ctx context.Context, store query.RelationalStore) error {
	_, err := store.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS tags (
			object_id TEXT NOT NULL,
			tag       TEXT NOT NULL,

			FOREIGN KEY (object_id)
				REFERENCES objects(id)
				ON DELETE CASCADE,

			UNIQUE (object_id, tag)
		);

		CREATE INDEX IF NOT EXISTS idx_tags_tag
			ON tags(tag);
	`)

	return err
}

func (SearchByTag) Run(ctx context.Context, store query.RelationalStore, params url.Values) (any, error) {
	tag := params.Get("tag")
	collection := params.Get("collection")

	if tag == "" {
		return nil, errors.New("tag is required")
	}

	rows, err := store.QueryContext(ctx, `
		SELECT t.object_id
		FROM tags t
		JOIN objects o ON o.id = t.object_id
		WHERE t.tag = ?
		  AND o.collection_id = ?
	`, tag, collection)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]string, 0)

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return map[string]any{"object_ids": ids}, nil
}

func (SearchByTag) Mutate(ctx context.Context, store query.RelationalStore, params url.Values) (any, error) {
	objectID := params.Get("object_id")
	tag := params.Get("tag")

	if objectID == "" || tag == "" {
		return nil, errors.New("object_id and tag are required")
	}

	_, err := store.ExecContext(ctx, `
		INSERT INTO tags (object_id, tag)
		VALUES (?, ?)
		ON CONFLICT(object_id, tag) DO NOTHING
	`, objectID, tag)
	if err != nil {
		return nil, err
	}

	return map[string]any{"added": true}, nil
}
