package model

import "time"

// EnrichmentSeed representa um registro da api_enrichments_seed (fonte simulada).
// JSON tags seguem o formato da API (snake_case).
type EnrichmentSeed struct {
	ID            string    `json:"id"`
	IDWorkspace   string    `json:"id_workspace"`
	WorkspaceName string    `json:"workspace_name"`
	TotalContacts int       `json:"total_contacts"`
	ContactType   string    `json:"contact_type"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PaginatedEnrichments resposta paginada do GET /people/v1/enrichments.
type PaginatedEnrichments struct {
	Data         []EnrichmentSeed `json:"data"`
	TotalItems   int              `json:"total_items"`
	TotalPages   int              `json:"total_pages"`
	CurrentPage  int              `json:"current_page"`
	Limit        int              `json:"limit"`
}

// GoldEnrichment representa um registro da dw_gold_enrichments.
type GoldEnrichment struct {
	IDEnriquecimento          string    `json:"id_enriquecimento"`
	IDWorkspace               string    `json:"id_workspace"`
	NomeWorkspace             string    `json:"nome_workspace"`
	TotalContatos             int       `json:"total_contatos"`
	TipoContato               string    `json:"tipo_contato"`
	StatusProcessamento       string    `json:"status_processamento"`
	DataCriacao               time.Time `json:"data_criacao"`
	DataAtualizacao           time.Time `json:"data_atualizacao"`
	DuracaoProcessamentoMin   *float64  `json:"duracao_processamento_minutos,omitempty"`
	TempoPorContatoMin        *float64  `json:"tempo_por_contato_minutos,omitempty"`
	ProcessamentoSucesso      bool      `json:"processamento_sucesso"`
	CategoriaTamanhoJob       string    `json:"categoria_tamanho_job"`
	NecessitaReprocessamento  bool      `json:"necessita_reprocessamento"`
	DataAtualizacaoDW         time.Time `json:"data_atualizacao_dw"`
}

// AnalyticsOverview KPIs agregados da Gold.
type AnalyticsOverview struct {
	TotalEnriquecimentos   int     `json:"total_enriquecimentos"`
	TotalSucesso           int     `json:"total_sucesso"`
	PercentualSucesso      float64 `json:"percentual_sucesso"`
	TempoMedioMinutos      float64 `json:"tempo_medio_minutos"`
	TotalPorCategoria      map[string]int `json:"total_por_categoria"`
}

// PaginatedGold resposta paginada do GET /analytics/enrichments.
type PaginatedGold struct {
	Data         []GoldEnrichment `json:"data"`
	TotalItems   int              `json:"total_items"`
	TotalPages   int              `json:"total_pages"`
	CurrentPage  int              `json:"current_page"`
	Limit        int              `json:"limit"`
}
