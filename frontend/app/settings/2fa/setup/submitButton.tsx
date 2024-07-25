"use client";

import { write } from "fs";
import { Session } from "next-auth";
import { useSession } from "next-auth/react";
import { revalidatePath } from "next/cache";

export default function SubmitButton({
  url,
  secret,
  code,
}: {
  url: string;
  secret: string;
  code: string;
}) {
  let session = useSession();

  return (
    <button
      className="btn"
      onClick={() => {
        getCodeVerification(session.data, url, secret, code);
      }}
    >
      Submit
    </button>
  );
}

function getCodeVerification(
  sessionData: Session | null,
  url: string,
  secret: string,
  code: string
) {
  let response = fetch(
    `${url}/validate?email=${sessionData?.user.email}&provider=${sessionData?.user.provider}&secret=${secret}&code=${code}`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${sessionData?.token}`,
      },
    }
  );

  response.then((res) => {
    if (res.status === 200) {
      res.json().then((body) => {
        if (body.verified) {
          writeSecret(sessionData, url, secret);
        } else {
          alert("Invalid code");
        }
      });
    }
  });
}

function writeSecret(sessionData: Session | null, url: string, secret: string) {
  let response = fetch(
    `${url}/2fa?email=${sessionData?.user.email}&provider=${sessionData?.user.provider}`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${sessionData?.token}`,
      },

      body: JSON.stringify({
        secret: secret,
      }),
    }
  );

  response.then((res) => {
    if (res.status === 200) {
      res.json().then((body) => {
        if (body.success) {
          revalidatePath("/settings");
        }
      });
    }
  });
}
