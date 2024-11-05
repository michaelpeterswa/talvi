import { redirect } from "next/navigation";
import { getAccount } from "./accounts";
import { get2FA } from "./2fa/2fa";
import TwoFactorStatusButton from "./2faStatusButton";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { auth } from "@/lib/auth/auth";

export default async function Settings() {
  const session = await auth();
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
      <div className="flex flex-col items-center justify-center min-h-[calc(100vh-80px)]">
        <Card className="w-[350px]">
          <CardHeader>
            <CardTitle>{account.name}</CardTitle>
            <CardDescription>{account.email}</CardDescription>
          </CardHeader>
          <CardContent>
            <form>
              <div className="grid w-full items-center gap-4">
                <div className="flex flex-col space-y-1.5">
                  <Label htmlFor="name">Created At:</Label>
                  <text>{}</text>
                </div>
                <div className="flex flex-col space-y-1.5">
                  <Label htmlFor="framework">Role</Label>
                  <Select>
                    <SelectTrigger id="framework">
                      <SelectValue placeholder={account.role} />
                    </SelectTrigger>
                    <SelectContent position="popper">
                      <SelectItem value="next">{account.role}</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
            </form>
          </CardContent>
          <CardFooter className="flex justify-between">
            <Button variant="outline">Cancel</Button>
            <Button>
              <TwoFactorStatusButton twoFactorStatus={twoFactor} />
            </Button>
            <Button>Save</Button>
          </CardFooter>
        </Card>
      </div>
    );
  } else {
    redirect("unauthorized");
  }
}
