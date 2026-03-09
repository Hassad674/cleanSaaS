package ai

type Model string

const (
	ModelClaude Model = "claude"
	ModelGPT    Model = "gpt"
	ModelGemini Model = "gemini"
)

func (m Model) IsValid() bool {
	return m == ModelClaude || m == ModelGPT || m == ModelGemini
}
