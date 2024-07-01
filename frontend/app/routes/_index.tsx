import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "~/components/ui/carousel";
import { Card, CardContent } from "~/components/ui/card";
import { H1 } from "~/components/ui/H1";

export async function loader({ request }: LoaderFunctionArgs) {
  return null;
}

export default function Index() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen">
      <H1>Rust Piscine</H1>
      <div className="w-full max-w-xl md:max-w-3xl h-96">
        <Carousel>
          <CarouselContent>
            {Array.from({ length: 5 }).map((_, index) => (
              <CarouselItem key={index}>
                <div className="p-1">
                  <Card>
                    <div className="px-4 py-2 bg-gray-200">
                      <h3 className="text-lg font-semibold">Day {index + 1}</h3>
                    </div>
                    <CardContent className="flex aspect-square items-center justify-center p-6">
                      <span className="text-4xl font-semibold">
                        {index + 1}
                      </span>
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
    </div>
  );
}
