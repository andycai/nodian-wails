<script lang="ts">
    import { onMount } from "svelte";
    import { createEventDispatcher } from "svelte";
    import {
        ListMarkdownFiles,
        CreateMarkdownFile,
        CreateFolder,
        RenameItem,
        SelectDirectory,
        GetCurrentDirectory,
        SetCurrentDirectory,
    } from "../wailsjs/go/main/App";
    import { truncateName } from "./utils";

    export let files: string[] = [];
    export let isDarkMode: boolean;

    let newItemName = "";
    let isCreatingFile = false;
    let isCreatingFolder = false;
    let creatingPath = "";
    let isRenaming = false;
    let renamingItem = "";
    let currentDirectory = "";

    const dispatch = createEventDispatcher();

    let expandedFolders: Set<string> = new Set();
    let allFoldersExpanded = false;
    let selectedItem: string | null = null;

    onMount(async () => {
        currentDirectory = await GetCurrentDirectory();
        await refreshFiles();
        // expandedFolders.add("."); // 默认展开根目录
        // expandedFolders = expandedFolders; // 触发更新
    });

    // async function refreshFiles() {
    //     files = await ListMarkdownFiles(currentDirectory);
    //     if (!files.includes(".")) {
    //         files = [".", ...files];
    //     }
    //     sortFiles();
    //     console.log("Files after refresh:", files);
    // }

    async function refreshFiles() {
        files = await ListMarkdownFiles(currentDirectory);
        sortFiles();
        // files.sort((a, b) => {
        //     const aIsDir = a.endsWith("/");
        //     const bIsDir = b.endsWith("/");
        //     if (aIsDir && !bIsDir) return -1;
        //     if (!bIsDir && aIsDir) return 1;
        //     return a.localeCompare(b);
        // });
    }

    function sortFiles() {
        files.sort((a, b) => {
            const aParts = a.split("/").filter(Boolean);
            const bParts = b.split("/").filter(Boolean);
            for (let i = 0; i < Math.min(aParts.length, bParts.length); i++) {
                if (aParts[i] !== bParts[i]) {
                    const aIsDir = aParts[i].endsWith("/");
                    const bIsDir = bParts[i].endsWith("/");
                    if (aIsDir && !bIsDir) return -1;
                    if (!aIsDir && bIsDir) return 1;
                    return aParts[i].localeCompare(bParts[i]);
                }
            }
            return aParts.length - bParts.length;
        });
    }

    // function sortFiles(files: string[]): string[] {
    //     return files.sort((a, b) => {
    //         const aDepth = getDepth(a);
    //         const bDepth = getDepth(b);
    //         if (aDepth !== bDepth) return aDepth - bDepth;
    //         const aIsDir = a.endsWith("/");
    //         const bIsDir = b.endsWith("/");
    //         if (aIsDir && !bIsDir) return -1;
    //         if (!aIsDir && bIsDir) return 1;
    //         return a.localeCompare(b);
    //     });
    // }

    // $: sortedFiles = sortFiles(files);

    async function selectDirectory() {
        const selectedDir = await SelectDirectory();
        if (selectedDir) {
            currentDirectory = selectedDir;
            await SetCurrentDirectory(currentDirectory);
            await refreshFiles();
        }
    }

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

    function cancelCreating() {
        isCreatingFile = false;
        isCreatingFolder = false;
        creatingPath = "";
        newItemName = "";
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === "Enter") {
            if (isRenaming) {
                renameItem();
            } else if (isCreatingFile || isCreatingFolder) {
                createItem();
            }
        } else if (event.key === "Escape") {
            cancelCreating();
        } else if (event.key === "F2" && selectedItem) {
            startRenaming(selectedItem);
        }
    }

    function selectFile(file: string) {
        dispatch("selectFile", file);
    }

    function handleContextMenu(event: MouseEvent, path: string) {
        event.preventDefault();
        if (path.endsWith("/")) {
            startCreatingFile(path.slice(0, -1));
        }
    }

    async function toggleAllFolders() {
        allFoldersExpanded = !allFoldersExpanded;
        if (allFoldersExpanded) {
            files.forEach((file) => {
                if (file.endsWith("/")) {
                    expandedFolders.add(file);
                }
            });
        } else {
            expandedFolders.clear();
        }
        expandedFolders = expandedFolders; // 触发更新

        await refreshFiles();
    }

    async function toggleFolder(folder: string) {
        if (expandedFolders.has(folder)) {
            expandedFolders.delete(folder);
        } else {
            expandedFolders.add(folder);
        }
        expandedFolders = expandedFolders; // 触发更新

        await refreshFiles();
    }

    function isExpanded(folder: string): boolean {
        return expandedFolders.has(folder);
    }

    function getDepth(path: string): number {
        return path.split("/").length - 1;
    }

    // function isVisible(file: string): boolean {
    //     if (file === ".") return true;
    //     const parts = file.split("/").filter(Boolean);
    //     let currentPath = ".";
    //     for (let i = 0; i < parts.length; i++) {
    //         if (i === parts.length - 1) return true; // 如果是最后一部分，总是显示
    //         currentPath += "/" + parts[i];
    //         if (!expandedFolders.has(currentPath)) {
    //             return false;
    //         }
    //     }
    //     return true;
    // }

    function isVisible(file: string): boolean {
        const parts = file.split("/").filter(Boolean);
        if (parts.length === 1) return true; // Root level items are always visible
        let currentPath = "";
        for (let i = 0; i < parts.length - 1; i++) {
            currentPath += parts[i] + "/";
            if (!expandedFolders.has(currentPath)) {
                return false;
            }
        }
        return true;
    }

    function selectItem(item: string) {
        selectedItem = item;
        if (item.endsWith("/")) {
            dispatch("selectFolder", item);
        } else {
            dispatch("selectFile", item);
        }
    }

    function startCreatingItem(isFile: boolean) {
        if (selectedItem && selectedItem.endsWith("/")) {
            creatingPath = selectedItem;
        } else {
            creatingPath = "";
        }
        isCreatingFile = isFile;
        isCreatingFolder = !isFile;
        newItemName = "";
    }

    async function createItem() {
        if (newItemName) {
            const fullPath = creatingPath
                ? `${creatingPath}${newItemName}`
                : `${currentDirectory}/${newItemName}`;

            // 检查并创建根目录
            await CreateFolder(currentDirectory);

            if (isCreatingFile) {
                await CreateMarkdownFile(fullPath, "");
            } else if (isCreatingFolder) {
                await CreateFolder(fullPath);
            }
            await refreshFiles();
            cancelCreating();
        }
    }

    function startRenaming(item: string) {
        isRenaming = true;
        renamingItem = item;
        newItemName = item.split("/").pop() || "";
    }

    async function renameItem() {
        if (newItemName && renamingItem !== newItemName) {
            const oldPath = `${currentDirectory}/${renamingItem}`;
            const newPath = `${currentDirectory}/${newItemName}`;
            await RenameItem(oldPath, newPath);
            await refreshFiles();
        }
        isRenaming = false;
        renamingItem = "";
        newItemName = "";
    }

    function getIndent(file: string): number {
        return file.split("/").filter(Boolean).length - 1;
    }

    function getFileName(file: string): string {
        const parts = file.split("/").filter(Boolean);
        return parts[parts.length - 1];
    }
