import type { Metadata } from "next";
import { Finlandica } from "next/font/google";
import { getServerSession } from "next-auth";
import SessionProvider from "../components/session/session-provider";
import "./globals.css";

const findlandica = Finlandica({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "talvi",
  description: "frontend for the talvi stack",
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const session = await getServerSession();

  return (
    <html lang="en">
      <head>
        <meta name="viewport" content="initial-scale=1, width=device-width" />
      </head>
      <body className={findlandica.className}>
        <SessionProvider session={session}>{children}</SessionProvider>
      </body>
    </html>
  );
}
