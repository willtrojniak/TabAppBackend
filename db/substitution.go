package db

import (
	"context"

	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/WilliamTrojniak/TabAppBackend/types"
	"github.com/jackc/pgx/v5"
)

func (s *PgxStore) CreateSubstitutionGroup(ctx context.Context, data *types.SubstitutionGroupCreate) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return handlePgxError(err)
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, `
    INSERT INTO item_substitution_groups (shop_id, name) VALUES (@shopId, @name) RETURNING id`,
		pgx.NamedArgs{
			"shopId": data.ShopId,
			"name":   data.Name,
		})
	var substitutionGroupId int
	err = row.Scan(&substitutionGroupId)

	if err != nil {
		return handlePgxError(err)
	}

	err = s.setSubstitutionGroupSubstitutions(ctx, tx, data.ShopId, substitutionGroupId, data.SubstitutionItemIds)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (s *PgxStore) UpdateSubstitutionGroup(ctx context.Context, shopId int, substitutionGroupId int, data *types.SubstitutionGroupUpdate) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return handlePgxError(err)
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx, `
    UPDATE item_substitution_groups SET name = @name
    WHERE id = @id AND shop_id = @shopId`,
		pgx.NamedArgs{
			"shopId": shopId,
			"id":     substitutionGroupId,
			"name":   data.Name,
		})

	if err != nil {
		return handlePgxError(err)
	}

	if result.RowsAffected() == 0 {
		return services.NewNotFoundServiceError(nil)
	}
	err = s.setSubstitutionGroupSubstitutions(ctx, tx, shopId, substitutionGroupId, data.SubstitutionItemIds)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}

func (s *PgxStore) GetSubstitutionGroups(ctx context.Context, shopId int) ([]types.SubstitutionGroup, error) {
	rows, err := s.pool.Query(ctx, `
    SELECT item_substitution_groups.name, item_substitution_groups.id,
    COALESCE(json_agg(items ORDER BY item_substitution_groups_to_items.index) FILTER (WHERE items.id IS NOT NULL), '[]') AS substitutions
    FROM item_substitution_groups
    LEFT JOIN item_substitution_groups_to_items ON
      item_substitution_groups.id = item_substitution_groups_to_items.substitution_group_id
      AND item_substitution_groups.shop_id = item_substitution_groups_to_items.shop_id
    LEFT JOIN items ON items.id = item_substitution_groups_to_items.item_id AND items.shop_id = item_substitution_groups_to_items.shop_id
    WHERE item_substitution_groups.shop_id = @shopId
    GROUP BY item_substitution_groups.shop_id, item_substitution_groups.id`,
		pgx.NamedArgs{
			"shopId": shopId,
		})

	if err != nil {
		return nil, handlePgxError(err)
	}

	data, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[types.SubstitutionGroup])
	if err != nil {
		return nil, handlePgxError(err)
	}

	return data, nil
}

func (s *PgxStore) DeleteSubstitutionGroup(ctx context.Context, shopId int, substitutionGroupId int) error {
	result, err := s.pool.Exec(ctx, `
    DELETE FROM item_substitution_groups 
    WHERE id = @id AND shop_id = @shopId`,
		pgx.NamedArgs{
			"shopId": shopId,
			"id":     substitutionGroupId,
		})

	if err != nil {
		return handlePgxError(err)
	}

	if result.RowsAffected() == 0 {
		return services.NewNotFoundServiceError(nil)
	}

	return nil
}

func (s *PgxStore) setSubstitutionGroupSubstitutions(ctx context.Context, tx pgx.Tx, shopId int, substitutionGroupId int, substitutionItemIds []int) error {
	_, err := tx.Exec(ctx, `
    CREATE TEMPORARY TABLE _temp_upsert_item_substitution_groups_to_items (LIKE item_substitution_groups_to_items INCLUDING ALL ) ON COMMIT DROP`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"_temp_upsert_item_substitution_groups_to_items"},
		[]string{"shop_id", "substitution_group_id", "item_id", "index"}, pgx.CopyFromSlice(len(substitutionItemIds), func(i int) ([]any, error) {
			return []any{shopId, substitutionGroupId, substitutionItemIds[i], i}, nil
		}))
	if err != nil {
		return handlePgxError(err)
	}

	_, err = tx.Exec(ctx, `
    INSERT INTO item_substitution_groups_to_items SELECT * FROM _temp_upsert_item_substitution_groups_to_items ON CONFLICT (shop_id, substitution_group_id, item_id) DO UPDATE
    SET index = excluded.index`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = tx.Exec(ctx,
		`DELETE FROM item_substitution_groups_to_items WHERE substitution_group_id = @substitutionGroupId AND shop_id = @shopId AND NOT (item_id = ANY (@substitutionItemIds))`,
		pgx.NamedArgs{
			"shopId":              shopId,
			"substitutionGroupId": substitutionGroupId,
			"substitutionItemIds": substitutionItemIds,
		})
	if err != nil {
		return handlePgxError(err)
	}

	return nil
}
