"use server";

import { cookies } from "next/headers";

export async function setTwoFactorVerifiedCookie() {
    cookies().set("two_factor_verified", "true", { secure: true, maxAge: 10*60 });
}

export async function getTwoFactorVerifiedCookie(): Promise<string> {
    let verified = cookies().get("two_factor_verified");
    if (verified) {
        return verified.value;
    }
    else {
        return "";
    }
}

export async function removeTwoFactorVerifiedCookie() {
    cookies().delete("two_factor_verified");
}