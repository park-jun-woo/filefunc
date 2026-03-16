//ff:type feature=chain type=model
//ff:what vLLM /v1/score 응답의 개별 점수 데이터
package chain

type scoreData struct {
	Index int     `json:"index"`
	Score float64 `json:"score"`
}
