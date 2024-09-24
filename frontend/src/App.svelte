<script lang="ts">
  import "./app.css";
  import { onMount } from "svelte";
  import MarkdownEditor from "./MarkdownEditor.svelte";
  import Calendar from "./Calendar.svelte";
  import TodoList from "./TodoList.svelte";
  import JsonFormatter from "./JsonFormatter.svelte";
  import TimestampConverter from "./TimestampConverter.svelte";
  import HashTool from "./HashTool.svelte";
  import EncodeTool from "./EncodeTool.svelte";
  import FileExplorer from "./FileExplorer.svelte";
  import {
    ListMarkdownFiles,
    ReadMarkdownFile,
    SaveMarkdownFile,
  } from "../wailsjs/go/main/App";
  import { truncateName } from "./utils";
  import { writable } from "svelte/store";

  let currentTool = "markdown";
  let isDarkMode = writable(true);
  let files: string[] = [];
  let openFiles: string[] = [];
  let selectedFile: string | null = null;
  let selectedFolder: string | null = null;
  let fileContents: { [key: string]: string } = {};
  let modifiedFiles: Set<string> = new Set();

  const MAX_OPEN_FILES = 10;

  function setCurrentTool(tool: string) {
    currentTool = tool;
  }

  function toggleDarkMode() {
    isDarkMode.update((value) => {
      const newValue = !value;
      localStorage.setItem("isDarkMode", JSON.stringify(newValue));
      if (newValue) {
        document.documentElement.classList.add("dark");
      } else {
        document.documentElement.classList.remove("dark");
      }
      return newValue;
    });
  }

  onMount(async () => {
    const storedDarkMode = localStorage.getItem("isDarkMode");
    if (storedDarkMode !== null) {
      const parsedDarkMode = JSON.parse(storedDarkMode);
      isDarkMode.set(parsedDarkMode);
      if (parsedDarkMode) {
        document.documentElement.classList.add("dark");
      } else {
        document.documentElement.classList.remove("dark");
      }
    }

    const storedOpenFiles = localStorage.getItem("openFiles");
    const storedSelectedFile = localStorage.getItem("selectedFile");
    if (storedOpenFiles !== null) {
      openFiles = JSON.parse(storedOpenFiles);
      for (const file of openFiles) {
        fileContents[file] = await ReadMarkdownFile(file);
      }
    }
    if (storedSelectedFile !== null) {
      selectedFile = storedSelectedFile;
    }

    // files = await ListMarkdownFiles("./nodian");
    // sortFiles();
  });

  function handleFileSelect(event: CustomEvent<string>) {
    const file = event.detail;
    if (!openFiles.includes(file)) {
      if (openFiles.length >= MAX_OPEN_FILES) {
        openFiles = [...openFiles.slice(1), file];
      } else {
        openFiles = [...openFiles, file];
      }
      localStorage.setItem("openFiles", JSON.stringify(openFiles));
    }
    selectedFile = file;
    localStorage.setItem("selectedFile", file);
    currentTool = "markdown";
  }

  function handlerFileSelectInTabs(file: string) {
    selectedFile = file;
    localStorage.setItem("selectedFile", file);
  }

  function handleFolderSelect(event: CustomEvent<string>) {
    selectedFolder = event.detail;
  }

  function closeFile(file: string) {
    openFiles = openFiles.filter((f) => f !== file);
    localStorage.setItem("openFiles", JSON.stringify(openFiles));
    if (selectedFile === file) {
      selectedFile = openFiles[openFiles.length - 1] || null;
      localStorage.setItem("selectedFile", selectedFile || "");
    }
    modifiedFiles.delete(file);
  }

  function sortFiles() {
    const sortedFiles: string[] = [];
    const fileTree: { [key: string]: string[] } = {};

    // 构建文件树
    files.forEach((file) => {
      const parts = file.split("/").filter(Boolean);
      let currentPath = "";
      for (let i = 0; i < parts.length; i++) {
        currentPath += (currentPath ? "/" : "") + parts[i];
        if (i === parts.length - 1) {
          const parentPath = currentPath.substring(
            0,
            currentPath.lastIndexOf("/"),
          );
          if (!fileTree[parentPath]) {
            fileTree[parentPath] = [];
          }
          fileTree[parentPath].push(currentPath);
        } else {
          if (!fileTree[currentPath]) {
            fileTree[currentPath] = [];
          }
        }
      }
    });

    // 递归排序和添加文件
    function addSortedFiles(path: string) {
      const items = fileTree[path] || [];
      const folders = items.filter((item) => fileTree[item]);
      const files = items.filter((item) => !fileTree[item]);

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

  function markFileAsModified(file: string) {
    modifiedFiles.add(file);
  }

  function markFileAsSaved(file: string) {
    modifiedFiles.delete(file);
  }
</script>

<main class="flex h-screen" class:dark={$isDarkMode}>
  <nav
    class="w-16 bg-gray-800 dark:bg-gray-900 flex flex-col items-center py-4 space-y-4"
  >
    <button
      on:click={() => setCurrentTool("markdown")}
      class="text-gray-400 hover:text-white"
    >
      <svg
        class="w-6 h-6"
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
    </button>
    <button
      on:click={() => setCurrentTool("calendar")}
      class="text-gray-400 hover:text-white"
    >
      <svg
        class="w-6 h-6"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
        ></path>
      </svg>
    </button>
    <button
      on:click={() => setCurrentTool("json")}
      class="text-gray-400 hover:text-white"
    >
      <svg
        class="w-6 h-6"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M4 6h16M4 10h16M4 14h16M4 18h16"
        ></path>
      </svg>
    </button>
    <button
      on:click={() => setCurrentTool("timestamp")}
      class="text-gray-400 hover:text-white"
    >
      <svg
        class="w-6 h-6"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M12 8v4l3 3m-6 3a9 9 0 11-18 0 9 9 0 0118 0z"
        ></path>
      </svg>
    </button>
    <button
      on:click={() => setCurrentTool("hash")}
      class="text-gray-400 hover:text-white"
    >
      <svg
        class="w-6 h-6"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
        ></path>
      </svg>
    </button>
    <button
      on:click={() => setCurrentTool("encode")}
      class="text-gray-400 hover:text-white"
    >
      <svg
        class="w-6 h-6"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
        ></path>
      </svg>
    </button>
    <button
      on:click={toggleDarkMode}
      class="mt-auto text-gray-400 hover:text-white"
    >
      {#if $isDarkMode}
        <svg
          class="w-6 h-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"
          ></path>
        </svg>
      {:else}
        <svg
          class="w-6 h-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
          ></path>
        </svg>
      {/if}
    </button>
  </nav>

  <div class="flex-1 flex">
    <FileExplorer
      {files}
      isDarkMode={$isDarkMode}
      on:selectFile={handleFileSelect}
      on:selectFolder={handleFolderSelect}
    />
    <div
      class="flex-1 flex flex-col bg-gray-100 dark:bg-gray-900 overflow-hidden"
    >
      {#if currentTool === "markdown"}
        <div class="flex bg-gray-200 dark:bg-gray-800 overflow-x-auto">
          {#each openFiles as file}
            <div
              class="px-4 py-2 flex items-center {selectedFile === file
                ? 'bg-white dark:bg-gray-700'
                : 'bg-gray-100 dark:bg-gray-800'} cursor-pointer"
              on:click={() => handlerFileSelectInTabs(file)}
            >
              <span class="mr-2 text-gray-900 dark:text-gray-100"
                >{truncateName(
                  file.split("/").pop() || "",
                  10,
                )}{modifiedFiles.has(file) ? "*" : ""}</span
              >
              <button
                on:click|stopPropagation={() => closeFile(file)}
                class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
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
                    d="M6 18L18 6M6 6l12 12"
                  ></path>
                </svg>
              </button>
            </div>
          {/each}
        </div>
      {/if}
      <div class="flex-1 p-4 overflow-y-auto">
        {#if currentTool === "markdown"}
          <MarkdownEditor
            {selectedFile}
            {fileContents}
            on:markModified={markFileAsModified}
            on:markSaved={markFileAsSaved}
          />
        {:else if currentTool === "calendar"}
          <div class="flex space-x-4">
            <Calendar />
            <TodoList />
          </div>
        {:else if currentTool === "json"}
          <JsonFormatter />
        {:else if currentTool === "timestamp"}
          <TimestampConverter />
        {:else if currentTool === "hash"}
          <HashTool />
        {:else if currentTool === "encode"}
          <EncodeTool />
        {/if}
      </div>
    </div>
  </div>
</main>

<style lang="postcss">
  @tailwind base;
  @tailwind components;
  @tailwind utilities;

  :global(body) {
    @apply bg-white text-gray-900 dark:bg-gray-900 dark:text-gray-100;
  }
</style>
