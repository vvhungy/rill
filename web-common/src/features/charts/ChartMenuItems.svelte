<script lang="ts">
  import { getFileAPIPathFromNameAndType } from "@rilldata/web-common/features/entity-management/entity-mappers";
  import NavigationMenuItem from "@rilldata/web-common/layout/navigation/NavigationMenuItem.svelte";
  import { runtime } from "../../runtime-client/runtime-store";
  import { deleteFileArtifact } from "../entity-management/actions";
  import { EntityType } from "../entity-management/types";
  import { useChartRoutes } from "./selectors";

  export let chartName: string;

  $: chartRoutes = useChartRoutes($runtime.instanceId);

  async function handleDeleteChart() {
    await deleteFileArtifact(
      $runtime.instanceId,
      getFileAPIPathFromNameAndType(chartName, EntityType.Chart),
      EntityType.Chart,
      $chartRoutes.data ?? [],
    );
  }
</script>

<NavigationMenuItem on:click={handleDeleteChart}>
  Delete chart
</NavigationMenuItem>
