"use client";

import { Bar, BarChart, XAxis } from "recharts";

import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "~/components/ui/card";
import {
	type ChartConfig,
	ChartContainer,
	ChartTooltip,
	ChartTooltipContent,
} from "./chart";

const chartData = [
	{ date: "2024-07-15", mandatory: 450, bonus: 300 },
	{ date: "2024-07-16", mandatory: 380, bonus: 420 },
	{ date: "2024-07-17", mandatory: 520, bonus: 120 },
	{ date: "2024-07-18", mandatory: 140, bonus: 550 },
	{ date: "2024-07-19", mandatory: 600, bonus: 350 },
	{ date: "2024-07-20", mandatory: 480, bonus: 400 },
	{ date: "2024-07-21", mandatory: 480, bonus: 400 },
];

const chartConfig = {
	mandatory: {
		label: "Mandatory",
		color: "hsl(var(--chart-1))",
	},
	bonus: {
		label: "Bonus",
		color: "hsl(var(--chart-2))",
	},
} satisfies ChartConfig;

export function TooltipChart() {
	return (
		<Card>
			<CardHeader>
				<CardTitle>Retention</CardTitle>
				<CardDescription>Retention</CardDescription>
			</CardHeader>
			<CardContent>
				<ChartContainer config={chartConfig}>
					<BarChart accessibilityLayer data={chartData}>
						<XAxis
							dataKey="date"
							tickLine={false}
							tickMargin={10}
							axisLine={false}
							tickFormatter={(value) => {
								return new Date(value).toLocaleDateString("en-US", {
									weekday: "short",
								});
							}}
						/>
						<Bar
							dataKey="mandatory"
							stackId="a"
							fill="var(--color-mandatory)"
							radius={[0, 0, 4, 4]}
						/>
						<Bar
							dataKey="bonus"
							stackId="a"
							fill="var(--color-bonus)"
							radius={[4, 4, 0, 0]}
						/>
						<ChartTooltip
							content={<ChartTooltipContent />}
							cursor={false}
							defaultIndex={1}
						/>
					</BarChart>
				</ChartContainer>
			</CardContent>
		</Card>
	);
}
