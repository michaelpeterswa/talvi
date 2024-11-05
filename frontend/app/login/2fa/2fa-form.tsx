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
import { useSearchParams } from "next/navigation";
import { setTwoFactorVerifiedCookie } from "./cookie";
import { SignInWithProvider } from "@/lib/auth/action";
import { toast, useToast } from "@/components/ui/use-toast";

const FormSchema = z.object({
  pin: z.string().min(6, {
    message: "Your one-time password must be 6 characters.",
  }),
});

export function InputOTPForm({ url }: { url: string }) {
  const { toast } = useToast();
  const searchParams = useSearchParams();
  const email = searchParams.get("email");
  const provider = searchParams.get("provider");

  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      pin: "",
    },
  });

  function onSubmit(data: z.infer<typeof FormSchema>) {
    if (email && provider) {
      getCodeVerification(url, email, provider, data.pin);
    }
  }

  return (
    <div>
      <Form {...form}>
        <form
          onSubmit={form.handleSubmit(onSubmit)}
          className="w-2/3 space-y-6"
        >
          <FormField
            control={form.control}
            name="pin"
            /* @ts-ignore */
            render={({ field }) => (
              <FormItem>
                <FormLabel>One-Time Password</FormLabel>
                <FormControl>
                  <InputOTP maxLength={6} {...field} data-1p-ignore>
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
    </div>
  );
}

function getCodeVerification(
  url: string,
  email: string,
  provider: string,
  code: string
) {
  let response = fetch(
    `${url}/verify?email=${email}&provider=${provider}&code=${code}`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    }
  );

  response.then((res) => {
    if (res.status === 200) {
      res.json().then((body) => {
        console.log(body);
        if (body.verified) {
          setTwoFactorVerifiedCookie();
          SignInWithProvider(provider, "/");
        } else if (body.verified === false) {
          console.log("Code verification failed");
          toast({
            variant: "destructive",
            title: "Incorrect Two Factor Code",
            description: "Please try again.",
          });
        }
      });
    }
  });
}
