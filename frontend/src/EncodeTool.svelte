<script lang="ts">
    import { EncodeString, DecodeString } from "../wailsjs/go/main/App";

    let input = "";
    let output = "";
    let encoding = "base64";

    const encodings = ["Base64", "URL"];

    async function encode() {
        output = await EncodeString(input, encoding);
    }

    async function decode() {
        try {
            output = await DecodeString(input, encoding);
        } catch (error) {
            // output = `Error: ${error.message}`;
        }
    }
</script>

<div class="encode-tool">
    <textarea bind:value={input} placeholder="Enter text to encode/decode"
    ></textarea>
    <div>
        <select bind:value={encoding}>
            {#each encodings as enc}
                <option value={enc.toLowerCase()}>{enc}</option>
            {/each}
        </select>
        <button on:click={encode}>Encode</button>
        <button on:click={decode}>Decode</button>
    </div>
    <textarea readonly value={output}></textarea>
</div>

<style>
    .encode-tool {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    textarea {
        width: 100%;
        height: 100px;
    }
</style>
