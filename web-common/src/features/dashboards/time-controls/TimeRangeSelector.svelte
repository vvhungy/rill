<script lang="ts">
  import IconSpaceFixer from "@rilldata/web-common/components/button/IconSpaceFixer.svelte";
  import { WithTogglableFloatingElement } from "@rilldata/web-common/components/floating-element";
  import Calendar from "@rilldata/web-common/components/icons/Calendar.svelte";
  import CaretDownIcon from "@rilldata/web-common/components/icons/CaretDownIcon.svelte";
  import { useMetricsView } from "@rilldata/web-common/features/dashboards/selectors/index";
  import { getStateManagers } from "@rilldata/web-common/features/dashboards/state-managers/state-managers";
  import DefaultTimeRangeMenuItem from "@rilldata/web-common/features/dashboards/time-controls/DefaultTimeRangeMenuItem.svelte";
  import TimeRangeScrubChip from "@rilldata/web-common/features/dashboards/time-controls/TimeRangeScrubChip.svelte";
  import { useTimeControlStore } from "@rilldata/web-common/features/dashboards/time-controls/time-control-store";
  import { getOrderedStartEnd } from "@rilldata/web-common/features/dashboards/time-series/utils";
  import {
    ALL_TIME,
    DEFAULT_TIME_RANGES,
  } from "@rilldata/web-common/lib/time/config";
  import { prettyFormatTimeRange } from "@rilldata/web-common/lib/time/ranges";
  import { humaniseISODuration } from "@rilldata/web-common/lib/time/ranges/iso-ranges";
  import {
    DashboardTimeControls,
    TimeComparisonOption,
    TimeRange,
    TimeRangePreset,
  } from "@rilldata/web-common/lib/time/types";
  import { createEventDispatcher } from "svelte";
  import { slide } from "svelte/transition";
  import { useDashboardStore } from "web-common/src/features/dashboards/stores/dashboard-stores";
  import { Menu, MenuItem } from "../../../components/menu";
  import Divider from "../../../components/menu/core/Divider.svelte";
  import { LIST_SLIDE_DURATION } from "../../../layout/config";
  import type { V1TimeGrain } from "../../../runtime-client";
  import CustomTimeRangeInput from "./CustomTimeRangeInput.svelte";
  import CustomTimeRangeMenuItem from "./CustomTimeRangeMenuItem.svelte";

  export let metricViewName: string;
  export let boundaryStart: Date;
  export let boundaryEnd: Date;
  export let minTimeGrain: V1TimeGrain;
  export let selectedRange: DashboardTimeControls;

  const dispatch = createEventDispatcher();

  $: dashboardStore = useDashboardStore(metricViewName);

  const ctx = getStateManagers();
  const timeControlsStore = useTimeControlStore(ctx);
  const metricsView = useMetricsView(ctx);
  const {
    selectors: {
      timeRangeSelectors: { timeRangeSelectorState },
    },
  } = ctx;

  let isCustomRangeOpen = false;
  let isCalendarRecentlyClosed = false;

  $: selectedSubRange =
    $dashboardStore?.selectedScrubRange?.start &&
    $dashboardStore?.selectedScrubRange?.end
      ? {
          start: $dashboardStore.selectedScrubRange.start,
          end: $dashboardStore.selectedScrubRange.end,
        }
      : null;

  function setIntermediateSelection(
    timeRangeName: TimeRangePreset | TimeComparisonOption | undefined,
  ) {
    return () => {
      if (timeRangeName) {
        intermediateSelection = timeRangeName;
      }
    };
  }

  function onSelectRelativeTimeRange(
    timeRange: TimeRange | undefined,
    closeMenu: () => void,
  ) {
    closeMenu();
    if (!timeRange) {
      return;
    }
    dispatch("select-time-range", {
      name: timeRange.name,
      start: timeRange.start,
      end: timeRange.end,
    });
  }

  function onSelectCustomTimeRange(
    startDate: string,
    endDate: string,
    closeMenu: () => void,
  ) {
    setIntermediateSelection(TimeRangePreset.CUSTOM)();
    closeMenu();
    dispatch("select-time-range", {
      name: TimeRangePreset.CUSTOM,
      start: startDate,
      end: endDate,
    });
  }

  function zoomScrub(toggleFloatingElement) {
    if (
      !$dashboardStore?.selectedScrubRange?.start ||
      !$dashboardStore?.selectedScrubRange?.end
    ) {
      return;
    }
    const { start, end } = getOrderedStartEnd(
      $dashboardStore?.selectedScrubRange?.start,
      $dashboardStore?.selectedScrubRange?.end,
    );
    onSelectRelativeTimeRange(
      {
        name: TimeRangePreset.CUSTOM,
        start,
        end,
      },
      toggleFloatingElement,
    );
    dispatch("remove-scrub");
  }

  function onClickOutside(closeMenu: () => void) {
    if (!isCalendarRecentlyClosed) {
      closeMenu();
    }
  }

  // If the user clicks outside to close the calendar, we don't want the `click-outside` event to close the whole menu.
  // A nice solution would be to check for `isCalendarOpen` in the `onClickOutside` function. However, the calendar
  // closes *before* the `click-outside` event is fired. A workaround is to check for `isCalendarRecentlyClosed`.
  function onCalendarClose() {
    isCalendarRecentlyClosed = true;

    setTimeout(() => {
      isCalendarRecentlyClosed = false;
    }, 300);
  }

  $: currentSelection = $dashboardStore?.selectedTimeRange?.name;
  $: intermediateSelection = currentSelection;

  const handleMenuOpen = () => {
    if (intermediateSelection !== TimeRangePreset.CUSTOM) {
      isCustomRangeOpen = false;
    }
  };

  $: defaultTimeRange = $metricsView.data?.defaultTimeRange
    ? ($metricsView.data?.defaultTimeRange as TimeRangePreset)
    : undefined;
