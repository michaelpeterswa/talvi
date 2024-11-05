"use client";

import { Dispatch, SetStateAction } from "react";

export default function InputField({
  code,
  setCode,
}: {
  code: string;
  setCode: Dispatch<SetStateAction<string>>;
}) {
  return (
    <input
      type="text"
      placeholder="Type here"
      className="input input-bordered w-full max-w-xs"
      value={code}
      onChange={(event) => setCode(event.target.value)}
    />
  );
}
