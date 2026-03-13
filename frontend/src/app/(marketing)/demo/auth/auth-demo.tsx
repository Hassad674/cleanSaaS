"use client";

import { useState, useCallback } from "react";
import Link from "next/link";

/* -------------------------------------------------------------------------- */
/*  Types                                                                      */
/* -------------------------------------------------------------------------- */

type Tab = "login" | "register";

interface FieldError {
  field: string;
  message: string;
}

/* -------------------------------------------------------------------------- */
/*  Helpers                                                                    */
/* -------------------------------------------------------------------------- */

function isValidEmail(email: string): boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

/* -------------------------------------------------------------------------- */
/*  Spinner                                                                    */
/* -------------------------------------------------------------------------- */

function Spinner() {
  return (
    <svg
      className="h-5 w-5 animate-spin"
      viewBox="0 0 24 24"
      fill="none"
    >
      <circle
        className="opacity-25"
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        strokeWidth="4"
      />
      <path
        className="opacity-75"
        fill="currentColor"
        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
      />
    </svg>
  );
}

/* -------------------------------------------------------------------------- */
/*  Google icon                                                                */
/* -------------------------------------------------------------------------- */

function GoogleIcon() {
  return (
    <svg className="h-5 w-5" viewBox="0 0 24 24">
      <path
        d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 01-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z"
        fill="#4285F4"
      />
      <path
        d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
        fill="#34A853"
      />
      <path
        d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
        fill="#FBBC05"
      />
      <path
        d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
        fill="#EA4335"
      />
    </svg>
  );
}

/* -------------------------------------------------------------------------- */
/*  Login form                                                                 */
/* -------------------------------------------------------------------------- */

function LoginForm() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [errors, setErrors] = useState<FieldError[]>([]);
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);

  const fieldError = useCallback(
    (field: string) => errors.find((e) => e.field === field)?.message,
    [errors],
  );

  function validate(): FieldError[] {
    const errs: FieldError[] = [];
    if (!email.trim()) {
      errs.push({ field: "email", message: "Email is required" });
    } else if (!isValidEmail(email)) {
      errs.push({ field: "email", message: "Please enter a valid email" });
    }
    if (!password) {
      errs.push({ field: "password", message: "Password is required" });
    }
    return errs;
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    const validationErrors = validate();
    setErrors(validationErrors);
    if (validationErrors.length > 0) return;

    setLoading(true);
    await sleep(1500);
    setLoading(false);
    setSuccess(true);
  }

  if (success) {
    return (
      <div className="text-center space-y-4 py-8">
        <div className="mx-auto h-14 w-14 rounded-full bg-primary/10 flex items-center justify-center">
          <svg
            className="h-7 w-7 text-primary"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            strokeWidth={2}
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M4.5 12.75l6 6 9-13.5"
            />
          </svg>
        </div>
        <h3 className="text-lg font-semibold text-foreground">
          Logged in successfully
        </h3>
        <p className="text-sm text-muted-foreground">
          Logged in as{" "}
          <span className="font-medium text-foreground">{email}</span>
        </p>
        <p className="text-xs text-muted-foreground">
          This is a simulated demo — no real session was created.
        </p>
        <button
          onClick={() => {
            setSuccess(false);
            setEmail("");
            setPassword("");
          }}
          className="mt-2 text-sm text-primary hover:underline font-medium"
        >
          Try again
        </button>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {/* Email */}
      <div>
        <label
          htmlFor="login-email"
          className="block text-sm font-medium mb-1.5 text-foreground"
        >
          Email
        </label>
        <input
          id="login-email"
          type="email"
          placeholder="demo@example.com"
          value={email}
          onChange={(e) => {
            setEmail(e.target.value);
            setErrors((prev) => prev.filter((err) => err.field !== "email"));
          }}
          className={`w-full px-3 py-2.5 border rounded-lg bg-background text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring transition-colors ${
            fieldError("email") ? "border-destructive" : "border-input"
          }`}
        />
        {fieldError("email") && (
          <p className="text-xs text-destructive mt-1">
            {fieldError("email")}
          </p>
        )}
      </div>

      {/* Password */}
      <div>
        <div className="flex items-center justify-between mb-1.5">
          <label
            htmlFor="login-password"
            className="block text-sm font-medium text-foreground"
          >
            Password
          </label>
          <button
            type="button"
            className="text-xs text-primary hover:underline"
            onClick={() =>
              alert(
                "In a real app, this would navigate to the forgot-password page.",
              )
            }
          >
            Forgot password?
          </button>
        </div>
        <input
          id="login-password"
          type="password"
          placeholder="Enter your password"
          value={password}
          onChange={(e) => {
            setPassword(e.target.value);
            setErrors((prev) =>
              prev.filter((err) => err.field !== "password"),
            );
          }}
          className={`w-full px-3 py-2.5 border rounded-lg bg-background text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring transition-colors ${
            fieldError("password") ? "border-destructive" : "border-input"
          }`}
        />
        {fieldError("password") && (
          <p className="text-xs text-destructive mt-1">
            {fieldError("password")}
          </p>
        )}
      </div>

      {/* Submit */}
      <button
        type="submit"
        disabled={loading}
        className="w-full bg-primary text-primary-foreground py-2.5 rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50 font-medium flex items-center justify-center gap-2"
      >
        {loading ? (
          <>
            <Spinner />
            Logging in...
          </>
        ) : (
          "Log in"
        )}
      </button>

      {/* Divider */}
      <div className="relative my-6">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t border-border" />
        </div>
        <div className="relative flex justify-center text-xs">
          <span className="bg-card px-3 text-muted-foreground">
            or continue with
          </span>
        </div>
      </div>

      {/* Google OAuth */}
      <button
        type="button"
        onClick={() =>
          alert("In a real app, this would redirect to Google OAuth.")
        }
        className="w-full border border-border bg-background text-foreground py-2.5 rounded-lg hover:bg-muted transition-colors font-medium flex items-center justify-center gap-3"
      >
        <GoogleIcon />
        Continue with Google
      </button>
    </form>
  );
}

