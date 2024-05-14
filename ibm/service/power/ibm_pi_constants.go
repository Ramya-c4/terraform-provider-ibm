package power

import "time"

const (
	// Arguments
	Arg_AffinityInstance                    = "pi_affinity_instance"
	Arg_AffinityPolicy                      = "pi_affinity_policy"
	Arg_AffinityVolume                      = "pi_affinity_volume"
	Arg_AntiAffinityInstances               = "pi_anti_affinity_instances"
	Arg_AntiAffinityVolumes                 = "pi_anti_affinity_volumes"
	Arg_CloudConnectionName                 = "pi_cloud_connection_name"
	Arg_CloudInstanceID                     = "pi_cloud_instance_id"
	Arg_DatacenterZone                      = "pi_datacenter_zone"
	Arg_Description                         = "pi_description"
	Arg_DhcpCidr                            = "pi_cidr"
	Arg_DhcpCloudConnectionID               = "pi_cloud_connection_id"
	Arg_DhcpDnsServer                       = "pi_dns_server"
	Arg_DhcpID                              = "pi_dhcp_id"
	Arg_DhcpName                            = "pi_dhcp_name"
	Arg_DhcpSnatEnabled                     = "pi_dhcp_snat_enabled"
	Arg_IBMiCSS                             = "pi_ibmi_css"
	Arg_IBMiPHA                             = "pi_ibmi_pha"
	Arg_IBMiRDSUsers                        = "pi_ibmi_rds_users"
	Arg_ImageName                           = "pi_image_name"
	Arg_InstanceName                        = "pi_instance_name"
	Arg_KeyName                             = "pi_key_name"
	Arg_LanguageCode                        = "pi_language_code"
	Arg_NetworkName                         = "pi_network_name"
	Arg_PIInstanceSharedProcessorPool       = "pi_shared_processor_pool"
	Arg_PlacementGroupName                  = "pi_placement_group_name"
	Arg_PlacementGroupPolicy                = "pi_placement_group_policy"
	Arg_PVMInstanceActionType               = "pi_action"
	Arg_PVMInstanceHealthStatus             = "pi_health_status"
	Arg_PVMInstanceId                       = "pi_instance_id"
	Arg_ReplicationEnabled                  = "pi_replication_enabled"
	Arg_SAP                                 = "sap"
	Arg_SAPProfileID                        = "pi_sap_profile_id"
	Arg_SharedProcessorPoolHostGroup        = "pi_shared_processor_pool_host_group"
	Arg_SharedProcessorPoolID               = "pi_shared_processor_pool_id"
	Arg_SharedProcessorPoolName             = "pi_shared_processor_pool_name"
	Arg_SharedProcessorPoolPlacementGroupID = "pi_shared_processor_pool_placement_group_id"
	Arg_SharedProcessorPoolReservedCores    = "pi_shared_processor_pool_reserved_cores"
	Arg_SnapshotID                          = "pi_snapshot_id"
	Arg_SnapShotName                        = "pi_snap_shot_name"
	Arg_SPPPlacementGroupID                 = "pi_spp_placement_group_id"
	Arg_SPPPlacementGroupName               = "pi_spp_placement_group_name"
	Arg_SPPPlacementGroupPolicy             = "pi_spp_placement_group_policy"
	Arg_SSHKey                              = "pi_ssh_key"
	Arg_StoragePool                         = "pi_storage_pool"
	Arg_StorageType                         = "pi_storage_type"
	Arg_VolumeGroupID                       = "pi_volume_group_id"
	Arg_VolumeID                            = "pi_volume_id"
	Arg_VolumeIDs                           = "pi_volume_ids"
	Arg_VolumeName                          = "pi_volume_name"
	Arg_VolumeOnboardingID                  = "pi_volume_onboarding_id"
	Arg_VolumePool                          = "pi_volume_pool"
	Arg_VolumeShareable                     = "pi_volume_shareable"
	Arg_VolumeSize                          = "pi_volume_size"
	Arg_VolumeType                          = "pi_volume_type"
	Arg_VTL                                 = "vtl"

	// Attributes
	Attr_AccessConfig                                = "access_config"
	Attr_Action                                      = "action"
	Attr_Addresses                                   = "addresses"
	Attr_AllocatedCores                              = "allocated_cores"
	Attr_Architecture                                = "architecture"
	Attr_Auxiliary                                   = "auxiliary"
	Attr_AuxiliaryChangedVolumeName                  = "auxiliary_changed_volume_name"
	Attr_AuxiliaryVolumeName                         = "auxiliary_volume_name"
	Attr_AvailabilityZone                            = "availability_zone"
	Attr_AvailableCores                              = "available_cores"
	Attr_AvailableIPCount                            = "available_ip_count"
	Attr_Bootable                                    = "bootable"
	Attr_BootVolumeID                                = "boot_volume_id"
	Attr_Capabilities                                = "capabilities"
	Attr_Capacity                                    = "capacity"
	Attr_Certified                                   = "certified"
	Attr_CIDR                                        = "cidr"
	Attr_ClassicEnabled                              = "classic_enabled"
	Attr_CloudConnectionID                           = "cloud_connection_id"
	Attr_CloudInstanceID                             = "cloud_instance_id"
	Attr_CloudInstances                              = "cloud_instances"
	Attr_Code                                        = "code"
	Attr_ConnectionMode                              = "connection_mode"
	Attr_Connections                                 = "connections"
	Attr_ConsistencyGroupName                        = "consistency_group_name"
	Attr_ConsoleLanguages                            = "console_languages"
	Attr_ContainerFormat                             = "container_format"
	Attr_CopyRate                                    = "copy_rate"
	Attr_CopyType                                    = "copy_type"
	Attr_CoreMemoryRatio                             = "core_memory_ratio"
	Attr_Cores                                       = "cores"
	Attr_CPUs                                        = "cpus"
	Attr_CreateTime                                  = "create_time"
	Attr_CreationDate                                = "creation_date"
	Attr_CRN                                         = "crn"
	Attr_CyclePeriodSeconds                          = "cycle_period_seconds"
	Attr_CyclingMode                                 = "cycling_mode"
	Attr_DatacenterCapabilities                      = "pi_datacenter_capabilities"
	Attr_DatacenterHref                              = "pi_datacenter_href"
	Attr_DatacenterLocation                          = "pi_datacenter_location"
	Attr_Datacenters                                 = "datacenters"
	Attr_DatacenterStatus                            = "pi_datacenter_status"
	Attr_DatacenterType                              = "pi_datacenter_type"
	Attr_Default                                     = "default"
	Attr_DeleteOnTermination                         = "delete_on_termination"
	Attr_DeploymentType                              = "deployment_type"
	Attr_Description                                 = "description"
	Attr_DhcpID                                      = "dhcp_id"
	Attr_DhcpLeaseInstanceIP                         = "instance_ip"
	Attr_DhcpLeaseInstanceMac                        = "instance_mac"
	Attr_DhcpLeases                                  = "leases"
	Attr_DhcpManaged                                 = "dhcp_managed"
	Attr_DhcpNetworkDeprecated                       = "network" // to deprecate
	Attr_DhcpNetworkID                               = "network_id"
	Attr_DhcpNetworkName                             = "network_name"
	Attr_DhcpServers                                 = "servers"
	Attr_DhcpStatus                                  = "status"
	Attr_DisasterRecoveryLocations                   = "disaster_recovery_locations"
	Attr_DiskFormat                                  = "disk_format"
	Attr_DiskType                                    = "disk_type"
	Attr_DNS                                         = "dns"
	Attr_Enabled                                     = "enabled"
	Attr_Endianness                                  = "endianness"
	Attr_ExternalIP                                  = "external_ip"
	Attr_FailureMessage                              = "failure_message"
	Attr_FlashCopyMappings                           = "flash_copy_mappings"
	Attr_FlashCopyName                               = "flash_copy_name"
	Attr_FreezeTime                                  = "freeze_time"
	Attr_Gateway                                     = "gateway"
	Attr_GlobalRouting                               = "global_routing"
	Attr_GreDestinationAddress                       = "gre_destination_address"
	Attr_GreSourceAddress                            = "gre_source_address"
	Attr_GroupID                                     = "group_id"
	Attr_HealthStatus                                = "health_status"
	Attr_HostID                                      = "host_id"
	Attr_Href                                        = "href"
	Attr_Hypervisor                                  = "hypervisor"
	Attr_HypervisorType                              = "hypervisor_type"
	Attr_IBMiCSS                                     = "ibmi_css"
	Attr_IBMIPAddress                                = "ibm_ip_address"
	Attr_IBMiPHA                                     = "ibmi_pha"
	Attr_IBMiRDS                                     = "ibmi_rds"
	Attr_IBMiRDSUsers                                = "ibmi_rds_users"
	Attr_ID                                          = "id"
	Attr_ImageID                                     = "image_id"
	Attr_ImageInfo                                   = "image_info"
	Attr_Images                                      = "images"
	Attr_ImageType                                   = "image_type"
	Attr_InputVolumes                                = "input_volumes"
	Attr_Instances                                   = "instances"
	Attr_InstanceSnapshots                           = "instance_snapshots"
	Attr_InstanceVolumes                             = "instance_volumes"
	Attr_IOThrottleRate                              = "io_throttle_rate"
	Attr_IP                                          = "ip"
	Attr_IPAddress                                   = "ipaddress"
	Attr_IPOctet                                     = "ipoctet"
	Attr_IsActive                                    = "is_active"
	Attr_Jumbo                                       = "jumbo"
	Attr_Key                                         = "key"
	Attr_KeyCreationDate                             = "creation_date"
	Attr_KeyID                                       = "key_id"
	Attr_KeyName                                     = "name"
	Attr_Keys                                        = "keys"
	Attr_Language                                    = "language"
	Attr_LastUpdateDate                              = "last_update_date"
	Attr_LastUpdatedDate                             = "last_updated_date"
	Attr_Leases                                      = "leases"
	Attr_LicenseRepositoryCapacity                   = "license_repository_capacity"
	Attr_Location                                    = "location"
	Attr_MacAddress                                  = "macaddress"
	Attr_MasterChangedVolumeName                     = "master_changed_volume_name"
	Attr_MasterVolumeName                            = "master_volume_name"
	Attr_Max                                         = "max"
	Attr_MaxAllocationSize                           = "max_allocation_size"
	Attr_MaxAvailable                                = "max_available"
	Attr_MaxCoresAvailable                           = "max_cores_available"
	Attr_MaximumStorageAllocation                    = "max_storage_allocation"
	Attr_MaxMem                                      = "maxmem"
	Attr_MaxMemoryAvailable                          = "max_memory_available"
	Attr_MaxProc                                     = "maxproc"
	Attr_MaxVirtualCores                             = "max_virtual_cores"
	Attr_Members                                     = "members"
	Attr_Memory                                      = "memory"
	Attr_Message                                     = "message"
	Attr_Metered                                     = "metered"
	Attr_MigrationStatus                             = "migration_status"
	Attr_Min                                         = "min"
	Attr_MinMem                                      = "minmem"
	Attr_MinProc                                     = "minproc"
	Attr_MinVirtualCores                             = "min_virtual_cores"
	Attr_MirroringState                              = "mirroring_state"
	Attr_MTU                                         = "mtu"
	Attr_Name                                        = "name"
	Attr_NetworkID                                   = "network_id"
	Attr_NetworkName                                 = "network_name"
	Attr_NetworkPorts                                = "network_ports"
	Attr_Networks                                    = "networks"
	Attr_NumberOfVolumes                             = "number_of_volumes"
	Attr_Onboardings                                 = "onboardings"
	Attr_OperatingSystem                             = "operating_system"
	Attr_PercentComplete                             = "percent_complete"
	Attr_PIInstanceSharedProcessorPool               = "shared_processor_pool"
	Attr_PIInstanceSharedProcessorPoolID             = "shared_processor_pool_id"
	Attr_PinPolicy                                   = "pin_policy"
	Attr_PlacementGroupID                            = "placement_group_id"
	Attr_PlacementGroups                             = "placement_groups"
	Attr_Policy                                      = "policy"
	Attr_Pool                                        = "pool"
	Attr_PoolName                                    = "pool_name"
	Attr_Port                                        = "port"
	Attr_PortID                                      = "portid"
	Attr_PowerEdgeRouter                             = "power_edge_router"
	Attr_PrimaryRole                                 = "primary_role"
	Attr_Processors                                  = "processors"
	Attr_ProcType                                    = "proctype"
	Attr_ProfileID                                   = "profile_id"
	Attr_Profiles                                    = "profiles"
	Attr_Progress                                    = "progress"
	Attr_PublicIP                                    = "public_ip"
	Attr_PVMInstanceID                               = "pvm_instance_id"
	Attr_PVMInstances                                = "pvm_instances"
	Attr_PVMSnapshots                                = "pvm_snapshots"
	Attr_Region                                      = "region"
	Attr_RemoteCopyID                                = "remote_copy_id"
	Attr_RemoteCopyRelationshipNames                 = "remote_copy_relationship_names"
	Attr_RemoteCopyRelationships                     = "remote_copy_relationships"
	Attr_ReplicationEnabled                          = "replication_enabled"
	Attr_ReplicationSites                            = "replication_sites"
	Attr_ReplicationStatus                           = "replication_status"
	Attr_ReplicationType                             = "replication_type"
	Attr_ReservedCores                               = "reserved_cores"
	Attr_ResultsOnboardedVolumes                     = "results_onboarded_volumes"
	Attr_ResultsVolumeOnboardingFailures             = "results_volume_onboarding_failures"
	Attr_ServerName                                  = "server_name"
	Attr_Shareable                                   = "shreable"
	Attr_SharedCoreRatio                             = "shared_core_ratio"
	Attr_SharedProcessorPool                         = "shared_processor_pool"
	Attr_SharedProcessorPoolAllocatedCores           = "allocated_cores"
	Attr_SharedProcessorPoolAvailableCores           = "available_cores"
	Attr_SharedProcessorPoolHostID                   = "host_id"
	Attr_SharedProcessorPoolID                       = "shared_processor_pool_id"
	Attr_SharedProcessorPoolInstanceAvailabilityZone = "availability_zone"
	Attr_SharedProcessorPoolInstanceCpus             = "cpus"
	Attr_SharedProcessorPoolInstanceId               = "id"
	Attr_SharedProcessorPoolInstanceMemory           = "memory"
	Attr_SharedProcessorPoolInstanceName             = "name"
	Attr_SharedProcessorPoolInstances                = "instances"
	Attr_SharedProcessorPoolInstanceStatus           = "status"
	Attr_SharedProcessorPoolInstanceUncapped         = "uncapped"
	Attr_SharedProcessorPoolInstanceVcpus            = "vcpus"
	Attr_SharedProcessorPoolName                     = "name"
	Attr_SharedProcessorPoolPlacementGroups          = "spp_placement_groups"
	Attr_SharedProcessorPoolReservedCores            = "reserved_cores"
	Attr_SharedProcessorPools                        = "shared_processor_pools"
	Attr_SharedProcessorPoolStatus                   = "status"
	Attr_SharedProcessorPoolStatusDetail             = "status_detail"
	Attr_Size                                        = "size"
	Attr_SnapshotID                                  = "snapshot_id"
	Attr_SourceVolumeName                            = "source_volume_name"
	Attr_Speed                                       = "speed"
	Attr_SPPPlacementGroupID                         = "spp_placement_group_id"
	Attr_SPPPlacementGroupMembers                    = "members"
	Attr_SPPPlacementGroupName                       = "name"
	Attr_SPPPlacementGroupPolicy                     = "policy"
	Attr_SPPPlacementGroups                          = "spp_placement_groups"
	Attr_SSHKey                                      = "ssh_key"
	Attr_StartTime                                   = "start_time"
	Attr_State                                       = "state"
	Attr_Status                                      = "status"
	Attr_StatusDescriptionErrors                     = "status_description_errors"
	Attr_StatusDetail                                = "status_detail"
	Attr_StoragePool                                 = "storage_pool"
	Attr_StoragePoolAffinity                         = "storage_pool_affinity"
	Attr_StoragePoolsCapacity                        = "storage_pools_capacity"
	Attr_StorageType                                 = "storage_type"
	Attr_StorageTypesCapacity                        = "storage_types_capacity"
	Attr_Synchronized                                = "synchronized"
	Attr_SystemPoolName                              = "system_pool_name"
	Attr_SystemPools                                 = "system_pools"
	Attr_Systems                                     = "systems"
	Attr_SysType                                     = "systype"
	Attr_TargetVolumeName                            = "target_volume_name"
	Attr_TenantID                                    = "tenant_id"
	Attr_TenantName                                  = "tenant_name"
	Attr_TotalCapacity                               = "total_capacity"
	Attr_TotalInstances                              = "total_instances"
	Attr_TotalMemoryConsumed                         = "total_memory_consumed"
	Attr_TotalProcessorsConsumed                     = "total_processors_consumed"
	Attr_TotalSSDStorageConsumed                     = "total_ssd_storage_consumed"
	Attr_TotalStandardStorageConsumed                = "total_standard_storage_consumed"
	Attr_Type                                        = "type"
	Attr_Uncapped                                    = "uncapped"
	Attr_URL                                         = "url"
	Attr_UsedIPCount                                 = "used_ip_count"
	Attr_UsedIPPercent                               = "used_ip_percent"
	Attr_UserIPAddress                               = "user_ip_address"
	Attr_VCPUs                                       = "vcpus"
	Attr_VirtualCoresAssigned                        = "virtual_cores_assigned"
	Attr_VLanID                                      = "vlan_id"
	Attr_VolumeGroupName                             = "volume_group_name"
	Attr_VolumeGroups                                = "volume_groups"
	Attr_VolumeID                                    = "volume_id"
	Attr_VolumeIDs                                   = "volume_ids"
	Attr_VolumePool                                  = "volume_pool"
	Attr_Volumes                                     = "volumes"
	Attr_VolumeSnapshots                             = "volume_snapshots"
	Attr_VolumeStatus                                = "volume_status"
	Attr_VPCCRNs                                     = "vpc_crns"
	Attr_VPCEnabled                                  = "vpc_enabled"
	Attr_WorkspaceCapabilities                       = "pi_workspace_capabilities"
	Attr_WorkspaceDetails                            = "pi_workspace_details"
	Attr_WorkspaceID                                 = "pi_workspace_id"
	Attr_WorkspaceLocation                           = "pi_workspace_location"
	Attr_WorkspaceName                               = "pi_workspace_name"
	Attr_Workspaces                                  = "workspaces"
	Attr_WorkspaceStatus                             = "pi_workspace_status"
	Attr_WorkspaceType                               = "pi_workspace_type"
	Attr_WWN                                         = "wwn"
	OS_IBMI                                          = "ibmi"

	// Affinty Values
	Affinity     = "affinity"
	AntiAffinity = "anti-affinity"

	// States
	State_Active              = "active"
	State_ACTIVE              = "ACTIVE"
	State_Added               = "added"
	State_Adding              = "adding"
	State_Available           = "available"
	State_BUILD               = "BUILD"
	State_Creating            = "creating"
	State_Deleted             = "deleted"
	State_Deleting            = "deleting"
	State_DELETING            = "DELETING"
	State_Failed              = "failed"
	State_Inactive            = "inactive"
	State_InProgress          = "in progress"
	State_InUse               = "in-use"
	State_NotFound            = "Not Found"
	State_PendingReclaimation = "pending_reclamation"
	State_Provisioning        = "provisioning"
	State_Removed             = "removed"
	State_Retry               = "retry"

	// Health
	Health_OK = "OK"

	// TODO: Second Half Cleanup, remove extra variables

	// SAP Profile
	PISAPProfiles         = "profiles"
	PISAPProfileCertified = "certified"
	PISAPProfileCores     = "cores"
	PISAPProfileMemory    = "memory"
	PISAPProfileID        = "profile_id"
	PISAPProfileType      = "type"

	PVMInstanceHealthOk      = "OK"
	PVMInstanceHealthWarning = "WARNING"

	//Added timeout values for warning  and active status
	warningTimeOut = 60 * time.Second
	activeTimeOut  = 2 * time.Minute
	// power service instance capabilities
	CUSTOM_VIRTUAL_CORES = "custom-virtualcores"

	PIConsoleLanguageCode               = "pi_language_code"
	PICloudConnectionId                 = "cloud_connection_id"
	PICloudConnectionStatus             = "status"
	PICloudConnectionIBMIPAddress       = "ibm_ip_address"
	PICloudConnectionUserIPAddress      = "user_ip_address"
	PICloudConnectionPort               = "port"
	PICloudConnectionClassicGreSource   = "gre_source_address"
	PICloudConnectionConnectionMode     = "connection_mode"
	PIInstanceDeploymentType            = "pi_deployment_type"
	PIInstanceMigratable                = "pi_migratable"
	PIInstanceNetwork                   = "pi_network"
	PIInstanceLicenseRepositoryCapacity = "pi_license_repository_capacity"
	PIInstanceStoragePool               = "pi_storage_pool"
	PIInstanceStorageType               = "pi_storage_type"
	PISAPInstanceProfileID              = "pi_sap_profile_id"
	PISAPInstanceDeploymentType         = "pi_sap_deployment_type"
	PIInstanceSharedProcessorPool       = "pi_shared_processor_pool"
	PIInstanceStorageConnection         = "pi_storage_connection"
	PIInstanceStoragePoolAffinity       = "pi_storage_pool_affinity"

	PIInstanceUserData  = "pi_user_data"
	PIInstanceVolumeIds = "pi_volume_ids"

	// Placement Group
	PIPlacementGroupID      = "placement_group_id"
	PIPlacementGroupMembers = "members"

	// Volume
	PIVolumeIds             = "pi_volume_ids"
	PIAffinityPolicy        = "pi_affinity_policy"
	PIAffinityVolume        = "pi_affinity_volume"
	PIAffinityInstance      = "pi_affinity_instance"
	PIAntiAffinityInstances = "pi_anti_affinity_instances"
	PIAntiAffinityVolumes   = "pi_anti_affinity_volumes"

	// Volume Clone
	PIVolumeCloneName   = "pi_volume_clone_name"
	PIVolumeCloneTaskID = "pi_volume_clone_task_id"
	PITargetStorageTier = "pi_target_storage_tier"

	// IBM PI Volume Group
	PIVolumeGroupName                 = "pi_volume_group_name"
	PIVolumeGroupConsistencyGroupName = "pi_consistency_group_name"
	PIVolumeGroupID                   = "pi_volume_group_id"
	PIVolumeGroupAction               = "pi_volume_group_action"
	PIVolumeOnboardingID              = "pi_volume_onboarding_id"

	// Disaster Recovery Location
	PIDRLocation = "location"

	// VPN
	PIVPNConnectionId                         = "connection_id"
	PIVPNConnectionStatus                     = "connection_status"
	PIVPNConnectionDeadPeerDetection          = "dead_peer_detections"
	PIVPNConnectionDeadPeerDetectionAction    = "action"
	PIVPNConnectionDeadPeerDetectionInterval  = "interval"
	PIVPNConnectionDeadPeerDetectionThreshold = "threshold"
	PIVPNConnectionLocalGatewayAddress        = "local_gateway_address"
	PIVPNConnectionVpnGatewayAddress          = "gateway_address"

	// Cloud Connections
	PICloudConnectionTransitEnabled = "pi_cloud_connection_transit_enabled"

	// status
	// common status states
	StatusShutoff = "SHUTOFF"
	StatusActive  = "ACTIVE"
	StatusResize  = "RESIZE"
	StatusError   = "ERROR"
	StatusBuild   = "BUILD"
	StatusPending = "PENDING"
	SctionStart   = "start"
	SctionStop    = "stop"

	// volume clone task status
	VolumeCloneCompleted = "completed"
	VolumeCloneRunning   = "running"

	// IBM PI Workspace
	PIWorkspaceName          = "pi_name"
	PIWorkspaceDatacenter    = "pi_datacenter"
	PIWorkspaceResourceGroup = "pi_resource_group_id"
	PIWorkspacePlan          = "pi_plan"
	PIVirtualOpticalDevice   = "pi_virtual_optical_device"
)
