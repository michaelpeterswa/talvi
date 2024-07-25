"use client";

import { useState } from "react";
import InputField from "./inputField";
import SubmitButton from "./submitButton";

export default function TwoFactorSubmitForm({
  url,
  secret,
}: {
  url: string;
  secret: string;
}) {
  const [twoFactorCode, setTwoFactorCode] = useState("");

  return (
    <div>
      <InputField code={twoFactorCode} setCode={setTwoFactorCode} />
      <SubmitButton secret={secret} code={twoFactorCode} url={url} />
    </div>
  );
}
