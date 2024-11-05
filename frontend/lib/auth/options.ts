import { NextAuthConfig } from 'next-auth';
import { JWT } from 'next-auth/jwt';
import { base64url, EncryptJWT } from "jose";
import { OAuthUserConfig, Provider } from 'next-auth/providers';
import GitHub, { GitHubProfile } from 'next-auth/providers/github';
import { cookies } from 'next/headers';

export const GitHubConfig: OAuthUserConfig<GitHubProfile> = {
    profile(profile) {
        return {
            id: profile.id.toString(),
            name: profile.name,
            email: profile.email,
            image: profile.avatar_url,
            role: "user",
        }
    }

}

export const providers: Provider[] = [
    GitHub(GitHubConfig)
];

export const authOptions: NextAuthConfig = {
    providers: providers,
    pages: {
        signIn: "/login",
      },
    callbacks: {
        async signIn({ account, profile }) {
          if (cookies().has('two_factor_verified')) {
            return true
          }

          let response = await fetch(
            `${process.env.BACKEND_HOST_URL}/api/v1/accounts/2fa/2fa?email=${profile?.email}&provider=${account?.provider}`,
            {
              method: "GET",
              headers: {
                "Content-Type": "application/json",
              },
            }
          );
  
          
          if (response.status === 200) {
            let body = await response.json();
            if (body.status === "enabled") {
              if (profile?.email && account?.provider) {
                return "/login/2fa?" + new URLSearchParams({ email: profile.email, provider: account.provider }).toString();
              }
            }
          }

          return true;
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
  const secret = base64url.decode(process.env.AUTH_SECRET ?? "")
  const jwt = await new EncryptJWT(token)
      .setProtectedHeader({ alg: 'dir', enc: 'A128CBC-HS256' })
      .setIssuedAt()
      .setExpirationTime('30d')
      .encrypt(secret)

  return jwt
}

async function getAccountOrCreate(token: JWT) {
  const jwe = await generateRequestJWE(token)
  if (jwe) {
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
}

async function createAccount(token: JWT , jwe: string) {
  const res = await fetch(`${process.env.BACKEND_HOST_URL}/api/v1/accounts/account?email=${token.email}&provider=${token.provider}`, {
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