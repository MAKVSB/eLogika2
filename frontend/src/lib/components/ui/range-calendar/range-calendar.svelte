<script lang="ts">
	import { RangeCalendar as RangeCalendarPrimitive } from "bits-ui";
	import * as RangeCalendar from "./index.js";
	import { cn, type WithoutChildrenOrChild } from "$lib/utils.js";
	import type { ButtonVariant } from "$lib/components/ui/button/index.js";
	import type { Snippet } from "svelte";
	import { isEqualMonth, type DateValue, ZonedDateTime } from "@internationalized/date";
	import RangeCalendarTime from "./range-calendar-time.svelte";

	export type Granularity = 'day' | 'hour' | 'minute' | 'second';

	let {
		ref = $bindable(null),
		value = $bindable(),
		placeholder = $bindable(),
		weekdayFormat = "short",
		class: className,
		buttonVariant = "ghost",
		captionLayout = "label",
		locale = "en-US",
		months: monthsProp,
		years,
		monthFormat: monthFormatProp,
		yearFormat = "numeric",
		day,
		disableDaysOutsideMonth = false,
		granularity = "day",
		hideTimeZone = false,
		disabled,
		...restProps
	}: WithoutChildrenOrChild<RangeCalendarPrimitive.RootProps> & {
		value: {
			start: ZonedDateTime,
			end: ZonedDateTime
		};
		buttonVariant?: ButtonVariant;
		captionLayout?: "dropdown" | "dropdown-months" | "dropdown-years" | "label";
		months?: RangeCalendarPrimitive.MonthSelectProps["months"];
		years?: RangeCalendarPrimitive.YearSelectProps["years"];
		monthFormat?: RangeCalendarPrimitive.MonthSelectProps["monthFormat"];
		yearFormat?: RangeCalendarPrimitive.YearSelectProps["yearFormat"];
		granularity?: Granularity
		hideTimeZone?: boolean;
		day?: Snippet<[{ day: DateValue; outsideMonth: boolean }]>;
	} = $props();

	const monthFormat = $derived.by(() => {
		if (monthFormatProp) return monthFormatProp;
		if (captionLayout.startsWith("dropdown")) return "long";
		return "long";
	});
</script>

<RangeCalendarPrimitive.Root
	bind:ref
	bind:value
	bind:placeholder
	{weekdayFormat}
	{disableDaysOutsideMonth}
	class={cn(
		"bg-background group/calendar p-3 [--cell-size:--spacing(8)] [[data-slot=card-content]_&]:bg-transparent [[data-slot=popover-content]_&]:bg-transparent",
		className
	)}
	{locale}
	{monthFormat}
	{yearFormat}
	{disabled}
	{...restProps}
>
	{#snippet children({ months, weekdays })}
		<RangeCalendar.Months>
			<RangeCalendar.Nav>
				<RangeCalendar.PrevButton variant={buttonVariant} />
				<RangeCalendar.NextButton variant={buttonVariant} />
			</RangeCalendar.Nav>
			{#each months as month, monthIndex (month)}
				<RangeCalendar.Month>
					<RangeCalendar.Header>
						<RangeCalendar.Caption
							{captionLayout}
							months={monthsProp}
							{monthFormat}
							{years}
							{yearFormat}
							month={month.value}
							bind:placeholder
							{locale}
							{monthIndex}
						/>
					</RangeCalendar.Header>

					<RangeCalendar.Grid>
						<RangeCalendar.GridHead>
							<RangeCalendar.GridRow class="select-none">
								{#each weekdays as weekday (weekday)}
									<RangeCalendar.HeadCell>
										{weekday.slice(0, 2)}
									</RangeCalendar.HeadCell>
								{/each}
							</RangeCalendar.GridRow>
						</RangeCalendar.GridHead>
						<RangeCalendar.GridBody>
							{#each month.weeks as weekDates (weekDates)}
								<RangeCalendar.GridRow class="w-full mt-2">
									{#each weekDates as date (date)}
										<RangeCalendar.Cell {date} month={month.value}>
											{#if day}
												{@render day({
													day: date,
													outsideMonth: !isEqualMonth(date, month.value),
												})}
											{:else}
												<RangeCalendar.Day />
											{/if}
										</RangeCalendar.Cell>
									{/each}
								</RangeCalendar.GridRow>
							{/each}
						</RangeCalendar.GridBody>
					</RangeCalendar.Grid>
				</RangeCalendar.Month>
			{/each}
		</RangeCalendar.Months>
		{#if granularity != "day"}
			<RangeCalendarTime {disabled} bind:value {locale} {hideTimeZone}></RangeCalendarTime>
		{/if}
	{/snippet}
</RangeCalendarPrimitive.Root>
