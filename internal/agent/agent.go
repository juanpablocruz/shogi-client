package agent

type AgentLevel int8

const (
	Begginer AgentLevel = iota
	Medium
	Pro
)

type Agent interface {
	AskHint(string) (string, error)
	AskMovement(string, AgentLevel) (string, error)
}
