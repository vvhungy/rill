<script lang="ts">
  import { onMount } from "svelte";
  import { slide } from "svelte/transition";
  import InfoCircle from "../icons/InfoCircle.svelte";
  import Tooltip from "../tooltip/Tooltip.svelte";
  import TooltipContent from "../tooltip/TooltipContent.svelte";

  export let id = "";
  export let label = "";
  export let error: string;
  export let value: string;
  export let placeholder = "";
  export let hint = "";
  export let claimFocusOnMount = false;

  let inputElement;

  if (claimFocusOnMount) {
    onMount(() => {
      inputElement.focus();
    });
  }
</script>

<div class="flex flex-row items-center pb-2 gap-x-1">
  <label for={id} class="text-gray-800 text-sm font-medium">{label}</label>
  {#if hint}
    <Tooltip location="right" alignment="middle" distance={8}>
      <div class="text-gray-500" style="transform:translateY(-.5px)">
        <InfoCircle size="13px" />
      </div>
      <TooltipContent maxWidth="400px" slot="tooltip-content">
        {@html hint}
      </TooltipContent>
    </Tooltip>
  {/if}
</div>
<input
  bind:this={inputElement}
  bind:value
  on:input
  on:change
  {id}
  name={id}
  type="text"
  {placeholder}
  autocomplete="off"
  class="border border-gray-300 rounded-sm px-3 py-1 cursor-pointer focus:outline-primary-500 w-full text-xs"
/>
{#if error}
  <div in:slide={{ duration: 200 }} class="pl-1 text-red-500 text-xs pt-1">
    {error}
  </div>
{/if}
