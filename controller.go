package autoscaler

var defaultEngine = NewInferenceEngine()

type ScalingController struct {
	engine *InferenceEngine
}

func NewScalingController() *ScalingController {
	return &ScalingController{engine: defaultEngine}
}

func (c *ScalingController) ComputeScalingAction(cpuLoad, trafficRate float64) int {
	return c.engine.Decide(cpuLoad, trafficRate)
}
