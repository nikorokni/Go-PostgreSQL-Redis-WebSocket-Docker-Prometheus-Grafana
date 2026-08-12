package risk

import (
	"math"
	"time"
)

type Snapshot struct {
	Symbol            string    `json:"symbol"`
	Price             float64   `json:"price"`
	LiquidityUSD      float64   `json:"liquidity_usd"`
	Volume24hUSD      float64   `json:"volume_24h_usd"`
	CollateralRatio   float64   `json:"collateral_ratio"`
	OracleAgeSeconds  float64   `json:"oracle_age_seconds"`
	OracleDispersion  float64   `json:"oracle_dispersion"`
	RedemptionRate24h float64   `json:"redemption_rate_24h"`
	Timestamp         time.Time `json:"timestamp"`
}

type Assessment struct {
	Symbol      string             `json:"symbol"`
	Score       float64            `json:"score"`
	Level       string             `json:"level"`
	DepegBPS    float64            `json:"depeg_bps"`
	Components  map[string]float64 `json:"components"`
	Alerts      []string           `json:"alerts"`
	AssessedAt  time.Time          `json:"assessed_at"`
}

func clamp(x float64) float64 { return math.Max(0, math.Min(100, x)) }

func Assess(s Snapshot) Assessment {
	depeg := math.Abs(s.Price-1) * 10000
	priceRisk := clamp(depeg / 5)
	liquidityRisk := clamp((1 - math.Min(1, s.LiquidityUSD/50_000_000)) * 100)
	turnoverRisk := clamp(s.Volume24hUSD / math.Max(s.LiquidityUSD, 1) * 60)
	collateralRisk := clamp((1.5-s.CollateralRatio)/0.5*100)
	oracleRisk := clamp(s.OracleAgeSeconds/3 + s.OracleDispersion*5000)
	runRisk := clamp(s.RedemptionRate24h * 400)
	c := map[string]float64{"price": priceRisk, "liquidity": liquidityRisk, "turnover": turnoverRisk, "collateral": collateralRisk, "oracle": oracleRisk, "bank_run": runRisk}
	score := .30*priceRisk + .20*liquidityRisk + .10*turnoverRisk + .15*collateralRisk + .15*oracleRisk + .10*runRisk
	level := "LOW"
	if score >= 75 { level = "CRITICAL" } else if score >= 50 { level = "HIGH" } else if score >= 25 { level = "MEDIUM" }
	alerts := []string{}
	if depeg >= 100 { alerts = append(alerts, "price deviation exceeds 1%") }
	if oracleRisk >= 60 { alerts = append(alerts, "oracle data is stale or divergent") }
	if runRisk >= 60 { alerts = append(alerts, "abnormal redemption pressure") }
	return Assessment{Symbol:s.Symbol, Score:math.Round(score*100)/100, Level:level, DepegBPS:math.Round(depeg*100)/100, Components:c, Alerts:alerts, AssessedAt:time.Now().UTC()}
}
