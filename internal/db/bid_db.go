package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/instinctG/tender/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
)

func (d *Database) GetUserBids(params model.GetUserBidsParams) ([]*model.Bid, error) {
	var bids []*model.Bid
	ctx := context.Background()
	conn, err := d.Client.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v", err)
		return nil, err
	}
	defer conn.Release()

	query := `SELECT id, name, status, author_type,author_id,version, created_at
		FROM bid
		WHERE author_id = (SELECT id FROM employee WHERE username = $1)
		ORDER BY created_at DESC LIMIT $2 OFFSET $3;
	`

	rows, err := conn.Query(ctx, query, params.Username, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bid model.Bid
		if err = rows.Scan(
			&bid.Id,
			&bid.Name,
			&bid.Status,
			&bid.AuthorType,
			&bid.AuthorId,
			&bid.Version,
			&bid.CreatedAt,
		); err != nil {
			return nil, err
		}
		bids = append(bids, &bid)
	}

	return bids, nil
}

func (d *Database) CreateBid(params model.CreateBidJSONBody) (*model.Bid, error) {
	var createdBid model.Bid
	ctx := context.Background()
	conn, err := d.Client.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v", err)
		return nil, err
	}
	defer conn.Release()

	var isAuthorized bool
	err = conn.QueryRow(ctx, `
        SELECT EXISTS (
            SELECT 1 
            FROM organization_responsible 
            WHERE user_id = $1
        )`,
		params.AuthorId).Scan(&isAuthorized)

	if err != nil {
		return nil, err
	}

	if !isAuthorized {
		return nil, errors.New("user is not authorized to create a bid for this organization")
	}

	query := `
        INSERT INTO bid (name, description, tender_id, author_type,author_id)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id,name,description,status,tender_id,author_type,author_id,version,created_at;
    `

	row := conn.QueryRow(ctx, query, params.Name, params.Description, params.TenderId, params.AuthorType, params.AuthorId)

	err = row.Scan(
		&createdBid.Id,
		&createdBid.Name,
		&createdBid.Description,
		&createdBid.Status,
		&createdBid.TenderId,
		&createdBid.AuthorType,
		&createdBid.AuthorId,
		&createdBid.Version,
		&createdBid.CreatedAt,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			log.Println(pgErr.Message)
		}
		log.Fatalf("Unable to scan from DB: %w", err)
		return nil, err
	}

	return &createdBid, nil
}

func (d *Database) EditBid(bidId string, params model.EditBidParams, body model.EditBidJSONBody) (*model.Bid, error) {
	var updatedBid model.Bid
	ctx := context.Background()
	conn, err := d.Client.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v", err)
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(ctx, `
        UPDATE bid
        SET name = COALESCE($1, name),
            description = COALESCE($2, description),
            version = version + 1
        WHERE id = $3 AND author_id = (SELECT id FROM employee WHERE username = $4)
        RETURNING id, name, status, author_type, author_id, version, created_at
    `,
		func() *string {
			if body.Name != "" {
				return &body.Name
			} else {
				return nil
			}
		}(),
		func() *string {
			if body.Description != "" {
				return &body.Description
			} else {
				return nil
			}
		}(),
		bidId, params.Username).Scan(
		&updatedBid.Id,
		&updatedBid.Name,
		&updatedBid.Status,
		&updatedBid.AuthorType,
		&updatedBid.AuthorId,
		&updatedBid.Version,
		&updatedBid.CreatedAt)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("bid not found")
		}
		return nil, err
	}

	return &updatedBid, nil
}

func (d *Database) SubmitBidFeedback(bidId string, params model.SubmitBidFeedbackParams) *model.Bid {
	return nil
}

func (d *Database) RollbackBid(bidId string, version int32, params model.RollbackBidParams) *model.Bid {
	return nil
}

func (d *Database) GetBidStatus(bidId string, params model.GetBidStatusParams) (string, error) {
	var status string
	ctx := context.Background()
	conn, err := d.Client.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v", err)
		return "", err
	}
	defer conn.Release()

	query := `SELECT status
				FROM bid 
				WHERE id = $1 AND author_id = (SELECT id FROM employee WHERE username = $2)`

	err = conn.QueryRow(ctx, query, bidId, params.Username).Scan(&status)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return "", errors.New("bid not found")
		}
		return "", err
	}
	return status, nil
}

func (d *Database) UpdateBidStatus(bidId string, params model.UpdateBidStatusParams) (*model.Bid, error) {
	var updatedBid model.Bid
	ctx := context.Background()
	conn, err := d.Client.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v", err)
		return nil, err
	}
	defer conn.Release()

	query := `UPDATE bid
				SET status = $1
				WHERE id = $2 AND author_id = (SELECT id FROM employee WHERE username = $3)
				RETURNING id, name, status, author_type, author_id, version, created_at;`

	err = conn.QueryRow(ctx, query, params.Status, bidId, params.Username).Scan(
		&updatedBid.Id,
		&updatedBid.Name,
		&updatedBid.Status,
		&updatedBid.AuthorType,
		&updatedBid.AuthorId,
		&updatedBid.Version,
		&updatedBid.CreatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, errors.New("bid not found")
		}
		return nil, err
	}

	return &updatedBid, nil
}

