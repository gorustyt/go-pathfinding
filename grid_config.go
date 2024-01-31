package path_finding

type PathFindingType string
type PathFindingCmd func(startX, startY, endX, endY int) (res []*PathPoint)

const (
	DescAStar               = "AStar"
	DescIdaStar             = "Ida star"
	DescBreadthFirstSearch  = "Breadth First Search"
	DescBestFirstSearch     = "Best First Search"
	DescDijkstra            = "Dijkstra"
	DescJumpPointSearch     = "Jump Point Search"
	DescOrthogonalJumpPoint = "Orthogonal Jump Point"
)

func GetPathFindingType(dest string, cfg *PathFindingConfig) PathFindingType {
	switch dest {
	case DescAStar:
		if cfg.BiDirectional {
			return BiAStar
		}
		return AStar
	case DescIdaStar:
		return IdAStar
	case DescBreadthFirstSearch:
		if cfg.BiDirectional {
			return BiBreadthFirst
		}
		return BreadthFirst
	case DescBestFirstSearch:
		if cfg.BiDirectional {
			return BiBestFirst
		}
		return BestFirst
	case DescDijkstra:
		if cfg.BiDirectional {
			return BiDijkstra
		}
		return Dijkstra
	case DescJumpPointSearch:
		cfg.DiagonalMovement = DiagonalMovementIfAtMostOneObstacle
		return JumpPoint
	case DescOrthogonalJumpPoint:
		cfg.DiagonalMovement = DiagonalMovementNever
		return JumpPoint
	}
	return Undefined
}

const (
	Undefined      PathFindingType = ""
	AStar          PathFindingType = "AStar"
	IdAStar        PathFindingType = "IdAStar"
	Dijkstra       PathFindingType = "Dijkstra"
	BestFirst      PathFindingType = "BestFirst"
	BreadthFirst   PathFindingType = "BreadthFirst"
	JumpPoint      PathFindingType = "JumpPoint"
	BiAStar        PathFindingType = "BiAStar"
	BiBestFirst    PathFindingType = "BiBestFirst"
	BiBreadthFirst PathFindingType = "BiBreadthFirst"
	BiDijkstra     PathFindingType = "BiDijkstra"
)

type PathFindingConfigOptions func(cfg *PathFindingConfig)
type PathFindingConfig struct {
	BiDirectional    bool             //是否是双向
	Weight           float64          //权重
	AllowDiagonal    bool             //允许对角线行走
	DiagonalMovement DiagonalMovement //对角线行走规则
	Heuristic        Heuristic        //估算函数
	DontCrossCorners bool             //是否跨越障碍物
	IdAStarTimeLimit int64            //最大搜索秒数，超过这个值也视为没有找到
	Trace            DebugTrace
}

func GetDefaultConfig() *PathFindingConfig {
	return &PathFindingConfig{
		AllowDiagonal: true,
		Heuristic:     manhattan,
	}
}

func (cfg *PathFindingConfig) check() {
	if cfg.Weight == 0 {
		cfg.Weight = 1
	}
	if cfg.Heuristic == nil {
		cfg.Heuristic = manhattan
	}
	if cfg.DiagonalMovement == DiagonalMovementNone {
		if !cfg.AllowDiagonal {
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
	if cfg.IdAStarTimeLimit == 0 {
		cfg.IdAStarTimeLimit = 10
	}

}

func WithWeight(weight float64) PathFindingConfigOptions {
	return func(cfg *PathFindingConfig) {
		cfg.Weight = weight
	}
}
func WithAllowDiagonal(allowDiagonal bool) PathFindingConfigOptions {
	return func(cfg *PathFindingConfig) {
		cfg.AllowDiagonal = allowDiagonal
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

func WithDebugTrace(trace DebugTrace) PathFindingConfigOptions {
	return func(cfg *PathFindingConfig) {
		cfg.Trace = trace
	}
}
