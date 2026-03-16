//ff:type feature=chain type=model
//ff:what vLLM /v1/score 응답 구조체
package chain

type scoreResponse struct {
	Data []scoreData `json:"data"`
}
