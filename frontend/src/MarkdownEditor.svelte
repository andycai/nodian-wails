<script lang="ts">
    import { onMount, createEventDispatcher } from "svelte";
    import { ReadMarkdownFile, SaveMarkdownFile } from "../wailsjs/go/main/App";
    import MarkdownIt from "markdown-it";
    import { writable } from "svelte/store";

    export let selectedFile: string | null;
    export let fileContents: { [key: string]: string };

    let content = "";
    let md = new MarkdownIt();
    let previousContent = "";
    const showPreview = writable(true);
    const dispatch = createEventDispatcher();

    $: if (selectedFile) {
        loadFile(selectedFile);
    }

    async function loadFile(file: string) {
        if (fileContents[file]) {
            content = fileContents[file];
        } else {
            content = await ReadMarkdownFile(file);
            fileContents[file] = content;
        }
        previousContent = content;
    }

    $: {
        if (selectedFile && content !== previousContent) {
            dispatch("markModified", [selectedFile, content]);
        }
    }

    function togglePreview() {
        showPreview.update((value) => {
            const newValue = !value;
            localStorage.setItem("showPreview", JSON.stringify(newValue));
            return newValue;
        });
    }

    function handleKeydown(event: KeyboardEvent) {
        // 支持 Windows 的 Ctrl+S 和 macOS 的 Cmd+S
        if ((event.ctrlKey || event.metaKey) && event.key === "s") {
            event.preventDefault();
            saveFile();
        }
    }

    async function saveFile() {
        if (selectedFile) {
            await SaveMarkdownFile(selectedFile, content);
            previousContent = content;
            dispatch("markSaved", selectedFile);
        }
    }

    onMount(() => {
        const storedShowPreview = localStorage.getItem("showPreview");
        if (storedShowPreview !== null) {
            showPreview.set(JSON.parse(storedShowPreview));
        }
        window.addEventListener("keydown", handleKeydown);
        return () => {
            window.removeEventListener("keydown", handleKeydown);
        };
    });

    $: renderedContent = md.render(content);
</script>

<div class="markdown-editor h-full flex flex-col">
    <div class="toolbar flex justify-end">
        <button
            on:click={togglePreview}
            class="rounded hover:bg-gray-200 dark:hover:bg-gray-700"
        >
            <svg
                class="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 22 22"
                xmlns="http://www.w3.org/2000/svg"
            >
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d={$showPreview
                        ? "M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                        : "M15 12a3 3 0 11-6 0 3 3 0 016 0z M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"}
                ></path>
            </svg>
        </button>
    </div>
    <div class="flex-1 flex">
        <div class="editor flex-1 px-2">
            <textarea
                bind:value={content}
                class="w-full h-full p-2 border rounded resize-none dark:bg-gray-800 dark:text-gray-100 custom-scrollbar"
            ></textarea>
        </div>
        {#if $showPreview}
            <div class="preview flex-1 px-2 overflow-y-auto dark:text-gray-100">
                {@html renderedContent}
            </div>
        {/if}
    </div>
</div>

<style>
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
