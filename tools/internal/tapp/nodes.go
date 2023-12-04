package tapp

import (
	"strconv"

	"github.com/andrescosta/goico/pkg/convertico"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type node struct {
	text string
	// true if this node has children and does not allow expansion
	expanded bool
	entity   any
	selected func(*TApp, *tview.TreeNode)
	// the handler recv the node getting the focus
	focus func(*TApp, *tview.TreeNode)
	// the handler recv the node loosing the focus and the one getting it
	blur         func(*TApp, *tview.TreeNode, *tview.TreeNode)
	children     []*node
	rootNodeType RootNodeType
	color        tcell.Color
}

type RootNodeType int

const (
	NoRootNode RootNodeType = iota
	RootNodeEnv
	RootNodePackage
	RootNodeFile
)

type sFile struct {
	tenant string
	file   string
}

type sServerNode struct {
	name string
	host *pb.Host
}

var rootNode = func(e *pb.Environment, j []*pb.JobPackage, f []*pb.TenantFiles) *node {
	return &node{
		text: "Jobico",
		children: []*node{
			{text: "Packages", entity: e, children: packageChildrenNodes(j), rootNodeType: RootNodePackage},
			{text: "Enviroment", entity: e, children: environmentChildrenNodes(e), rootNodeType: RootNodeEnv},
			{text: "Files", entity: e, children: fileChildrenNodes(f), rootNodeType: RootNodeFile},
			{text: "(*) Job Results", color: tcell.ColorGreen, expanded: true,
				children: []*node{
					{text: "<< start >>", entity: e,
						selected: onSelectedGettingJobResults,
						focus:    func(c *TApp, _ *tview.TreeNode) { switchToPageIfExists(c.mainView, "results") },
					},
				}},
		},
	}
}

var packageChildrenNodes = func(j []*pb.JobPackage) []*node {
	return convertico.SliceWithFunc(j, jobPackageNode)
}

var environmentChildrenNodes = func(e *pb.Environment) []*node {
	return []*node{
		{text: e.ID, entity: e, children: []*node{
			{text: "Services", children: convertico.SliceWithFunc(e.Services, serviceNode)},
		}},
	}
}

var fileChildrenNodes = func(r []*pb.TenantFiles) []*node {
	return convertico.SliceWithFunc(r, tenantFileNode)
}

var serviceNode = func(e *pb.Service) *node {
	return &node{
		text: e.ID, entity: e,
		children: []*node{
			{text: "Servers", children: convertico.SliceWithFuncName(e.ID, e.Servers, serverNode)},
			{text: "Storages", children: convertico.SliceWithFunc(e.Storages, storageNode)},
		},
	}
}

var jobPackageNode = func(e *pb.JobPackage) *node {
	return &node{
		text: e.ID, entity: e, focus: onFocusJobPackageNode,
	}
}

var serverNode = func(name string, e *pb.Host) *node {
	return &node{
		text: e.Ip + ":" + strconv.Itoa(int(e.Port)), entity: &sServerNode{name, e},
		focus: onFocusServerNode,
	}
}

var storageNode = func(s *pb.Storage) *node {
	return &node{
		text: s.ID, entity: s,
	}
}

var tenantFileNode = func(e *pb.TenantFiles) *node {
	return &node{
		text: e.TenantId, entity: e,
		children: convertico.SliceWithFuncName(e.TenantId, e.Files, fileNode),
	}
}

var fileNode = func(tenant string, file string) *node {
	return &node{
		text: file, entity: &sFile{tenant, file},
		focus: onFocusFileNode,
	}
}
