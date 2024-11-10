import { Session } from "next-auth";
import { redirect } from "next/navigation";

export type TwoFactorStatus = {
  status: string;
};

export async function get2FA(session: Session): Promise<TwoFactorStatus> {
  if (!session) {
    redirect("/unauthorized");
  }

  const res = await fetch(
    `${process.env.BACKEND_HOST_URL}/api/v1/accounts/2fa/2fa?email=${session.user.email}&provider=${session.user.provider}`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session.token}`,
      },
    }
  );

  if (res.status === 404) {
    throw new Error("not found");
  } else if (res.status != 200) {
    throw new Error("failed to get account");
  }

  const body = res.json();

  return body as Promise<TwoFactorStatus>;
}

export type GeneratedTwoFactorConfig = {
  secret: string;
  image: string;
};

export async function generate2FAConfig(
  session: Session
): Promise<GeneratedTwoFactorConfig> {
  if (!session) {
    redirect("/unauthorized");
  }

  const res = await fetch(
    `${process.env.BACKEND_HOST_URL}/api/v1/accounts/2fa/generate?email=${session.user.email}&provider=${session.user.provider}`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session.token}`,
      },
    }
  );

  if (res.status != 200) {
    throw new Error("failed to generate 2fa");
  }

  const body = res.json();

  return body as Promise<GeneratedTwoFactorConfig>;
}
