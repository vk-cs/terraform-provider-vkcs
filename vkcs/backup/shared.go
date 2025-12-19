package backup

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/plans"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/providers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/triggers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/instances"
	"golang.org/x/exp/maps"
)

const (
	ProviderNameNova  = "cloud_servers"
	ProviderNameTrove = "dbaas"
	Nova              = "OS::Nova"
	Trove             = "OS::Trove"
	NovaInstance      = "OS::Nova::Server"
	TroveInstance     = "OS::Trove::Instance"
	TroveCluster      = "OS::Trove::Cluster"

	CinderVolume = "OS::Cinder::Volume"

	RetentionFull = "max_backups"
	RetentionGFS  = "gfs"

	TimeFormatWithoutZone = "15:04"
	TimeFormatWithZone    = "15:04-07"
)

var providerNameMapping = map[string]string{
	Nova:  ProviderNameNova,
	Trove: ProviderNameTrove,
}

func getProviderNames() []string {
	names := maps.Values(providerNameMapping)
	sort.Strings(names)
	return names
}

func getResourcesInfo(config clients.Config, region string, instancesID []types.String, resourceType string) ([]*plans.BackupPlanResource, error) {
	if resourceType == ProviderNameNova {
		return getNovaResourceInfo(config, region, instancesID)
	}
	if resourceType == ProviderNameTrove {
		return getTroveResourceInfo(config, region, instancesID)
	}

	return nil, fmt.Errorf("error getting resources info: unknown resource type")
}

func getNovaResourceInfo(config clients.Config, region string, instancesID []types.String) ([]*plans.BackupPlanResource, error) {
	computeClient, err := config.ComputeV2Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating VKCS compute client: %s", err.Error())
	}

	allPages, err := servers.List(computeClient, servers.ListOpts{}).AllPages()
	if err != nil {
		return nil, fmt.Errorf("error getting servers info: %s", err)
	}

	allServers, err := servers.ExtractServers(allPages)
	if err != nil {
		return nil, fmt.Errorf("error getting servers info: %s", err)
	}

	serversMap := make(map[string]servers.Server)
	for _, server := range allServers {
		serversMap[server.ID] = server
	}

	resourcesInfo := make([]*plans.BackupPlanResource, 0)
	missingResources := make([]string, 0)

	for _, instanceID := range instancesID {
		if serv, ok := serversMap[instanceID.ValueString()]; ok {
			resourceInfo := plans.BackupPlanResource{
				ID:   serv.ID,
				Type: NovaInstance,
				Name: serv.Name,
			}
			resourcesInfo = append(resourcesInfo, &resourceInfo)
		} else {
			missingResources = append(missingResources, instanceID.ValueString())
		}
	}

	if len(missingResources) > 0 {
		return nil, fmt.Errorf("error getting resources info: could not find resources: %s", strings.Join(missingResources, ", "))
	}

	return resourcesInfo, nil
}

func getTroveResourceInfo(config clients.Config, region string, instancesID []types.String) ([]*plans.BackupPlanResource, error) {
	dbClient, err := config.DatabaseV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating VKCS database client: %s", err.Error())
	}

	allInstancesPages, err := instances.List(dbClient).AllPages()
	if err != nil {
		return nil, fmt.Errorf("error getting database instances info: %s", err)
	}

	allInstances, err := instances.ExtractInstances(allInstancesPages)
	if err != nil {
		return nil, fmt.Errorf("error getting database instances info: %s", err)
	}

	allClustersPages, err := clusters.List(dbClient).AllPages()
	if err != nil {
		return nil, fmt.Errorf("error getting database clusters info: %s", err)
	}

	allClusters, err := clusters.ExtractClusters(allClustersPages)
	if err != nil {
		return nil, fmt.Errorf("error getting database clusters info: %s", err)
	}

	instancesMap := make(map[string]instances.InstanceResp, len(allInstances))
	for _, inst := range allInstances {
		instancesMap[inst.ID] = inst
	}

	clustersMap := make(map[string]clusters.ClusterResp)
	for _, cluster := range allClusters {
		clustersMap[cluster.ID] = cluster
	}

	resourcesInfo := make([]*plans.BackupPlanResource, 0)
	missingResources := make([]string, 0)

	for _, instanceID := range instancesID {
		id := instanceID.ValueString()
		if inst, ok := instancesMap[id]; ok {
			resourcesInfo = append(resourcesInfo, &plans.BackupPlanResource{
				ID:   inst.ID,
				Type: TroveInstance,
				Name: inst.Name,
			})
		} else if cl, ok := clustersMap[id]; ok {
			resourcesInfo = append(resourcesInfo, &plans.BackupPlanResource{
				ID:   cl.ID,
				Type: TroveCluster,
				Name: cl.Name,
			})
		} else {
			missingResources = append(missingResources, id)
		}
	}

	if len(missingResources) > 0 {
		return nil, fmt.Errorf("error getting resources info: could not find resources: %s", strings.Join(missingResources, ", "))
	}

	return resourcesInfo, nil
}

