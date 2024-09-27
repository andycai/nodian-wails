<script lang="ts">
    import { onMount } from "svelte";
    import { createEventDispatcher } from "svelte";
    import {
        ListMarkdownFiles,
        CreateMarkdownFile,
        CreateFolder,
        CreateRootFolder,
        RenameItem,
        SelectDirectory,
        GetCurrentDirectory,
        SetCurrentDirectory,
        DeleteItem,
    } from "../wailsjs/go/main/App";
    import { truncateName } from "./utils";
    import ConfirmDialog from "./ConfirmDialog.svelte";

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
    let hoveredItem: string | null = null;

    let showConfirmDialog = false;
    let itemToDelete: string | null = null;
    let confirmMessage = "";

    onMount(async () => {
        currentDirectory = await GetCurrentDirectory();
        await refreshFiles();
    });

    async function refreshFiles() {
        files = await ListMarkdownFiles(currentDirectory);
        // 确保所有文件夹都以 '/' 结尾
        files = files.map((file) => {
            if (isFolder(file) && !file.endsWith("/")) {
                return file + "/";
            }
            return file;
        });
        sortFiles();

        // 保留展开状态
        const newExpandedFolders = new Set<string>();
        expandedFolders.forEach((folder) => {
            if (files.includes(folder)) {
                newExpandedFolders.add(folder);
            }
        });
        expandedFolders = newExpandedFolders;
    }

    function sortFiles() {
        const sortedFiles: string[] = [];
        const fileTree: { [key: string]: string[] } = {};

        // 构建文件树
        files.forEach((file) => {
            const parts = file.split("/"); //.filter(Boolean);
            // const parts = file.split("/").filter(Boolean);
            // console.log("parts", file, parts);
            let currentPath = "";
            for (let i = 0; i < parts.length; i++) {
                currentPath += (currentPath ? "/" : "") + parts[i];
                // 路径中最后一个元素
                if (i === parts.length - 1) {
                    // 文件夹
                    if (parts[i] === "") {
                        //
                    } else {
                        //
                    }

                    const parentPath = currentPath.substring(
                        0,
                        currentPath.lastIndexOf("/"),
                    );
                    if (!fileTree[parentPath]) {
                        fileTree[parentPath] = [];
                    }
                    // console.log(
                    //     "currentPath",
                    //     parentPath,
                    //     currentPath,
                    //     parts[i],
                    // );
                    // fileTree[parentPath].push(currentPath);
                    if (parts[i] !== "") {
                        fileTree[parentPath].push(currentPath);
                    } else {
                        const parentParts = parentPath
                            .split("/")
                            .filter(Boolean);

                        if (parentParts.length === 1) {
                            fileTree[""].push(parentPath);
                        } else {
                            fileTree[parentParts[0]].push(parentPath);
                        }
                    }
                } else {
                    if (!fileTree[currentPath]) {
                        fileTree[currentPath] = [];
                    }
                }
            }
        });

        console.log("File tree:", fileTree);
        // 递归排序和添加文件
        function addSortedFiles(path: string) {
            const items = fileTree[path] || [];
            // const folders = items.filter(
            //     (item) => fileTree[item] !== undefined,
            // );
            const folders = items.filter(
                (item) => !(fileTree[item] === undefined),
            );
            const files = items.filter((item) => fileTree[item] === undefined);
            // console.log("items:", path, items);
            console.log("folders:", path, folders);
            // console.log("files:", path, files);

            // 文件夹按字母顺序排序
            folders.sort((a, b) => a.localeCompare(b));
            // 文件按字母顺序排序
            files.sort((a, b) => a.localeCompare(b));

            // 先添加文件夹
            folders.forEach((folder) => {
                sortedFiles.push(folder + "/");
                addSortedFiles(folder);
            });

            // 再添加文件
            files.forEach((file) => {
                sortedFiles.push(file);
            });
        }

        addSortedFiles("");

        files = sortedFiles;
        console.log("Sorted files:", files);
    }

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
            if (isRenaming) {
                cancelRenaming();
            } else {
                cancelCreating();
            }
        } else if (event.key === "F2" && selectedItem) {
            // startRenaming(selectedItem);
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
        expandedFolders = new Set(expandedFolders); // 触发更新

        await refreshFiles();
    }

    async function toggleFolder(folder: string) {
        if (expandedFolders.has(folder)) {
            expandedFolders.delete(folder);
        } else {
            expandedFolders.add(folder);
        }
        expandedFolders = new Set(expandedFolders); // 触发更新
        await refreshFiles();
    }

    function isVisible(file: string): boolean {
        if (file === "" || file === "/") return true;
        const parts = file.split("/").filter(Boolean);
        let currentPath = "";
        for (let i = 0; i < parts.length; i++) {
            if (i === parts.length - 1) return true; // 如果是最后一部分，总是显示
            currentPath += (currentPath ? "/" : "") + parts[i];
            if (!expandedFolders.has(currentPath + "/")) {
                return false;
            }
        }
        return true;
    }

    function selectItem(item: string) {
        selectedItem = item;
        if (isFolder(item)) {
            dispatch("selectFolder", item);
        } else {
            dispatch("selectFile", item);
        }
    }

    function startCreatingItem(isFile: boolean) {
        if (selectedItem && isFolder(selectedItem)) {
            creatingPath = selectedItem;
        } else {
            creatingPath = currentDirectory + "/";
        }
        isCreatingFile = isFile;
        isCreatingFolder = !isFile;
        newItemName = "";
    }

    async function createItem() {
        if (newItemName) {
            let fullPath;
            if (creatingPath) {
                fullPath = `${creatingPath}${isFolder(creatingPath) ? "" : "/"}${newItemName}`;
            } else {
                fullPath = `${currentDirectory}/${newItemName}`;
            }

            console.log("Creating item:", fullPath, currentDirectory);
            // 检查并创建根目录
            await CreateRootFolder(currentDirectory);

            if (isCreatingFile) {
                await CreateMarkdownFile(fullPath, "");
            } else if (isCreatingFolder) {
                await CreateFolder(fullPath);
            }
            await refreshFiles();
            cancelCreating();
        }
    }

    function startRenaming(item: string, event: MouseEvent) {
        event.stopPropagation();
        isRenaming = true;
        renamingItem = item;
        newItemName = getFileName(item);
    }

    function cancelRenaming() {
        isRenaming = false;
        renamingItem = "";
        newItemName = "";
    }

    async function renameItem() {
        if (newItemName && renamingItem !== newItemName) {
            const oldPath = `${currentDirectory}/${renamingItem}`;
            let newItemPath = getFilePath(renamingItem);
            const newPath = `${currentDirectory}/${newItemPath}/${newItemName}`;
            await RenameItem(oldPath, newPath);
            await refreshFiles();
        }
        isRenaming = false;
        renamingItem = "";
        newItemName = "";
    }

    function showDeleteConfirmation(item: string, event: MouseEvent) {
        event.stopPropagation();
        itemToDelete = item;
        confirmMessage = `Are you sure you want to delete ${getFileName(item)}?`;
        showConfirmDialog = true;
    }

    async function handleDeleteConfirm() {
        if (itemToDelete) {
            await DeleteItem(`${currentDirectory}/${itemToDelete}`);
            await refreshFiles();
        }
        itemToDelete = null;
    }

    function handleDeleteCancel() {
        itemToDelete = null;
    }

    function getIndent(file: string): number {
        return file.split("/").filter(Boolean).length - 1;
    }

    function getFileName(file: string): string {
        const parts = file.split("/").filter(Boolean);
        return parts[parts.length - 1];
    }

    function getFilePath(file: string): string {
        const parts = file.split("/").filter(Boolean);
        return parts.slice(0, -1).join("/");
    }

    function isFolder(file: string): boolean {
        // return !file.endsWith(".mk");
        return (
            file.endsWith("/") || files.some((f) => f.startsWith(file + "/"))
        );
    }
