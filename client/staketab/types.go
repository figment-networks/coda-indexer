package staketab

type Providers struct {
	ProvidersCount   uint64            `json:"providers_count"`
	StakingProviders []StakingProvider `json:"staking_providers"`
}

type StakingProvider struct {
	ProviderId      int     `json:"provider_id"`
	ProviderAddress string  `json:"provider_address"`
	ProviderFee     float64 `json:"provider_fee"`
}
