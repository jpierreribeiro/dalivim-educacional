package services

type BehaviorAnalysis struct {
	AuthorshipScore float64  `json:"authorship_score"`
	Confidence      string   `json:"confidence"`
	Signals         []string `json:"signals"`
}

func analyzeBehavior(features map[string]interface{}) BehaviorAnalysis {
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

	// Delete ratio check (low deletion is suspicious)
	deleteRatio := getFloat(features, "deleteRatio")
	if deleteRatio < 0.02 {
		signals = append(signals, "low_edit_ratio")
		suspicionScore += 0.25
	}

	// Linear editing check (very linear is suspicious)
	linearScore := getFloat(features, "linearEditingScore")
	if linearScore > 0.9 {
		signals = append(signals, "highly_linear_editing")
		suspicionScore += 0.2
	}

	// Large paste events
	pasteEvents := getInt(features, "pasteEvents")
	if pasteEvents > 3 {
		signals = append(signals, "multiple_paste_events")
		suspicionScore += 0.15
	}

	// Fast completion with no executions
	execCount := getInt(features, "executionCount")
	totalTime := getFloat(features, "totalTime")
	if execCount == 0 && totalTime < 120 {
		signals = append(signals, "fast_completion_no_testing")
		suspicionScore += 0.2
	}

	// Suspicious focus loss
	focusLoss := getInt(features, "focusLossCount")
	if focusLoss > 5 {
		signals = append(signals, "frequent_focus_loss")
		suspicionScore += 0.1
	}

	// Low keystroke variance (robotic typing)
	burstiness := getFloat(features, "burstiness")
	if burstiness < 0.3 {
		signals = append(signals, "low_typing_variance")
		suspicionScore += 0.15
	}

	// Calculate final authorship score (inverse of suspicion)
	authorshipScore := 1.0 - min(suspicionScore, 1.0)

	// Determine confidence
	confidence := "low"
	if len(signals) >= 4 {
		confidence = "high"
	} else if len(signals) >= 2 {
		confidence = "medium"
	}

	return BehaviorAnalysis{
		AuthorshipScore: authorshipScore,
		Confidence:      confidence,
		Signals:         signals,
	}
}
