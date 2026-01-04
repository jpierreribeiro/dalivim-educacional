package service

type AnalysisResult struct {
	AuthorshipScore float64  `json:"authorship_score"`
	Confidence      string   `json:"confidence"`
	Signals         []string `json:"signals"`
}

type AnalysisService interface {
	Analyze(features map[string]interface{}) AnalysisResult
}

type analysisService struct{}

func NewAnalysisService() AnalysisService {
	return &analysisService{}
}

func (s *analysisService) Analyze(features map[string]interface{}) AnalysisResult {
	signals := []string{}
	suspicionScore := 0.0

	// Paste ratio check
	pasteRatio := getFloat(features, "pasteCharRatio")
	if pasteRatio > 0.6 {
		signals = append(signals, "high_paste_ratio")
		suspicionScore += 0.3
	} else if pasteRatio > 0.3 {
		signals = append(signals, "moderate_paste_ratio")
		suspicionScore += 0.15
	}

	// Delete ratio check
	deleteRatio := getFloat(features, "deleteRatio")
	if deleteRatio < 0.02 {
		signals = append(signals, "low_edit_ratio")
		suspicionScore += 0.25
	}

	// Linear editing check
	linearScore := getFloat(features, "linearEditingScore")
	if linearScore > 0.9 {
		signals = append(signals, "highly_linear_editing")
		suspicionScore += 0.2
	}

	// Paste events check
	pasteEvents := getInt(features, "pasteEvents")
	if pasteEvents > 3 {
		signals = append(signals, "multiple_paste_events")
		suspicionScore += 0.15
	}

	// Fast completion check
	execCount := getInt(features, "executionCount")
	totalTime := getFloat(features, "totalTime")
	if execCount == 0 && totalTime < 120 {
		signals = append(signals, "fast_completion_no_testing")
		suspicionScore += 0.2
	}

	// Focus loss check
	focusLoss := getInt(features, "focusLossCount")
	if focusLoss > 5 {
		signals = append(signals, "frequent_focus_loss")
		suspicionScore += 0.1
	}

	// Keystroke variance check
	burstiness := getFloat(features, "burstiness")
	if burstiness < 0.3 {
		signals = append(signals, "low_typing_variance")
		suspicionScore += 0.15
	}

	// Calculate authorship score
	authorshipScore := 1.0 - min(suspicionScore, 1.0)

	// Determine confidence
	confidence := "low"
	if len(signals) >= 4 {
		confidence = "high"
	} else if len(signals) >= 2 {
		confidence = "medium"
	}

	return AnalysisResult{
		AuthorshipScore: authorshipScore,
		Confidence:      confidence,
		Signals:         signals,
	}
}

func getFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 0.0
}

func getInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		if f, ok := v.(float64); ok {
			return int(f)
		}
	}
	return 0
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
