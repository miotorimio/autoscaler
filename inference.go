package autoscaler

import "math"

type LinguisticTerm struct {
	Name     string
	Function MembershipFunction
}

type OutputTerm struct {
	Name     string
	Function MembershipFunction
	Centroid float64
}

type Rule struct {
	CPU     *LinguisticTerm
	Traffic *LinguisticTerm
	Output  *OutputTerm
}

type InferenceEngine struct {
	cpuTerms     []*LinguisticTerm
	trafficTerms []*LinguisticTerm
	outputTerms  []*OutputTerm
	rules        []*Rule
}

func NewInferenceEngine() *InferenceEngine {
	cpuLow := &LinguisticTerm{Name: "Низкая", Function: Trapezoid{A: 0, B: 0, C: 30, D: 50}}
	cpuMed := &LinguisticTerm{Name: "Средняя", Function: Triangle{A: 30, B: 50, C: 70}}
	cpuHigh := &LinguisticTerm{Name: "Высокая", Function: Trapezoid{A: 50, B: 70, C: 100, D: 100}}

	trafficStable := &LinguisticTerm{Name: "Стабильный", Function: Trapezoid{A: 0, B: 0, C: 400, D: 700}}
	trafficExplosive := &LinguisticTerm{Name: "Взрывной", Function: Trapezoid{A: 400, B: 700, C: 1000, D: 1000}}

	remove := &OutputTerm{Name: "Удалить", Function: Trapezoid{A: -5, B: -5, C: -4, D: -1}, Centroid: -4}
	keep := &OutputTerm{Name: "Оставить", Function: Triangle{A: -2, B: 0, C: 2}, Centroid: 0}
	add := &OutputTerm{Name: "Добавить", Function: Trapezoid{A: 1, B: 4, C: 5, D: 5}, Centroid: 4}

	rules := []*Rule{
		{CPU: cpuLow, Traffic: trafficStable, Output: keep},
		{CPU: cpuLow, Traffic: trafficExplosive, Output: add},
		{CPU: cpuMed, Traffic: trafficStable, Output: keep},
		{CPU: cpuMed, Traffic: trafficExplosive, Output: add},
		{CPU: cpuHigh, Traffic: trafficStable, Output: add},
		{CPU: cpuHigh, Traffic: trafficExplosive, Output: add},
	}

	return &InferenceEngine{
		cpuTerms:     []*LinguisticTerm{cpuLow, cpuMed, cpuHigh},
		trafficTerms: []*LinguisticTerm{trafficStable, trafficExplosive},
		outputTerms:  []*OutputTerm{remove, keep, add},
		rules:        rules,
	}
}

func (e *InferenceEngine) fuzzifyCPU(value float64) map[string]float64 {
	value = clamp(value, 0, 100)
	result := make(map[string]float64, len(e.cpuTerms))
	for _, term := range e.cpuTerms {
		result[term.Name] = term.Function.Membership(value)
	}
	return result
}

func (e *InferenceEngine) fuzzifyTraffic(value float64) map[string]float64 {
	if value < 0 {
		value = 0
	}
	if value > 1000 {
		value = 1000
	}
	result := make(map[string]float64, len(e.trafficTerms))
	for _, term := range e.trafficTerms {
		result[term.Name] = term.Function.Membership(value)
	}
	return result
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func (e *InferenceEngine) infer(cpuValue, trafficValue float64) map[string]float64 {
	cpuDegrees := e.fuzzifyCPU(cpuValue)
	trafficDegrees := e.fuzzifyTraffic(trafficValue)
	activation := map[string]float64{"Удалить": 0, "Оставить": 0, "Добавить": 0}

	for _, rule := range e.rules {
		strength := min(cpuDegrees[rule.CPU.Name], trafficDegrees[rule.Traffic.Name])
		if strength > activation[rule.Output.Name] {
			activation[rule.Output.Name] = strength
		}
	}

	return activation
}

func (e *InferenceEngine) defuzzify(activation map[string]float64) int {
	numerator := 0.0
	denominator := 0.0

	for _, term := range e.outputTerms {
		weight := activation[term.Name]
		numerator += term.Centroid * weight
		denominator += weight
	}

	if denominator == 0 {
		return 0
	}

	result := numerator / denominator
	result = math.Round(result)
	if result < -5 {
		result = -5
	}
	if result > 5 {
		result = 5
	}
	return int(result)
}

func (e *InferenceEngine) Decide(cpuLoad, trafficRate float64) int {
	activation := e.infer(cpuLoad, trafficRate)
	return e.defuzzify(activation)
}
