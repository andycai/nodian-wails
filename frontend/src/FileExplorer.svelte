<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import {
        ListMarkdownFiles,
        CreateMarkdownFile,
        CreateFolder,
    } from "../wailsjs/go/main/App";

    export let files: string[] = [];

    let newItemName = "";
    let isCreatingFile = false;
    let isCreatingFolder = false;
    let creatingPath = "";

    const dispatch = createEventDispatcher();

    async function startCreatingFile(path = "") {
        isCreatingFile = true;
        isCreatingFolder = false;
        creatingPath = path;
        newItemName = "";
    }

    async function startCreatingFolder(path = "") {
        isCreatingFolder = true;
        isCreatingFile = false;
        creatingPath = path;
        newItemName = "";
    }

    async function createItem() {
        if (newItemName) {
            const fullPath = creatingPath
                ? `${creatingPath}/${newItemName}`
                : newItemName;
            if (isCreatingFile) {
                await CreateMarkdownFile(fullPath, "");
            } else if (isCreatingFolder) {
                await CreateFolder(fullPath);
            }
            files = await ListMarkdownFiles(".");
            cancelCreating();
        }
    }

    function cancelCreating() {
        isCreatingFile = false;
        isCreatingFolder = false;
        creatingPath = "";
        newItemName = "";
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === "Enter") {
            createItem();
        } else if (event.key === "Escape") {
            cancelCreating();
        }
    }

    function selectFile(file: string) {
        dispatch("selectFile", file);
    }

    function handleContextMenu(event: MouseEvent, path: string) {
        event.preventDefault();
        // Here you can implement a context menu
        // For now, we'll just use it to create a file in the folder
        if (path.endsWith("/")) {
            startCreatingFile(path.slice(0, -1));
        }
    }
</script>

<div class="p-4">
    <div class="flex justify-between items-center mb-4">
        <h2 class="text-lg font-semibold text-gray-700 dark:text-gray-300">
            Files
        </h2>
        <div class="flex space-x-2">
            <button
                on:click={() => startCreatingFile()}
                class="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-600"
            >
                <svg
                    class="w-5 h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M9 13h6m-3-3v6m5 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                    ></path>
                </svg>
            </button>
            <button
                on:click={() => startCreatingFolder()}
                class="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-600"
            >
                <svg
                    class="w-5 h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z"
                    ></path>
                </svg>
            </button>
        </div>
    </div>
    <ul class="space-y-1">
        {#each files as file}
            <li>
                {#if (isCreatingFile || isCreatingFolder) && creatingPath === file.slice(0, -1)}
                    <input
                        type="text"
                        bind:value={newItemName}
                        on:keydown={handleKeydown}
                        class="w-full p-1 text-sm border rounded dark:bg-gray-700 dark:text-white"
                        autofocus
                    />
                {:else}
                    <button
                        on:click={() => selectFile(file)}
                        on:contextmenu={(e) => handleContextMenu(e, file)}
                        class="w-full text-left p-1 text-sm hover:bg-gray-200 dark:hover:bg-gray-600 rounded flex items-center"
                    >
                        {#if file.endsWith("/")}
                            <svg
                                class="w-4 h-4 mr-2"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                                xmlns="http://www.w3.org/2000/svg"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
                                ></path>
                            </svg>
                        {:else}
                            <svg
                                class="w-4 h-4 mr-2"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                                xmlns="http://www.w3.org/2000/svg"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                                ></path>
                            </svg>
                        {/if}
                        {file}
                    </button>
                {/if}
            </li>
        {/each}
        {#if (isCreatingFile || isCreatingFolder) && !creatingPath}
            <li>
                <input
                    type="text"
                    bind:value={newItemName}
                    on:keydown={handleKeydown}
                    class="w-full p-1 text-sm border rounded dark:bg-gray-700 dark:text-white"
                    autofocus
                />
            </li>
        {/if}
    </ul>
</div>
