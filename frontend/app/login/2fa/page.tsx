import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { InputOTPForm } from "./2fa-form";

export default function TwoFactorVerification() {
  var validateUrl = `${process.env.BACKEND_HOST_URL}/api/v1/accounts/2fa`;

  return (
    <div className="flex flex-col items-center justify-center min-h-[calc(100vh-80px)]">
      <Card className="w-[350px]">
        <CardHeader>
          <CardTitle>2-Factor Login</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid w-full items-center gap-4">
            <div className="flex flex-col space-y-1.5">
              <InputOTPForm url={validateUrl} />
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
