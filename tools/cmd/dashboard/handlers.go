package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/andrescosta/goico/pkg/service"
	info "github.com/andrescosta/goico/pkg/service/info/grpc"
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/rivo/tview"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func onFocusFileNode(c *CliApp, file string, tenant string) {
	if strings.HasSuffix(file, ".json") {
		pageName := tenant + "/" + file
		if c.mainView.HasPage(pageName) {
			c.mainView.SwitchToPage(pageName)
		} else {
			f, err := c.repoClient.GetFile(context.Background(), tenant, file)
			if err == nil {
				cv := createContentView(string(f))
				c.mainView.AddAndSwitchToPage(pageName, cv, true)
			}
		}
	}
}

func onFocusServerNode(c *CliApp, name string, e *pb.Host) {
	addr := e.Ip + ":" + strconv.Itoa(int(e.Port))
	pageName := addr
	if c.mainView.HasPage(pageName) {
		c.mainView.SwitchToPage(pageName)
	} else {
		helthCheckClient := c.helthCheckClients[addr]
		infoClient, ok := c.infoClients[addr]
		if !ok {
			var err error
			infoClient, err = service.NewGrpcServerInfoClient(addr)
			if err != nil {
				//
				return
			}
			c.infoClients[addr] = infoClient
			helthCheckClient, err = service.NewGrpcServerHelthCheckClient(addr)
			if err != nil {
				//
				return
			}
			c.helthCheckClients[addr] = helthCheckClient
		}
		infos, err := infoClient.Info(context.Background(), &info.InfoRequest{})
		if err != nil {
			//
			return
		}
		s, err := helthCheckClient.Check(context.Background(), name)
		if err != nil {
			s = healthpb.HealthCheckResponse_NOT_SERVING
		}
		view := renderTableServers(infos, s)
		c.mainView.AddAndSwitchToPage(pageName, view, true)
	}
}

func onSelectedGettingJobResults(ca *CliApp, n *tview.TreeNode) {
	n.SetText("<< stop >>")
	nl := n.GetReference().(*node)
	nl.selected = onSelectedStopGettingJobResults
	startGettingJobResults(ca)
}

func onSelectedStopGettingJobResults(ca *CliApp, n *tview.TreeNode) {
	ca.cancelJobResultsGetter()
	nl := n.GetReference().(*node)
	n.SetText("<< start >>")
	nl.selected = onSelectedGettingJobResults
}

func startGettingJobResults(ca *CliApp) {
	var textView *tview.TextView
	lines := int32(5)
	if ca.mainView.HasPage("results") {
		ca.mainView.SwitchToPage("results")
		_, i := ca.mainView.GetFrontPage()
		textView = i.(*tview.TextView)
		lines = 0
	} else {
		textView = tview.NewTextView().
			SetTextAlign(tview.AlignLeft).
			//SetDynamicColors(true).
			//SetRegions(true).
			SetScrollable(true).
			SetWordWrap(false).
			SetWrap(false)
		//textView.SetBorder(true)
		ca.mainView.AddAndSwitchToPage("results", textView, true)
	}
	// fmt.Fprintln(textView, `{"level":"info2","Event":"ev1","Queue":"queue1","Code":1,"Result":"ok","time":"2023-11-30T22:52:02-05:00"}`)
	// fmt.Fprintln(textView, `{"level":"info1","Event":"ev1","Queue":"queue1","Code":1,"Result":"ok","time":"2023-11-30T22:52:02-05:00"}`)
	// print(lines)
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
						//textView.Write([]byte(l))
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
			//SetTextColor(color).
			SetAlign(tview.AlignCenter))
	status := "No ok"
	if s == healthpb.HealthCheckResponse_SERVING {
		status = "Ok"
	}
	table.SetCell(0, 1,
		tview.NewTableCell(status).
			//SetTextColor(color).
			SetAlign(tview.AlignCenter))
	for ix, info := range infos {
		table.SetCell(ix+1, 0,
			tview.NewTableCell(info.Key).
				//SetTextColor(color).
				SetAlign(tview.AlignCenter))
		table.SetCell(ix+1, 1,
			tview.NewTableCell(info.Value).
				//SetTextColor(color).
				SetAlign(tview.AlignCenter))
	}
	return table
}
