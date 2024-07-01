import { useState } from "react";
import { Button } from "./ui/button";
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "./ui/sheet";
import {
  MenuIcon,
  Home,
  Info,
  TriangleAlert,
  MessageCircle,
} from "lucide-react";
import { NavLink } from "@remix-run/react";
import classNames from "classnames";

const navItems = [
  {
    label: "Home",
    href: "/",
    icon: Home,
  },
  {
    label: "My Team",
    href: "/about",
    icon: Info,
  },
];

function MainNav() {
  return (
    <div className="mr-4 hidden gap-2 md:flex">
      {navItems.map((item, index) => (
        <NavLink
          key={item.href}
          to={item.href}
          className={({ isActive, isPending }) => {
            return classNames(
              "text-primary underline-offset-4 hover:underline",
              "px-2 py-2",
              {
                "text-violet-600 underline font-bold": isActive,
              },
            );
          }}
        >
          {item.label}
        </NavLink>
      ))}
    </div>
  );
}

function MobileNav() {
  const [open, setOpen] = useState(false);

  return (
    <div className="md:hidden w-full">
      <Sheet open={open} onOpenChange={setOpen}>
        <div className="flex flex-row items-center justify-between w-full">
          <SheetTrigger asChild>
            <Button variant="ghost" size="icon" className="md:hidden">
              <MenuIcon />
            </Button>
          </SheetTrigger>
          <div className="font-bold uppercase text-violet-600">
            Student Council
          </div>
        </div>

        <SheetContent side="left">
          <div className="flex flex-col items-start">
            <div className="mb-6 font-bold uppercase text-gray-600">
              Student Council
            </div>
            {navItems.map((item, index) => (
              <NavLink
                key={index}
                to={item.href}
                className={({ isActive, isPending }) => {
                  return classNames("mb-4 flex flex-row items-center", {
                    "text-violet-600 font-bold": isActive,
                  });
                }}
                onClick={() => {
                  setOpen(false);
                }}
              >
                <item.icon className="size-5 mr-4" /> <p>{item.label}</p>
              </NavLink>
            ))}
          </div>
        </SheetContent>
      </Sheet>
    </div>
  );
}

export default function NavBar() {
  return (
    <header className="w-full border-b">
      <div className="flex h-14 items-center px-4">
        <MainNav />
        <MobileNav />
      </div>
    </header>
  );
}
