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
    const dispatch = createEventDispatcher();
    let editorElement: HTMLTextAreaElement;
    let previewElement: HTMLDivElement;

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

    function handleInput(event: Event) {
        const target = event.target as HTMLTextAreaElement;
        content = target.value;
    }

    onMount(() => {
        window.addEventListener("keydown", handleKeydown);
        return () => {
            window.removeEventListener("keydown", handleKeydown);
        };
    });

    $: renderedContent = md.render(content);

    function syncScroll() {
        if (previewElement) {
            const percentage =
                editorElement.scrollTop /
                (editorElement.scrollHeight - editorElement.clientHeight);
            previewElement.scrollTop =
                percentage *
                (previewElement.scrollHeight - previewElement.clientHeight);
        }
    }
</script>

<div class="markdown-editor h-full flex">
    <textarea
        bind:this={editorElement}
        bind:value={content}
        on:input={handleInput}
        on:scroll={syncScroll}
        class="editor flex-1 p-4 overflow-y-auto dark:bg-gray-800 dark:text-gray-100 custom-scrollbar resize-none"
    ></textarea>
    <div
        bind:this={previewElement}
        class="preview flex-1 p-4 overflow-y-auto dark:bg-gray-800 dark:text-gray-100 custom-scrollbar"
    >
        {@html renderedContent}
    </div>
</div>

<style>
    .markdown-editor {
        display: flex;
    }

    .editor,
    .preview {
        width: 50%;
        height: 100%;
    }

    .editor {
        font-family: monospace;
        white-space: pre-wrap;
        word-break: break-word;
        border: none;
        outline: none;
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
