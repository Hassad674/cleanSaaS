import type { Metadata } from "next";
import { AuthDemo } from "./auth-demo";

export const metadata: Metadata = {
  title: "Auth Demo — CleanSaaS",
  description:
    "Interactive demo of CleanSaaS authentication: login, register, forgot password, and OAuth. No account required.",
};

export default function AuthDemoPage() {
  return <AuthDemo />;
}
