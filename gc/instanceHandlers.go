package gc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

// GetInstanceStatus gets the instance state given the instance id
func GetInstanceStatus(writer http.ResponseWriter, request *http.Request) {
	InitHeaders(writer)
	log.Print("Instance status request invoked!")

	params := mux.Vars(request)
	if len(params) != 0 {
		instanceID := getInstance(params["project"], params["zone"], params["instance"])
		json.NewEncoder(writer).Encode(fmt.Sprintf("Instance status is %#v", instanceID))
	} else {
		log.Fatal("Failed to retrieve query arguments for status")
	}
}

func getInstance(project string, zone string, instance string) uint64 {
	ctx := context.Background()

	c, err := google.DefaultClient(ctx, compute.ComputeScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := computeService.Instances.Get(project, zone, instance).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	return resp.Id
}

// CreateInstance provisisons new instance in GCP
func CreateInstance(writer http.ResponseWriter, request *http.Request) {
	InitHeaders(writer)
	log.Print("Instance create/provsion requested!")

	params := mux.Vars(request)
	log.Print(fmt.Printf("POST params are %#v", params))
	if len(params) != 0 {
		instanceID := provisionInstance(params["project"], params["region"], params["zone"], params["username"], params["userpass"])
		json.NewEncoder(writer).Encode(fmt.Sprintf("Instance provisioned and instance ID is %#v", instanceID))
	} else {
		log.Fatal("Failed to provision new instance!")
	}

}

func provisionInstance(project string, region string, zone string, username string, userpass string) uint64 {
	// customize the attributes based on the instance to be provisioned
	var diskSize int64
	instance := "demo-stack"
	instanceType := "g1-small"
	diskType := "pd-standard"
	image := "projects/ubuntu-os-cloud/global/images/ubuntu-1604-xenial-v20200611"
	diskSize = 10
	accEmail := "gcp service account"
	enAutoRestart := true

	ctx := context.Background()
	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}

	disks := compute.AttachedDisk{
		AutoDelete: true,
		Boot:       true,
		DeviceName: instance,
		DiskSizeGb: diskSize,
		Mode:       "READ_WRITE",
		Type:       "PERSISTENT",
		InitializeParams: &compute.AttachedDiskInitializeParams{
			DiskName:    instance,
			DiskSizeGb:  diskSize,
			DiskType:    fmt.Sprintf("projects/%s/zones/%s/diskTypes/%s", project, zone, diskType),
			SourceImage: image,
		},
	}

	netwInterf := compute.NetworkInterface{
		Subnetwork: fmt.Sprintf("projects/%s/regions/%s/subnetworks/default", project, region),
		AccessConfigs: []*compute.AccessConfig{&compute.AccessConfig{
			Kind:        "compute#accessConfig",
			Name:        "External NAT",
			Type:        "ONE_TO_ONE_NAT",
			NetworkTier: "STANDARD",
		},
		},
	}

	servAcc := compute.ServiceAccount{
		Email: accEmail,
		Scopes: []string{"https://www.googleapis.com/auth/cloud-platform",
			"https://www.googleapis.com/auth/devstorage.read_only",
			"https://www.googleapis.com/auth/logging.write",
			"https://www.googleapis.com/auth/monitoring.write",
			"https://www.googleapis.com/auth/servicecontrol",
			"https://www.googleapis.com/auth/service.management.readonly",
			"https://www.googleapis.com/auth/trace.append",
		},
	}

	secInstConfig := compute.ShieldedInstanceConfig{
		EnableIntegrityMonitoring: true,
		EnableSecureBoot:          false,
		EnableVtpm:                true,
	}

	shell := "/bin/bash"
	startupScript := "#!/bin/bash\npwfile=\"/etc/passwd\"\nshadowfile=\"/etc/shadow\"\ngfile=\"/etc/group\"\nhdir=\"/home\"\n\n# custom metadata to be set during provision\nlogin=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/login -H \"Metadata-Flavor: Google\")\n\nshell=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/shell -H \"Metadata-Flavor: Google\")\n\npass=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/pass -H \"Metadata-Flavor: Google\")\n\nif [ \"$(id -un)\" != \"root\" ]; then\n  echo \"Error: Must be root to run this command.\" >&2\n  exit 1\nfi\n\necho \"Adding new user account $1 to $(hostname)\"\n\nfullname=$login\n\nuid=\"$(awk -F: '{ if (big < $3 && $3 < 5000) big=$3 } END { print big + 1 }' $pwfile)\"\nhomedir=$hdir/$login\ngid=$uid\n\necho ${login}:x:${uid}:${gid}:${fullname}:${homedir}:$shell >>$pwfile\necho ${login}:*:11647:0:99999:7::: >>$shadowfile\necho \"${login}:x:${gid}:$login\" >>$gfile\nmkdir $homedir\ncp -R /etc/skel/.[a-zA-Z]* $homedir\nchmod 755 $homedir\nchown -R ${login}:${login} $homedir\n\necho $login:$pass | chpasswd\n\nusermod -aG sudo $login"

	md := compute.Metadata{
		Kind: "compute#metadata",
		Items: []*compute.MetadataItems{
			{
				Key:   "login",
				Value: &username,
			},
			{
				Key:   "shell",
				Value: &shell,
			},
			{
				Key:   "pass",
				Value: &userpass,
			},
			{
				Key:   "startup-script",
				Value: &startupScript,
			},
		},
	}

	schedConfig := compute.Scheduling{
		AutomaticRestart:  &enAutoRestart,
		Preemptible:       false,
		OnHostMaintenance: "MIGRATE",
		NodeAffinities:    []*compute.SchedulingNodeAffinity{},
	}

	rb := &compute.Instance{
		Name:                   instance,
		Zone:                   fmt.Sprintf("projects/%s/zones/%s", project, zone),
		MachineType:            fmt.Sprintf("projects/%s/zones/%s/machineTypes/%s", project, zone, instanceType),
		CanIpForward:           false,
		Metadata:               &md,
		Disks:                  []*compute.AttachedDisk{&disks},
		NetworkInterfaces:      []*compute.NetworkInterface{&netwInterf},
		ServiceAccounts:        []*compute.ServiceAccount{&servAcc},
		ShieldedInstanceConfig: &secInstConfig,
		Scheduling:             &schedConfig,
	}

	resp, err := computeService.Instances.Insert(project, zone, rb).Context(ctx).Do()
	log.Printf("Response received from GCP req %#v", resp)
	if err != nil {
		log.Fatal(err)
	}

	if resp.HttpErrorStatusCode == 0 && resp.Status == "RUNNING" {
		fmt.Printf("Newly provisioned instance id is: %#v\n", resp.Id)
		return resp.Id
	}

	return 0
}