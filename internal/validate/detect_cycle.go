//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what import 방향 그래프에서 DFS 기반 순환 경로를 모두 탐지
package validate

// detectCycle finds all cycles in a directed graph using DFS.
// Each cycle is a list of node keys forming the cycle (last -> first closes it).
func detectCycle(graph map[string][]string) [][]string {
	var cycles [][]string
	visited := make(map[string]bool)
	inStack := make(map[string]bool)
	stackIndex := make(map[string]int)
	var stack []string

	var dfs func(node string)
	dfs = func(node string) {
		visited[node] = true
		inStack[node] = true
		stackIndex[node] = len(stack)
		stack = append(stack, node)

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				dfs(neighbor)
			} else if inStack[neighbor] {
				idx := stackIndex[neighbor]
				cycle := make([]string, len(stack)-idx)
				copy(cycle, stack[idx:])
				cycles = append(cycles, cycle)
			}
		}

		stack = stack[:len(stack)-1]
		inStack[node] = false
	}

	for node := range graph {
		if !visited[node] {
			dfs(node)
		}
	}
	return cycles
}
