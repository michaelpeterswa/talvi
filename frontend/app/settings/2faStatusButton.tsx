"use client";
import Link from "next/link";
import { TwoFactorStatus } from "./2fa/2fa";
import { Session } from "next-auth";
import { useSession } from "next-auth/react";
import { redirectServerSide } from "@/lib/redirect/redirect";

export default function TwoFactorStatusButton({
  twoFactorStatus,
}: {
  twoFactorStatus: TwoFactorStatus;
}) {
  const { data: session } = useSession();
  if (!session) {
    return null;
  }

  if (twoFactorStatus.status === "enabled") {
    return (
      <button
        className="btn"
        onClick={() => {
          remove2FARequest(session);
        }}
      >
        Remove
      </button>
    );
  } else {
    return (
      <Link className="btn" href="/settings/2fa/setup">
        Setup 2FA
      </Link>
    );
  }
}

async function remove2FARequest(session: Session) {
  const res = await fetch(
    `${process.env.NEXT_PUBLIC_BACKEND_HOST_URL}/api/v1/accounts/2fa/2fa?email=${session.user.email}&provider=${session.user.provider}`,
    {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session.token}`,
      },
    }
  );
  if (res.status === 200) {
    redirectServerSide("/settings");
  } else {
    return false;
  }
}
