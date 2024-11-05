import { removeTwoFactorVerifiedCookie } from "@/app/login/2fa/cookie";
import { SignIn, SignOut } from "@/lib/auth/action";
import { Button } from "../ui/button";

export function SignInButton() {
  return (
    <form
      action={async () => {
        await SignIn();
      }}
    >
      <Button type="submit">sign in</Button>
    </form>
  );
}

export function SignOutButton() {
  return (
    <form
      action={async () => {
        await removeTwoFactorVerifiedCookie();
        await SignOut();
      }}
    >
      <Button type="submit">sign out</Button>
    </form>
  );
}
