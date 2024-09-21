<script lang="ts">
    import { HashString } from "../wailsjs/go/main/App";

    let input = "";
    let algorithm = "md5";
    let output = "";

    const algorithms = [
        "CRC-32",
        "MD5",
        "SHA1",
        "SHA256",
        "SHA512",
        "SHA3-256",
        "SHA3-512",
    ];

    async function hash() {
        output = await HashString(input, algorithm);
    }
</script>

<div class="hash-tool">
    <textarea bind:value={input} placeholder="Enter text to hash"></textarea>
    <div>
        <select bind:value={algorithm}>
            {#each algorithms as algo}
                <option value={algo.toLowerCase()}>{algo}</option>
            {/each}
        </select>
        <button on:click={hash}>Hash</button>
    </div>
    <input readonly value={output} />
</div>

<style>
    .hash-tool {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    textarea {
        width: 100%;
        height: 100px;
    }
</style>
