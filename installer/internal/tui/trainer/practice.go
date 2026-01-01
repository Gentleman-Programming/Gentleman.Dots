package trainer

import (
	"math/rand"
	"sort"
	"time"
)

// MasteryThreshold is the number of consecutive correct answers needed to master an exercise
const MasteryThreshold = 3

// init seeds the random number generator
func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetExerciseStats returns stats for a specific exercise, creating if needed
func (mp *ModuleProgress) GetExerciseStats(exerciseID string) *ExerciseStats {
	if mp.ExerciseStats == nil {
		mp.ExerciseStats = make(map[string]*ExerciseStats)
	}
	if _, ok := mp.ExerciseStats[exerciseID]; !ok {
		mp.ExerciseStats[exerciseID] = &ExerciseStats{}
	}
	return mp.ExerciseStats[exerciseID]
}

// RecordPracticeResult records the result of a practice attempt
func (mp *ModuleProgress) RecordPracticeResult(exerciseID string, correct bool) {
	stats := mp.GetExerciseStats(exerciseID)
	stats.TotalAttempts++
	stats.LastAttempted = time.Now().Unix()

	if correct {
		stats.TotalCorrect++
		stats.ConsecutiveRight++

		// Check for mastery
		if stats.ConsecutiveRight >= MasteryThreshold {
			stats.Mastered = true
		}
	} else {
		stats.TotalWrong++
		stats.ConsecutiveRight = 0
		// Un-master if they get it wrong
		stats.Mastered = false
	}

	// Update overall progress stats
	mp.PracticeAttempts++
	if correct {
		mp.PracticeCorrect++
	}
	if mp.PracticeAttempts > 0 {
		mp.PracticeAccuracy = float64(mp.PracticeCorrect) / float64(mp.PracticeAttempts)
	}
	mp.LastPracticed = time.Now()
}

// GetPracticeWeight calculates the weight for an exercise (higher = more likely to appear)
func (stats *ExerciseStats) GetPracticeWeight() int {
	if stats.Mastered {
		return 0 // Mastered exercises don't appear in practice
	}

	// Base weight
	weight := 10

	// Add weight based on wrong answers (more wrong = higher weight)
	weight += stats.TotalWrong * 5

	// Reduce weight based on consecutive correct (getting better)
	weight -= stats.ConsecutiveRight * 2

	// Never practiced = medium-high priority
	if stats.TotalAttempts == 0 {
		weight = 15
	}

	// Ensure minimum weight of 1 for non-mastered
	if weight < 1 {
		weight = 1
	}

	return weight
}

// exerciseWeight pairs an exercise with its weight for sorting
type exerciseWeight struct {
	Exercise Exercise
	Weight   int
}

// GetWeightedPracticeExercises returns exercises ordered by practice need
// Exercises with more errors appear more frequently
func GetWeightedPracticeExercises(module ModuleID, progress *ModuleProgress) []Exercise {
	lessons := GetLessons(module)
	if len(lessons) == 0 {
		return []Exercise{}
	}

	// Calculate weights for each exercise
	weights := make([]exerciseWeight, 0, len(lessons))
	totalWeight := 0
	unmasteredCount := 0

	for _, lesson := range lessons {
		stats := progress.GetExerciseStats(lesson.ID)
		weight := stats.GetPracticeWeight()

		if !stats.Mastered {
			unmasteredCount++
		}

		weights = append(weights, exerciseWeight{
			Exercise: lesson,
			Weight:   weight,
		})
		totalWeight += weight
	}

	// If all mastered, return empty (practice complete!)
	if unmasteredCount == 0 {
		return []Exercise{}
	}

	// Sort by weight descending (highest priority first for debugging/transparency)
	sort.Slice(weights, func(i, j int) bool {
		return weights[i].Weight > weights[j].Weight
	})

	// Build result slice with exercises that have weight > 0
	result := make([]Exercise, 0, unmasteredCount)
	for _, ew := range weights {
		if ew.Weight > 0 {
			ex := ew.Exercise
			ex.Type = ExercisePractice
			result = append(result, ex)
		}
	}

	return result
}