/* -------------------------------------------------------------------------- */
/*  Register form                                                              */
/* -------------------------------------------------------------------------- */

function RegisterForm() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [errors, setErrors] = useState<FieldError[]>([]);
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);

  const fieldError = useCallback(
    (field: string) => errors.find((e) => e.field === field)?.message,
    [errors],
  );

  function validate(): FieldError[] {
    const errs: FieldError[] = [];
    if (!name.trim()) {
      errs.push({ field: "name", message: "Name is required" });
    }
    if (!email.trim()) {
      errs.push({ field: "email", message: "Email is required" });
    } else if (!isValidEmail(email)) {
      errs.push({ field: "email", message: "Please enter a valid email" });
    }
    if (!password) {
      errs.push({ field: "password", message: "Password is required" });
    } else if (password.length < 8) {
      errs.push({
        field: "password",
        message: "Password must be at least 8 characters",
      });
    }
    if (!confirmPassword) {
      errs.push({
        field: "confirmPassword",
        message: "Please confirm your password",
      });
    } else if (password !== confirmPassword) {
      errs.push({
        field: "confirmPassword",
        message: "Passwords do not match",
      });
    }
    return errs;
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    const validationErrors = validate();
    setErrors(validationErrors);
    if (validationErrors.length > 0) return;

    setLoading(true);
    await sleep(1500);
    setLoading(false);
    setSuccess(true);
  }

  if (success) {
    return (
      <div className="text-center space-y-4 py-8">
        <div className="mx-auto h-14 w-14 rounded-full bg-primary/10 flex items-center justify-center">
          <svg
            className="h-7 w-7 text-primary"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            strokeWidth={2}
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M21.75 6.75v10.5a2.25 2.25 0 01-2.25 2.25h-15a2.25 2.25 0 01-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25m19.5 0v.243a2.25 2.25 0 01-1.07 1.916l-7.5 4.615a2.25 2.25 0 01-2.36 0L3.32 8.91a2.25 2.25 0 01-1.07-1.916V6.75"
            />
          </svg>
        </div>
        <h3 className="text-lg font-semibold text-foreground">
          Check your email
        </h3>
        <p className="text-sm text-muted-foreground">
          We sent a verification link to{" "}
          <span className="font-medium text-foreground">{email}</span>
        </p>
        <p className="text-xs text-muted-foreground">
          This is a simulated demo — no real email was sent.
        </p>
        <button
          onClick={() => {
            setSuccess(false);
            setName("");
            setEmail("");
            setPassword("");
            setConfirmPassword("");
          }}
          className="mt-2 text-sm text-primary hover:underline font-medium"
        >
          Try again
        </button>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {/* Name */}
      <div>
        <label
          htmlFor="register-name"
          className="block text-sm font-medium mb-1.5 text-foreground"
        >
          Full name
        </label>
        <input
          id="register-name"
          type="text"
          placeholder="Jane Doe"
          value={name}
          onChange={(e) => {
            setName(e.target.value);
            setErrors((prev) => prev.filter((err) => err.field !== "name"));
          }}
          className={`w-full px-3 py-2.5 border rounded-lg bg-background text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring transition-colors ${
            fieldError("name") ? "border-destructive" : "border-input"
          }`}
        />
        {fieldError("name") && (
          <p className="text-xs text-destructive mt-1">
            {fieldError("name")}
          </p>
        )}
      </div>

      {/* Email */}
      <div>
        <label
          htmlFor="register-email"
          className="block text-sm font-medium mb-1.5 text-foreground"
        >
          Email
        </label>
        <input
          id="register-email"
          type="email"
          placeholder="jane@example.com"
          value={email}
          onChange={(e) => {
            setEmail(e.target.value);
            setErrors((prev) => prev.filter((err) => err.field !== "email"));
          }}
          className={`w-full px-3 py-2.5 border rounded-lg bg-background text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring transition-colors ${
            fieldError("email") ? "border-destructive" : "border-input"
          }`}
        />
        {fieldError("email") && (
          <p className="text-xs text-destructive mt-1">
            {fieldError("email")}
          </p>
        )}
      </div>

      {/* Password */}
      <div>
        <label
          htmlFor="register-password"
          className="block text-sm font-medium mb-1.5 text-foreground"
        >
          Password
        </label>
        <input
          id="register-password"
          type="password"
          placeholder="At least 8 characters"
          value={password}
          onChange={(e) => {
            setPassword(e.target.value);
            setErrors((prev) =>
              prev.filter((err) => err.field !== "password"),
            );
          }}
          className={`w-full px-3 py-2.5 border rounded-lg bg-background text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring transition-colors ${
            fieldError("password") ? "border-destructive" : "border-input"
          }`}
        />
        {fieldError("password") && (
          <p className="text-xs text-destructive mt-1">
            {fieldError("password")}
          </p>
        )}
      </div>

      {/* Confirm password */}
      <div>
        <label
          htmlFor="register-confirm"
          className="block text-sm font-medium mb-1.5 text-foreground"
        >
          Confirm password
        </label>
        <input
          id="register-confirm"
          type="password"
          placeholder="Re-enter your password"
          value={confirmPassword}
          onChange={(e) => {
            setConfirmPassword(e.target.value);
            setErrors((prev) =>
              prev.filter((err) => err.field !== "confirmPassword"),
            );
          }}
          className={`w-full px-3 py-2.5 border rounded-lg bg-background text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring transition-colors ${
            fieldError("confirmPassword")
              ? "border-destructive"
              : "border-input"
          }`}
        />
        {fieldError("confirmPassword") && (
          <p className="text-xs text-destructive mt-1">
            {fieldError("confirmPassword")}
          </p>
        )}
      </div>

      {/* Submit */}
      <button
        type="submit"
        disabled={loading}
        className="w-full bg-primary text-primary-foreground py-2.5 rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50 font-medium flex items-center justify-center gap-2"
      >
        {loading ? (
          <>
            <Spinner />
            Creating account...
          </>
        ) : (
          "Create account"
        )}
      </button>
    </form>
  );
}

