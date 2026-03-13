import type { Metadata } from "next";

export const metadata: Metadata = { title: "Forgot password" };

export default function ForgotPasswordPage() {
  return (
    <div className="text-center">
      <h1 className="text-2xl font-bold mb-2">Forgot password</h1>
      <p className="text-muted-foreground">Password reset will be implemented here.</p>
    </div>
  );
}
