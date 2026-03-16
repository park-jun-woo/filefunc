//ff:type feature=chain type=model
//ff:what vLLM /v1/score 요청 구조체
package chain

type scoreRequest struct {
	Model string `json:"model"`
	Text1 string `json:"text_1"`
	Text2 string `json:"text_2"`
}
