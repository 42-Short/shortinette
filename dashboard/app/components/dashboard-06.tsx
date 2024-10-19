import { Link } from "@remix-run/react";

import { LineChart, Package2, Settings, Users2 } from "lucide-react";

import {
	Card,
	CardDescription,
	CardHeader,
	CardTitle,
} from "components/ui/card";

import { Tabs, TabsContent } from "components/ui/tabs";
import {
	Tooltip,
	TooltipContent,
	TooltipTrigger,
	TooltipProvider,
} from "components/ui/tooltip";

import { TooltipChart } from "./ui/retention";
import { AverageScore } from "./ui/average-score";
import { CardContent } from "./ui/card";

export default function Dashboard() {
	return (
		<TooltipProvider>
			<div className="flex min-h-screen w-full flex-col bg-muted/40">
				<aside className="fixed inset-y-0 left-0 z-10 hidden w-14 flex-col border-r bg-background sm:flex">
					<nav className="flex flex-col items-center gap-4 px-2 sm:py-5">
						<Link
							to="/"
							className="group flex h-9 w-9 shrink-0 items-center justify-center gap-2 rounded-full bg-primary text-lg font-semibold text-primary-foreground md:h-8 md:w-8 md:text-base"
						>
							<Package2 className="h-4 w-4 transition-all group-hover:scale-110" />
							<span className="sr-only">Dashboard</span>
						</Link>
						<Tooltip>
							<TooltipTrigger asChild>
								<Link
									to="/participants"
									className="flex h-9 w-9 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:text-foreground md:h-8 md:w-8"
								>
									<Users2 className="h-5 w-5" />
									<span className="sr-only">Participants</span>
								</Link>
							</TooltipTrigger>
							<TooltipContent side="right">Participants</TooltipContent>
						</Tooltip>
						<Tooltip>
							<TooltipTrigger asChild>
								<Link
									to="/analytics"
									className="flex h-9 w-9 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:text-foreground md:h-8 md:w-8"
								>
									<LineChart className="h-5 w-5" />
									<span className="sr-only">Analytics</span>
								</Link>
							</TooltipTrigger>
							<TooltipContent side="right">Analytics</TooltipContent>
						</Tooltip>
					</nav>
					<nav className="mt-auto flex flex-col items-center gap-4 px-2 sm:py-5">
						<Tooltip>
							<TooltipTrigger asChild>
								<Link
									to="/settings"
									className="flex h-9 w-9 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:text-foreground md:h-8 md:w-8"
								>
									<Settings className="h-5 w-5" />
									<span className="sr-only">Settings</span>
								</Link>
							</TooltipTrigger>
							<TooltipContent side="right">Settings</TooltipContent>
						</Tooltip>
					</nav>
				</aside>
				<div className="flex flex-col sm:gap-4 sm:py-4 sm:pl-14">
					<header className="sticky top-0 z-30 flex h-14 items-center gap-4 border-b bg-background px-4 sm:static sm:h-auto sm:border-0 sm:bg-transparent sm:px-6" />
					<main className="grid flex-1 items-start gap-4 p-4 sm:px-6 sm:py-0 md:gap-8">
						<Card>
							<CardHeader>
								<CardTitle>Dashboard</CardTitle>
								<CardDescription>Manage your Short here</CardDescription>
							</CardHeader>
							<CardContent>
								<div className="grid grid-cols-1 md:grid-cols-2 gap-4">
									<div className="col-span-1">
										<TooltipChart />
									</div>
									<div className="col-span-1">
										<AverageScore />
									</div>
								</div>
							</CardContent>
						</Card>
					</main>
				</div>
			</div>
		</TooltipProvider>
	);
}
