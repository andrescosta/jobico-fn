package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/andrescosta/goico/pkg/service"
	info "github.com/andrescosta/goico/pkg/service/info/grpc"
	"github.com/andrescosta/goico/pkg/yamlico"
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/rivo/tview"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func onFocusFileNode(c *App, n *tview.TreeNode) {
	f := (n.GetReference().(*node)).entity.(*sFile)
	if strings.HasSuffix(f.file, ".json") {
		pageName := f.tenant + "/" + f.file
		trySwitchToPage(pageName, c.mainView, func() (tview.Primitive, error) {
			f, err := c.repoClient.GetFile(context.Background(), f.tenant, f.file)
			if err != nil {
				return nil, err
			} else {
				cv := createContentView(string(f))
				return cv, nil
			}
		})
	}
}

func onFocusServerNode(c *App, n *tview.TreeNode) {
	h := (n.GetReference().(*node)).entity.(*sServerNode)
	addr := h.host.Ip + ":" + strconv.Itoa(int(h.host.Port))
	trySwitchToPage(addr, c.mainView, func() (tview.Primitive, error) {
		helthCheckClient := c.helthCheckClients[addr]
		infoClient, ok := c.infoClients[addr]
		if !ok {
			var err error
			infoClient, err = service.NewGrpcServerInfoClient(addr)
			if err != nil {
				return nil, err
			}
			c.infoClients[addr] = infoClient
			helthCheckClient, err = service.NewGrpcServerHelthCheckClient(addr)
			if err != nil {
				return nil, err
			}
			c.helthCheckClients[addr] = helthCheckClient
		}
		infos, err := infoClient.Info(context.Background(), &info.InfoRequest{})
		if err != nil {
			return nil, err
		}
		s, err := helthCheckClient.Check(context.Background(), h.name)
		if err != nil {
			s = healthpb.HealthCheckResponse_NOT_SERVING
		}
		view := renderTableServers(infos, s)
		return view, nil
	})
}

func onSelectedGettingJobResults(ca *App, n *tview.TreeNode) {
	n.SetText("<< stop >>")
	nl := n.GetReference().(*node)
	nl.selected = onSelectedStopGettingJobResults
	startGettingJobResults(ca)
}

func onSelectedStopGettingJobResults(ca *App, n *tview.TreeNode) {
	ca.cancelJobResultsGetter()
	nl := n.GetReference().(*node)
	n.SetText("<< start >>")
	nl.selected = onSelectedGettingJobResults
}

func onFocusJobPackageNode(ca *App, n *tview.TreeNode) {
	p := (n.GetReference().(*node)).entity.(*pb.JobPackage)
	pn := "package/" + p.TenantId + "/" + p.ID
	trySwitchToPage(pn, ca.mainView, func() (tview.Primitive, error) {
		pkg, err := ca.controlClient.GetPackage(context.Background(), p.TenantId, &p.ID)
		if err != nil {
			return nil, err
		}
		yaml, err := yamlico.Encode(pkg[0])
		if err != nil {
			return nil, err
		}
		return createContentView(*yaml), nil
	})
}

func startGettingJobResults(ca *App) {
	var textView *tview.TextView
	lines := int32(5)
	if ca.mainView.HasPage("results") {
		ca.mainView.SwitchToPage("results")
		_, fp := ca.mainView.GetFrontPage()
		textView = fp.(*tview.TextView)
		lines = 0
	} else {
		textView = tview.NewTextView().
			SetTextAlign(tview.AlignLeft).
			SetScrollable(true).
			SetWordWrap(false).
			SetWrap(false).
			SetMaxLines(100)
		ca.mainView.AddAndSwitchToPage("results", textView, true)
	}
	ch := make(chan string)
	ca.ctxJobResultsGetter, ca.cancelJobResultsGetter = context.WithCancel(context.Background())
	go func(mc <-chan string) {
		for {
			select {
			case <-ca.ctxJobResultsGetter.Done():
				return
			case l, ok := <-mc:
				if ok {
					ca.app.QueueUpdateDraw(func() {
						fmt.Fprintln(textView, l)
					})
				}
			}
		}
	}(ch)
	go func() {
		defer close(ch)
		ca.recorderClient.GetJobExecutions(ca.ctxJobResultsGetter, "", lines, ch)
	}()
}

func renderTableServers(infos []*info.Info, s healthpb.HealthCheckResponse_ServingStatus) *tview.Table {
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
