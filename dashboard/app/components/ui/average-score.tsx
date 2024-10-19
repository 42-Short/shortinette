"use client";

import { TrendingUp } from "lucide-react";
import { Bar, BarChart, CartesianGrid, Rectangle, XAxis } from "recharts";

import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "~/components/ui/card";
import {
	type ChartConfig,
	ChartContainer,
	ChartTooltip,
	ChartTooltipContent,
} from "~/components/ui/chart";

export const description = "A bar chart with an active bar";

const chartData = [
	{ day: "Mon", averageScore: 50, fill: "var(--color-Mon)" },
	{ day: "Tue", averageScore: 20, fill: "var(--color-Tue)" },
	{ day: "Wed", averageScore: 30, fill: "var(--color-Wed)" },
	{ day: "Thu", averageScore: 60, fill: "var(--color-Thu)" },
	{ day: "Fri", averageScore: 60, fill: "var(--color-Fri)" },
	{ day: "Sat", averageScore: 20, fill: "var(--color-Sat)" },
	{ day: "Sun", averageScore: 70, fill: "var(--color-Sun)" },
];

const chartConfig = {
	averageScore: {
		label: "Average Score",
	},
	Mon: {
		label: "Mon",
		color: "hsl(var(--chart-1))",
	},
	Tue: {
		label: "Tue",
		color: "hsl(var(--chart-2))",
	},
	Wed: {
		label: "Wed",
		color: "hsl(var(--chart-3))",
	},
	Thu: {
		label: "Thu",
		color: "hsl(var(--chart-4))",
	},
	Fri: {
		label: "Fri",
		color: "hsl(var(--chart-5))",
	},
	Sat: {
		label: "Sat",
		color: "hsl(var(--chart-1))",
	},
	Sun: {
		label: "Sun",
		color: "hsl(var(--chart-2))",
	},
} satisfies ChartConfig;

export function AverageScore() {
	return (
		<Card>
			<CardHeader>
				<CardTitle>Average Score</CardTitle>
				<CardDescription>January - June 2024</CardDescription>
			</CardHeader>
			<CardContent>
				<ChartContainer config={chartConfig}>
					<BarChart accessibilityLayer data={chartData}>
						<CartesianGrid vertical={false} />
						<XAxis
							dataKey="day"
							tickLine={false}
							tickMargin={10}
							axisLine={false}
							tickFormatter={(value) =>
								chartConfig[value as keyof typeof chartConfig]?.label
							}
						/>
						<ChartTooltip
							cursor={false}
							content={<ChartTooltipContent hideLabel />}
						/>
						<Bar
							dataKey="averageScore"
							strokeWidth={2}
							radius={8}
							activeIndex={2}
							activeBar={({ ...props }) => {
								return (
									<Rectangle
										{...props}
										fillOpacity={0.8}
										stroke={props.payload.fill}
										strokeDasharray={4}
										strokeDashoffset={4}
									/>
								);
							}}
						/>
					</BarChart>
				</ChartContainer>
			</CardContent>
			<CardFooter className="flex-col items-start gap-2 text-sm">
				<div className="flex gap-2 font-medium leading-none">
					Trending up by 5.2% this month <TrendingUp className="h-4 w-4" />
				</div>
				<div className="leading-none text-muted-foreground">
					Showing total averageScore for the last 6 months
				</div>
			</CardFooter>
		</Card>
	);
}
