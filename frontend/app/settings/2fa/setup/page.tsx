import { getServerSession } from "next-auth";
import { generate2FAConfig } from "../2fa";
import { authOptions } from "@/app/api/auth/[...nextauth]/options";
import { redirect } from "next/navigation";
import Image from "next/image";
import TwoFactorSubmitForm from "./twoFactorSubmitForm";

export default async function Setup2FA() {
  const session = await getServerSession(authOptions);
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
        <div>Secret: {twoFactorConfig.secret}</div>
        <div>Scan the QR code below to setup 2FA</div>
        <Image src={qrCode} alt="2fa qr code" width={200} height={200} />
        <TwoFactorSubmitForm
          secret={twoFactorConfig.secret}
          url={validateUrl}
        />
      </div>
    );
  } else {
    redirect("/unauthorized");
  }
}
