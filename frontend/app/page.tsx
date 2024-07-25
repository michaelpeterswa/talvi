"use client";
import { signIn, signOut, useSession } from "next-auth/react";

export default function Home() {
  const { data: session } = useSession();

  return (
    <main>
      <div className="hero min-h-screen">
        <div className="hero-content text-center">
          <div className="max-w-xl">
            <h1 className="text-7xl">talvi</h1>
            <p className="py-6">
              Provident cupiditate voluptatem et in. Quaerat fugiat ut assumenda
              excepturi exercitationem quasi. In deleniti eaque aut repudiandae
              et a id nisi.
            </p>
            <button
              onClick={() => {
                signIn(undefined, { callbackUrl: "/" });
              }}
            >
              sign in
            </button>
            <br />
            {session?.user?.name?.toLowerCase()}
            <br />
            <button
              onClick={() => {
                signOut({ callbackUrl: "/" });
              }}
            >
              sign out
            </button>
          </div>
        </div>
      </div>
    </main>
  );
}
