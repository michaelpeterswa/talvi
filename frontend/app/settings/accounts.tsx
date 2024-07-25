import { Session } from "next-auth";

type Account = {
  id: string;
  name: string;
  email: string;
  role: string;
  provider: string;
  created_at: string;
};

export async function getAccount(session: Session): Promise<Account> {
  const res = await fetch(
    `${process.env.BACKEND_HOST_URL}/api/v1/accounts/account?email=${session.user.email}&provider=${session.user.provider}`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session.token}`,
      },
    }
  );

  if (res.status === 404) {
    throw new Error("not found");
  } else if (res.status != 200) {
    throw new Error("failed to get account");
  }

  const body = res.json();

  return body as Promise<Account>;
}
