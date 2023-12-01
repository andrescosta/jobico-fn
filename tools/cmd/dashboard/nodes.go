package main

import (
	"strconv"

	pb "github.com/andrescosta/workflew/api/types"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type node struct {
	text string
	// true if this node has children and does not allow expansion
	expanded bool
	entity   any
	selected func(*CliApp, *tview.TreeNode)
	// the handler recv the node getting the focus
	focus func(*CliApp, *tview.TreeNode)
	// the handler recv the node loosing the focus and the one getting it
	blur     func(*CliApp, *tview.TreeNode, *tview.TreeNode)
	children []*node
	color    tcell.Color
}

type sFile struct {
	tenant string
	file   string
}

type sServerNode struct {
	name string
	host *pb.Host
}

var rootNode = func(e *pb.Environment, j []*pb.JobPackage, r []*pb.TenantFiles) *node {
	return &node{
		text: "Jobico Manager",
		children: []*node{
			{text: "Packages", entity: e, children: generator(j, jobPackageNode)},
			{text: "Enviroment", entity: e, children: []*node{
				{text: e.ID, entity: e, children: []*node{
					{text: "Services", children: generator(e.Services, serviceNode)},
				}},
			}},
			{text: "Files", entity: e, children: generator(r, tenantFileNode)},
			{text: "(*) Job Results", color: tcell.ColorGreen, expanded: true,
				children: []*node{
					{text: "<< start >>", entity: e,
						selected: onSelectedGettingJobResults,
						focus:    func(c *CliApp, _ *tview.TreeNode) { switchToPageIfExists(c.mainView, "results") },
					},
				}},
		},
	}
}

var serviceNode = func(e *pb.Service) *node {
	return &node{
		text: e.ID, entity: e,
		children: []*node{
			{text: "Servers", children: generatorNamed(e.ID, e.Servers, serverNode)},
			{text: "Storages", children: generator(e.Storages, storageNode)},
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
		children: generatorNamed(e.TenantId, e.Files, fileNode),
	}
}

var fileNode = func(tenant string, file string) *node {
	return &node{
		text: file, entity: &sFile{tenant, file},
		focus: onFocusFileNode,
	}
}
