package repository

import (
	"context"
	"fmt"

	"github.com/driva/api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SeedRepository acessa api_enrichments_seed.
type SeedRepository struct {
	pool *pgxpool.Pool
}

// NewSeedRepository cria um novo repositório.
func NewSeedRepository(pool *pgxpool.Pool) *SeedRepository {
	return &SeedRepository{pool: pool}
}

// ListPaginated retorna enrichments paginados.
func (r *SeedRepository) ListPaginated(ctx context.Context, page, limit int) (*model.PaginatedEnrichments, error) {
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	// Total de registros
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM api_enrichments_seed`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count seed: %w", err)
	}

	totalPages := (total + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	// Dados
	rows, err := r.pool.Query(ctx, `
		SELECT id, id_workspace, workspace_name, total_contacts, contact_type, status, created_at, updated_at
		FROM api_enrichments_seed
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list seed: %w", err)
	}
	defer rows.Close()

	var data []model.EnrichmentSeed
	for rows.Next() {
		var e model.EnrichmentSeed
		err := rows.Scan(&e.ID, &e.IDWorkspace, &e.WorkspaceName, &e.TotalContacts, &e.ContactType, &e.Status, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan seed: %w", err)
		}
		data = append(data, e)
	}

	return &model.PaginatedEnrichments{
		Data:        data,
		TotalItems:  total,
		TotalPages:  totalPages,
		CurrentPage: page,
		Limit:       limit,
	}, nil
}