func enrichWithVolumes(ctx context.Context, resources []*plans.BackupPlanResource, backupTargets []PlanResourceBackupTargetModel) ([]*plans.BackupPlanResource, diag.Diagnostics) {
	var diags diag.Diagnostics

	resourceByID := make(map[string]*plans.BackupPlanResource, len(resources))
	for _, r := range resources {
		resourceByID[r.ID] = r
	}

	for _, backupTarget := range backupTargets {
		if backupTarget.VolumeIDs.IsNull() || backupTarget.VolumeIDs.IsUnknown() {
			continue
		}

		res, ok := resourceByID[backupTarget.InstanceID.ValueString()]
		if !ok {
			diags.AddError("Invalid backup_targets", fmt.Sprintf("Instance %s not found in plan resources", backupTarget.InstanceID.ValueString()))
			return nil, diags
		}

		var volumeIDs []string
		d := backupTarget.VolumeIDs.ElementsAs(ctx, &volumeIDs, false)
		diags.Append(d...)
		if d.HasError() {
			return nil, diags
		}

		for _, volumeID := range volumeIDs {
			res.Resources = append(res.Resources, &plans.BackupPlanResource{
				ID:   volumeID,
				Type: CinderVolume,
			})
		}
	}

	return resources, diags
}

func dayToNumber(day string) int {
	days := map[string]int{
		"Mo": 0,
		"Tu": 1,
		"We": 2,
		"Th": 3,
		"Fr": 4,
		"Sa": 5,
		"Su": 6,
	}
	return days[day]
}

func numberToDay(number int) string {
	days := []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}
	return days[number]
}

func expandIncrementalFullDay(plan PlanResourceModel) (int, error) {
	scheduleDates := plan.Schedule.Date
	if len(scheduleDates) > 1 {
		return 0, fmt.Errorf("invalid resource schema. There should be only one date for incremental_backups")
	}
	fullDay := dayToNumber(scheduleDates[0].ValueString())
	return fullDay, nil
}

func expandGFS(plan PlanResourceModel) *plans.BackupPlanGFS {
	gfs := &plans.BackupPlanGFS{
		Son: int(plan.GFSRetention.GFSWeekly.ValueInt64()),
	}
	if !plan.GFSRetention.GFSMonthly.IsNull() {
		gfs.Father = int(plan.GFSRetention.GFSMonthly.ValueInt64())
	}
	if !plan.GFSRetention.GFSYearly.IsNull() {
		gfs.Grandfather = int(plan.GFSRetention.GFSYearly.ValueInt64())
	}
	return gfs
}

func expandTriggerSchedule(plan PlanResourceModel) (string, error) {
	var triggerSchedule string
	if !plan.Schedule.EveryHours.IsNull() {
		everyHours := plan.Schedule.EveryHours.ValueInt64()
		triggerSchedule = fmt.Sprintf("0 */%d * * *", everyHours)
	} else {
		incrementalBackups := plan.IncrementalBackup.ValueBool()
		parsedTime, err := parseTime(plan.Schedule.Time.ValueString())
		if err != nil {
			return "", fmt.Errorf("invalid resource schema. Invalid time: %s", err)
		}
		planTime := parsedTime.In(time.UTC)
		if incrementalBackups {
			triggerSchedule = fmt.Sprintf("%d %d * * *", planTime.Minute(), planTime.Hour())
		} else {
			scheduleDates := plan.Schedule.Date
			days := make([]string, 0)
			for _, date := range scheduleDates {
				day := dayToNumber(date.ValueString())
				days = append(days, strconv.Itoa(day))
			}
			triggerSchedule = fmt.Sprintf("%d %d * * %s", planTime.Minute(), planTime.Hour(), strings.Join(days, ","))
		}
	}
	return triggerSchedule, nil
}

