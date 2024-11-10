import NextAuth from "next-auth";
import { authOptions, providers } from "@/lib/auth/options";

export const { handlers, signIn, signOut, auth } = NextAuth(authOptions);

export const providerMap = providers.map((provider, index) => {
    if (typeof provider === "function") {
      const providerData = provider()
      return { id: providerData.id, name: providerData.name, key: index }
    } else {
      return { id: provider.id, name: provider.name, key: index }
    }
  })