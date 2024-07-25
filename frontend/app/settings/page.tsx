import { getServerSession } from "next-auth/next";
import { authOptions } from "../api/auth/[...nextauth]/options";
import { redirect } from "next/navigation";
import Image from "next/image";
import { getAccount } from "./accounts";
import { generate2FAConfig, get2FA } from "./2fa/2fa";
import TwoFactorStatusButton from "./2faStatusButton";

export default async function Settings() {
  const session = await getServerSession(authOptions);
  if (session) {
    const account = await getAccount(session);
    if (!account) {
      return (
        <div className="flex flex-col items-center justify-center h-screen">
          no profile
        </div>
      );
    }
    const twoFactor = await get2FA(session);
    if (!twoFactor) {
      return (
        <div className="flex flex-col items-center justify-center h-screen">
          no 2fa
        </div>
      );
    }

    return (
      <div className="flex flex-col items-center justify-center h-screen">
        <div>{account.id}</div>
        <div>{account.name}</div>
        <div>{account.created_at}</div>
        <div>{account.email}</div>
        <div>{account.provider}</div>
        <div>{account.role}</div>
        <TwoFactorStatusButton twoFactorStatus={twoFactor} />
      </div>
    );
  } else {
    redirect("unauthorized");
  }
}
