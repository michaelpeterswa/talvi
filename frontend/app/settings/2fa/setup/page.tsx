import { generate2FAConfig } from "../2fa";
import { redirect } from "next/navigation";
import Image from "next/image";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { InputOTPForm } from "./2fa-form";
import { auth } from "@/lib/auth/auth";

export default async function Setup2FA() {
  const session = await auth();
  if (session) {
    const twoFactorConfig = await generate2FAConfig(session);
    if (!twoFactorConfig) {
      return (
        <div className="flex flex-col items-center justify-center h-screen">
          failed to generate 2fa config
        </div>
      );
    }

    var qrCode = `data:image/png;base64,${twoFactorConfig.image}`;

    var validateUrl = `${process.env.BACKEND_HOST_URL}/api/v1/accounts/2fa`;

    return (
      <div className="flex flex-col items-center justify-center h-screen">
        <Card className="w-[350px]">
          <CardHeader>
            <CardTitle>2-Factor Auth Setup</CardTitle>
            <CardDescription>{twoFactorConfig.secret}</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid w-full items-center gap-4">
              <div className="flex flex-col space-y-1.5 items-center">
                <Label>QR Code:</Label>
                <Image
                  src={qrCode}
                  alt="2fa qr code"
                  width={200}
                  height={200}
                />
              </div>
              <div className="flex flex-col space-y-1.5">
                <InputOTPForm
                  secret={twoFactorConfig.secret}
                  url={validateUrl}
                />
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  } else {
    redirect("/unauthorized");
  }
}
