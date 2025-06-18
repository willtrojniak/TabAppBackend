package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/willtrojniak/TabAppBackend/models"
	"github.com/willtrojniak/TabAppBackend/services"
)

func (q *PgxQueries) CreateSubstitutionGroup(ctx context.Context, data *models.SubstitutionGroupCreate) error {
	return q.WithTx(ctx, func(q *PgxQueries) error {
		row := q.tx.QueryRow(ctx, `
    INSERT INTO item_substitution_groups (shop_id, name) VALUES (@shopId, @name) RETURNING id`,
			pgx.NamedArgs{
				"shopId": data.ShopId,
				"name":   data.Name,
			})
		var substitutionGroupId int
		err := row.Scan(&substitutionGroupId)

		if err != nil {
			return handlePgxError(err)
		}

		err = q.setSubstitutionGroupSubstitutions(ctx, data.ShopId, substitutionGroupId, data.SubstitutionItemIds)
		if err != nil {
			return err
		}

		return nil
	})
}

func (q *PgxQueries) UpdateSubstitutionGroup(ctx context.Context, shopId int, substitutionGroupId int, data *models.SubstitutionGroupUpdate) error {
	return q.WithTx(ctx, func(q *PgxQueries) error {
		result, err := q.tx.Exec(ctx, `
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
		err = q.setSubstitutionGroupSubstitutions(ctx, shopId, substitutionGroupId, data.SubstitutionItemIds)
		if err != nil {
			return err
		}
		return nil
	})
}

func (q *PgxQueries) GetSubstitutionGroups(ctx context.Context, shopId int) ([]models.SubstitutionGroup, error) {
	rows, err := q.tx.Query(ctx, `
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

	data, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.SubstitutionGroup])
	if err != nil {
		return nil, handlePgxError(err)
	}

	return data, nil
}

func (q *PgxQueries) DeleteSubstitutionGroup(ctx context.Context, shopId int, substitutionGroupId int) error {
	result, err := q.tx.Exec(ctx, `
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

func (q *PgxQueries) setSubstitutionGroupSubstitutions(ctx context.Context, shopId int, substitutionGroupId int, substitutionItemIds []int) error {
	_, err := q.tx.Exec(ctx, `
    CREATE TEMPORARY TABLE _temp_upsert_item_substitution_groups_to_items (LIKE item_substitution_groups_to_items INCLUDING ALL ) ON COMMIT DROP`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = q.tx.CopyFrom(ctx, pgx.Identifier{"_temp_upsert_item_substitution_groups_to_items"},
		[]string{"shop_id", "substitution_group_id", "item_id", "index"}, pgx.CopyFromSlice(len(substitutionItemIds), func(i int) ([]any, error) {
			return []any{shopId, substitutionGroupId, substitutionItemIds[i], i}, nil
		}))
	if err != nil {
		return handlePgxError(err)
	}

	_, err = q.tx.Exec(ctx, `
    INSERT INTO item_substitution_groups_to_items SELECT * FROM _temp_upsert_item_substitution_groups_to_items ON CONFLICT (shop_id, substitution_group_id, item_id) DO UPDATE
    SET index = excluded.index`)
	if err != nil {
		return handlePgxError(err)
	}

	_, err = q.tx.Exec(ctx,
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
