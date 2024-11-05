"use server";

import { redirect } from "next/navigation";

export async function redirectServerSide(location: string) {
  redirect(location);
}