/* -------------------------------------------------------------------------- */
/*  Tab switcher                                                               */
/* -------------------------------------------------------------------------- */

function TabSwitcher({
  active,
  onChange,
}: {
  active: Tab;
  onChange: (tab: Tab) => void;
}) {
  return (
    <div className="flex rounded-lg bg-muted p-1">
      <button
        type="button"
        onClick={() => onChange("login")}
        className={`flex-1 py-2 text-sm font-medium rounded-md transition-colors ${
          active === "login"
            ? "bg-card text-foreground shadow-sm"
            : "text-muted-foreground hover:text-foreground"
        }`}
      >
        Log in
      </button>
      <button
        type="button"
        onClick={() => onChange("register")}
        className={`flex-1 py-2 text-sm font-medium rounded-md transition-colors ${
          active === "register"
            ? "bg-card text-foreground shadow-sm"
            : "text-muted-foreground hover:text-foreground"
        }`}
      >
        Register
      </button>
    </div>
  );
}

/* -------------------------------------------------------------------------- */
/*  Feature list (right panel on desktop)                                      */
/* -------------------------------------------------------------------------- */

const features = [
  {
    title: "Email & password",
    description: "Classic credential-based login with bcrypt hashing.",
    icon: (
      <svg
        className="h-5 w-5"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        strokeWidth={1.5}
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M21.75 6.75v10.5a2.25 2.25 0 01-2.25 2.25h-15a2.25 2.25 0 01-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25m19.5 0v.243a2.25 2.25 0 01-1.07 1.916l-7.5 4.615a2.25 2.25 0 01-2.36 0L3.32 8.91a2.25 2.25 0 01-1.07-1.916V6.75"
        />
      </svg>
    ),
  },
  {
    title: "Google OAuth",
    description: "One-click sign-in with Google. Easily add more providers.",
    icon: (
      <svg
        className="h-5 w-5"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        strokeWidth={1.5}
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M15.75 5.25a3 3 0 013 3m3 0a6 6 0 01-7.029 5.912c-.563-.097-1.159.026-1.563.43L10.5 17.25H8.25v2.25H6v2.25H2.25v-2.818c0-.597.237-1.17.659-1.591l6.499-6.499c.404-.404.527-1 .43-1.563A6 6 0 1121.75 8.25z"
        />
      </svg>
    ),
  },
  {
    title: "Email verification",
    description: "Verify new accounts via a token-based email link.",
    icon: (
      <svg
        className="h-5 w-5"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        strokeWidth={1.5}
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
        />
      </svg>
    ),
  },
  {
    title: "Forgot & reset password",
    description: "Secure password recovery with time-limited tokens.",
    icon: (
      <svg
        className="h-5 w-5"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        strokeWidth={1.5}
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z"
        />
      </svg>
    ),
  },
  {
    title: "JWT sessions",
    description: "Stateless auth with signed JWTs and proper expiration.",
    icon: (
      <svg
        className="h-5 w-5"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        strokeWidth={1.5}
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M7.864 4.243A7.5 7.5 0 0119.5 10.5c0 2.92-.556 5.709-1.568 8.268M5.742 6.364A7.465 7.465 0 004.5 10.5a48.667 48.667 0 00-1.488 8.546m4.07-2.163a48.11 48.11 0 01-1.653 4.12M14.1 14.1l-2.1 2.1m0 0l-2.1 2.1M12 16.2l2.1 2.1M12 16.2l-2.1-2.1"
        />
      </svg>
    ),
  },
  {
    title: "Rate limiting",
    description: "Protect auth endpoints from brute-force attacks.",
    icon: (
      <svg
        className="h-5 w-5"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        strokeWidth={1.5}
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M12 9v3.75m0-10.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285z"
        />
      </svg>
    ),
  },
];

