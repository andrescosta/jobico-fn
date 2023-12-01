package main

import (
	"strconv"

	pb "github.com/andrescosta/workflew/api/types"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type node struct {
	text     string
	expanded bool
	selected func(*CliApp, *tview.TreeNode)
	focus    func(*CliApp)
	blur     func(*CliApp)
	children []*node
	color    tcell.Color
}

var rootNode = func(e *pb.Environment, j []*pb.JobPackage, r []*pb.TenantFiles) *node {
	return &node{
		text: "Jobico Manager",
		children: []*node{
			{text: "Packages", children: jobPackagesNode(j)},
			{text: "Enviroment", children: []*node{
				{text: e.ID, children: []*node{
					{text: "Services", children: servicesNode(e.Services)},
				}},
			}},
			{text: "Files", children: tenantFilesNode(r)},
			{text: "(*) Job Results", color: tcell.ColorGreen, expanded: true,
				children: []*node{
					{text: "<< start >>",
						selected: onSelectedGettingJobResults,
						focus:    func(c *CliApp) { switchToPageIfExists(c.mainView, "results") },
						blur:     func(c *CliApp) { switchToEmptyPage(c.mainView) }},
				}},
		},
	}
}

var serviceNode = func(e *pb.Service) *node {
	return &node{
		text: e.ID,
		children: []*node{
			{text: "Servers", children: serversNode(e.ID, e.Servers)},
			{text: "Storages", children: storagesNode(e.Storages)},
		},
	}
}

var jobPackageNode = func(e *pb.JobPackage) *node {
	return &node{
		text: e.ID,
	}
}
var serverNode = func(name string, e *pb.Host) *node {
	return &node{
		text:  e.Ip + ":" + strconv.Itoa(int(e.Port)),
		blur:  func(c *CliApp) { switchToEmptyPage(c.mainView) },
		focus: func(c *CliApp) { onFocusServerNode(c, name, e) },
	}
}

var storageNode = func(s *pb.Storage) *node {
	return &node{
		text: s.ID,
	}
}

var tenantFileNode = func(e *pb.TenantFiles) *node {
	return &node{
		text:     e.TenantId,
		children: filesNode(e.TenantId, e.Files),
	}
}

var fileNode = func(tenant string, file string) *node {
	return &node{
		text: file, focus: func(c *CliApp) { onFocusFileNode(c, file, tenant) },
		blur: func(c *CliApp) { switchToEmptyPage(c.mainView) },
	}
}

var serversNode = func(name string, e []*pb.Host) []*node {
	r := make([]*node, 0)
	for _, ee := range e {
		r = append(r, serverNode(name, ee))
	}
	return r
}

var storagesNode = func(e []*pb.Storage) []*node {
	r := make([]*node, 0)
	for _, ee := range e {
		r = append(r, storageNode(ee))
	}
	return r
}

var servicesNode = func(e []*pb.Service) []*node {
	r := make([]*node, 0)
	for _, ee := range e {
		r = append(r, serviceNode(ee))
	}
	return r
}

var jobPackagesNode = func(e []*pb.JobPackage) []*node {
	r := make([]*node, 0)
	for _, ee := range e {
		r = append(r, jobPackageNode(ee))
	}
	return r
}

var filesNode = func(merchant string, e []string) []*node {
	r := make([]*node, 0)
	for _, ee := range e {
		r = append(r, fileNode(merchant, ee))
	}
	return r

}

var tenantFilesNode = func(e []*pb.TenantFiles) []*node {
	r := make([]*node, 0)
	for _, ee := range e {
		r = append(r, tenantFileNode(ee))
	}
	return r

}
