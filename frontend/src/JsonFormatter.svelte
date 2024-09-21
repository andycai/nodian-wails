<script lang="ts">
    import { FormatJSON } from "../wailsjs/go/main/App";

    let input = "";
    let output = "";
    let compact = false;

    async function formatJSON() {
        try {
            output = await FormatJSON(input, compact);
        } catch (error) {
            output = `Error: ${error.message}`;
        }
    }
</script>

<div class="json-formatter">
    <textarea bind:value={input} placeholder="Enter JSON here"></textarea>
    <div>
        <button on:click={formatJSON}>Format</button>
        <label>
            <input type="checkbox" bind:checked={compact} />
            Compact
        </label>
    </div>
    <pre>{output}</pre>
</div>

<style>
    .json-formatter {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    textarea,
    pre {
        width: 100%;
        height: 200px;
        font-family: monospace;
    }
</style>
