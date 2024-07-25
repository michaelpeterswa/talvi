import Link from "next/link";
import { TwoFactorStatus } from "./2fa/2fa";

export default function TwoFactorStatusButton({
  twoFactorStatus,
}: {
  twoFactorStatus: TwoFactorStatus;
}) {
  if (twoFactorStatus.status === "enabled") {
    return <button className="btn">Disable</button>;
  } else if (twoFactorStatus.status === "disabled") {
    return <button className="btn">Enable</button>;
  } else {
    return (
      <Link className="btn" href="/settings/2fa/setup">
        Setup 2FA
      </Link>
    );
  }
}
