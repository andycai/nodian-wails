<script lang="ts">
  import { onMount } from "svelte";
  import MarkdownEditor from "./MarkdownEditor.svelte";
  import Calendar from "./Calendar.svelte";
  import TodoList from "./TodoList.svelte";
  import JsonFormatter from "./JsonFormatter.svelte";
  import TimestampConverter from "./TimestampConverter.svelte";
  import HashTool from "./HashTool.svelte";
  import EncodeTool from "./EncodeTool.svelte";

  let currentTool = "markdown";

  function setCurrentTool(tool: string) {
    currentTool = tool;
  }

  onMount(() => {
    // Initialize any necessary data or state here
  });
</script>

<main>
  <nav>
    <button on:click={() => setCurrentTool("markdown")}>Markdown Editor</button>
    <button on:click={() => setCurrentTool("calendar")}>Calendar & TODO</button>
    <button on:click={() => setCurrentTool("json")}>JSON Formatter</button>
    <button on:click={() => setCurrentTool("timestamp")}
      >Timestamp Converter</button
    >
    <button on:click={() => setCurrentTool("hash")}>Hash Tool</button>
    <button on:click={() => setCurrentTool("encode")}>Encode/Decode Tool</button
    >
  </nav>

  <div class="content">
    {#if currentTool === "markdown"}
      <MarkdownEditor />
    {:else if currentTool === "calendar"}
      <div class="calendar-todo">
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
</main>

<style>
  main {
    display: flex;
    flex-direction: column;
    height: 100vh;
  }

  nav {
    display: flex;
    justify-content: space-around;
    padding: 1rem;
    background-color: #f0f0f0;
  }

  .content {
    flex: 1;
    padding: 1rem;
    overflow-y: auto;
  }

  .calendar-todo {
    display: flex;
    gap: 1rem;
  }
</style>
