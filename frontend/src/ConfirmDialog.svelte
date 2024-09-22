<script lang="ts">
    import { createEventDispatcher } from "svelte";

    export let isOpen = false;
    export let message = "";

    const dispatch = createEventDispatcher();

    function confirm() {
        dispatch("confirm");
        isOpen = false;
    }

    function cancel() {
        dispatch("cancel");
        isOpen = false;
    }
</script>

{#if isOpen}
    <div
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
        <div
            class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-xl max-w-sm w-full"
        >
            <h2 class="text-xl font-bold mb-4 text-gray-900 dark:text-gray-100">
                Confirm Action
            </h2>
            <p class="mb-6 text-gray-700 dark:text-gray-300">{message}</p>
            <div class="flex justify-end space-x-4">
                <button
                    on:click={cancel}
                    class="px-4 py-2 bg-gray-300 text-gray-800 rounded hover:bg-gray-400 transition-colors"
                >
                    Cancel
                </button>
                <button
                    on:click={confirm}
                    class="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600 transition-colors"
                >
                    Confirm
                </button>
            </div>
        </div>
    </div>
{/if}
