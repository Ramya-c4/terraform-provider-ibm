// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

func ResourceIBMPIInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPIInstanceCreate,
		ReadContext:   resourceIBMPIInstanceRead,
		UpdateContext: resourceIBMPIInstanceUpdate,
		DeleteContext: resourceIBMPIInstanceDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "This is the Power Instance id that is assigned to the account",
			},
			helpers.PIInstanceLicenseRepositoryCapacity: {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Deprecated:  "This field is deprecated.",
				Description: "The VTL license repository capacity TB value",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "PI instance status",
			},
			"min_processors": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Minimum number of the CPUs",
			},
			"min_memory": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Minimum memory",
			},
			"max_processors": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Maximum number of processors",
			},
			"max_memory": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Maximum memory size",
			},
			helpers.PIInstanceVolumeIds: {
				Type:             schema.TypeSet,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Set:              schema.HashString,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "List of PI volumes",
			},
			helpers.PIInstanceUserData: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Base64 encoded data to be passed in for invoking a cloud init script",
			},
			helpers.PIInstanceStorageType: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Storage type for server deployment; if pi_storage_type is not provided the storage type will default to tier3",
			},
			PIInstanceStoragePool: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Storage Pool for server deployment; if provided then pi_storage_pool_affinity will be ignored; Only valid when you deploy one of the IBM supplied stock images. Storage pool for a custom image (an imported image or an image that is created from a VM capture) defaults to the storage pool the image was created in",
			},
			PIAffinityPolicy: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Affinity policy for pvm instance being created; ignored if pi_storage_pool provided; for policy affinity requires one of pi_affinity_instance or pi_affinity_volume to be specified; for policy anti-affinity requires one of pi_anti_affinity_instances or pi_anti_affinity_volumes to be specified",
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"affinity", "anti-affinity"}),
			},
			PIAffinityVolume: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Volume (ID or Name) to base storage affinity policy against; required if requesting affinity and pi_affinity_instance is not provided",
				ConflictsWith: []string{PIAffinityInstance},
			},
			PIAffinityInstance: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "PVM Instance (ID or Name) to base storage affinity policy against; required if requesting storage affinity and pi_affinity_volume is not provided",
				ConflictsWith: []string{PIAffinityVolume},
			},
			PIAntiAffinityVolumes: {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "List of volumes to base storage anti-affinity policy against; required if requesting anti-affinity and pi_anti_affinity_instances is not provided",
				ConflictsWith: []string{PIAntiAffinityInstances},
			},
			PIAntiAffinityInstances: {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "List of pvmInstances to base storage anti-affinity policy against; required if requesting anti-affinity and pi_anti_affinity_volumes is not provided",
				ConflictsWith: []string{PIAntiAffinityVolumes},
			},
			helpers.PIInstanceStorageConnection: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"vSCSI"}),
				Description:  "Storage Connectivity Group for server deployment",
			},
			PIInstanceStoragePoolAffinity: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates if all volumes attached to the server must reside in the same storage pool",
			},
			PIInstanceNetwork: {
				Type:             schema.TypeList,
				DiffSuppressFunc: flex.ApplyOnce,
				Required:         true,
				Description:      "List of one or more networks to attach to the instance",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"network_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			helpers.PIPlacementGroupID: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Placement group ID",
			},
			Arg_PIInstanceSharedProcessorPool: {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{PISAPInstanceProfileID},
				Description:   "Shared Processor Pool the instance is deployed on",
			},
			Attr_PIInstanceSharedProcessorPoolID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Shared Processor Pool ID the instance is deployed on",
			},
			"health_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "PI Instance health status",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID",
			},
			"pin_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "PIN Policy of the Instance",
			},
			helpers.PIInstanceImageId: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "PI instance image id",
				DiffSuppressFunc: flex.ApplyOnce,
			},
			helpers.PIInstanceProcessors: {
				Type:          schema.TypeFloat,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{PISAPInstanceProfileID},
				Description:   "Processors count",
			},
			helpers.PIInstanceName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PI Instance name",
			},
			helpers.PIInstanceProcType: {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validate.ValidateAllowedStringValues([]string{"dedicated", "shared", "capped"}),
				ConflictsWith: []string{PISAPInstanceProfileID},
				Description:   "Instance processor type",
			},
			helpers.PIInstanceSSHKeyName: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "SSH key name",
			},
			helpers.PIInstanceMemory: {
				Type:          schema.TypeFloat,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{PISAPInstanceProfileID},
				Description:   "Memory size",
			},
			PIInstanceDeploymentType: {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"EPIC", "VMNoStorage"}),
				Description:  "Custom Deployment Type Information",
			},
			PISAPInstanceProfileID: {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{helpers.PIInstanceProcessors, helpers.PIInstanceMemory, helpers.PIInstanceProcType},
				Description:   "SAP Profile ID for the amount of cores and memory",
			},
			PISAPInstanceDeploymentType: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Custom SAP Deployment Type Information",
			},
			PIVirtualOpticalDevice: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"attach"}),
				Description:  "Virtual Machine's Cloud Initialization Virtual Optical Device",
			},
			helpers.PIInstanceSystemType: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "PI Instance system type",
			},
			helpers.PIInstanceReplicants: {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Optional:    true,
				Default:     1,
				Description: "PI Instance replicas count",
			},
			helpers.PIInstanceReplicationPolicy: {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"affinity", "anti-affinity", "none"}),
				Default:      "none",
				Description:  "Replication policy for the PI Instance",
			},
			helpers.PIInstanceReplicationScheme: {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"prefix", "suffix"}),
				Default:      "suffix",
				Description:  "Replication scheme",
			},
			helpers.PIInstanceProgress: {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Progress of the operation",
			},
			helpers.PIInstancePinPolicy: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Pin Policy of the instance",
				Default:      "none",
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"none", "soft", "hard"}),
			},
			"operating_system": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operating System",
			},
			"os_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OS Type",
			},
			helpers.PIInstanceHealthStatus: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{helpers.PIInstanceHealthOk, helpers.PIInstanceHealthWarning}),
				Default:      "OK",
				Description:  "Allow the user to set the status of the lpar so that they can connect to it faster",
			},
			helpers.PIVirtualCoresAssigned: {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Virtual Cores Assigned to the PVMInstance",
			},
			"max_virtual_cores": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum Virtual Cores Assigned to the PVMInstance",
			},
			"min_virtual_cores": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Minimum Virtual Cores Assigned to the PVMInstance",
			},
			Arg_IBMiCSS: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "IBM i Cloud Storage Solution",
			},
			Arg_IBMiPHA: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "IBM i Power High Availability",
			},
			Attr_IBMiRDS: {
				Type:        schema.TypeBool,
				Optional:    false,
				Required:    false,
				Computed:    true,
				Description: "IBM i Rational Dev Studio",
			},
			Arg_IBMiRDSUsers: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "IBM i Rational Dev Studio Number of User Licenses",
			},
		},
	}
}

func resourceIBMPIInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("Now in the PowerVMCreate")
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	client := st.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)
	sapClient := st.NewIBMPISAPInstanceClient(ctx, sess, cloudInstanceID)
	imageClient := st.NewIBMPIImageClient(ctx, sess, cloudInstanceID)

	var pvmList *models.PVMInstanceList
	if _, ok := d.GetOk(PISAPInstanceProfileID); ok {
		pvmList, err = createSAPInstance(d, sapClient)
	} else {
		pvmList, err = createPVMInstance(d, client, imageClient)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	var instanceReadyStatus string
	if r, ok := d.GetOk(helpers.PIInstanceHealthStatus); ok {
		instanceReadyStatus = r.(string)
	}

	// id is a combination of the cloud instance id and all of the pvm instance ids
	id := cloudInstanceID
	for _, pvm := range *pvmList {
		id += "/" + *pvm.PvmInstanceID
	}

	d.SetId(id)

	for _, s := range *pvmList {
		if dt, ok := d.GetOk(PIInstanceDeploymentType); ok && dt.(string) == "VMNoStorage" {
			_, err = isWaitForPIInstanceShutoff(ctx, client, *s.PvmInstanceID, instanceReadyStatus)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			_, err = isWaitForPIInstanceAvailable(ctx, client, *s.PvmInstanceID, instanceReadyStatus)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// If Storage Pool Affinity is given as false we need to update the vm instance.
	// Default value is true which indicates that all volumes attached to the server
	// must reside in the same storage pool.
	storagePoolAffinity := d.Get(PIInstanceStoragePoolAffinity).(bool)
	if !storagePoolAffinity {
		for _, s := range *pvmList {
			body := &models.PVMInstanceUpdate{
				StoragePoolAffinity: &storagePoolAffinity,
			}
			// This is a synchronous process hence no need to check for health status
			_, err = client.Update(*s.PvmInstanceID, body)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	// If virtual optical device provided then update cloud initialization
	if vod, ok := d.GetOk(PIVirtualOpticalDevice); ok {
		for _, s := range *pvmList {
			body := &models.PVMInstanceUpdate{
				CloudInitialization: &models.CloudInitialization{
					VirtualOpticalDevice: vod.(string),
				},
			}
			_, err = client.Update(*s.PvmInstanceID, body)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceIBMPIInstanceRead(ctx, d, meta)
}

func resourceIBMPIInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	idArr, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := idArr[0]
	instanceID := idArr[1]

	client := st.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)
	powervmdata, err := client.Get(instanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(helpers.PIInstanceMemory, powervmdata.Memory)
	d.Set(helpers.PIInstanceProcessors, powervmdata.Processors)
	if powervmdata.Status != nil {
		d.Set("status", powervmdata.Status)
	}
	d.Set(helpers.PIInstanceProcType, powervmdata.ProcType)
	d.Set("min_processors", powervmdata.Minproc)
	d.Set(helpers.PIInstanceProgress, powervmdata.Progress)
	if powervmdata.StorageType != nil && *powervmdata.StorageType != "" {
		d.Set(helpers.PIInstanceStorageType, powervmdata.StorageType)
	}
	d.Set(PIInstanceStoragePool, powervmdata.StoragePool)
	d.Set(PIInstanceStoragePoolAffinity, powervmdata.StoragePoolAffinity)
	d.Set(helpers.PICloudInstanceId, cloudInstanceID)
	d.Set("instance_id", powervmdata.PvmInstanceID)
	d.Set(helpers.PIInstanceName, powervmdata.ServerName)
	d.Set(helpers.PIInstanceImageId, powervmdata.ImageID)
	if *powervmdata.PlacementGroup != "none" {
		d.Set(helpers.PIPlacementGroupID, powervmdata.PlacementGroup)
	}
	d.Set(Arg_PIInstanceSharedProcessorPool, powervmdata.SharedProcessorPool)
	d.Set(Attr_PIInstanceSharedProcessorPoolID, powervmdata.SharedProcessorPoolID)

	networksMap := []map[string]interface{}{}
	if powervmdata.Networks != nil {
		for _, n := range powervmdata.Networks {
			if n != nil {
				v := map[string]interface{}{
					"ip_address":   n.IPAddress,
					"mac_address":  n.MacAddress,
					"network_id":   n.NetworkID,
					"network_name": n.NetworkName,
					"type":         n.Type,
					"external_ip":  n.ExternalIP,
				}
				networksMap = append(networksMap, v)
			}
		}
	}
	d.Set(PIInstanceNetwork, networksMap)

	if powervmdata.SapProfile != nil && powervmdata.SapProfile.ProfileID != nil {
		d.Set(PISAPInstanceProfileID, powervmdata.SapProfile.ProfileID)
	}
	d.Set(helpers.PIInstanceSystemType, powervmdata.SysType)
	d.Set("min_memory", powervmdata.Minmem)
	d.Set("max_processors", powervmdata.Maxproc)
	d.Set("max_memory", powervmdata.Maxmem)
	d.Set("pin_policy", powervmdata.PinPolicy)
	d.Set("operating_system", powervmdata.OperatingSystem)
	d.Set("os_type", powervmdata.OsType)

	if powervmdata.Health != nil {
		d.Set("health_status", powervmdata.Health.Status)
	}
	if powervmdata.VirtualCores != nil {
		d.Set(helpers.PIVirtualCoresAssigned, powervmdata.VirtualCores.Assigned)
		d.Set("max_virtual_cores", powervmdata.VirtualCores.Max)
		d.Set("min_virtual_cores", powervmdata.VirtualCores.Min)
	}
	d.Set(helpers.PIInstanceLicenseRepositoryCapacity, powervmdata.LicenseRepositoryCapacity)
	d.Set(PIInstanceDeploymentType, powervmdata.DeploymentType)
	if powervmdata.SoftwareLicenses != nil {
		d.Set(Arg_IBMiCSS, powervmdata.SoftwareLicenses.IbmiCSS)
		d.Set(Arg_IBMiPHA, powervmdata.SoftwareLicenses.IbmiPHA)
		d.Set(Attr_IBMiRDS, powervmdata.SoftwareLicenses.IbmiRDS)
		if *powervmdata.SoftwareLicenses.IbmiRDS {
			d.Set(Arg_IBMiRDSUsers, powervmdata.SoftwareLicenses.IbmiRDSUsers)
		} else {
			d.Set(Arg_IBMiRDSUsers, 0)
		}
	}
	return nil
}

func resourceIBMPIInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	name := d.Get(helpers.PIInstanceName).(string)
	mem := d.Get(helpers.PIInstanceMemory).(float64)
	procs := d.Get(helpers.PIInstanceProcessors).(float64)
	processortype := d.Get(helpers.PIInstanceProcType).(string)
	assignedVirtualCores := int64(d.Get(helpers.PIVirtualCoresAssigned).(int))

	if d.Get("health_status") == "WARNING" {
		return diag.Errorf("the operation cannot be performed when the lpar health in the WARNING State")
	}

	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.Errorf("failed to get the session from the IBM Cloud Service")
	}

	cloudInstanceID, instanceID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := st.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)

	// Check if cloud instance is capable of changing virtual cores
	cloudInstanceClient := st.NewIBMPICloudInstanceClient(ctx, sess, cloudInstanceID)
	cloudInstance, err := cloudInstanceClient.Get(cloudInstanceID)
	if err != nil {
		return diag.FromErr(err)
	}
	cores_enabled := checkCloudInstanceCapability(cloudInstance, CUSTOM_VIRTUAL_CORES)

	if d.HasChanges(helpers.PIInstanceName, PIVirtualOpticalDevice) {
		body := &models.PVMInstanceUpdate{}
		if d.HasChange(helpers.PIInstanceName) {
			body.ServerName = name
		}
		if d.HasChange(PIVirtualOpticalDevice) {
			body.CloudInitialization.VirtualOpticalDevice = d.Get(PIVirtualOpticalDevice).(string)
		}
		_, err = client.Update(instanceID, body)
		if err != nil {
			return diag.Errorf("failed to update the lpar: %v", err)
		}
		_, err = isWaitForPIInstanceAvailable(ctx, client, instanceID, "OK")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange(helpers.PIInstanceProcType) {
		// Stop the lpar
		if d.Get("status") == "SHUTOFF" {
			log.Printf("the lpar is in the shutoff state. Nothing to do . Moving on ")
		} else {
			err := stopLparForResourceChange(ctx, client, instanceID)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		// Modify
		log.Printf("At this point the lpar should be off. Executing the Processor Update Change")
		updatebody := &models.PVMInstanceUpdate{ProcType: processortype}
		if cores_enabled {
			log.Printf("support for %s is enabled", CUSTOM_VIRTUAL_CORES)
			updatebody.VirtualCores = &models.VirtualCores{Assigned: &assignedVirtualCores}
		} else {
			log.Printf("no virtual cores support enabled for this customer..")
		}
		_, err = client.Update(instanceID, updatebody)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = isWaitForPIInstanceStopped(ctx, client, instanceID)
		if err != nil {
			return diag.FromErr(err)
		}

		// Start the lpar
		err := startLparAfterResourceChange(ctx, client, instanceID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Virtual core will be updated only if service instance capability is enabled
	if d.HasChange(helpers.PIVirtualCoresAssigned) {
		body := &models.PVMInstanceUpdate{
			VirtualCores: &models.VirtualCores{Assigned: &assignedVirtualCores},
		}
		_, err = client.Update(instanceID, body)
		if err != nil {
			return diag.Errorf("failed to update the lpar with the change for virtual cores: %v", err)
		}
		_, err = isWaitForPIInstanceAvailable(ctx, client, instanceID, "OK")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Start of the change for Memory and Processors
	if d.HasChange(helpers.PIInstanceMemory) || d.HasChange(helpers.PIInstanceProcessors) {

		maxMemLpar := d.Get("max_memory").(float64)
		maxCPULpar := d.Get("max_processors").(float64)
		//log.Printf("the required memory is set to [%d] and current max memory is set to  [%d] ", int(mem), int(maxMemLpar))

		if mem > maxMemLpar || procs > maxCPULpar {
			log.Printf("Will require a shutdown to perform the change")
		} else {
			log.Printf("maxMemLpar is set to %f", maxMemLpar)
			log.Printf("maxCPULpar is set to %f", maxCPULpar)
		}

		//if d.GetOkExists("reboot_for_resource_change")

		instanceState := d.Get("status")
		log.Printf("the instance state is %s", instanceState)

		if (mem > maxMemLpar || procs > maxCPULpar) && instanceState != "SHUTOFF" {
			err = performChangeAndReboot(ctx, client, instanceID, cloudInstanceID, mem, procs)
			if err != nil {
				return diag.FromErr(err)
			}

		} else {
			body := &models.PVMInstanceUpdate{
				Memory:     mem,
				Processors: procs,
			}
			if cores_enabled {
				log.Printf("support for %s is enabled", CUSTOM_VIRTUAL_CORES)
				body.VirtualCores = &models.VirtualCores{Assigned: &assignedVirtualCores}
			} else {
				log.Printf("no virtual cores support enabled for this customer..")
			}

			_, err = client.Update(instanceID, body)
			if err != nil {
				return diag.Errorf("failed to update the lpar with the change %v", err)
			}
			if instanceState == "SHUTOFF" {
				_, err = isWaitforPIInstanceUpdate(ctx, client, instanceID)
				if err != nil {
					return diag.FromErr(err)
				}
			} else {
				_, err = isWaitForPIInstanceAvailable(ctx, client, instanceID, "OK")
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	// License repository capacity will be updated only if service instance is a vtl instance
	// might need to check if lrc was set
	if d.HasChange(helpers.PIInstanceLicenseRepositoryCapacity) {
		lrc := d.Get(helpers.PIInstanceLicenseRepositoryCapacity).(int64)
		body := &models.PVMInstanceUpdate{
			LicenseRepositoryCapacity: lrc,
		}
		_, err = client.Update(instanceID, body)
		if err != nil {
			return diag.Errorf("failed to update the lpar with the change for license repository capacity %s", err)
		}
		_, err = isWaitForPIInstanceAvailable(ctx, client, instanceID, "OK")
		if err != nil {
			diag.FromErr(err)
		}
	}

	if d.HasChange(PISAPInstanceProfileID) {
		// Stop the lpar
		if d.Get("status") == "SHUTOFF" {
			log.Printf("the lpar is in the shutoff state. Nothing to do... Moving on ")
		} else {
			err := stopLparForResourceChange(ctx, client, instanceID)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		// Update the profile id
		profileID := d.Get(PISAPInstanceProfileID).(string)
		body := &models.PVMInstanceUpdate{
			SapProfileID: profileID,
		}
		_, err = client.Update(instanceID, body)
		if err != nil {
			return diag.Errorf("failed to update the lpar with the change for sap profile: %v", err)
		}

		// Wait for the resize to complete and status to reset
		_, err = isWaitForPIInstanceStopped(ctx, client, instanceID)
		if err != nil {
			return diag.FromErr(err)
		}

		// Start the lpar
		err := startLparAfterResourceChange(ctx, client, instanceID)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange(PIInstanceStoragePoolAffinity) {
		storagePoolAffinity := d.Get(PIInstanceStoragePoolAffinity).(bool)
		body := &models.PVMInstanceUpdate{
			StoragePoolAffinity: &storagePoolAffinity,
		}
		// This is a synchronous process hence no need to check for health status
		_, err = client.Update(instanceID, body)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange(helpers.PIPlacementGroupID) {
		pgClient := st.NewIBMPIPlacementGroupClient(ctx, sess, cloudInstanceID)

		oldRaw, newRaw := d.GetChange(helpers.PIPlacementGroupID)
		old := oldRaw.(string)
		new := newRaw.(string)

		if len(strings.TrimSpace(old)) > 0 {
			placementGroupID := old
			//remove server from old placement group
			body := &models.PlacementGroupServer{
				ID: &instanceID,
			}
			pgID, err := pgClient.DeleteMember(placementGroupID, body)
			if err != nil {
				// ignore delete member error where the server is already not in the PG
				if !strings.Contains(err.Error(), "is not part of placement-group") {
					return diag.FromErr(err)
				}
			} else {
				_, err = isWaitForPIInstancePlacementGroupDelete(ctx, pgClient, *pgID.ID, instanceID)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		if len(strings.TrimSpace(new)) > 0 {
			placementGroupID := new
			// add server to a new placement group
			body := &models.PlacementGroupServer{
				ID: &instanceID,
			}
			pgID, err := pgClient.AddMember(placementGroupID, body)
			if err != nil {
				return diag.FromErr(err)
			} else {
				_, err = isWaitForPIInstancePlacementGroupAdd(ctx, pgClient, *pgID.ID, instanceID)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}
	if d.HasChanges(Arg_IBMiCSS, Arg_IBMiPHA, Arg_IBMiRDSUsers) {
		if d.Get("status") == "ACTIVE" {
			log.Printf("the lpar is in the Active state, continuing with update")
		} else {
			_, err = isWaitForPIInstanceAvailable(ctx, client, instanceID, "OK")
			if err != nil {
				return diag.FromErr(err)
			}
		}

		sl := &models.SoftwareLicenses{}
		sl.IbmiCSS = flex.PtrToBool(d.Get(Arg_IBMiCSS).(bool))
		sl.IbmiPHA = flex.PtrToBool(d.Get(Arg_IBMiPHA).(bool))
		ibmrdsUsers := d.Get(Arg_IBMiRDSUsers).(int)
		if ibmrdsUsers < 0 {
			return diag.Errorf("request with  IBM i Rational Dev Studio property requires IBM i Rational Dev Studio number of users")
		}
		sl.IbmiRDS = flex.PtrToBool(ibmrdsUsers > 0)
		sl.IbmiRDSUsers = int64(ibmrdsUsers)

		updatebody := &models.PVMInstanceUpdate{SoftwareLicenses: sl}
		_, err = client.Update(instanceID, updatebody)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = isWaitForPIInstanceSoftwareLicenses(ctx, client, instanceID, sl)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceIBMPIInstanceRead(ctx, d, meta)
}

func resourceIBMPIInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	idArr, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := idArr[0]
	client := st.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)
	for _, instanceID := range idArr[1:] {
		err = client.Delete(instanceID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	for _, instanceID := range idArr[1:] {
		_, err = isWaitForPIInstanceDeleted(ctx, client, instanceID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")
	return nil
}

func isWaitForPIInstanceDeleted(ctx context.Context, client *st.IBMPIInstanceClient, id string) (interface{}, error) {

	log.Printf("Waiting for  (%s) to be deleted.", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"retry", helpers.PIInstanceDeleting},
		Target:     []string{helpers.PIInstanceNotFound},
		Refresh:    isPIInstanceDeleteRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
		Timeout:    10 * time.Minute,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIInstanceDeleteRefreshFunc(client *st.IBMPIInstanceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pvm, err := client.Get(id)
		if err != nil {
			log.Printf("The power vm does not exist")
			return pvm, helpers.PIInstanceNotFound, nil
		}
		return pvm, helpers.PIInstanceDeleting, nil
	}
}

func isWaitForPIInstanceAvailable(ctx context.Context, client *st.IBMPIInstanceClient, id string, instanceReadyStatus string) (interface{}, error) {
	log.Printf("Waiting for PIInstance (%s) to be available and active ", id)

	queryTimeOut := activeTimeOut
	if instanceReadyStatus == helpers.PIInstanceHealthWarning {
		queryTimeOut = warningTimeOut
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"PENDING", helpers.PIInstanceBuilding, helpers.PIInstanceHealthWarning},
		Target:     []string{helpers.PIInstanceAvailable, helpers.PIInstanceHealthOk, "ERROR", "", "SHUTOFF"},
		Refresh:    isPIInstanceRefreshFunc(client, id, instanceReadyStatus),
		Delay:      30 * time.Second,
		MinTimeout: queryTimeOut,
		Timeout:    120 * time.Minute,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIInstanceRefreshFunc(client *st.IBMPIInstanceClient, id, instanceReadyStatus string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {

		pvm, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}
		// Check for `instanceReadyStatus` health status and also the final health status "OK"
		if *pvm.Status == helpers.PIInstanceAvailable && (pvm.Health.Status == instanceReadyStatus || pvm.Health.Status == helpers.PIInstanceHealthOk) {
			return pvm, helpers.PIInstanceAvailable, nil
		}
		if *pvm.Status == "ERROR" {
			if pvm.Fault != nil {
				err = fmt.Errorf("failed to create the lpar: %s", pvm.Fault.Message)
			} else {
				err = fmt.Errorf("failed to create the lpar")
			}
			return pvm, *pvm.Status, err
		}

		return pvm, helpers.PIInstanceBuilding, nil
	}
}

func isWaitForPIInstancePlacementGroupAdd(ctx context.Context, client *st.IBMPIPlacementGroupClient, pgID string, id string) (interface{}, error) {
	log.Printf("Waiting for PIInstance Placement Group (%s) to be updated ", id)

	queryTimeOut := activeTimeOut

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Adding},
		Target:     []string{State_Added},
		Refresh:    isPIInstancePlacementGroupAddRefreshFunc(client, pgID, id),
		Delay:      30 * time.Second,
		MinTimeout: queryTimeOut,
		Timeout:    10 * time.Minute,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIInstancePlacementGroupAddRefreshFunc(client *st.IBMPIPlacementGroupClient, pgID string, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pg, err := client.Get(pgID)
		if err != nil {
			return nil, "", err
		}
		for _, x := range pg.Members {
			if x == id {
				return pg, State_Added, nil
			}
		}
		return pg, State_Adding, nil
	}
}

func isWaitForPIInstancePlacementGroupDelete(ctx context.Context, client *st.IBMPIPlacementGroupClient, pgID string, id string) (interface{}, error) {
	log.Printf("Waiting for PIInstance Placement Group (%s) to be updated ", id)

	queryTimeOut := activeTimeOut

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Deleting},
		Target:     []string{State_Deleted},
		Refresh:    isPIInstancePlacementGroupDeleteRefreshFunc(client, pgID, id),
		Delay:      30 * time.Second,
		MinTimeout: queryTimeOut,
		Timeout:    10 * time.Minute,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIInstancePlacementGroupDeleteRefreshFunc(client *st.IBMPIPlacementGroupClient, pgID string, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pg, err := client.Get(pgID)
		if err != nil {
			return nil, "", err
		}
		for _, x := range pg.Members {
			if x == id {
				return pg, State_Deleting, nil
			}
		}
		return pg, State_Deleted, nil
	}
}

func isWaitForPIInstanceSoftwareLicenses(ctx context.Context, client *st.IBMPIInstanceClient, id string, softwareLicenses *models.SoftwareLicenses) (interface{}, error) {
	log.Printf("Waiting for PIInstance Software Licenses (%s) to be updated ", id)

	queryTimeOut := activeTimeOut

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"notdone"},
		Target:     []string{"done"},
		Refresh:    isPIInstanceSoftwareLicensesRefreshFunc(client, id, softwareLicenses),
		Delay:      90 * time.Second,
		MinTimeout: queryTimeOut,
		Timeout:    120 * time.Minute,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIInstanceSoftwareLicensesRefreshFunc(client *st.IBMPIInstanceClient, id string, softwareLicenses *models.SoftwareLicenses) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {

		pvm, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		// Check that each software license we modified has been updated
		if softwareLicenses.IbmiCSS != nil {
			if *softwareLicenses.IbmiCSS != *pvm.SoftwareLicenses.IbmiCSS {
				return pvm, "notdone", nil
			}
		}

		if softwareLicenses.IbmiPHA != nil {
			if *softwareLicenses.IbmiPHA != *pvm.SoftwareLicenses.IbmiPHA {
				return pvm, "notdone", nil
			}
		}

		if softwareLicenses.IbmiRDS != nil {
			// If the update set IBMiRDS to false, don't check IBMiRDSUsers as it will be updated on the terraform side on the read
			if !*softwareLicenses.IbmiRDS {
				if *softwareLicenses.IbmiRDS != *pvm.SoftwareLicenses.IbmiRDS {
					return pvm, "notdone", nil
				}
			} else if (*softwareLicenses.IbmiRDS != *pvm.SoftwareLicenses.IbmiRDS) || (softwareLicenses.IbmiRDSUsers != pvm.SoftwareLicenses.IbmiRDSUsers) {
				return pvm, "notdone", nil
			}
		}

		return pvm, "done", nil
	}
}

func isWaitForPIInstanceShutoff(ctx context.Context, client *st.IBMPIInstanceClient, id string, instanceReadyStatus string) (interface{}, error) {
	log.Printf("Waiting for PIInstance (%s) to be shutoff and health active ", id)

	queryTimeOut := activeTimeOut
	if instanceReadyStatus == helpers.PIInstanceHealthWarning {
		queryTimeOut = warningTimeOut
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{StatusPending, helpers.PIInstanceBuilding, helpers.PIInstanceHealthWarning},
		Target:     []string{helpers.PIInstanceHealthOk, StatusError, "", StatusShutoff},
		Refresh:    isPIInstanceShutoffRefreshFunc(client, id, instanceReadyStatus),
		Delay:      30 * time.Second,
		MinTimeout: queryTimeOut,
		Timeout:    120 * time.Minute,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIInstanceShutoffRefreshFunc(client *st.IBMPIInstanceClient, id, instanceReadyStatus string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {

		pvm, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}
		if *pvm.Status == StatusShutoff && (pvm.Health.Status == instanceReadyStatus || pvm.Health.Status == helpers.PIInstanceHealthOk) {
			return pvm, StatusShutoff, nil
		}
		if *pvm.Status == StatusError {
			if pvm.Fault != nil {
				err = fmt.Errorf("failed to create the lpar: %s", pvm.Fault.Message)
			} else {
				err = fmt.Errorf("failed to create the lpar")
			}
			return pvm, *pvm.Status, err
		}

		return pvm, helpers.PIInstanceBuilding, nil
	}
}

// This function takes the input string and encodes into base64 if isn't already encoded
func encodeBase64(userData string) string {
	_, err := base64.StdEncoding.DecodeString(userData)
	if err != nil {
		return base64.StdEncoding.EncodeToString([]byte(userData))
	}
	return userData
}

func isWaitForPIInstanceStopped(ctx context.Context, client *st.IBMPIInstanceClient, id string) (interface{}, error) {
	log.Printf("Waiting for PIInstance (%s) to be stopped and powered off ", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"STOPPING", "RESIZE", "VERIFY_RESIZE", helpers.PIInstanceHealthWarning},
		Target:     []string{"OK", "SHUTOFF"},
		Refresh:    isPIInstanceRefreshFuncOff(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 2 * time.Minute, // This is the time that the client will execute to check the status of the request
		Timeout:    30 * time.Minute,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIInstanceRefreshFuncOff(client *st.IBMPIInstanceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {

		log.Printf("Calling the check Refresh status of the pvm instance %s", id)
		pvm, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}
		if *pvm.Status == "SHUTOFF" && pvm.Health.Status == helpers.PIInstanceHealthOk {
			return pvm, "SHUTOFF", nil
		}
		return pvm, "STOPPING", nil
	}
}

func stopLparForResourceChange(ctx context.Context, client *st.IBMPIInstanceClient, id string) error {
	body := &models.PVMInstanceAction{
		//Action: flex.PtrToString("stop"),
		Action: flex.PtrToString("immediate-shutdown"),
	}
	err := client.Action(id, body)
	if err != nil {
		return fmt.Errorf("failed to perform the stop action on the pvm instance %v", err)
	}

	_, err = isWaitForPIInstanceStopped(ctx, client, id)

	return err
}

// Start the lpar
func startLparAfterResourceChange(ctx context.Context, client *st.IBMPIInstanceClient, id string) error {
	body := &models.PVMInstanceAction{
		Action: flex.PtrToString("start"),
	}
	err := client.Action(id, body)
	if err != nil {
		return fmt.Errorf("failed to perform the start action on the pvm instance %v", err)
	}

	_, err = isWaitForPIInstanceAvailable(ctx, client, id, "OK")

	return err
}

// Stop / Modify / Start only when the lpar is off limits
func performChangeAndReboot(ctx context.Context, client *st.IBMPIInstanceClient, id, cloudInstanceID string, mem, procs float64) error {
	/*
		These are the steps
		1. Stop the lpar - Check if the lpar is SHUTOFF
		2. Once the lpar is SHUTOFF - Make the cpu / memory change - DUring this time , you can check for RESIZE and VERIFY_RESIZE as the transition states
		3. If the change is successful , the lpar state will be back in SHUTOFF
		4. Once the LPAR state is SHUTOFF , initiate the start again and check for ACTIVE + OK
	*/
	//Execute the stop

	log.Printf("Calling the stop lpar for Resource Change code ..")
	err := stopLparForResourceChange(ctx, client, id)
	if err != nil {
		return err
	}

	body := &models.PVMInstanceUpdate{
		Memory:     mem,
		Processors: procs,
	}

	_, updateErr := client.Update(id, body)
	if updateErr != nil {
		return fmt.Errorf("failed to update the lpar with the change, %s", updateErr)
	}

	_, err = isWaitforPIInstanceUpdate(ctx, client, id)
	if err != nil {
		return fmt.Errorf("failed to get an update from the Service after the resource change, %s", err)
	}

	// Now we can start the lpar
	log.Printf("Calling the start lpar After the  Resource Change code ..")
	err = startLparAfterResourceChange(ctx, client, id)
	if err != nil {
		return err
	}

	return nil

}

func isWaitforPIInstanceUpdate(ctx context.Context, client *st.IBMPIInstanceClient, id string) (interface{}, error) {
	log.Printf("Waiting for PIInstance (%s) to be ACTIVE or SHUTOFF AFTER THE RESIZE Due to DLPAR Operation ", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"RESIZE", "VERIFY_RESIZE"},
		Target:     []string{"ACTIVE", "SHUTOFF", helpers.PIInstanceHealthOk},
		Refresh:    isPIInstanceShutAfterResourceChange(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Minute,
		Timeout:    60 * time.Minute,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIInstanceShutAfterResourceChange(client *st.IBMPIInstanceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {

		pvm, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if *pvm.Status == "SHUTOFF" && pvm.Health.Status == helpers.PIInstanceHealthOk {
			log.Printf("The lpar is now off after the resource change...")
			return pvm, "SHUTOFF", nil
		}

		return pvm, "RESIZE", nil
	}
}

func expandPVMNetworks(networks []interface{}) []*models.PVMInstanceAddNetwork {
	pvmNetworks := make([]*models.PVMInstanceAddNetwork, 0, len(networks))
	for _, v := range networks {
		network := v.(map[string]interface{})
		pvmInstanceNetwork := &models.PVMInstanceAddNetwork{
			IPAddress: network["ip_address"].(string),
			NetworkID: flex.PtrToString(network["network_id"].(string)),
		}
		pvmNetworks = append(pvmNetworks, pvmInstanceNetwork)
	}
	return pvmNetworks
}

func checkCloudInstanceCapability(cloudInstance *models.CloudInstance, custom_capability string) bool {
	log.Printf("Checking for the following capability %s", custom_capability)
	log.Printf("the instance features are %s", cloudInstance.Capabilities)
	for _, v := range cloudInstance.Capabilities {
		if v == custom_capability {
			return true
		}
	}
	return false
}

func createSAPInstance(d *schema.ResourceData, sapClient *st.IBMPISAPInstanceClient) (*models.PVMInstanceList, error) {

	name := d.Get(helpers.PIInstanceName).(string)
	profileID := d.Get(PISAPInstanceProfileID).(string)
	imageid := d.Get(helpers.PIInstanceImageId).(string)

	pvmNetworks := expandPVMNetworks(d.Get(PIInstanceNetwork).([]interface{}))

	var replicants int64
	if r, ok := d.GetOk(helpers.PIInstanceReplicants); ok {
		replicants = int64(r.(int))
	}
	var replicationpolicy string
	if r, ok := d.GetOk(helpers.PIInstanceReplicationPolicy); ok {
		replicationpolicy = r.(string)
	}
	var replicationNamingScheme string
	if r, ok := d.GetOk(helpers.PIInstanceReplicationScheme); ok {
		replicationNamingScheme = r.(string)
	}
	instances := &models.PVMInstanceMultiCreate{
		AffinityPolicy: &replicationpolicy,
		Count:          replicants,
		Numerical:      &replicationNamingScheme,
	}

	body := &models.SAPCreate{
		ImageID:   &imageid,
		Instances: instances,
		Name:      &name,
		Networks:  pvmNetworks,
		ProfileID: &profileID,
	}

	if v, ok := d.GetOk(PISAPInstanceDeploymentType); ok {
		body.DeploymentType = v.(string)
	}
	if v, ok := d.GetOk(helpers.PIInstanceVolumeIds); ok {
		volids := flex.ExpandStringList((v.(*schema.Set)).List())
		if len(volids) > 0 {
			body.VolumeIDs = volids
		}
	}
	if p, ok := d.GetOk(helpers.PIInstancePinPolicy); ok {
		pinpolicy := p.(string)
		if d.Get(helpers.PIInstancePinPolicy) == "soft" || d.Get(helpers.PIInstancePinPolicy) == "hard" {
			body.PinPolicy = models.PinPolicy(pinpolicy)
		}
	}

	if v, ok := d.GetOk(helpers.PIInstanceSSHKeyName); ok {
		sshkey := v.(string)
		body.SSHKeyName = sshkey
	}
	if u, ok := d.GetOk(helpers.PIInstanceUserData); ok {
		userData := u.(string)
		body.UserData = encodeBase64(userData)
	}
	if sys, ok := d.GetOk(helpers.PIInstanceSystemType); ok {
		body.SysType = sys.(string)
	}

	if st, ok := d.GetOk(helpers.PIInstanceStorageType); ok {
		body.StorageType = st.(string)
	}
	if sp, ok := d.GetOk(PIInstanceStoragePool); ok {
		body.StoragePool = sp.(string)
	}

	if ap, ok := d.GetOk(PIAffinityPolicy); ok {
		policy := ap.(string)
		affinity := &models.StorageAffinity{
			AffinityPolicy: &policy,
		}

		if policy == "affinity" {
			if av, ok := d.GetOk(PIAffinityVolume); ok {
				afvol := av.(string)
				affinity.AffinityVolume = &afvol
			}
			if ai, ok := d.GetOk(PIAffinityInstance); ok {
				afins := ai.(string)
				affinity.AffinityPVMInstance = &afins
			}
		} else {
			if avs, ok := d.GetOk(PIAntiAffinityVolumes); ok {
				afvols := flex.ExpandStringList(avs.([]interface{}))
				affinity.AntiAffinityVolumes = afvols
			}
			if ais, ok := d.GetOk(PIAntiAffinityInstances); ok {
				afinss := flex.ExpandStringList(ais.([]interface{}))
				affinity.AntiAffinityPVMInstances = afinss
			}
		}
		body.StorageAffinity = affinity
	}

	if pg, ok := d.GetOk(helpers.PIPlacementGroupID); ok {
		body.PlacementGroup = pg.(string)
	}

	pvmList, err := sapClient.Create(body)
	if err != nil {
		return nil, fmt.Errorf("failed to provision: %v", err)
	}
	if pvmList == nil {
		return nil, fmt.Errorf("failed to provision")
	}

	return pvmList, nil
}

func createPVMInstance(d *schema.ResourceData, client *st.IBMPIInstanceClient, imageClient *st.IBMPIImageClient) (*models.PVMInstanceList, error) {

	name := d.Get(helpers.PIInstanceName).(string)
	imageid := d.Get(helpers.PIInstanceImageId).(string)

	var mem, procs float64
	var systype, processortype string
	if v, ok := d.GetOk(helpers.PIInstanceMemory); ok {
		mem = v.(float64)
	} else {
		return nil, fmt.Errorf("%s is required for creating pvm instances", helpers.PIInstanceMemory)
	}
	if v, ok := d.GetOk(helpers.PIInstanceProcessors); ok {
		procs = v.(float64)
	} else {
		return nil, fmt.Errorf("%s is required for creating pvm instances", helpers.PIInstanceProcessors)
	}
	if v, ok := d.GetOk(helpers.PIInstanceSystemType); ok {
		systype = v.(string)
	} else {
		return nil, fmt.Errorf("%s is required for creating pvm instances", helpers.PIInstanceSystemType)
	}
	if v, ok := d.GetOk(helpers.PIInstanceProcType); ok {
		processortype = v.(string)
	} else {
		return nil, fmt.Errorf("%s is required for creating pvm instances", helpers.PIInstanceProcType)
	}

	pvmNetworks := expandPVMNetworks(d.Get(PIInstanceNetwork).([]interface{}))

	var volids []string
	if v, ok := d.GetOk(helpers.PIInstanceVolumeIds); ok {
		volids = flex.ExpandStringList((v.(*schema.Set)).List())
	}
	var replicants float64
	if r, ok := d.GetOk(helpers.PIInstanceReplicants); ok {
		replicants = float64(r.(int))
	}
	var replicationpolicy string
	if r, ok := d.GetOk(helpers.PIInstanceReplicationPolicy); ok {
		replicationpolicy = r.(string)
	}
	var replicationNamingScheme string
	if r, ok := d.GetOk(helpers.PIInstanceReplicationScheme); ok {
		replicationNamingScheme = r.(string)
	}
	var pinpolicy string
	if p, ok := d.GetOk(helpers.PIInstancePinPolicy); ok {
		pinpolicy = p.(string)
		if pinpolicy == "" {
			pinpolicy = "none"
		}
	}

	var userData string
	if u, ok := d.GetOk(helpers.PIInstanceUserData); ok {
		userData = u.(string)
	}

	body := &models.PVMInstanceCreate{
		Processors:              &procs,
		Memory:                  &mem,
		ServerName:              flex.PtrToString(name),
		SysType:                 systype,
		ImageID:                 flex.PtrToString(imageid),
		ProcType:                flex.PtrToString(processortype),
		Replicants:              replicants,
		UserData:                encodeBase64(userData),
		ReplicantNamingScheme:   flex.PtrToString(replicationNamingScheme),
		ReplicantAffinityPolicy: flex.PtrToString(replicationpolicy),
		Networks:                pvmNetworks,
	}
	if s, ok := d.GetOk(helpers.PIInstanceSSHKeyName); ok {
		sshkey := s.(string)
		body.KeyPairName = sshkey
	}
	if len(volids) > 0 {
		body.VolumeIDs = volids
	}
	if d.Get(helpers.PIInstancePinPolicy) == "soft" || d.Get(helpers.PIInstancePinPolicy) == "hard" {
		body.PinPolicy = models.PinPolicy(pinpolicy)
	}

	var assignedVirtualCores int64
	if a, ok := d.GetOk(helpers.PIVirtualCoresAssigned); ok {
		assignedVirtualCores = int64(a.(int))
		body.VirtualCores = &models.VirtualCores{Assigned: &assignedVirtualCores}
	}

	if st, ok := d.GetOk(helpers.PIInstanceStorageType); ok {
		body.StorageType = st.(string)
	}
	if sp, ok := d.GetOk(PIInstanceStoragePool); ok {
		body.StoragePool = sp.(string)
	}

	if dt, ok := d.GetOk(PIInstanceDeploymentType); ok {
		body.DeploymentType = dt.(string)
	}

	if ap, ok := d.GetOk(PIAffinityPolicy); ok {
		policy := ap.(string)
		affinity := &models.StorageAffinity{
			AffinityPolicy: &policy,
		}

		if policy == "affinity" {
			if av, ok := d.GetOk(PIAffinityVolume); ok {
				afvol := av.(string)
				affinity.AffinityVolume = &afvol
			}
			if ai, ok := d.GetOk(PIAffinityInstance); ok {
				afins := ai.(string)
				affinity.AffinityPVMInstance = &afins
			}
		} else {
			if avs, ok := d.GetOk(PIAntiAffinityVolumes); ok {
				afvols := flex.ExpandStringList(avs.([]interface{}))
				affinity.AntiAffinityVolumes = afvols
			}
			if ais, ok := d.GetOk(PIAntiAffinityInstances); ok {
				afinss := flex.ExpandStringList(ais.([]interface{}))
				affinity.AntiAffinityPVMInstances = afinss
			}
		}
		body.StorageAffinity = affinity
	}

	if sc, ok := d.GetOk(helpers.PIInstanceStorageConnection); ok {
		body.StorageConnection = sc.(string)
	}

	if pg, ok := d.GetOk(helpers.PIPlacementGroupID); ok {
		body.PlacementGroup = pg.(string)
	}

	if spp, ok := d.GetOk(Arg_PIInstanceSharedProcessorPool); ok {
		body.SharedProcessorPool = spp.(string)
	}
	imageData, err := imageClient.GetStockImage(imageid)
	if err != nil {
		// check if vtl image is cloud instance image
		imageData, err = imageClient.Get(imageid)
		if err != nil {
			return nil, fmt.Errorf("image doesn't exist. %e", err)
		}
	}
	if lrc, ok := d.GetOk(helpers.PIInstanceLicenseRepositoryCapacity); ok {

		if imageData.Specifications.ImageType == "stock-vtl" {
			body.LicenseRepositoryCapacity = int64(lrc.(int))
		} else {
			return nil, fmt.Errorf("pi_license_repository_capacity should only be used when creating VTL instances. %e", err)
		}
	}

	if imageData.Specifications.OperatingSystem == OS_IBMI {
		// Default value
		falseBool := false
		sl := &models.SoftwareLicenses{
			IbmiCSS:      &falseBool,
			IbmiPHA:      &falseBool,
			IbmiRDS:      &falseBool,
			IbmiRDSUsers: 0,
		}
		if ibmiCSS, ok := d.GetOk(Arg_IBMiCSS); ok {
			sl.IbmiCSS = flex.PtrToBool(ibmiCSS.(bool))
		}
		if ibmiPHA, ok := d.GetOk(Arg_IBMiPHA); ok {
			sl.IbmiPHA = flex.PtrToBool(ibmiPHA.(bool))
		}
		if ibmrdsUsers, ok := d.GetOk(Arg_IBMiRDSUsers); ok {
			if ibmrdsUsers.(int) < 0 {
				return nil, fmt.Errorf("request with IBM i Rational Dev Studio property requires IBM i Rational Dev Studio number of users")
			}
			sl.IbmiRDS = flex.PtrToBool(ibmrdsUsers.(int) > 0)
			sl.IbmiRDSUsers = int64(ibmrdsUsers.(int))
		}
		body.SoftwareLicenses = sl
	}

	pvmList, err := client.Create(body)

	if err != nil {
		return nil, fmt.Errorf("failed to provision: %v", err)
	}
	if pvmList == nil {
		return nil, fmt.Errorf("failed to provision")
	}

	return pvmList, nil
}

func splitID(id string) (id1, id2 string, err error) {
	parts, err := flex.IdParts(id)
	if err != nil {
		return
	}
	id1 = parts[0]
	id2 = parts[1]
	return
}
