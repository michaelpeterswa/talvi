"use client";
import { pacifico } from "./fonts";
import { cn } from "@/lib/utils";

export default function Home() {
  return (
    <main>
      <div className="flex flex-col items-center justify-center min-h-[calc(100vh-80px)]">
        <div className="text-center">
          <div className="max-w-xl">
            <h1
              className={cn(
                "text-8xl p-4 bg-gradient-to-bl from-fuchsia-500 to-cyan-500 overflow-visible	inline-block text-transparent bg-clip-text",
                pacifico.className
              )}
            >
              hello!
            </h1>
            <p className="py-6">
              Provident cupiditate voluptatem et in. Quaerat fugiat ut assumenda
              excepturi exercitationem quasi. In deleniti eaque aut repudiandae
              et a id nisi.
            </p>
          </div>
        </div>
      </div>
    </main>
  );
}
