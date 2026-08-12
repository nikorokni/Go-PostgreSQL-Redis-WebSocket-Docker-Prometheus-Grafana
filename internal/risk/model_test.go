package risk

import "testing"

func TestHealthyStablecoin(t *testing.T) {
	a := Assess(Snapshot{Symbol:"USDX", Price:1, LiquidityUSD:100e6, Volume24hUSD:5e6, CollateralRatio:1.8, OracleAgeSeconds:5})
	if a.Level != "LOW" || a.Score >= 25 { t.Fatalf("unexpected assessment: %+v", a) }
}

func TestDepegRaisesRisk(t *testing.T) {
	healthy := Assess(Snapshot{Price:1, LiquidityUSD:100e6, CollateralRatio:1.8})
	depeg := Assess(Snapshot{Price:.94, LiquidityUSD:3e6, Volume24hUSD:10e6, CollateralRatio:1.05, OracleAgeSeconds:300, RedemptionRate24h:.3})
	if depeg.Score <= healthy.Score || len(depeg.Alerts) == 0 { t.Fatalf("risk did not increase") }
}
