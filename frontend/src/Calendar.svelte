<script lang="ts">
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();

    let currentDate = new Date();
    $: year = currentDate.getFullYear();
    $: month = currentDate.getMonth();

    function getDaysInMonth(year, month) {
        return new Date(year, month + 1, 0).getDate();
    }

    function getFirstDayOfMonth(year, month) {
        return new Date(year, month, 1).getDay();
    }

    function prevMonth() {
        currentDate = new Date(year, month - 1, 1);
    }

    function nextMonth() {
        currentDate = new Date(year, month + 1, 1);
    }

    function selectDate(day) {
        dispatch("selectDate", new Date(year, month, day));
    }

    $: daysInMonth = getDaysInMonth(year, month);
    $: firstDay = getFirstDayOfMonth(year, month);
</script>

<div class="calendar">
    <div class="header">
        <button on:click={prevMonth}>&lt;</button>
        <h2>
            {currentDate.toLocaleString("default", {
                month: "long",
                year: "numeric",
            })}
        </h2>
        <button on:click={nextMonth}>&gt;</button>
    </div>
    <div class="days">
        {#each ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"] as day}
            <div class="day-name">{day}</div>
        {/each}
        {#each Array(firstDay) as _}
            <div class="day empty"></div>
        {/each}
        {#each Array(daysInMonth) as _, i}
            <div class="day" on:click={() => selectDate(i + 1)}>{i + 1}</div>
        {/each}
    </div>
</div>

<style>
    .calendar {
        width: 300px;
    }

    .header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1rem;
    }

    .days {
        display: grid;
        grid-template-columns: repeat(7, 1fr);
        gap: 5px;
    }

    .day,
    .day-name {
        text-align: center;
        padding: 5px;
    }

    .day {
        cursor: pointer;
    }

    .day:hover {
        background-color: #f0f0f0;
    }

    .empty {
        visibility: hidden;
    }
</style>
