import { LoaderFunction, json } from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";

import { Card, CardContent } from "./components/card";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "./components/carousel";

export const loader: LoaderFunction = async () => {
  return json({ items: Array.from({ length: 5 }) });
};

export default function Root() {
  const data = useLoaderData() as { items: any[] };

  return (
    <div>
      <h1>Welcome to My Remix App</h1>
      <Carousel
        opts={{
          align: "start",
        }}
        className="w-full max-w-sm"
      >
        <CarouselContent>
          {data.items.map((_, index: number) => (
            <CarouselItem key={index} className="md:basis-1/2 lg:basis-1/3">
              <div className="p-1">
                <Card>
                  <CardContent className="flex aspect-square items-center justify-center p-6">
                    <span className="text-3xl font-semibold">{index + 1}</span>
                  </CardContent>
                </Card>
              </div>
            </CarouselItem>
          ))}
        </CarouselContent>
        <CarouselPrevious />
        <CarouselNext />
      </Carousel>
    </div>
  );
}
