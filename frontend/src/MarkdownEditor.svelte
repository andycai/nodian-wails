<script lang="ts">
    import { onMount } from "svelte";
    import {
        CreateMarkdownFile,
        ReadMarkdownFile,
        SaveMarkdownFile,
    } from "../wailsjs/go/main/App";
    import MarkdownIt from "markdown-it";

    export let selectedFile: string | null;

    let content = "";
    let md = new MarkdownIt();

    $: if (selectedFile) {
        loadFile(selectedFile);
    }

    async function loadFile(file: string) {
        content = await ReadMarkdownFile(file);
    }

    async function saveFile() {
        if (selectedFile) {
            await SaveMarkdownFile(selectedFile, content);
        }
    }

    $: renderedContent = md.render(content);
</script>

<div class="markdown-editor">
    <div class="editor">
        <textarea
            bind:value={content}
            class="w-full h-64 p-2 border rounded dark:bg-gray-800 dark:text-white"
        ></textarea>
        <button
            on:click={saveFile}
            class="mt-2 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 dark:bg-blue-700 dark:hover:bg-blue-800"
            >Save</button
        >
    </div>
    <div class="preview">
        {@html renderedContent}
    </div>
</div>

<style>
    .markdown-editor {
        display: flex;
        height: 100%;
    }

    .editor,
    .preview {
        flex: 1;
        padding: 1rem;
    }

    .preview :global(h1) {
        @apply text-2xl font-bold mb-4;
    }

    .preview :global(h2) {
        @apply text-xl font-bold mb-3;
    }

    .preview :global(p) {
        @apply mb-2;
    }

    .preview :global(ul),
    .preview :global(ol) {
        @apply ml-4 mb-2;
    }

    .preview :global(li) {
        @apply mb-1;
    }
</style>