function FeatureList() {
  return (
    <div className="space-y-4">
      <h3 className="text-sm font-semibold text-foreground uppercase tracking-wider">
        What&apos;s included
      </h3>
      <ul className="space-y-3">
        {features.map((feature) => (
          <li key={feature.title} className="flex gap-3 items-start">
            <div className="flex-shrink-0 h-9 w-9 rounded-lg bg-primary/10 text-primary flex items-center justify-center">
              {feature.icon}
            </div>
            <div className="min-w-0">
              <p className="text-sm font-medium text-foreground">
                {feature.title}
              </p>
              <p className="text-xs text-muted-foreground leading-relaxed">
                {feature.description}
              </p>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}

/* -------------------------------------------------------------------------- */
/*  Main demo component                                                        */
/* -------------------------------------------------------------------------- */

export function AuthDemo() {
  const [activeTab, setActiveTab] = useState<Tab>("login");

  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-12 sm:py-16">
      {/* Back link */}
      <Link
        href="/demo"
        className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors mb-8"
      >
        <svg
          className="h-4 w-4"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={2}
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18"
          />
        </svg>
        Back to demos
      </Link>

      {/* Header */}
      <div className="max-w-2xl mb-10">
        <span className="inline-block bg-primary/10 text-primary text-sm font-medium px-3 py-1 rounded-full mb-4">
          Authentication
        </span>
        <h1 className="text-3xl sm:text-4xl font-bold text-foreground tracking-tight">
          Secure auth, out of the box
        </h1>
        <p className="text-muted-foreground mt-3 text-base sm:text-lg leading-relaxed">
          Login, register, email verification, password reset, and OAuth — all
          pre-built and ready to customise. Try the forms below with simulated
          data.
        </p>
      </div>

      {/* Content: form + features side by side on desktop */}
      <div className="grid grid-cols-1 lg:grid-cols-5 gap-8 lg:gap-12 max-w-5xl">
        {/* Left: form card */}
        <div className="lg:col-span-3">
          <div className="bg-card border border-border rounded-xl p-6 sm:p-8 shadow-sm">
            {/* Mobile tab switcher */}
            <div className="mb-6 lg:hidden">
              <TabSwitcher active={activeTab} onChange={setActiveTab} />
            </div>

            {/* Desktop: side-by-side tabs */}
            <div className="hidden lg:block mb-6">
              <TabSwitcher active={activeTab} onChange={setActiveTab} />
            </div>

            {/* Forms */}
            {activeTab === "login" ? <LoginForm /> : <RegisterForm />}

            {/* Demo hint */}
            <div className="mt-6 rounded-lg bg-muted/50 border border-border p-3">
              <p className="text-xs text-muted-foreground text-center">
                <span className="font-medium text-foreground">Demo mode</span>{" "}
                — Fill in any values and submit to see the full flow. No real
                requests are made.
              </p>
            </div>
          </div>
        </div>

        {/* Right: feature list */}
        <div className="lg:col-span-2">
          <div className="bg-card border border-border rounded-xl p-6 shadow-sm sticky top-24">
            <FeatureList />

            {/* Code hint */}
            <div className="mt-6 pt-5 border-t border-border">
              <p className="text-xs text-muted-foreground leading-relaxed">
                All auth logic lives in{" "}
                <code className="bg-muted px-1.5 py-0.5 rounded text-xs font-mono">
                  features/auth/
                </code>{" "}
                on the frontend and{" "}
                <code className="bg-muted px-1.5 py-0.5 rounded text-xs font-mono">
                  internal/auth/
                </code>{" "}
                on the backend. Fully removable — delete the folder and its
                wiring, zero compilation errors.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
