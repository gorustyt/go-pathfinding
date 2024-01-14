package path_finding

type PathFindingType string
type PathFindingCmd func(startX, startY, endX, endY int) (res []*GridNodeInfo)

const (
	AStar                                PathFindingType = "AStar"
	IdAStar                              PathFindingType = "IdAStar"
	Dijkstra                             PathFindingType = "Dijkstra"
	BestFirst                            PathFindingType = "BestFirst"
	BreadthFirst                         PathFindingType = "BreadthFirst"
	JumpPoint                            PathFindingType = "JumpPoint"
	JPFNeverMoveDiagonally               PathFindingType = "JPFNeverMoveDiagonally"
	JPFMoveDiagonallyIfNoObstacles       PathFindingType = "JPFMoveDiagonallyIfNoObstacles"
	JPFMoveDiagonallyIfAtMostOneObstacle PathFindingType = "JPFMoveDiagonallyIfAtMostOneObstacle"
	JPFAlwaysMoveDiagonally              PathFindingType = "JPFAlwaysMoveDiagonally"

	BiAStar        PathFindingType = "BiAStar"
	BiBestFirst    PathFindingType = "BiBestFirst"
	BiBreadthFirst PathFindingType = "BiBreadthFirst"
	BiDijkstra     PathFindingType = "BiDijkstra"
)

type PathFindingConfigOptions func(cfg *PathFindingConfig)
type PathFindingConfig struct {
	weight           float64          //权重
	allowDiagonal    bool             //允许对角线行走
	DiagonalMovement DiagonalMovement //对角线行走规则
	Heuristic        Heuristic        //估算函数
	DontCrossCorners bool             //是否跨越障碍物
	*IdAStarConfig
}

type IdAStarConfig struct {
	IdAStarTimeLimit int64 //最大搜索秒数，超过这个值也视为没有找到
}

func (cfg *PathFindingConfig) check() {
	if cfg.weight == 0 {
		cfg.weight = 1
	}
	if cfg.Heuristic == nil {
		cfg.Heuristic = manhattan
	}
	if cfg.DiagonalMovement == DiagonalMovementNone {
		if !cfg.allowDiagonal {
			cfg.DiagonalMovement = DiagonalMovementNever
		} else {
			if cfg.DontCrossCorners {
				cfg.DiagonalMovement = DiagonalMovementOnlyWhenNoObstacles
			} else {
				cfg.DiagonalMovement = DiagonalMovementIfAtMostOneObstacle
			}
		}
	}
	if cfg.DiagonalMovement == DiagonalMovementNever {
		cfg.Heuristic = manhattan
	} else {
		cfg.Heuristic = octile
	}
	if cfg.IdAStarConfig == nil {
		cfg.IdAStarConfig = &IdAStarConfig{
			IdAStarTimeLimit: 10,
		}
	}

}

func WithWeight(weight float64) PathFindingConfigOptions {
	return func(cfg *PathFindingConfig) {
		cfg.weight = weight
	}
}
func WithAllowDiagonal(allowDiagonal bool) PathFindingConfigOptions {
	return func(cfg *PathFindingConfig) {
		cfg.allowDiagonal = allowDiagonal
	}
}

func WithDiagonalMovement(DiagonalMovement DiagonalMovement) PathFindingConfigOptions {
	return func(cfg *PathFindingConfig) {
		cfg.DiagonalMovement = DiagonalMovement
	}
}
func WithHeuristic(Heuristic Heuristic) PathFindingConfigOptions {
	return func(cfg *PathFindingConfig) {
		cfg.Heuristic = Heuristic
	}
}
func WithDontCrossCorners(DontCrossCorners bool) PathFindingConfigOptions {
	return func(cfg *PathFindingConfig) {
		cfg.DontCrossCorners = DontCrossCorners
	}
}
