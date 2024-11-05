"use client";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
} from "@/components/ui/input-otp";
import { revalidatePath } from "next/cache";
import { Session } from "next-auth";
import { redirectServerSide } from "@/lib/redirect/redirect";
import { useSession } from "next-auth/react";

const FormSchema = z.object({
  pin: z.string().min(6, {
    message: "Your one-time password must be 6 characters.",
  }),
});

export function InputOTPForm({ url, secret }: { url: string; secret: string }) {
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      pin: "",
    },
  });

  let session = useSession();

  function onSubmit(data: z.infer<typeof FormSchema>) {
    getCodeVerification(session.data, url, secret, data.pin);
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="w-2/3 space-y-6">
        <FormField
          control={form.control}
          name="pin"
          /* @ts-ignore */
          render={({ field }) => (
            <FormItem>
              <FormLabel>One-Time Password</FormLabel>
              <FormControl>
                <InputOTP maxLength={6} {...field}>
                  <InputOTPGroup>
                    <InputOTPSlot index={0} />
                    <InputOTPSlot index={1} />
                    <InputOTPSlot index={2} />
                    <InputOTPSlot index={3} />
                    <InputOTPSlot index={4} />
                    <InputOTPSlot index={5} />
                  </InputOTPGroup>
                </InputOTP>
              </FormControl>
              <FormDescription>
                Please enter the one-time password sent to your phone.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type="submit">Submit</Button>
      </form>
    </Form>
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
          redirectServerSide("/settings");
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
