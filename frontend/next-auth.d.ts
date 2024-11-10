import { DefaultSession, DefaultUser } from "next-auth"
import { DefaultJWT } from "next-auth/jwt"

declare module "next-auth" {
    interface Session {
        user: {
            id: string,
            role: string,
            name: string,
            email: string,
            provider: string,
            image: string,
        } & DefaultSession
        token: string
    }

    interface User extends DefaultUser {
        role: string,
    }
}

declare module "next-auth/jwt" {
    interface JWT extends DefaultJWT {
        role: string,
        provider: string,
    }
}