func parseTime(value string) (*time.Time, error) {
	var parsedTime time.Time
	parsedTime, err := time.Parse(TimeFormatWithZone, value)
	if err != nil {
		parsedTime, err = time.Parse(TimeFormatWithoutZone, value)
		if err != nil {
			return nil, err
		}
	}
	return &parsedTime, nil
}

func flattenGFS(planResp *plans.PlanResponse) *PlanResourceGFSRetentionModel {
	gfsRetention := PlanResourceGFSRetentionModel{
		GFSWeekly: types.Int64Value(int64(planResp.GFS.Son)),
	}
	if planResp.GFS.Father != 0 {
		gfsRetention.GFSMonthly = types.Int64Value(int64(planResp.GFS.Father))
	}
	if planResp.GFS.Grandfather != 0 {
		gfsRetention.GFSYearly = types.Int64Value(int64(planResp.GFS.Grandfather))
	}
	return &gfsRetention
}

func flattenSchedule(planResp *plans.PlanResponse, triggerResp *triggers.TriggerResponse, location *time.Location) *PlanResourceScheduleModel {
	planSchedule := PlanResourceScheduleModel{}
	scheduleParts := strings.Split(triggerResp.Properties.Pattern, " ")
	if strings.HasPrefix(scheduleParts[1], "*/") {
		planSchedule.Date = nil
		planSchedule.Time = types.StringNull()

		everyHoursParts := strings.Split(scheduleParts[1], "/")
		everyHours, _ := strconv.Atoi(everyHoursParts[1])
		planSchedule.EveryHours = types.Int64Value(int64(everyHours))
	} else {
		timeString := fmt.Sprintf("%02s:%02s", scheduleParts[1], scheduleParts[0])
		if location != time.UTC {
			timeParsed, _ := time.Parse(TimeFormatWithoutZone, timeString)
			timeWithLocation := timeParsed.In(location)
			timeFormatted := timeWithLocation.Format(TimeFormatWithZone)
			planSchedule.Time = types.StringValue(timeFormatted)
		} else {
			planSchedule.Time = types.StringValue(timeString)
		}

		if planResp.FullDay == nil {
			days := make([]types.String, 0)
			daysParts := strings.Split(scheduleParts[4], ",")
			for _, dayStr := range daysParts {
				dayNum, _ := strconv.Atoi(dayStr)
				days = append(days, types.StringValue(numberToDay(dayNum)))
			}
			planSchedule.Date = days
		} else {
			day := numberToDay(*planResp.FullDay)
			planSchedule.Date = []types.String{types.StringValue(day)}
		}
		planSchedule.EveryHours = types.Int64Null()
	}
	return &planSchedule
}

func findProvider(client *gophercloud.ServiceClient, providerID string, providerName string) (*providers.Provider, error) {
	allProviders, err := providers.List(client).Extract()
	if err != nil {
		return nil, fmt.Errorf("error retrieving backup providers")
	}
	var foundProvider *providers.Provider
	for _, provider := range allProviders {
		if providerID != "" && provider.ID != providerID {
			continue
		}
		if providerName != "" && providerNameMapping[provider.Name] != providerName {
			continue
		}
		foundProvider = provider
		foundProvider.Name = providerNameMapping[provider.Name]
		break
	}
	if foundProvider == nil {
		return nil, fmt.Errorf("error retrieving backup provider: no suitable providers found")
	}
	return foundProvider, nil
}

func findTrigger(client *gophercloud.ServiceClient, planID string) (*triggers.TriggerResponse, error) {
	allPages, err := triggers.List(client).AllPages()
	if err != nil {
		return nil, fmt.Errorf("error getting backup triggers info: %s", err.Error())
	}
	allTriggers, err := triggers.ExtractTriggers(allPages)
	if err != nil {
		return nil, fmt.Errorf("error getting backup triggers info: %s", err.Error())
	}
	var triggerResp *triggers.TriggerResponse
	for _, tr := range allTriggers {
		if tr.PlanID == planID {
			triggerResp = &tr
			break
		}
	}
	if triggerResp == nil {
		return nil, fmt.Errorf("backup trigger not found for plan %s", planID)
	}
	return triggerResp, nil
}
