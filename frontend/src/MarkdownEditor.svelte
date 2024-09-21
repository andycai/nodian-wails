<script lang="ts">
    import { onMount } from "svelte";
    import {
        CreateMarkdownFile,
        ReadMarkdownFile,
        SaveMarkdownFile,
        ListMarkdownFiles,
    } from "../wailsjs/go/main/App";
    import MarkdownIt from "markdown-it";

    let content = "";
    let filename = "";
    let files = [];
    let md = new MarkdownIt();

    onMount(async () => {
        files = await ListMarkdownFiles(".");
    });

    async function createFile() {
        if (filename) {
            await CreateMarkdownFile(filename, content);
            files = await ListMarkdownFiles(".");
        }
    }

    async function loadFile(file) {
        content = await ReadMarkdownFile(file);
        filename = file;
    }

    async function saveFile() {
        if (filename) {
            await SaveMarkdownFile(filename, content);
        }
    }

    $: renderedContent = md.render(content);
</script>

<div class="markdown-editor">
    <div class="file-list">
        <h3>Files</h3>
        <ul>
            {#each files as file}
                <li on:click={() => loadFile(file)}>{file}</li>
            {/each}
        </ul>
        <input bind:value={filename} placeholder="New file name" />
        <button on:click={createFile}>Create</button>
    </div>
    <div class="editor">
        <textarea bind:value={content}></textarea>
        <button on:click={saveFile}>Save</button>
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

    .file-list {
        width: 200px;
        border-right: 1px solid #ccc;
        padding: 1rem;
    }

    .editor,
    .preview {
        flex: 1;
        padding: 1rem;
    }

    textarea {
        width: 100%;
        height: 100%;
        resize: none;
    }
</style>