</script>

<WithTogglableFloatingElement
  alignment="start"
  distance={8}
  let:active
  let:toggleFloatingElement
  let:handleClose
  on:open={handleMenuOpen}
>
  <button on:click={toggleFloatingElement}>
    {#if selectedSubRange}
      <TimeRangeScrubChip
        on:remove={() => {
          handleClose();
          dispatch("remove-scrub");
        }}
        {active}
        start={selectedSubRange.start}
        end={selectedSubRange.end}
        zone={$dashboardStore?.selectedTimezone}
      />
    {:else}
      <div
        class:bg-gray-200={active}
        class="px-3 py-2 rounded flex flex-row gap-x-2 hover:bg-gray-200 hover:dark:bg-gray-600 items-baseline"
        aria-label="Select time range"
      >
        <div class="flex flew-row gap-x-3">
          <div class="font-bold flex flex-row items-center gap-x-3">
            <span class="ui-copy-icon"><Calendar size="16px" /></span>
            <span style:transform="translateY(1px)">
              <!-- This conditional shouldn't be necessary because there should always be a selected (at least default) time range -->
              {#if intermediateSelection === TimeRangePreset.CUSTOM}
                Custom range
              {:else if currentSelection}
                {#if currentSelection in DEFAULT_TIME_RANGES}
                  {DEFAULT_TIME_RANGES[currentSelection].label}
                {:else}
                  Last {humaniseISODuration(currentSelection)}
                {/if}
              {:else}
                Select a time range
              {/if}
            </span>
          </div>
          <span style:transform="translateY(1px)">
            {prettyFormatTimeRange(
              $timeControlsStore?.selectedTimeRange?.start,
              $timeControlsStore?.selectedTimeRange?.end,
              $timeControlsStore?.selectedTimeRange?.name,
              $dashboardStore?.selectedTimezone,
            )}
          </span>
        </div>
        <IconSpaceFixer pullRight>
          <div class="transition-transform" class:-rotate-180={active}>
            <CaretDownIcon size="14px" />
          </div>
        </IconSpaceFixer>
      </div>
    {/if}
  </button>
  <Menu
    label="Time range selector"
    let:handleClose
    maxWidth="300px"
    on:click-outside={() => onClickOutside(handleClose)}
    on:escape={handleClose}
    slot="floating-element"
  >
    {@const allTime = {
      name: TimeRangePreset.ALL_TIME,
      label: ALL_TIME.label,
      start: boundaryStart,
      end: new Date(boundaryEnd.getTime() + 1), // end is exclusive
    }}
    {#if selectedSubRange}
      <MenuItem
        on:before-select={setIntermediateSelection(TimeRangePreset.CUSTOM)}
        on:select={() => zoomScrub(handleClose)}
      >
        <span> Zoom to subrange </span>
        <span slot="right">Z</span>
      </MenuItem>
      <Divider />
    {/if}
    <MenuItem
      on:before-select={setIntermediateSelection(allTime.name)}
      on:select={() => onSelectRelativeTimeRange(allTime, handleClose)}
    >
      <span class:font-bold={intermediateSelection === allTime.name}>
        {allTime.label}
      </span>
    </MenuItem>
    {#if $timeRangeSelectorState.showDefaultItem && defaultTimeRange}
      <DefaultTimeRangeMenuItem
        on:before-select={setIntermediateSelection(defaultTimeRange)}
        on:select={() =>
          onSelectRelativeTimeRange(
            $timeControlsStore.defaultTimeRange,
            handleClose,
          )}
        selected={intermediateSelection === defaultTimeRange}
        isoDuration={defaultTimeRange}
      />
    {/if}
    {#if $timeRangeSelectorState.latestWindowTimeRanges?.length}
      <Divider />
      {#each $timeRangeSelectorState.latestWindowTimeRanges as timeRange}
        <MenuItem
          on:before-select={setIntermediateSelection(timeRange.name)}
          on:select={() => onSelectRelativeTimeRange(timeRange, handleClose)}
        >
          <span class:font-bold={intermediateSelection === timeRange.name}>
            {timeRange.label}
          </span>
        </MenuItem>
      {/each}
    {/if}
    {#if $timeRangeSelectorState.periodToDateRanges?.length}
      <Divider />
      {#each $timeRangeSelectorState.periodToDateRanges as timeRange}
        <MenuItem
          on:before-select={setIntermediateSelection(timeRange.name)}
          on:select={() => onSelectRelativeTimeRange(timeRange, handleClose)}
        >
          <span class:font-bold={intermediateSelection === timeRange.name}>
            {timeRange.label}
          </span>
        </MenuItem>
      {/each}
    {/if}
    {#if $timeRangeSelectorState.previousCompleteDateRanges?.length}
      <Divider />
      {#each $timeRangeSelectorState.previousCompleteDateRanges as timeRange}
        <MenuItem
          on:before-select={setIntermediateSelection(timeRange.name)}
          on:select={() => onSelectRelativeTimeRange(timeRange, handleClose)}
        >
          <span class:font-bold={intermediateSelection === timeRange.name}>
            {timeRange.label}
          </span>
        </MenuItem>
      {/each}
    {/if}
    <Divider />
    <CustomTimeRangeMenuItem
      on:select={() => {
        isCustomRangeOpen = !isCustomRangeOpen;
      }}
      open={isCustomRangeOpen}
      selected={intermediateSelection === TimeRangePreset.CUSTOM}
    />
    {#if isCustomRangeOpen}
      <div transition:slide={{ duration: LIST_SLIDE_DURATION }}>
        <CustomTimeRangeInput
          {boundaryStart}
          {boundaryEnd}
          {minTimeGrain}
          zone={$dashboardStore?.selectedTimezone}
          defaultDate={selectedRange}
          on:apply={(e) =>
            onSelectCustomTimeRange(
              e.detail.startDate,
              e.detail.endDate,
              handleClose,
            )}
          on:close-calendar={onCalendarClose}
        />
      </div>
    {/if}
  </Menu>
</WithTogglableFloatingElement>
