import type { Metadata } from "next";
import "./globals.css";
import { cn } from "@/lib/utils";
import { finlandica } from "./fonts";
import { auth } from "@/lib/auth/auth";
import AuthProvider from "@/lib/auth/authprovider";
import { ThemeProvider } from "@/components/theme/theme-provider";
import Navigation from "@/components/navigation/navigation";
import { Toaster } from "@/components/ui/toaster";
export const metadata: Metadata = {
  title: "talvi",
  description: "frontend for the talvi stack",
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const session = await auth();

  return (
    <html lang="en" suppressHydrationWarning={true}>
      <head>
        <meta name="viewport" content="initial-scale=1, width=device-width" />
      </head>
      <body className={cn("", finlandica.className)}>
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          <AuthProvider session={session}>
            <Navigation />
            {children}
            <Toaster />
          </AuthProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
