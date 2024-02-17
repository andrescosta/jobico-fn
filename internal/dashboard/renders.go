package dashboard

import (
	"fmt"
	"strings"

	"github.com/andrescosta/goico/pkg/service/grpc/svcmeta"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var iconPreffixNodesMap = map[bool]string{
	true:  iconExpanded,
	false: iconContracted,
}

func renderNode(target *node) *tview.TreeNode {
	if len(target.children) > 0 {
		if !target.expanded {
			target.text = renderNodeText(iconContracted, target.text)
		}
		target.color = tcell.ColorGreen
	} else {
		target.color = tcell.ColorWhite
	}
	node := tview.NewTreeNode(target.text).
		SetExpanded(target.expanded).
		SetReference(target).
		SetColor(target.color)
	for _, child := range target.children {
		node.AddChild(renderNode(child))
	}
	return node
}

func reRenderNode(target *node, tn *tview.TreeNode) {
	if len(target.children) > 0 {
		if !target.expanded {
			target.text = renderNodeText(iconContracted, target.text)
		}
		target.color = tcell.ColorGreen
	} else {
		newText, f := strings.CutPrefix(target.text, iconContracted)
		if !f {
			newText, _ = strings.CutPrefix(target.text, iconExpanded)
		}
		target.text = newText
		target.color = tcell.ColorWhite
		target.expanded = false
	}
	tn.SetText(target.text).
		SetExpanded(target.expanded).
		SetReference(target).
		SetColor(target.color)
}

func renderNodeText(icon, text string) string {
	return fmt.Sprintf("%s%s", icon, text)
}

func renderHTTPTableServer(info map[string]string) *tview.Table {
	table := tview.NewTable().
		SetBorders(true)
	table.SetCell(0, 0,
		tview.NewTableCell("Status").
			SetAlign(tview.AlignCenter))
	status := "Unknown"
	table.SetCell(0, 1,
		tview.NewTableCell(status).
			SetAlign(tview.AlignCenter))
	ix := 0
	for k, v := range info {
		table.SetCell(ix+1, 0,
			tview.NewTableCell(k).
				SetAlign(tview.AlignCenter))
		table.SetCell(ix+1, 1,
			tview.NewTableCell(v).
				SetAlign(tview.AlignCenter))
		ix++
	}
	return table
}

func renderGrpcTableServer(infos []*svcmeta.GrpcServerMetadata, s healthpb.HealthCheckResponse_ServingStatus) *tview.Table {
	table := tview.NewTable().
		SetBorders(true)
	table.SetCell(0, 0,
		tview.NewTableCell("Status").
			SetAlign(tview.AlignCenter))
	status := "Not ok"
	if s == healthpb.HealthCheckResponse_SERVING {
		status = "Ok"
	}
	table.SetCell(0, 1,
		tview.NewTableCell(status).
			SetAlign(tview.AlignCenter))
	for ix, info := range infos {
		table.SetCell(ix+1, 0,
			tview.NewTableCell(info.Key).
				SetAlign(tview.AlignCenter))
		table.SetCell(ix+1, 1,
			tview.NewTableCell(info.Value).
				SetAlign(tview.AlignCenter))
	}
	return table
}
