import { AuthOptions, Session, getServerSession } from "next-auth";
import GitHubProvider from "next-auth/providers/github";
import { GithubProfile } from 'next-auth/providers/github'
import { base64url, EncryptJWT } from "jose";
import { JWT } from "next-auth/jwt";


export const authOptions: AuthOptions = {
    providers: [
      GitHubProvider({
        profile(profile: GithubProfile) {
            return {
                ...profile,
                role: profile.role ?? "user",
                id: profile.id.toString(),
                image: profile.avatar_url,
            }
        },
        clientId: process.env.GITHUB_ID ?? "",
        clientSecret: process.env.GITHUB_SECRET ?? "",
      }),
    ],
    callbacks: {
      async signIn({ account, profile }) {
        return true
        // if (account?.provider === 'github' &&
        //   profile?.email?.endsWith(`@${process.env.NEXTAUTH_GITHUB_DOMAIN_RESTRICTION}`)) {
        //   return true
        // } else {
        //   return '/unauthorized'
        // }
      },
      async jwt({ token, account, user}) {
        // Persist the OAuth access_token and or the user id to the token right after signin
        if (account) {
          token.id = account.providerAccountId
          token.provider = account.provider
        }

        if (user) token.role = user.role

        getAccountOrCreate(token)
  
        return token
      },
      async session({ session, token, user }) {
        session.token = await generateRequestJWE(token)
        session.user.provider = token.provider
        return session
      },
    },
  };

  async function generateRequestJWE(token: JWT) {
    const secret = base64url.decode(process.env.NEXTAUTH_SECRET ?? "")
    const jwt = await new EncryptJWT(token)
        .setProtectedHeader({ alg: 'dir', enc: 'A128CBC-HS256' })
        .setIssuedAt()
        .setExpirationTime('30d')
        .encrypt(secret)
    
    return jwt
  }

  async function getAccountOrCreate(token: JWT) {
    const jwe = await generateRequestJWE(token)
    if (token) {
      const res = await fetch(`${process.env.BACKEND_HOST_URL}/api/v1/accounts/account?email=${token.email}&provider=${token.provider}`, {
        method: 'GET',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${jwe}`,
        }});
      if (res.status === 404) {
        await createAccount(token, jwe)
      } else if (res.status != 200){
        return false
      } else {
        return true
      }}
    else {
      return false
    }
  };

  async function createAccount(token: JWT , jwe: string) {
    const res = await fetch(`${process.env.BACKEND_HOST_URL}/api/v1/accounts/account`, {
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${jwe}`,
      },
      body: JSON.stringify({
        name: token.name,
        email: token.email,
        role: token.role,
        provider: token.provider,
    })});
    if (res.status === 201) {
      return true
    } else {
      return false
    }
  }