</script>

<div
    class="w-64 bg-gray-100 dark:bg-gray-800 p-4 overflow-y-auto custom-scrollbar"
>
    <div class="flex justify-between items-center mb-4">
        <h2
            class="text-lg font-semibold text-gray-900 dark:text-gray-100 cursor-pointer"
            on:click={selectDirectory}
        >
            {currentDirectory.split("/").pop()?.charAt(0).toUpperCase() +
                currentDirectory.split("/").pop()?.slice(1) || "Root"}
        </h2>
        <div class="flex space-x-2">
            <button
                on:click={() => startCreatingItem(true)}
                class="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-700"
                title="New File"
            >
                <svg
                    class="w-4 h-4 text-gray-900 dark:text-gray-100"
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
                class="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-700"
                title="New Folder"
            >
                <svg
                    class="w-4 h-4 text-gray-900 dark:text-gray-100"
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
                class="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-700"
                title="Refresh"
            >
                <svg
                    class="w-4 h-4 text-gray-900 dark:text-gray-100"
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
                class="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-700"
                title={allFoldersExpanded ? "Collapse All" : "Expand All"}
            >
                <svg
                    class="w-4 h-4 text-gray-900 dark:text-gray-100"
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
                <li
                    style="padding-left: {getIndent(file) * 16}px;"
                    on:mouseenter={() => (hoveredItem = file)}
                    on:mouseleave={() => (hoveredItem = null)}
                >
                    {#if isRenaming && file === renamingItem}
                        <input
                            type="text"
                            bind:value={newItemName}
                            on:keydown={handleKeydown}
                            class="w-full p-1 text-sm border rounded dark:bg-gray-700 dark:text-white"
                            autofocus
                        />
                    {:else}
                        <div class="flex items-center justify-between">
                            <button
                                on:click={() => {
                                    selectItem(file);
                                    if (isFolder(file)) toggleFolder(file);
                                }}
                                class="flex-grow text-left p-1 text-sm hover:bg-gray-200 dark:hover:bg-gray-700 rounded flex items-center text-gray-900 dark:text-gray-100 {selectedItem ===
                                file
                                    ? 'bg-gray-200 dark:bg-gray-700'
                                    : ''}"
                            >
                                {#if isFolder(file)}
                                    <svg
                                        class="w-4 h-4 mr-2 text-gray-900 dark:text-gray-100 flex-shrink-0"
                                        fill="none"
                                        stroke="currentColor"
                                        viewBox="0 0 24 24"
                                        xmlns="http://www.w3.org/2000/svg"
                                    >
                                        <path
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                            stroke-width="2"
                                            d={expandedFolders.has(file)
                                                ? "M19 9l-7 7-7-7"
                                                : "M9 5l7 7-7 7"}
                                        ></path>
                                    </svg>
                                {:else}
                                    <svg
                                        class="w-4 h-4 mr-2 text-gray-900 dark:text-gray-100 flex-shrink-0"
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
                                        isFolder(file) ? 8 : 12,
                                    )}
                                </span>
                            </button>
                            {#if hoveredItem === file}
                                <div class="flex">
                                    <button
                                        on:click={(event) =>
                                            startRenaming(file, event)}
                                        class="p-1 text-gray-900 dark:text-gray-100 hover:text-white"
                                    >
                                        <svg
                                            class="w-4 h-4"
                                            fill="none"
                                            stroke="currentColor"
                                            viewBox="0 0 24 24"
                                            xmlns="http://www.w3.org/2000/svg"
                                        >
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="2"
                                                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                                            ></path>
                                        </svg>
                                    </button>
                                    <button
                                        on:click={(event) =>
                                            showDeleteConfirmation(file, event)}
                                        class="p-1 text-gray-900 dark:text-gray-100 hover:text-white"
                                    >
                                        <svg
                                            class="w-4 h-4"
                                            fill="none"
                                            stroke="currentColor"
                                            viewBox="0 0 24 24"
                                            xmlns="http://www.w3.org/2000/svg"
                                        >
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="2"
                                                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                                            ></path>
                                        </svg>
                                    </button>
                                </div>
                            {/if}
                        </div>
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

<ConfirmDialog
    bind:isOpen={showConfirmDialog}
    message={confirmMessage}
    on:confirm={handleDeleteConfirm}
    on:cancel={handleDeleteCancel}
/>

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
