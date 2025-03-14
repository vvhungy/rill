<!-- @component 

Implementation for showing a Chip with a title + a list of values. Shows the first value in selectedValues, while enabling
a way to select values from a menu popover.

The RemovableListChip has a few features that are worth noting:
- the remove toggle is on the left side, rather than the right side, which is more traditional
with chips. The main reason for this is the user should not have to look to the left side of a longer chip to see
the name and then move the cursor to the right to cancel it.
- clicking the chip body will expand out a the RemovableListMenu. This component will be in charge of both selecting / de-selecting
existing elements in the lib as well as changing the type (include, exclude) and enabling list search. The implementation of these parts
are details left to the consumer of the component; this component should remain pure-ish (only internal state) if possible.
-->

<script context="module" lang="ts">
  import { createEventDispatcher, onMount } from "svelte";
  import { fly } from "svelte/transition";
  import WithTogglableFloatingElement from "../../floating-element/WithTogglableFloatingElement.svelte";
  import Tooltip from "../../tooltip/Tooltip.svelte";
  import TooltipContent from "../../tooltip/TooltipContent.svelte";
  import TooltipTitle from "../../tooltip/TooltipTitle.svelte";
  import { ChipColors, defaultChipColors } from "../chip-types";
  import { Chip } from "../index";
  import RemovableListBody from "./RemovableListBody.svelte";
  import RemovableListMenu from "./RemovableListMenu.svelte";
</script>

<script lang="ts">
  export let name: string;
  export let selectedValues: string[];
  export let allValues: string[] | null;

  /** an optional type label that will appear in the tooltip */
  export let typeLabel: string;
  export let excludeMode: boolean;
  export let colors: ChipColors = defaultChipColors;
  export let label: string | undefined = undefined;

  let active = !selectedValues.length;

  const dispatch = createEventDispatcher();

  onMount(() => {
    dispatch("mount");
  });

  function handleDismiss() {
    if (!selectedValues.length) {
      dispatch("remove");
    } else {
      active = false;
    }
  }
</script>

<WithTogglableFloatingElement
  bind:active
  let:toggleFloatingElement
  distance={8}
  alignment="start"
>
  <Tooltip
    location="bottom"
    alignment="start"
    distance={8}
    activeDelay={60}
    suppress={active}
  >
    <Chip
      removable
      on:click={() => {
        toggleFloatingElement();
        dispatch("click");
      }}
      on:remove={() => dispatch("remove")}
      {active}
      {...colors}
      {label}
      outline
    >
      <!-- remove button tooltip -->
      <svelte:fragment slot="remove-tooltip">
        <slot name="remove-tooltip-content">
          remove {selectedValues.length}
          value{#if selectedValues.length !== 1}s{/if} for {name}</slot
        >
      </svelte:fragment>
      <!-- body -->
      <RemovableListBody
        slot="body"
        label={name}
        values={selectedValues}
        show={1}
        {active}
      />
    </Chip>
    <div slot="tooltip-content" transition:fly={{ duration: 100, y: 4 }}>
      <TooltipContent maxWidth="400px">
        <TooltipTitle>
          <svelte:fragment slot="name">{name}</svelte:fragment>
          <svelte:fragment slot="description">{typeLabel || ""}</svelte:fragment
          >
        </TooltipTitle>
        {#if $$slots["body-tooltip-content"]}
          <slot name="body-tooltip-content">click to edit the values</slot>
        {/if}
      </TooltipContent>
    </div>
  </Tooltip>
  <RemovableListMenu
    slot="floating-element"
    {excludeMode}
    {allValues}
    {selectedValues}
    on:escape={handleDismiss}
    on:click-outside={handleDismiss}
    on:apply
    on:search
    on:toggle
  />
</WithTogglableFloatingElement>
