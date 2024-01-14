package path_finding

import "math"

/**
 *  PF.Heuristic
 *  A collection of heuristic functions.
 */

/**
 * Manhattan distance.
 * @param {number} dx - Difference in x.
 * @param {number} dy - Difference in y.
 * @return {number} dx + dy
 */
func manhattan(dx, dy int) int {
	return dx + dy
}

/**
 * Euclidean distance.
 * @param {number} dx - Difference in x.
 * @param {number} dy - Difference in y.
 * @return {number} sqrt(dx * dx + dy * dy)
 */
func euclidean(dx, dy int) int {
	return int(math.Sqrt(float64(dx)*float64(dx) + float64(dy)*float64(dy)))
}

/**
 * Octile distance.
 * @param {number} dx - Difference in x.
 * @param {number} dy - Difference in y.
 * @return {number} sqrt(dx * dx + dy * dy) for grids
 */
func octile(dx, dy int) int {
	var F = math.Sqrt2 - 1
	if dx < dy {
		return int(F*float64(dx) + float64(dy))
	}
	return int(F*float64(dy) + float64(dx))
}

/**
 * Chebyshev distance.
 *  dx - Difference in x.
 * dy - Difference in y.
 * return {number} max(dx, dy)
 */
func chebyshev(dx, dy int) int {
	return int(math.Max(float64(dx), float64(dy)))
}
