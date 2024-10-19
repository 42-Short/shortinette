"use client";

import { useLoaderData } from "@remix-run/react";

export const description = "A collection of health charts.";

export async function loader() {
	return null;
}

import Dashboard from "~/components/dashboard-06";

export default function Index() {
	const data = useLoaderData<typeof loader>();
	return <Dashboard />;
}
