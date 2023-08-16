package repository

import (
	"context"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"project/api"
)

type PostgresSQLItemRepository struct {
	db *sqlx.DB
}

type ItemRow struct {
	Uuid    uuid.UUID `db:"uuid"`
	Payload api.Item  `db:"payload"`
}

func ItemPostgresSQLItemRepository(db *sqlx.DB) *PostgresSQLItemRepository {
	return &PostgresSQLItemRepository{
		db: db,
	}
}

func (p PostgresSQLItemRepository) GetItem(ctx context.Context, uuid uuid.UUID) (api.Item, error) {
	item := api.Item{}

	sqlQuery, _, _ := sq.
		Select("payload").
		From("task.items").
		Where(sq.Eq{"uuid": uuid.String()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	//sqlQuery := `SELECT payload FROM task.items where uuid = $1 LIMIT 1`

	err := p.db.GetContext(ctx, &item, sqlQuery)
	if err != nil {
		return item,
			fmt.Errorf("can't get item from postgressql repository uuid=%v psql_db error:%w", uuid.String(), err)
	}

	return item, nil
}

// GetItems don't use after cache initialisation. Not transactional because rowsCount could change
func (p PostgresSQLItemRepository) GetItems(ctx context.Context) ([]ItemRow, error) {
	//sqlQuery := `SELECT count(*) FROM task.items`
	var rowsCount int
	countSqlQuery := `SELECT count(*) FROM task.items`
	err := p.db.Get(&rowsCount, countSqlQuery)
	if err != nil {
		return []ItemRow{},
			fmt.Errorf("can't get rows count from task.items postgressql in psql_db repository error: %w", err)
	}

	//sqlQuery := `SELECT * FROM task.items`
	sqlRowsQuery, _, _ := sq.
		Select("*").
		From("task.items").
		ToSql()
	rows, err := p.db.QueryxContext(ctx, sqlRowsQuery)
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	if err != nil {
		return []ItemRow{},
			fmt.Errorf("can't get rows from task.items postgressql in psql_db repository error: %w", err)
	}

	itemRows := make([]ItemRow, 0, rowsCount)
	for rows.Next() {
		/*var rowUUIDBytes string
		var rowItemBytes []byte
		err = rows.Scan(&rowUUIDBytes, &rowItemBytes)*/

		type ItemRowDb struct {
			Uuid    uuid.UUID `db:"uuid"`
			Payload []byte    `db:"payload"`
		}
		var itemRowDb ItemRowDb
		err = rows.StructScan(&itemRowDb)
		if err != nil {
			return nil,
				fmt.Errorf("can't serialize row from task.items postgressql in psql_db repository error: %w", err)
		}

		var payloadDb api.Item
		err := json.Unmarshal(itemRowDb.Payload, &payloadDb)
		if err != nil {
			return nil, fmt.Errorf("can't unmarshal payload:jsonb from from task.items postgressql in psql_db repository error: %w", err)
		}
		itemRow := ItemRow{itemRowDb.Uuid, payloadDb}
		itemRows = append(itemRows, itemRow)
	}

	return itemRows, nil
}

func (p PostgresSQLItemRepository) AddItem(ctx context.Context, uuid uuid.UUID, item api.Item) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sqlQuery, _, _ := psql.
		Insert("task.items").
		Columns("uuid", "payload").
		Values(uuid, item).
		Suffix(`ON CONFLICT (uuid) DO UPDATE SET uuid = excluded.uuid, payload = excluded.payload;`).
		ToSql()

	/*sqlQuery :=
	`INSERT INTO task.items (uuid, payload)
	VALUES ($1, $2)
	ON CONFLICT (uuid) DO UPDATE
	SET uuid = excluded.uuid,
		payload = excluded.payload;`*/

	jsonItem, _ := json.Marshal(item)

	p.db.MustExecContext(ctx, sqlQuery, uuid, jsonItem)
}