</script>

<div
    class="w-64 bg-gray-700 dark:bg-gray-800 p-4 overflow-y-auto custom-scrollbar"
>
    <div class="flex justify-between items-center mb-4">
        <h2
            class="text-lg font-semibold text-gray-300 cursor-pointer"
            on:click={selectDirectory}
        >
            {currentDirectory.split("/").pop()?.charAt(0).toUpperCase() +
                currentDirectory.split("/").pop()?.slice(1) || "Root"}
        </h2>
        <div class="flex space-x-2">
            <button
                on:click={() => startCreatingItem(true)}
                class="p-1 rounded hover:bg-gray-600"
                title="New File"
            >
                <svg
                    class="w-4 h-4 text-gray-300"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M12 6v6m0 0v6m0-6h6m-6 0H6"
                    ></path>
                </svg>
            </button>
            <button
                on:click={() => startCreatingItem(false)}
                class="p-1 rounded hover:bg-gray-600"
                title="New Folder"
            >
                <svg
                    class="w-4 h-4 text-gray-300"
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
            <button
                on:click={refreshFiles}
                class="p-1 rounded hover:bg-gray-600"
                title="Refresh"
            >
                <svg
                    class="w-4 h-4 text-gray-300"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                    ></path>
                </svg>
            </button>
            <button
                on:click={toggleAllFolders}
                class="p-1 rounded hover:bg-gray-600"
                title={allFoldersExpanded ? "Collapse All" : "Expand All"}
            >
                <svg
                    class="w-4 h-4 text-gray-300"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d={allFoldersExpanded
                            ? "M19 9l-7 7-7-7"
                            : "M9 5l7 7-7 7"}
                    ></path>
                </svg>
            </button>
        </div>
    </div>
    <ul class="space-y-1">
        {#each files as file}
            {#if isVisible(file)}
                <li style="padding-left: {getIndent(file) * 16}px;">
                    {#if isRenaming && file === renamingItem}
                        <input
                            type="text"
                            bind:value={newItemName}
                            on:keydown={handleKeydown}
                            class="w-full p-1 text-sm border rounded dark:bg-gray-700 dark:text-white"
                            autofocus
                        />
                    {:else}
                        <button
                            on:click={() => {
                                selectItem(file);
                                if (file.endsWith("/") || file === ".")
                                    toggleFolder(file);
                            }}
                            on:keydown={handleKeydown}
                            class="w-full text-left p-1 text-sm hover:bg-gray-600 rounded flex items-center text-gray-300 {selectedItem ===
                            file
                                ? 'bg-gray-600'
                                : ''}"
                        >
                            {#if file.endsWith("/") || file === "."}
                                <svg
                                    class="w-4 h-4 mr-2 text-gray-300 flex-shrink-0"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                    xmlns="http://www.w3.org/2000/svg"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d={isExpanded(file)
                                            ? "M19 9l-7 7-7-7"
                                            : "M9 5l7 7-7 7"}
                                    ></path>
                                </svg>
                            {:else}
                                <svg
                                    class="w-4 h-4 mr-2 text-gray-300 flex-shrink-0"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                    xmlns="http://www.w3.org/2000/svg"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z"
                                    ></path>
                                </svg>
                            {/if}
                            <span class="truncate">
                                {truncateName(
                                    getFileName(file),
                                    file.endsWith("/") || file === "." ? 8 : 12,
                                )}
                            </span>
                        </button>
                    {/if}
                </li>
            {/if}
        {/each}
        {#if (isCreatingFile || isCreatingFolder) && creatingPath}
            <li style="padding-left: {(getIndent(creatingPath) + 1) * 16}px;">
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

<style>
    /* 自定义滚动条样式 */
    .custom-scrollbar {
        scrollbar-width: thin;
        scrollbar-color: rgba(156, 163, 175, 0.5) transparent;
    }

    .custom-scrollbar::-webkit-scrollbar {
        width: 6px;
        height: 6px;
    }

    .custom-scrollbar::-webkit-scrollbar-track {
        background: transparent;
    }

    .custom-scrollbar::-webkit-scrollbar-thumb {
        background-color: rgba(156, 163, 175, 0.5);
        border-radius: 3px;
    }
</style>
