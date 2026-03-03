package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/driva/api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GoldRepository acessa dw_gold_enrichments.
type GoldRepository struct {
	pool *pgxpool.Pool
}

// NewGoldRepository cria um novo repositório.
func NewGoldRepository(pool *pgxpool.Pool) *GoldRepository {
	return &GoldRepository{pool: pool}
}

// Overview retorna KPIs agregados da Gold.
func (r *GoldRepository) Overview(ctx context.Context) (*model.AnalyticsOverview, error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM dw_gold_enrichments`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count gold: %w", err)
	}

	var totalSucesso int
	err = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM dw_gold_enrichments WHERE processamento_sucesso = true`).Scan(&totalSucesso)
	if err != nil {
		return nil, fmt.Errorf("count success: %w", err)
	}

	percentual := 0.0
	if total > 0 {
		percentual = float64(totalSucesso) / float64(total) * 100
	}

	var tempoMedio *float64
	err = r.pool.QueryRow(ctx, `
		SELECT AVG(duracao_processamento_minutos) FROM dw_gold_enrichments WHERE duracao_processamento_minutos IS NOT NULL
	`).Scan(&tempoMedio)
	if err != nil {
		return nil, fmt.Errorf("avg duration: %w", err)
	}

	tm := 0.0
	if tempoMedio != nil {
		tm = *tempoMedio
	}

	// Total por categoria_tamanho_job
	totalPorCategoria := make(map[string]int)
	rows, err := r.pool.Query(ctx, `
		SELECT categoria_tamanho_job, COUNT(*) FROM dw_gold_enrichments GROUP BY categoria_tamanho_job
	`)
	if err != nil {
		return nil, fmt.Errorf("group by categoria: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var cat string
		var cnt int
		if err := rows.Scan(&cat, &cnt); err != nil {
			return nil, fmt.Errorf("scan categoria: %w", err)
		}
		totalPorCategoria[cat] = cnt
	}

	return &model.AnalyticsOverview{
		TotalEnriquecimentos: total,
		TotalSucesso:         totalSucesso,
		PercentualSucesso:    percentual,
		TempoMedioMinutos:    tm,
		TotalPorCategoria:    totalPorCategoria,
	}, nil
}

// ListPaginated retorna Gold enrichments paginados com filtros opcionais.
func (r *GoldRepository) ListPaginated(ctx context.Context, page, limit int, status, categoria string) (*model.PaginatedGold, error) {
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	var where []string
	var filterArgs []interface{}
	n := 1
	if status != "" {
		where = append(where, fmt.Sprintf("status_processamento = $%d", n))
		filterArgs = append(filterArgs, status)
		n++
	}
	if categoria != "" {
		where = append(where, fmt.Sprintf("categoria_tamanho_job = $%d", n))
		filterArgs = append(filterArgs, categoria)
		n++
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = " WHERE " + strings.Join(where, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM dw_gold_enrichments" + whereClause
	var total int
	if len(filterArgs) > 0 {
		if err := r.pool.QueryRow(ctx, countQuery, filterArgs...).Scan(&total); err != nil {
			return nil, fmt.Errorf("count gold: %w", err)
		}
	} else {
		if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
			return nil, fmt.Errorf("count gold: %w", err)
		}
	}

	totalPages := (total + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	listArgs := append(filterArgs, limit, offset)
	limitPlaceholder := n
	offsetPlaceholder := n + 1
	listQuery := fmt.Sprintf(`
		SELECT id_enriquecimento, id_workspace, nome_workspace, total_contatos, tipo_contato, status_processamento,
		       data_criacao, data_atualizacao, duracao_processamento_minutos, tempo_por_contato_minutos,
		       processamento_sucesso, categoria_tamanho_job, necessita_reprocessamento, data_atualizacao_dw
		FROM dw_gold_enrichments
		%s
		ORDER BY data_criacao DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, limitPlaceholder, offsetPlaceholder)

	rows, err := r.pool.Query(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, fmt.Errorf("list gold: %w", err)
	}
	defer rows.Close()

	var data []model.GoldEnrichment
	for rows.Next() {
		var e model.GoldEnrichment
		err := rows.Scan(&e.IDEnriquecimento, &e.IDWorkspace, &e.NomeWorkspace, &e.TotalContatos, &e.TipoContato,
			&e.StatusProcessamento, &e.DataCriacao, &e.DataAtualizacao, &e.DuracaoProcessamentoMin, &e.TempoPorContatoMin,
			&e.ProcessamentoSucesso, &e.CategoriaTamanhoJob, &e.NecessitaReprocessamento, &e.DataAtualizacaoDW)
		if err != nil {
			return nil, fmt.Errorf("scan gold: %w", err)
		}
		data = append(data, e)
	}

	return &model.PaginatedGold{
		Data:        data,
		TotalItems:  total,
		TotalPages:  totalPages,
		CurrentPage: page,
		Limit:       limit,
	}, nil
}
