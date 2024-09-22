<script lang="ts">
  import { onMount } from "svelte";
  import MarkdownEditor from "./MarkdownEditor.svelte";
  import Calendar from "./Calendar.svelte";
  import TodoList from "./TodoList.svelte";
  import JsonFormatter from "./JsonFormatter.svelte";
  import TimestampConverter from "./TimestampConverter.svelte";
  import HashTool from "./HashTool.svelte";
  import EncodeTool from "./EncodeTool.svelte";
  import FileExplorer from "./FileExplorer.svelte";
  import { ListMarkdownFiles } from "../wailsjs/go/main/App";

  let currentTool = "markdown";
  let isDarkMode = false;
  let files: string[] = [];
  let selectedFile: string | null = null;

  function setCurrentTool(tool: string) {
    currentTool = tool;
  }

  function toggleDarkMode() {
    isDarkMode = !isDarkMode;
    if (isDarkMode) {
      document.documentElement.classList.add("dark");
    } else {
      document.documentElement.classList.remove("dark");
    }
  }

  function handleFileSelect(event: CustomEvent<string>) {
    selectedFile = event.detail;
    currentTool = "markdown";
  }

  onMount(async () => {
    files = await ListMarkdownFiles("./assets");
  });
</script>

<main class={`flex h-screen ${isDarkMode ? "dark" : ""}`}>
  <nav
    class="w-16 bg-gray-100 dark:bg-gray-800 flex flex-col items-center py-4 space-y-4"
  >
    <button
      on:click={() => setCurrentTool("markdown")}
      class="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white"
    >
      <svg
        class="w-8 h-8"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
        ></path></svg
      >
    </button>
    <button
      on:click={() => setCurrentTool("calendar")}
      class="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white"
    >
      <svg
        class="w-8 h-8"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
        ></path></svg
      >
    </button>
    <button
      on:click={() => setCurrentTool("json")}
      class="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white"
    >
      <svg
        class="w-8 h-8"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M4 6h16M4 10h16M4 14h16M4 18h16"
        ></path></svg
      >
    </button>
    <button
      on:click={() => setCurrentTool("timestamp")}
      class="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white"
    >
      <svg
        class="w-8 h-8"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
        ></path></svg
      >
    </button>
    <button
      on:click={() => setCurrentTool("hash")}
      class="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white"
    >
      <svg
        class="w-8 h-8"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
        ></path></svg
      >
    </button>
    <button
      on:click={() => setCurrentTool("encode")}
      class="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white"
    >
      <svg
        class="w-8 h-8"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z"
        ></path></svg
      >
    </button>
    <button
      on:click={toggleDarkMode}
      class="mt-auto text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white"
    >
      {#if isDarkMode}
        <svg
          class="w-8 h-8"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"
          ></path></svg
        >
      {:else}
        <svg
          class="w-8 h-8"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
          ></path></svg
        >
      {/if}
    </button>
  </nav>

  <div class="flex-1 flex">
    <FileExplorer
      {files}
      {isDarkMode}
      on:selectFile={handleFileSelect}
      class="w-64 bg-gray-50 dark:bg-gray-700 p-4 overflow-y-auto"
    />
    <div class="flex-1 p-4 bg-white dark:bg-gray-900 overflow-y-auto">
      {#if currentTool === "markdown"}
        <MarkdownEditor {selectedFile} />
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
</main>

<style global>
  /* 全局样式 */
  :global(body) {
    background-color: white;
    color: #1a202c;
  }

  :global(.dark body) {
    background-color: #1a202c;
    color: #f7fafc;
  }

  :global(input),
  :global(textarea),
  :global(select) {
    background-color: white;
    color: #1a202c;
    border-color: #e2e8f0;
  }

  :global(.dark input),
  :global(.dark textarea),
  :global(.dark select) {
    background-color: #2d3748;
    color: #f7fafc;
    border-color: #4a5568;
  }

  :global(button) {
    background-color: #edf2f7;
    color: #1a202c;
  }

  :global(button:hover) {
    background-color: #e2e8f0;
  }

  :global(.dark button) {
    background-color: #4a5568;
    color: #f7fafc;
  }

  :global(.dark button:hover) {
    background-color: #2d3748;
  }
</style>
