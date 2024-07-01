import { Link } from "@remix-run/react";
import { GitHub } from "./icon/GitHub";

export function Footer() {
  return (
    <footer
      className="bg-white rounded-lg shadow dark:bg-gray-800 fixed bottom-0 w-full"
      style={{
        height: "fit-content",
      }}
    >
      <div className="w-full mx-auto max-w-screen-xl p-4 flex flex-col sm:flex-row items-center justify-start sm:justify-between gap-x-7 gap-y-3">
        <Link
          to="https://github.com/42-Short/shortinette"
          target="_blank"
          className="hover:underline"
        >
          Found a bug? Open an issue! <GitHub className="size-5" />
        </Link>
      </div>
    </footer>
  );
}
