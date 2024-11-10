"use client";
import { Sheet, SheetTrigger, SheetContent } from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import {
  NavigationMenu,
  NavigationMenuList,
  NavigationMenuLink,
} from "@/components/ui/navigation-menu";
import { Home, Menu } from "lucide-react";
import { useSession } from "next-auth/react";
import { ModeToggle } from "../theme/mode-toggle";
import { useEffect, useState } from "react";
import { SignInButton, SignOutButton } from "../sign-in/sign-in";
import { NavigationLink } from "./link";

export default function Navigation() {
  const { data: session } = useSession();
  const [hasSession, setHasSession] = useState(false);
  useEffect(() => {
    if (session) {
      setHasSession(true);
    }
  }, [session]);

  let sheetLinks: NavigationLink[] = [
    { id: 1, title: "Example Link 1", href: "#" },
    { id: 2, title: "Example Link 2", href: "#" },
    { id: 3, title: "Example Link 3", href: "#" },
    { id: 4, title: "Example Link 4", href: "#" },
  ];

  let barLinks: NavigationLink[] = [
    { id: 1, title: "Example Link 1", href: "#" },
    { id: 2, title: "Example Link 2", href: "#" },
    { id: 3, title: "Example Link 3", href: "#" },
  ];

  return (
    <header className="flex h-20 w-full shrink-0 items-center px-4 md:px-6">
      <Sheet>
        <SheetTrigger asChild>
          <Button variant="outline" size="icon" className="lg:hidden">
            <Menu />
            <span className="sr-only">Toggle navigation menu</span>
          </Button>
        </SheetTrigger>
        <SheetContent side="left">
          <Link href="/" className="flex items-center gap-2" prefetch={false}>
            <Home />
            <span className="text-lg font-semibold">talvi</span>
          </Link>
          <div className="grid gap-4 py-6">
            {sheetLinks.map((link, index) => (
              <Link
                key={index}
                href={link.href}
                className="flex w-full items-center py-2 text-lg font-semibold"
                prefetch={false}
              >
                {link.title}
              </Link>
            ))}
          </div>
        </SheetContent>
      </Sheet>
      <div className="w-[150px] hidden lg:flex">
        <Link href="/" className="flex items-center gap-2" prefetch={false}>
          <Home />
          <span className="text-lg font-semibold">talvi</span>
        </Link>
      </div>
      <div className="flex w-full justify-center">
        <NavigationMenu className="hidden lg:flex">
          <NavigationMenuList>
            {barLinks.map((link) => (
              <NavigationMenuLink key={link.id} asChild>
                <Link
                  href={link.href}
                  className="group inline-flex h-9 w-max items-center justify-center rounded-md bg-background px-4 py-2 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground focus:outline-none disabled:pointer-events-none disabled:opacity-50 data-[active]:bg-accent/50 data-[state=open]:bg-accent/50"
                  prefetch={false}
                >
                  {link.title}
                </Link>
              </NavigationMenuLink>
            ))}
          </NavigationMenuList>
        </NavigationMenu>
      </div>
      <div className="ml-auto pr-8">
        <ModeToggle />
      </div>
      {hasSession && (
        <div className="ml-auto px-1">
          <Link href="/settings">
            <Button>settings</Button>
          </Link>
        </div>
      )}
      <div className="ml-auto px-1">
        {hasSession ? SignOutButton() : SignInButton()}
      </div>
    </header>
  );
}
