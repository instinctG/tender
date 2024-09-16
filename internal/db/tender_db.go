package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/instinctG/tender/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
)

func (d *Database) GetTenders(params model.GetTendersParams) ([]*model.Tender, error) {
	var tenders []*model.Tender
	ctx := context.Background()
	conn, err := d.Client.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v", err)
		return nil, err
	}
	defer conn.Release()

	query := `
        SELECT id, name, description, status, service_type, version, created_at
        FROM tender
        WHERE status = 'Published'
    `

	var args []interface{}

	if len(params.ServiceType) > 0 {
		query += " AND service_type = ANY($1)"
		args = append(args, params.ServiceType)
		query += " ORDER BY name ASC LIMIT $2 OFFSET $3"
	} else {
		query += " ORDER BY name ASC LIMIT $1 OFFSET $2"
	}

	args = append(args, params.Limit, params.Offset)

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tender model.Tender
		err = rows.Scan(
			&tender.Id,
			&tender.Name,
			&tender.Description,
			&tender.Status,
			&tender.ServiceType,
			&tender.Version,
			&tender.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tenders = append(tenders, &tender)
	}

	return tenders, nil
}

func (d *Database) GetUserTenders(params model.GetUserTendersParams) ([]*model.Tender, error) {
	var tenders []*model.Tender
	ctx := context.Background()
	conn, err := d.Client.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v", err)
		return nil, err
	}
	defer conn.Release()

	query := `SELECT id, name, description, service_type, status, version, created_at
        FROM tender
        WHERE creator_username=$1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
        `

	rows, err := conn.Query(ctx, query, params.Username, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tender model.Tender
		if err = rows.Scan(
			&tender.Id,
			&tender.Name,
			&tender.Description,
			&tender.ServiceType,
			&tender.Status,
			&tender.Version,
			&tender.CreatedAt,
		); err != nil {
			return nil, err
		}
		tenders = append(tenders, &tender)
	}

	return tenders, nil
}

func (d *Database) CreateTender(params model.CreateTenderJSONBody) (*model.Tender, error) {
	var createdTender model.Tender
	ctx := context.Background()
	conn, err := d.Client.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v", err)
		return nil, err
	}
	defer conn.Release()

	var isResponsible bool
	err = conn.QueryRow(ctx, `
        SELECT EXISTS (
            SELECT 1 
            FROM organization_responsible 
            WHERE organization_id = $1 
            AND user_id = (SELECT id FROM employee WHERE username = $2)
        )`,
		params.OrganizationId, params.CreatorUsername).Scan(&isResponsible)

	if err != nil {
		return nil, err
	}

	if !isResponsible {
		return nil, errors.New("user is not responsible for this organization")
	}

	query := `
        INSERT INTO tender (name, description, service_type,organization_id, creator_username)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id,name,description,status,service_type,version,created_at;
    `

	row := conn.QueryRow(ctx, query, params.Name, params.Description, params.ServiceType, params.OrganizationId, params.CreatorUsername)

	err = row.Scan(
		&createdTender.Id,
		&createdTender.Name,
		&createdTender.Description,
		&createdTender.Status,
		&createdTender.ServiceType,
		&createdTender.Version,
		&createdTender.CreatedAt,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			log.Println(pgErr.Message)
		}
		log.Fatalf("Unable to scan JSON body: %w", err)
		return nil, err
	}

	return &createdTender, nil
}

func (d *Database) EditTender(tenderId string, par model.EditTenderParams, params model.EditTenderJSONBody) (*model.Tender, error) {
	var updatedTender model.Tender
	ctx := context.Background()
	conn, err := d.Client.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v", err)
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(ctx, `
        UPDATE tender
        SET name = COALESCE($1, name),
            description = COALESCE($2, description),
            service_type = COALESCE($3, service_type),
            version = version + 1
        WHERE id = $4 AND creator_username = $5
        RETURNING id, name, description, service_type, status, version, created_at
    `,
		func() *string {
			if params.Name != "" {
				return &params.Name
			} else {
				return nil
			}
		}(),
		func() *string {
			if params.Description != "" {
				return &params.Description
			} else {
				return nil
			}
		}(),
		func() *string {
			if params.ServiceType != "" {
				return &params.ServiceType
			} else {
				return nil
			}
		}(),
		tenderId, par.Username).Scan(
		&updatedTender.Id,
		&updatedTender.Name,
		&updatedTender.Description,
		&updatedTender.ServiceType,
		&updatedTender.Status,
		&updatedTender.Version,
		&updatedTender.CreatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("tender not found")
		}
		return nil, err
	}

	return &updatedTender, nil
}

func (d *Database) RollbackTender(tenderId string, version int32, params model.RollbackTenderParams) (*model.Tender, error) {
	return nil, nil
}

func (d *Database) GetTenderStatus(tenderId string, params model.GetTenderStatusParams) (string, error) {
	var status string
	ctx := context.Background()
	conn, err := d.Client.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v", err)
		return "", err
	}
	defer conn.Release()

	query := `SELECT status 
			  FROM tender
              WHERE id = $1
              `

	if params.Username != "" {
		query += ` AND creator_username = $2`
		err = conn.QueryRow(ctx, query, tenderId, params.Username).Scan(&status)
	} else {
		err = conn.QueryRow(ctx, query, tenderId).Scan(&status)
	}

	if err != nil {
		if err.Error() == "no rows in result set" {
			return "", errors.New("tender not found")
		}
		return "", err
	}
	return status, nil
}

func (d *Database) UpdateTenderStatus(tenderId string, params model.UpdateTenderStatusParams) (*model.Tender, error) {
	var updatedTender model.Tender
	ctx := context.Background()
	conn, err := d.Client.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v", err)
		return nil, err
	}
	defer conn.Release()

	query := `UPDATE tender SET
			  status = $1
			  WHERE id = $2 AND creator_username = $3
			  RETURNING id,name,description,status,service_type,version,created_at
              `

	err = conn.QueryRow(ctx, query, params.Status, tenderId, params.Username).Scan(
		&updatedTender.Id,
		&updatedTender.Name,
		&updatedTender.Description,
		&updatedTender.Status,
		&updatedTender.ServiceType,
		&updatedTender.Version,
		&updatedTender.CreatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, errors.New("tender not found")
		}
		return nil, err
	}

	return &updatedTender, nil
}
