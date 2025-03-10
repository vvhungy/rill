import { expect } from "@playwright/test";
import type { Page } from "playwright";
import {
  clickMenuButton,
  openEntityMenu,
  renameEntityUsingTitle,
  TestEntityType,
} from "./commonHelpers";
import { waitForEntity } from "./waitHelpers";

export async function createModel(page: Page, name: string) {
  // add model button
  await page.locator("button#create-model-button").click();
  await waitForEntity(page, TestEntityType.Model, "model", true);
  await renameEntityUsingTitle(page, name);
  await waitForEntity(page, TestEntityType.Model, name, true);
}

export async function createModelFromSource(page: Page, source: string) {
  await openEntityMenu(page, source);
  await clickMenuButton(page, "Create New Model");
}

export async function modelHasError(page: Page, hasError: boolean, error = "") {
  const errorLocator = page.locator(".editor-pane .error");
  try {
    await errorLocator.waitFor({
      timeout: 100,
    });
  } catch (err) {
    // assertions not needed
  }

  if (hasError) {
    const actualError = await errorLocator.textContent();
    expect(actualError).toMatch(error);
  } else {
    expect(await errorLocator.count()).toBe(0);
  }
}
