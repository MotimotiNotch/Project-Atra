package analyzer

import (
	"fmt"
)

// State は特定の時点における各ノードのエネルギー（活性）レベルを保持します。
type State map[string]float64

// Evaluator は時間ステップ t に沿って宇宙の動的な変化を計算します。
type Evaluator struct {
	space *Space
}

func NewEvaluator(s *Space) *Evaluator {
	return &Evaluator{space: s}
}

// Step は時間 t から t+1 への状態遷移を計算します。
// アトラクタへの「滑落」：活性が高いノードから引力のあるノードへエネルギーが移動する。
func (e *Evaluator) Step(currentState State, deltaTime float64) State {
	nextState := make(State)
	for k, v := range currentState {
		nextState[k] = v
	}

	// 各接続（引力・斥力）に基づくエネルギーの移動
	for _, conn := range e.space.Connections {
		srcEnergy := currentState[conn.Source]

		// 引力 (Strength > 0) の場合、エネルギーは Target へと「滑り落ちる」
		if conn.Strength > 0 {
			// 移動量 = 転がるエネルギー × 張力 × 時間
			flow := srcEnergy * (conn.Strength / 100.0) * deltaTime
			nextState[conn.Source] -= flow
			nextState[conn.Target] += flow
		}

		// 斥力 (Strength < 0) の場合、エネルギーは反発して遠ざかる
		if conn.Strength < 0 {
			// 斥力はソースのエネルギーを減衰させるか、あるいは不安定化させる
			repel := srcEnergy * (mathAbs(conn.Strength) / 100.0) * deltaTime
			nextState[conn.Source] -= repel
		}
	}

	return nextState
}

func mathAbs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func (s State) Dump(t int) {
	fmt.Printf("Step [%d]:\n", t)
	for name, energy := range s {
		if energy > 0.01 {
			fmt.Printf("  %-12s: %.2f\n", name, energy)
		}
	}
}
