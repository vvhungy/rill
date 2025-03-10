<script lang="ts">
  import { Switch } from "@rilldata/web-common/components/button";
  import Button from "@rilldata/web-common/components/button/Button.svelte";
  import Close from "@rilldata/web-common/components/icons/Close.svelte";
  import Column from "@rilldata/web-common/components/icons/Column.svelte";
  import Row from "@rilldata/web-common/components/icons/Row.svelte";
  import SearchIcon from "@rilldata/web-common/components/icons/Search.svelte";
  import { Search } from "@rilldata/web-common/components/search";
  import SearchableFilterChip from "@rilldata/web-common/components/searchable-filter-menu/SearchableFilterChip.svelte";
  import Shortcut from "@rilldata/web-common/components/tooltip/Shortcut.svelte";
  import Tooltip from "@rilldata/web-common/components/tooltip/Tooltip.svelte";
  import TooltipContent from "@rilldata/web-common/components/tooltip/TooltipContent.svelte";
  import TooltipShortcutContainer from "@rilldata/web-common/components/tooltip/TooltipShortcutContainer.svelte";
  import TooltipTitle from "@rilldata/web-common/components/tooltip/TooltipTitle.svelte";
  import SelectAllButton from "@rilldata/web-common/features/dashboards/dimension-table/SelectAllButton.svelte";
  import ReplacePivotDialog from "@rilldata/web-common/features/dashboards/pivot/ReplacePivotDialog.svelte";
  import { useMetricsView } from "@rilldata/web-common/features/dashboards/selectors/index";
  import { getStateManagers } from "@rilldata/web-common/features/dashboards/state-managers/state-managers";
  import {
    metricsExplorerStore,
    useDashboardStore,
  } from "@rilldata/web-common/features/dashboards/stores/dashboard-stores";
  import ComparisonSelector from "@rilldata/web-common/features/dashboards/time-controls/ComparisonSelector.svelte";
  import Spinner from "@rilldata/web-common/features/entity-management/Spinner.svelte";
  import { EntityStatus } from "@rilldata/web-common/features/entity-management/types";
  import { TIME_GRAIN } from "@rilldata/web-common/lib/time/config";
  import type { TimeGrain } from "@rilldata/web-common/lib/time/types";
  import { slideRight } from "@rilldata/web-common/lib/transitions";
  import { createEventDispatcher } from "svelte";
  import { fly } from "svelte/transition";
  import { featureFlags } from "../../feature-flags";
  import { PivotChipType } from "../pivot/types";
  import TDDExportButton from "./TDDExportButton.svelte";
  import type { TDDComparison } from "./types";

  export let metricViewName: string;
  export let dimensionName: string;
  export let isFetching = false;
  export let comparing: TDDComparison | undefined;
  export let areAllTableRowsSelected = false;
  export let isRowsEmpty = false;

  const dispatch = createEventDispatcher();
  const { adminServer } = featureFlags;

  const {
    actions: {
      dimensionsFilter: { toggleDimensionFilterMode },
    },
  } = getStateManagers();

  $: metricsView = useMetricsView(getStateManagers());
  $: dashboardStore = useDashboardStore(metricViewName);

  $: expandedMeasureName = $dashboardStore?.expandedMeasureName;
  $: allMeasures = $metricsView?.data?.measures ?? [];

  $: selectableMeasures = allMeasures
    ?.filter((m) => m.name !== undefined || m.label !== undefined)
    .map((m) =>
      // Note: undefined values are filtered out above, so the
      // empty string fallback is unreachable.
      ({
        name: m.name ?? "",
        label: m.label ?? "",
      }),
    );

  $: selectedItems = allMeasures?.map((m) => m.name === expandedMeasureName);

  $: selectedMeasureLabel =
    allMeasures?.find((m) => m.name === expandedMeasureName)?.label ??
    expandedMeasureName ??
    "";

  $: excludeMode =
    $dashboardStore?.dimensionFilterExcludeMode.get(dimensionName) ?? false;

  $: filterKey = excludeMode ? "exclude" : "include";
  $: otherFilterKey = excludeMode ? "include" : "exclude";

  let searchToggle = false;

  let searchText = "";
  function onSearch() {
    dispatch("search", searchText);
  }

  function closeSearchBar() {
    searchText = "";
    searchToggle = false;
    dispatch("search", searchText);
  }

  function onSubmit() {
    if (!areAllTableRowsSelected) {
      dispatch("toggle-all-search-items");
      closeSearchBar();
    }
  }

  function toggleFilterMode() {
    toggleDimensionFilterMode(dimensionName);
  }

  function switchMeasure(event) {
    metricsExplorerStore.setExpandedMeasureName(metricViewName, event.detail);
  }

  let showReplacePivotModal = false;
  function startPivotForTDD() {
    const pivot = $dashboardStore?.pivot;

    if (
      pivot.rows.dimension.length ||
      pivot.columns.measure.length ||
      pivot.columns.dimension.length
    ) {
      showReplacePivotModal = true;
    } else {
      createPivot();
    }
  }

  function createPivot() {
    showReplacePivotModal = false;
    const dashboardGrain = $dashboardStore?.selectedTimeRange?.interval;
    if (!dashboardGrain || !expandedMeasureName) return;

    const timeGrain: TimeGrain = TIME_GRAIN[dashboardGrain];
    const rowDimensions = dimensionName
      ? [
          {
            id: dimensionName,
            title: dimensionName,
            type: PivotChipType.Dimension,
          },
        ]
      : [];
    metricsExplorerStore.createPivot(
      metricViewName,
      { dimension: rowDimensions },
      {
        dimension: [
          {
            id: dashboardGrain,
            title: timeGrain.label,
            type: PivotChipType.Time,
          },
        ],
        measure: [
          {
            id: expandedMeasureName,
            title: expandedMeasureName,
            type: PivotChipType.Measure,
          },
        ],
      },
    );
  }
