import { Suspense } from "react";
import { VerifyEmailForm } from "@/features/auth/components/verify-email-form";

export const metadata = {
  title: "Verify Email — CleanSaaS",
  description: "Verify your email address",
};

export default function VerifyEmailPage() {
  return (
    <Suspense fallback={<p className="text-center text-muted-foreground">Loading...</p>}>
      <VerifyEmailForm />
    </Suspense>
  );
}
