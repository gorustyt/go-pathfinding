package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	path_finding "github.com/gorustyt/go-pathfinding"
	"github.com/gorustyt/go-pathfinding/example/demo/grid_map"
)

type OptionsMenu struct {
	uid string
	widget.BaseWidget
	cfg            *grid_map.Config
	heuristicLabel *widget.Label
	heuristic      *widget.RadioGroup

	options      *widget.CheckGroup
	optionsLabel *widget.Label

	visualizeRecursion *widget.CheckGroup

	weight      *fyne.Container
	secondLimit *fyne.Container
}

func NewOptionsMenu(cfg *grid_map.Config) fyne.CanvasObject {
	o := &OptionsMenu{cfg: cfg}
	o.BaseWidget.ExtendBaseWidget(o)
	o.initHeuristic()
	o.initOptions()
	o.initVisualizeRecursion()
	o.initEntryWeight()
	o.initSecondLimit()
	return o
}

func (m *OptionsMenu) initHeuristic() {
	g := widget.NewRadioGroup([]string{
		path_finding.DescHeuristicManhattan,
		path_finding.DescHeuristicEuclidean,
		path_finding.DescHeuristicOctile,
		path_finding.DescHeuristicChebyshev,
	}, func(s string) {
		m.cfg.Heuristic = path_finding.GetHeuristicByDesc(s)
	})
	g.Selected = path_finding.DescHeuristicManhattan
	m.heuristic = g
	m.heuristicLabel = widget.NewLabel("Heuristic")
}
func (m *OptionsMenu) SetUid(uid string) {
	m.uid = uid
	m.Refresh()
}

func (m *OptionsMenu) Refresh() {
	switch m.uid {
	case path_finding.DescAStar:
		m.heuristicLabel.Show()
		m.heuristic.Show()
		m.options.Options = []string{
			optionsAllowDiagonal,
			optionsBiDirectional,
			optionsDontCrossCorners,
		}
		m.options.Selected = []string{
			optionsAllowDiagonal,
		}
		m.options.Show()

		m.weight.Show()

		m.secondLimit.Hide()
		m.visualizeRecursion.Hide()
	case path_finding.DescIdaStar:
		m.heuristicLabel.Show()
		m.heuristic.Show()

		m.options.Options = []string{
			optionsAllowDiagonal,
			optionsDontCrossCorners,
		}
		m.options.Selected = []string{
			optionsAllowDiagonal,
		}
		m.options.Show()

		m.weight.Show()
		m.secondLimit.Show()
		m.visualizeRecursion.Show()
	case path_finding.DescBestFirstSearch:
		m.heuristicLabel.Show()
		m.heuristic.Show()
		m.options.Options = []string{
			optionsAllowDiagonal,
			optionsBiDirectional,
			optionsDontCrossCorners,
		}
		m.options.Selected = []string{
			optionsAllowDiagonal,
		}
		m.options.Show()
		m.weight.Hide()
		m.secondLimit.Hide()
		m.visualizeRecursion.Hide()
	case path_finding.DescDijkstra, path_finding.DescBreadthFirstSearch:
		m.heuristicLabel.Hide()
		m.heuristic.Hide()

		m.options.Options = []string{
			optionsAllowDiagonal,
			optionsBiDirectional,
			optionsDontCrossCorners,
		}
		m.options.Selected = []string{
			optionsAllowDiagonal,
		}
		m.options.Show()

		m.weight.Hide()
		m.secondLimit.Hide()
		m.visualizeRecursion.Hide()
	case path_finding.DescJumpPointSearch, path_finding.DescOrthogonalJumpPoint:
		m.heuristicLabel.Show()
		m.heuristic.Show()

		m.options.Hide()
		m.weight.Hide()
		m.secondLimit.Hide()

		m.visualizeRecursion.Show()
	case path_finding.DescTrace:
		m.heuristicLabel.Show()
		m.heuristic.Show()

		m.options.Options = []string{
			optionsAllowDiagonal,
			optionsDontCrossCorners,
		}
		m.options.Selected = []string{
			optionsAllowDiagonal,
		}
		m.options.Show()
		m.weight.Hide()
		m.secondLimit.Hide()
		m.visualizeRecursion.Hide()
	}
}

func (m *OptionsMenu) initOptions() {
	g := widget.NewCheckGroup([]string{
		optionsAllowDiagonal,
		optionsBiDirectional,
		optionsDontCrossCorners,
	}, func(s []string) {
		parseOptions(s, m.cfg)
	})
	if m.cfg.AllowDiagonal {
		g.Selected = append(g.Selected, optionsAllowDiagonal)
	}
	if m.cfg.BiDirectional {
		g.Selected = append(g.Selected, optionsBiDirectional)
	}
	if m.cfg.DontCrossCorners {
		g.Selected = append(g.Selected, optionsDontCrossCorners)
	}
	m.options = g
	m.optionsLabel = widget.NewLabel("Options")
}

func (m *OptionsMenu) initVisualizeRecursion() {
	g := widget.NewCheckGroup([]string{"Visualize recursion"}, func(strings []string) {

	})
	g.Selected = []string{"Visualize recursion"}
	m.visualizeRecursion = g

}

func (m *OptionsMenu) initEntryWeight() {
	entry := widget.NewEntry()
	entry.Text = "1"
	entry.Bind(m.cfg.Weight)
	m.weight = container.NewGridWithColumns(2,
		entry,
		widget.NewLabel("Weight"))
}

func (m *OptionsMenu) initSecondLimit() {
	entry := widget.NewEntry()
	entry.Bind(m.cfg.SecondLimit)
	entry.Text = "10"
	m.secondLimit = container.NewGridWithColumns(2,
		entry,
		widget.NewLabel("Seconds limit"))
}

type OptionsMenuRender struct {
	box  *fyne.Container
	menu *OptionsMenu
}

func (o OptionsMenuRender) Destroy() {

}

func (o OptionsMenuRender) Layout(size fyne.Size) {
	o.box.Resize(size)
}

func (o OptionsMenuRender) MinSize() fyne.Size {
	return o.box.MinSize()
}

func (o OptionsMenuRender) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{o.box}
}

func (o OptionsMenuRender) Refresh() {
	o.menu.Refresh()
}

func (m *OptionsMenu) CreateRenderer() fyne.WidgetRenderer {
	return newOptionsMenuRender(m)
}

func newOptionsMenuRender(menu *OptionsMenu) fyne.WidgetRenderer {
	m := &OptionsMenuRender{
		menu: menu,
		box: container.NewVBox(
			menu.heuristicLabel,
			menu.heuristic,
			menu.optionsLabel,
			menu.options,
			menu.weight,
			menu.secondLimit,
			menu.visualizeRecursion,
		),
	}
	return m
}