</script>

<div
  class="grid grid-auto-cols justify-between grid-flow-col items-center p-1 pb-3 h-11 w-full"
>
  <div class="flex gap-x-3 items-center font-normal text-gray-500">
    <div class="flex items-center gap-x-2">
      <div class="flex items-center gap-x-1">
        <Row size="16px" /> Rows
      </div>

      <ComparisonSelector chipStyle {metricViewName} />
    </div>

    <div class="flex items-center gap-x-2 pl-2">
      <div class="flex items-center gap-x-1">
        <Column size="16px" /> Columns
      </div>
      <SearchableFilterChip
        label={selectedMeasureLabel}
        on:item-clicked={switchMeasure}
        selectableItems={selectableMeasures}
        {selectedItems}
        tooltipText="Choose a measure to display"
      />
    </div>

    {#if isFetching}
      <Spinner size="18px" status={EntityStatus.Running} />
    {/if}
  </div>

  {#if comparing === "dimension"}
    <div class="flex items-center mr-4 gap-x-3" style:cursor="pointer">
      {#if !isRowsEmpty}
        <SelectAllButton {areAllTableRowsSelected} on:toggle-all-search-items />
      {/if}

      {#if !searchToggle}
        <button
          class="flex items-center ui-copy-icon"
          in:fly|global={{ x: 10, duration: 300 }}
          style:grid-column-gap=".2rem"
          on:click={() => (searchToggle = !searchToggle)}
        >
          <SearchIcon size="16px" />
          <span> Search </span>
        </button>
      {:else}
        <div
          transition:slideRight={{ leftOffset: 8 }}
          class="flex items-center gap-x-1"
        >
          <Search
            bind:value={searchText}
            on:input={onSearch}
            on:submit={onSubmit}
          />
          <button
            class="ui-copy-icon"
            style:cursor="pointer"
            on:click={() => closeSearchBar()}
          >
            <Close />
          </button>
        </div>
      {/if}

      <Tooltip distance={16} location="left">
        <div class="ui-copy-icon" style:grid-column-gap=".4rem">
          <Switch checked={excludeMode} on:click={() => toggleFilterMode()}>
            Exclude
          </Switch>
        </div>
        <TooltipContent slot="tooltip-content">
          <TooltipTitle>
            <svelte:fragment slot="name">
              Output {filterKey}s selected values
            </svelte:fragment>
          </TooltipTitle>
          <TooltipShortcutContainer>
            <div>Toggle to {otherFilterKey} values</div>
            <Shortcut>Click</Shortcut>
          </TooltipShortcutContainer>
        </TooltipContent>
      </Tooltip>

      <TDDExportButton {metricViewName} includeScheduledReport={$adminServer} />

      <Button
        compact
        type="text"
        on:click={() => {
          startPivotForTDD();
        }}
      >
        Start Pivot
      </Button>
    </div>
  {/if}
</div>

<ReplacePivotDialog
  open={showReplacePivotModal}
  on:close={() => {
    showReplacePivotModal = false;
  }}
  on:replace={() => createPivot()}
/>
