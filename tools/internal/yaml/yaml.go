package yaml

import (
	"fmt"
	"os"

	pb "github.com/andrescosta/workflew/api/types"
	"gopkg.in/yaml.v3"
)

func Decode(file string) (*pb.JobPackage, error) {
	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	d := yaml.NewDecoder(r)
	p := pb.JobPackage{}
	if err = d.Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func Encode(j *pb.JobPackage) (*string, error) {
	o, err := yaml.Marshal(j)
	if err != nil {
		return nil, err
	} else {
		s := string(o)
		return &s, err
	}

}

func Debug() {
	id := uint64(1)
	k := pb.JobPackage{
		ID:       &id,
		Name:     "packname",
		TenantId: "merch1",
	}
	o, err := yaml.Marshal(&k)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(o))
	}
}