// SelectRandomPracticeExercise selects an exercise using weighted random selection
func SelectRandomPracticeExercise(module ModuleID, progress *ModuleProgress) *Exercise {
	lessons := GetLessons(module)
	if len(lessons) == 0 {
		return nil
	}

	// Build weighted pool
	type weightedEx struct {
		exercise Exercise
		weight   int
	}

	pool := make([]weightedEx, 0, len(lessons))
	totalWeight := 0

	for _, lesson := range lessons {
		stats := progress.GetExerciseStats(lesson.ID)
		weight := stats.GetPracticeWeight()

		if weight > 0 {
			pool = append(pool, weightedEx{exercise: lesson, weight: weight})
			totalWeight += weight
		}
	}

	// All mastered
	if totalWeight == 0 {
		return nil
	}

	// Weighted random selection
	r := rand.Intn(totalWeight)
	cumulative := 0

	for _, we := range pool {
		cumulative += we.weight
		if r < cumulative {
			ex := we.exercise
			ex.Type = ExercisePractice
			return &ex
		}
	}

	// Fallback (shouldn't happen)
	if len(pool) > 0 {
		ex := pool[0].exercise
		ex.Type = ExercisePractice
		return &ex
	}

	return nil
}

// GetPracticeStats returns summary stats for practice mode
type PracticeStats struct {
	TotalExercises   int
	MasteredCount    int
	RemainingCount   int
	OverallAccuracy  float64
	WeakestExercises []string // IDs of exercises with most errors
	PracticeComplete bool
}

// GetPracticeStats calculates practice statistics for a module
func GetPracticeStatsForModule(module ModuleID, progress *ModuleProgress) PracticeStats {
	lessons := GetLessons(module)

	stats := PracticeStats{
		TotalExercises:   len(lessons),
		WeakestExercises: make([]string, 0),
	}

	if len(lessons) == 0 {
		return stats
	}

	// Collect exercise stats
	type exError struct {
		id     string
		errors int
	}
	errorList := make([]exError, 0, len(lessons))

	for _, lesson := range lessons {
		exStats := progress.GetExerciseStats(lesson.ID)
		if exStats.Mastered {
			stats.MasteredCount++
		}
		if exStats.TotalWrong > 0 {
			errorList = append(errorList, exError{id: lesson.ID, errors: exStats.TotalWrong})
		}
	}

	stats.RemainingCount = stats.TotalExercises - stats.MasteredCount
	stats.PracticeComplete = stats.RemainingCount == 0
	stats.OverallAccuracy = progress.PracticeAccuracy

	// Sort by errors descending
	sort.Slice(errorList, func(i, j int) bool {
		return errorList[i].errors > errorList[j].errors
	})

	// Take top 3 weakest
	for i := 0; i < len(errorList) && i < 3; i++ {
		stats.WeakestExercises = append(stats.WeakestExercises, errorList[i].id)
	}

	return stats
}

// ResetModulePractice resets all practice data for a module
func (mp *ModuleProgress) ResetModulePractice() {
	mp.ExerciseStats = make(map[string]*ExerciseStats)
	mp.PracticeAttempts = 0
	mp.PracticeCorrect = 0
	mp.PracticeAccuracy = 0
	mp.WeakExercises = []string{}
	mp.LastPracticed = time.Time{}
}

// IsPracticeComplete returns true if all exercises are mastered
func (mp *ModuleProgress) IsPracticeComplete(module ModuleID) bool {
	lessons := GetLessons(module)
	if len(lessons) == 0 {
		return false
	}

	for _, lesson := range lessons {
		stats := mp.GetExerciseStats(lesson.ID)
		if !stats.Mastered {
			return false
		}
	}

	return true
}
