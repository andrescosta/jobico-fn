package tapp

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/andrescosta/goico/pkg/service"
	info "github.com/andrescosta/goico/pkg/service/info/grpc"
	"github.com/andrescosta/goico/pkg/yamlico"
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func onFocusFileNode(c *TApp, n *tview.TreeNode) {
	f := (n.GetReference().(*node)).entity.(*sFile)
	if strings.HasSuffix(f.file, ".json") {
		pageName := f.tenant + "/" + f.file
		trySwitchToPage(pageName, c.mainView, c, func() (tview.Primitive, error) {
			f, err := c.repoClient.GetFile(context.Background(), f.tenant, f.file)
			if err != nil {
				return nil, errors.Join(errors.New(`"Repo" service down`), err)
			} else {
				cv := createContentView(string(f))
				return cv, nil
			}
		})
	}
}

func onFocusServerNode(c *TApp, n *tview.TreeNode) {
	h := (n.GetReference().(*node)).entity.(*sServerNode)
	addr := h.host.Ip + ":" + strconv.Itoa(int(h.host.Port))
	trySwitchToPage(addr, c.mainView, c, func() (tview.Primitive, error) {
		helthCheckClient := c.helthCheckClients[addr]
		infoClient, ok := c.infoClients[addr]
		if !ok {
			var err error
			infoClient, err = service.NewGrpcServerInfoClient(addr)
			if err != nil {
				return nil, errors.Join(errors.New(`"Server Info" service down`), err)
			}
			c.infoClients[addr] = infoClient
			helthCheckClient, err = service.NewGrpcServerHelthCheckClient(addr)
			if err != nil {
				return nil, errors.Join(errors.New(`"Healcheck" service down`), err)
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

func onSelectedGettingJobResults(ca *TApp, n *tview.TreeNode) {
	n.SetText("<< stop >>")
	nl := n.GetReference().(*node)
	nl.selected = onSelectedStopGettingJobResults
	startGettingJobResults(ca, n)
}

func onSelectedStopGettingJobResults(ca *TApp, n *tview.TreeNode) {
	ca.cancelJobResultsGetter()
	nl := n.GetReference().(*node)
	n.SetText("<< start >>")
	nl.selected = onSelectedGettingJobResults
}

func onFocusJobPackageNode(c *TApp, n *tview.TreeNode) {
	p := (n.GetReference().(*node)).entity.(*pb.JobPackage)
	pn := "package/" + p.TenantId + "/" + p.ID
	trySwitchToPage(pn, c.mainView, c, func() (tview.Primitive, error) {
		pkg, err := c.controlClient.GetPackage(context.Background(), p.TenantId, &p.ID)
		if err != nil {
			return nil, errors.Join(errors.New(`"Ctl" service down`), err)
		}
		yaml, err := yamlico.Encode(pkg[0])
		if err != nil {
			return nil, errors.Join(errors.New(`package cannot displayed`), err)
		}
		textView := createContentView(decorate(*yaml))
		textView.SetRegions(true)
		textView.SetDynamicColors(true)
		return textView, nil
	})
}

func decorate(yaml string) string {

	reAttributes := regexp.MustCompile(`(?:^|\n).*:`)

	yaml = reAttributes.ReplaceAllStringFunc(yaml, func(match string) string {
		return "[#ff8282]" + match[:len(match)-1] + "[white:black]:"
	})

	reValues := regexp.MustCompile(`: .+\n`)
	yaml = reValues.ReplaceAllStringFunc(yaml, func(match string) string {
		return ": [#d1ffbd]" + match[2:]
	})

	return yaml
}

func startGettingJobResults(ca *TApp, n *tview.TreeNode) {
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
				} else {
					return
				}
			}
		}
	}(ch)
	go func() {
		defer close(ch)
		err := ca.recorderClient.GetJobExecutions(ca.ctxJobResultsGetter, "", lines, ch)
		if err != nil {
			log.Err(err)
			showText(ca.status, "Error getting results", tcell.ColorRed, 3*time.Second, ca)
			ca.app.QueueUpdateDraw(func() {
				onSelectedStopGettingJobResults(ca, n)
				disableTreeNode(n)
			})
		}
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
