package dashboard

import (
	"context"
	"errors"
	"strconv"

	"github.com/andrescosta/goico/pkg/service"
	"github.com/andrescosta/goico/pkg/service/grpc"
	"github.com/andrescosta/goico/pkg/service/grpc/svcmeta"

	"github.com/andrescosta/goico/pkg/yamlutil"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/rivo/tview"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func onFocusFileNode(_ context.Context, c *TApp, n *tview.TreeNode) {
	f := (n.GetReference().(*node)).entity.(*sFile)
	if f.file.Type == pb.File_JsonSchema {
		pageName := f.tenant + "/" + f.file.Name
		trySwitchToPage(pageName, c.mainView, c, func() (tview.Primitive, error) {
			f, err := c.repoCli.File(context.Background(), f.tenant, f.file.Name)
			if err != nil {
				return nil, errors.Join(errors.New(`"Repo" service down`), err)
			}
			cv := buildTextView(string(f))
			return cv, nil
		})
	} else {
		switchToEmptyPage(c)
	}
}

func onFocusServerNode(ctx context.Context, c *TApp, n *tview.TreeNode) {
	h := (n.GetReference().(*node)).entity.(*sServerNode)
	addr := h.host.Ip + ":" + strconv.Itoa(int(h.host.Port))
	trySwitchToPage(addr, c.mainView, c, func() (tview.Primitive, error) {
		switch h.host.Type {
		case pb.Host_Undefined:
			c.debugInfo("undefined hos type")
		case pb.Host_Grpc:
			helthCheckClient := c.helthCheckClients[addr]
			infoClient, ok := c.infoClients[addr]
			if !ok {
				var err error
				infoClient, err = svcmeta.NewInfoClient(ctx, addr, service.DefaultGrpcDialer)
				if err != nil {
					return nil, errors.Join(errors.New(`"Server Info" service down`), err)
				}
				c.infoClients[addr] = infoClient
				helthCheckClient, err = grpc.NewHelthCheckClient(ctx, addr, h.name, service.DefaultGrpcDialer)
				if err != nil {
					return nil, errors.Join(errors.New(`"Healcheck" service down`), err)
				}
				c.helthCheckClients[addr] = helthCheckClient
			}
			info, err := infoClient.Info(context.Background(), &svcmeta.GrpcMetadataRequest{})
			if err != nil {
				return nil, err
			}
			s, err := helthCheckClient.Check(context.Background())
			if err != nil {
				s = healthpb.HealthCheckResponse_NOT_SERVING
			}
			view := renderGrpcTableServer(info, s)
			return view, nil
		case pb.Host_Http, pb.Host_Headless:
			metadata, err := c.metadataCli.Metadata(ctx, h.name)
			if err != nil {
				return nil, err
			}
			view := renderHTTPTableServer(metadata)
			return view, nil
		default:
			return nil, nil
		}
		return nil, nil
	})
}

func onSelectedGettingJobResults(_ context.Context, ca *TApp, n *tview.TreeNode) {
	n.SetText("<< stop >>")
	nl := n.GetReference().(*node)
	nl.selected = onSelectedStopGettingJobResults
	ca.startGettingJobResults(n)
}

func onSelectedStopGettingJobResults(_ context.Context, ca *TApp, n *tview.TreeNode) {
	ca.cancelJobResultsGetter()
	nl := n.GetReference().(*node)
	n.SetText("<< start >>")
	nl.selected = onSelectedGettingJobResults
}

func onFocusJobPackageNode(_ context.Context, c *TApp, n *tview.TreeNode) {
	p := (n.GetReference().(*node)).entity.(*pb.JobPackage)
	pn := "package/" + p.Tenant + "/" + p.ID
	trySwitchToPage(pn, c.mainView, c, func() (tview.Primitive, error) {
		pkg, err := c.controlCli.Package(context.Background(), p.Tenant, &p.ID)
		if err != nil {
			return nil, errors.Join(errors.New(`"Ctl" service down`), err)
		}
		yaml, err := yamlutil.Marshal(pkg[0])
		if err != nil {
			return nil, errors.Join(errors.New(`package cannot displayed`), err)
		}
		textView := buildTextView(syntaxHighlightYaml(*yaml))
		return textView, nil
	})
}