func (d *Database) SubmitBidDecision(bidId string, params model.SubmitBidDecisionParams) (*model.Bid, error) {
	/*
		var bid model.Bid
			ctx := context.Background()
			conn, err := d.db.Acquire(ctx)
			if err != nil {
				log.Fatalf("Unable to acquire a database connection: %v", err)
				return nil, err
			}
			defer conn.Release()

			var isResponsible bool
			query := `
				SELECT EXISTS (
					SELECT 1
					FROM organization_responsible r
					JOIN bid b ON b.organization_id = r.organization_id
					WHERE r.user_id = (SELECT id FROM employee WHERE username = $1)
					  AND b.id = $2
				)
			`
			err = conn.QueryRow(ctx, query, params.Username, bidId).Scan(&isResponsible)
			if err != nil {
				return &bid, err
			}
			if !isResponsible {
				return &bid, errors.New("user is not responsible for the organization")
			}


			query = `
				INSERT INTO bid (id, user_id, decision)
				VALUES ($1, (SELECT id FROM employee WHERE username = $2), $3, $4)
				ON CONFLICT (bid_id, user_id) DO UPDATE
				SET decision = EXCLUDED.decision
			`


			_, err = conn.Exec(ctx, query, bidId, params.Username, params.Decision)
			if err != nil {
				return &bid, err
			}


			var totalResponsible int
			query = `
				SELECT COUNT(*)
				FROM organization_responsible r
				JOIN bids b ON b.organization_id = r.organization_id
				WHERE b.id = $1
			`
			err = db.QueryRow(ctx, query, bidId).Scan(&totalResponsible)
			if err != nil {
				return bid, err
			}

			// Получаем решения всех ответственных
			var totalAccepted, totalRejected int
			query = `
				SELECT
					SUM(CASE WHEN decision = 'accept' THEN 1 ELSE 0 END),
					SUM(CASE WHEN decision = 'reject' THEN 1 ELSE 0 END)
				FROM decisions
				WHERE bid_id = $1
			`
			err = db.QueryRow(ctx, query, bidId).Scan(&totalAccepted, &totalRejected)
			if err != nil {
				return bid, err
			}

			// Логика отклонения предложения, если есть хотя бы один "reject"
			if totalRejected > 0 {
				query = `
					UPDATE bids
					SET status = 'rejected'
					WHERE id = $1
					RETURNING id, name, description, status, version, created_at
				`
				err = db.QueryRow(ctx, query, bidId).Scan(
					&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.Version, &bid.CreatedAt,
				)
				if err != nil {
					return bid, err
				}
				return bid, nil
			}

			// Проверка кворума (min(3, количество ответственных))
			quorum := min(3, totalResponsible)
			if totalAccepted >= quorum {
				// Согласование предложения, обновление статуса
				query = `
					UPDATE bids
					SET status = 'accepted'
					WHERE id = $1
					RETURNING id, name, description, status, version, created_at
				`
				err = db.QueryRow(ctx, query, bidId).Scan(
					&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.Version, &bid.CreatedAt,
				)
				if err != nil {
					return bid, err
				}

				// Закрываем тендер, если предложение было согласовано
				query = `
					UPDATE tenders
					SET status = 'closed'
					WHERE id = (SELECT tender_id FROM bids WHERE id = $1)
				`
				_, err = db.Exec(ctx, query, bidId)
				if err != nil {
					return bid, err
				}
			}

			return bid, nil
	*/
	return nil, nil
}

func (d *Database) GetBidsForTender(tenderId string, params model.GetBidsForTenderParams) ([]*model.Bid, error) {
	var bids []*model.Bid
	ctx := context.Background()
	conn, err := d.Client.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v", err)
		return nil, err
	}
	defer conn.Release()

	query := `
        SELECT id, name, status, author_type, author_id, version, created_at
        FROM bid 
        WHERE tender_id = $1 and author_id = (SELECT id FROM employee WHERE username = $2)
        LIMIT $3 OFFSET $4
    `

	rows, err := conn.Query(ctx, query, tenderId, params.Username, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bid model.Bid
		if err = rows.Scan(
			&bid.Id,
			&bid.Name,
			&bid.Status,
			&bid.AuthorType,
			&bid.AuthorId,
			&bid.Version,
			&bid.CreatedAt,
		); err != nil {
			return nil, err
		}
		bids = append(bids, &bid)
	}

	return bids, nil
}

func (d *Database) GetBidReviews(tenderId string, params model.GetBidReviewsParams) []*model.BidReview {
	return []*model.BidReview{}
}
