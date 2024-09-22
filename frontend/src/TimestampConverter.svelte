<script lang="ts">
    import { TimestampToDate, DateToTimestamp } from "../wailsjs/go/main/App";

    let timestamp: string | null = null;
    let date: string | null = null;
    let isMilliseconds: boolean = true;

    async function convertTimestamp() {
        if (timestamp) {
            date = await TimestampToDate(parseInt(timestamp), isMilliseconds);
        }
    }

    async function convertDate() {
        if (date) {
            timestamp = await DateToTimestamp(date, isMilliseconds).toString();
        }
    }
</script>

<div class="timestamp-converter">
    <div>
        <input type="number" bind:value={timestamp} placeholder="Timestamp" />
        <button on:click={convertTimestamp}>To Date</button>
    </div>
    <div>
        <input
            type="text"
            bind:value={date}
            placeholder="YYYY-MM-DD HH:MM:SS"
        />
        <button on:click={convertDate}>To Timestamp</button>
    </div>
    <label>
        <input type="checkbox" bind:checked={isMilliseconds} />
        Use milliseconds
    </label>
</div>

<style>
    .timestamp-converter {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }
</style>
