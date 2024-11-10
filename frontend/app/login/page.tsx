import { redirect } from "next/navigation";
import { providerMap } from "@/lib/auth/auth";
import { AuthError } from "next-auth";
import {
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Github } from "lucide-react";
import { SignInWithProvider } from "@/lib/auth/action";

export default async function SignInPage() {
  return (
    <div className="flex flex-col items-center justify-center min-h-[calc(100vh-80px)]">
      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle>talvi</CardTitle>
          <CardDescription>login</CardDescription>
        </CardHeader>
        <CardFooter>
          {Object.values(providerMap).map((provider) => (
            <div className="w-full" key={provider.id}>
              <form
                action={async () => {
                  "use server";
                  try {
                    await SignInWithProvider(provider.id, "/");
                  } catch (error) {
                    // Signin can fail for a number of reasons, such as the user
                    // not existing, or the user not having the correct role.
                    // In some cases, you may want to redirect to a custom error
                    if (error instanceof AuthError) {
                      return redirect(`/unauthorized`);
                    }

                    // Otherwise if a redirects happens NextJS can handle it
                    // so you can just re-thrown the error and let NextJS handle it.
                    // Docs:
                    // https://nextjs.org/docs/app/api-reference/functions/redirect#server-component
                    throw error;
                  }
                }}
              >
                <Button className="w-full">
                  <Github className="mr-2 h-4 w-4" /> GitHub
                </Button>
              </form>
            </div>
          ))}
        </CardFooter>
      </Card>
    </div>
  );
}
