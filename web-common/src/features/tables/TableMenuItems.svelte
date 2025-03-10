<script lang="ts">
  import Explore from "@rilldata/web-common/components/icons/Explore.svelte";
  import Model from "@rilldata/web-common/components/icons/Model.svelte";
  import { useModelFileNames } from "@rilldata/web-common/features/models/selectors";
  import { appScreen } from "@rilldata/web-common/layout/app-store";
  import { behaviourEvent } from "@rilldata/web-common/metrics/initMetrics";
  import { BehaviourEventMedium } from "@rilldata/web-common/metrics/service/BehaviourEventTypes";
  import {
    MetricsEventScreenName,
    MetricsEventSpace,
  } from "@rilldata/web-common/metrics/service/MetricsTypes";
  import { WandIcon } from "lucide-svelte";
  import { runtime } from "../../runtime-client/runtime-store";
  import { useCreateDashboardFromTableUIAction } from "../metrics-views/ai-generation/generateMetricsView";
  import { createModelFromSource } from "../sources/createModel";
  import { useIsModelingSupportedForCurrentOlapDriver } from "./selectors";
  import NavigationMenuItem from "@rilldata/web-common/layout/navigation/NavigationMenuItem.svelte";

  export let fullyQualifiedTableName: string;

  $: isModelingSupportedForCurrentOlapDriver =
    useIsModelingSupportedForCurrentOlapDriver($runtime.instanceId);
  $: tableName = fullyQualifiedTableName.split(".")[1];
  $: runtimeInstanceId = $runtime.instanceId;
  $: modelNames = useModelFileNames($runtime.instanceId);

  const handleCreateModel = async () => {
    try {
      const previousActiveEntity = $appScreen?.type;
      const newModelName = await createModelFromSource(
        runtimeInstanceId,
        $modelNames.data ?? [],
        tableName,
        tableName,
      );

      await behaviourEvent.fireNavigationEvent(
        newModelName,
        BehaviourEventMedium.Menu,
        MetricsEventSpace.LeftPanel,
        previousActiveEntity,
        MetricsEventScreenName.Model,
      );
    } catch (err) {
      console.error(err);
    }
  };

  $: createDashboardFromTable = useCreateDashboardFromTableUIAction(
    $runtime.instanceId,
    tableName,
    BehaviourEventMedium.Menu,
    MetricsEventSpace.LeftPanel,
  );
</script>

{#if $isModelingSupportedForCurrentOlapDriver.data}
  <NavigationMenuItem on:click={() => handleCreateModel()}>
    <Model slot="icon" />
    Create new model
  </NavigationMenuItem>
{/if}

<NavigationMenuItem on:click={createDashboardFromTable}>
  <Explore slot="icon" />
  <div class="flex gap-x-2 items-center">
    Generate dashboard with AI
    <WandIcon class="w-3 h-3" />
  </div>
</NavigationMenuItem>
