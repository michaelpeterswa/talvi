"use client";
import { Session } from "next-auth";
import { SessionProvider } from "next-auth/react";

// @ts-ignore
export default function AuthProvider({
  session,
  children,
}: {
  session: Session | null;
  children: any;
}) {
  return <SessionProvider session={session}>{children}</SessionProvider>;
}
