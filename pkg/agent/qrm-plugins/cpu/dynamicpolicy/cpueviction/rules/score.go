/*
Copyright 2022 The Katalyst Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rules

import (
	"fmt"
	"math"
	"sort"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/kubewharf/katalyst-core/pkg/consts"
	"github.com/kubewharf/katalyst-core/pkg/metrics"
	"github.com/kubewharf/katalyst-core/pkg/util/general"
)

const (
	DeploymentEvictionFrequencyScorerName  = "DeploymentEvictionFrequency"
	UsageGapScorerName                     = "Usage"
	targetMetric                           = consts.MetricCPUUsageContainer
	metricNameEvictPodScore                = "pressure_numa_evict_pod_score"
	metricNameEvictPodLimitRatio           = "pressure_numa_evict_pod_limit_ratio"
	metricNameEvictDeploymentEvictionRatio = "pressure_numa_evict_deployment_eviction_ratio"

	score = 100.0

	deploymentEvictionFrequencyScoreWeight      = 1.0
	deploymentEvictionFrequencyLimitRatioWeight = 0.5

	usageGapScoreWeight = 0.5
	usageGapLimit       = 0.3
)

var DefaultEnabledScorers = []string{
	DeploymentEvictionFrequencyScorerName,
	UsageGapScorerName,
}

// ScoreFunc defines the function signature for scoring candidate pods
type ScoreFunc func(pod *CandidatePod, params interface{}) int

// Scorer manages a collection of scoring functions and their parameters
type Scorer struct {
	scorers      map[string]ScoreFunc
	scorerParams map[string]interface{}
	emitter      metrics.MetricEmitter
}

// allScores contains all available scoring functions that can be enabled
var allScores = map[string]ScoreFunc{
	DeploymentEvictionFrequencyScorerName: DeploymentEvictionFrequencyScorer,
	UsageGapScorerName:                    UsageGapScorer,
}

func NewScorer(enabledScorers []string, emitter metrics.MetricEmitter, scorerParams map[string]interface{}) (*Scorer, error) {
	if scorerParams == nil {
		general.Warningf("scorerParams is nil, using empty parameters for all scorers")
		scorerParams = make(map[string]interface{})
	}
	enabled := sets.NewString(enabledScorers...)
	scorers := make(map[string]ScoreFunc)
	params := make(map[string]interface{})
	for name := range enabled {
		scorer, exists := allScores[name]
		if !exists {
			return nil, fmt.Errorf("scorer %s not exists", name)
		}
		scorers[name] = scorer
		params[name] = scorerParams[name]
	}
	general.Infof("initialized scorer with %d enabled scorers", len(scorers))
	return &Scorer{
		scorers:      scorers,
		scorerParams: params,
		emitter:      emitter,
	}, nil
}

// score return: []*CandidatePod - Sorted list of pods by total score (ascending)
func (s *Scorer) Score(pods []*CandidatePod) []*CandidatePod {
	if pods == nil {
		general.Warningf("pods is nil, returning empty list")
		return []*CandidatePod{}
	}
	if len(s.scorers) == 0 {
		general.Warningf("no scorers enabled, returning pods without scoring")
		return pods
	}
	var validPods []*CandidatePod
	for _, pod := range pods {
		if pod == nil {
			general.Warningf("nil pod found in scoring list, skipping")
			continue
		}
		validPods = append(validPods, pod)
	}

	for _, pod := range validPods {
		if pod == nil {
			general.Warningf("nil pod found in scoring list, skipping")
			continue
		}
		pod.TotalScore = 0
		pod.Scores = make(map[string]int)
		for name, scorer := range s.scorers {
			score := scorer(pod, s.scorerParams[name])
			pod.Scores[name] = score
			pod.TotalScore += score
		}
	}
	sort.Slice(validPods, func(i, j int) bool {
		return validPods[i].TotalScore > validPods[j].TotalScore
	})
	if validPods[0].WorkloadsEvictionInfo != nil {
		evictionInfo, exists := validPods[0].WorkloadsEvictionInfo[workloadName]
		if exists {
			_ = s.emitter.StoreInt64(metricNameEvictPodScore, int64(validPods[0].TotalScore), metrics.MetricTypeNameRaw)
			limitRatio := evictionInfo.Limit / evictionInfo.Replicas
			_ = s.emitter.StoreFloat64(metricNameEvictPodLimitRatio, float64(limitRatio), metrics.MetricTypeNameRaw)
			// PerHourEvictionRatio := evictionRatio / window
			for _, stats := range evictionInfo.StatsByWindow {
				_ = s.emitter.StoreFloat64(metricNameEvictDeploymentEvictionRatio, float64(stats.EvictionRatio), metrics.MetricTypeNameRaw)
				break
			}
		}
	}
	general.Infof("scored %d pods, top score: %s, %d, bottom score: %s, %d", len(validPods), validPods[0].Pod.Name, validPods[0].TotalScore, validPods[len(validPods)-1].Pod.Name, validPods[len(validPods)-1].TotalScore)
	for scorerName, score := range validPods[0].Scores {
		general.Infof("scorer name: %s, score: %d", scorerName, score)
		_ = s.emitter.StoreInt64(metricNameEvictPodScore+"_"+scorerName, int64(score), metrics.MetricTypeNameRaw)
	}
	return validPods
}

func (s *Scorer) SetScorerParam(key string, value interface{}) {
	if key == "" {
		general.Warningf("scorer key is empty, will not set scorer param")
		return
	}
	s.scorerParams[key] = value
}

func (s *Scorer) SetScorer(key string, scorer ScoreFunc) {
	if key == "" {
		general.Warningf("scorer key is empty, will not set scorer")
		return
	}
	s.scorers[key] = scorer
}

func DeploymentEvictionFrequencyScorer(pod *CandidatePod, params interface{}) int {
	if len(pod.WorkloadsEvictionInfo) == 0 {
		general.Warningf("no eviction info for pod %s", pod.Pod.Name)
		return 0
	}
	var totalScore float64
	var workloadCount int
	for _, workloadInfo := range pod.WorkloadsEvictionInfo {
		if len(workloadInfo.StatsByWindow) == 0 {
			general.Warningf("no eviction info for workload %s", workloadInfo.WorkloadName)
			continue
		}
		var windowScore float64
		var weightSum float64
		for window, stats := range workloadInfo.StatsByWindow {
			weight := 1 / window
			// perHourCount := float64(stats.EvictionCount) / window
			// general.Infof("limit: %v, perHourCount: %v", workloadInfo.Limit, perHourCount)
			// countScore := normalizeCount(perHourCount, workloadInfo.Limit)
			freqScore := normalizeDeploymentEvictionFrequencyScore(workloadInfo, window, stats)
			// windowContribution := countScore * stats.EvictionRatio * 10
			// general.Infof("window: %v, countScore: %v, ratio: %v, windowContribution: %v", window, countScore, stats.EvictionRatio, windowContribution)
			windowContribution := freqScore * weight
			windowScore += windowContribution
			weightSum += weight
		}
		if weightSum > 0 {
			workloadScore := windowScore / weightSum
			totalScore += workloadScore
			workloadCount++
		}
	}
	if workloadCount == 0 {
		return 0
	}
	// general.Infof("workloadCount: %d, frequency totalScore:  %v", workloadCount, totalScore)
	avgScore := totalScore / float64(workloadCount)
	general.Infof("DeploymentEvictionFrequencyScorer,  pod: %v, frequencyScore:  %v", pod.Pod.Name, avgScore)
	return int(avgScore)
}

// UsageGapScorer params: numaOverStat []NumaOverStat
func UsageGapScorer(pod *CandidatePod, params interface{}) int {
	if pod == nil || pod.Pod == nil {
		general.Warningf("nil pod or pod spec passed to UsageGapScorer")
		return 0
	}
	numaOverStats, ok := params.([]NumaOverStat)
	if !ok {
		general.Errorf("UsageGapScorer params type error, params: %v", params)
		return 0
	}
	if len(numaOverStats) == 0 {
		general.Errorf("UsageGapScorer received empty numaOverStats")
		return 0
	}
	numaID := numaOverStats[0].NumaID
	metricsHistory := numaOverStats[0].MetricsHistory
	numaHis, ok := metricsHistory.Inner[numaID]
	if !ok {
		general.Warningf("no metrics history for numa %d", numaID)
		return 0
	}
	podUID := string(pod.Pod.UID)
	podHis, existMetric := numaHis[podUID]
	if !existMetric {
		general.Warningf("no metric history for pod %s on numa %d", podUID, numaID)
		return 0
	}
	metricRing, ok := podHis[targetMetric]

	if !ok {
		general.Warningf("no %s metric history for pod %s", targetMetric, podUID)
		return 0
	}
	avgUsageRatio := metricRing.Avg()

	pod.UsageRatio = avgUsageRatio
	gap := numaOverStats[0].Gap
	// choose low score
	usageGap := math.Abs(avgUsageRatio - math.Abs(gap))
	usageGapScore := normalizeUsageGapScore(usageGap)
	general.Infof("UsageGapScorer,  pod: %v, usageGapScore:  %v, numaGap: %v", pod.Pod.Name, usageGapScore, gap)
	return int(usageGapScore)
}

func normalizeUsageGapScore(usageGap float64) float64 {
	if 0 <= usageGap && usageGap < usageGapLimit {
		return usageGapScoreWeight * score * (1 - usageGap/usageGapLimit)
	}
	return 0
}

func normalizeDeploymentEvictionFrequencyScore(workloadInfo *WorkloadEvictionInfo, window float64, stats *EvictionStats) float64 {
	normalizedPerHourEvictionRatio := stats.EvictionRatio / window
	normalizedLimitRatio := deploymentEvictionFrequencyLimitRatioWeight * (float64(workloadInfo.Limit) / float64(workloadInfo.Replicas))
	if normalizedPerHourEvictionRatio >= normalizedLimitRatio {
		return 0
	} else {
		return deploymentEvictionFrequencyScoreWeight * score * (1 - normalizedPerHourEvictionRatio/normalizedLimitRatio)
	}
}

// _ = emitter.StoreInt64(metricNameNumaOverloadNumaCount, int64(overloadNumaCount), metrics.MetricTypeNameRaw)
