"use client";

import { useState, useCallback } from "react";
import Link from "next/link";

/* -------------------------------------------------------------------------- */
/*  Types                                                                      */
/* -------------------------------------------------------------------------- */

type View = "login" | "register" | "forgot-password";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081";

interface AuthUser {
  id: string;
  email: string;
  name: string;
  avatar_url: string;
  role: "admin" | "member";
  email_verified: boolean;
}

interface AuthApiResponse {
  token: string;
  user: AuthUser;
}

interface AuthSession {
  user: AuthUser;
  token: string;
}

/* -------------------------------------------------------------------------- */
/*  API helpers                                                                */
/* -------------------------------------------------------------------------- */

async function apiRegister(
  name: string,
  email: string,
  password: string
): Promise<{ data: AuthApiResponse | null; error: string | null }> {
  try {
    const res = await fetch(`${API_URL}/auth/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name, email, password }),
    });

    const json = await res.json().catch(() => null);

    if (res.ok) {
      return { data: json as AuthApiResponse, error: null };
    }

    const errorMessage =
      (json as { error?: string })?.error || "Registration failed";
    return { data: null, error: errorMessage };
  } catch {
    return { data: null, error: "Network error — is the backend running?" };
  }
}

async function apiLogin(
  email: string,
  password: string
): Promise<{ data: AuthApiResponse | null; error: string | null }> {
  try {
    const res = await fetch(`${API_URL}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });

    const json = await res.json().catch(() => null);

    if (res.ok) {
      return { data: json as AuthApiResponse, error: null };
    }

    const errorMessage =
      (json as { error?: string })?.error || "Login failed";
    return { data: null, error: errorMessage };
  } catch {
    return { data: null, error: "Network error — is the backend running?" };
  }
}

/* -------------------------------------------------------------------------- */
/*  Demo Login form — mirrors features/auth/components/login-form.tsx          */
/* -------------------------------------------------------------------------- */

function DemoLoginForm({
  onNavigate,
  onLogin,
}: {
  onNavigate: (view: View) => void;
  onLogin: (session: AuthSession) => void;
}) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setLoading(true);

    const result = await apiLogin(email, password);

    if (result.error || !result.data) {
      setError(result.error || "Login failed");
      setLoading(false);
      return;
    }

    setLoading(false);
    onLogin({ user: result.data.user, token: result.data.token });
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <h1 className="text-2xl font-bold text-center text-card-foreground">
        Log in
      </h1>

      {error && (
        <p className="text-sm text-destructive text-center">{error}</p>
      )}

      <div>
        <label
          htmlFor="demo-login-email"
          className="block text-sm font-medium mb-1 text-card-foreground"
        >
          Email
        </label>
        <input
          id="demo-login-email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          className="w-full px-3 py-2 border border-input rounded-lg bg-background text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
      </div>

      <div>
        <label
          htmlFor="demo-login-password"
          className="block text-sm font-medium mb-1 text-card-foreground"
        >
          Password
        </label>
        <input
          id="demo-login-password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          className="w-full px-3 py-2 border border-input rounded-lg bg-background text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
      </div>

      <button
        type="submit"
        disabled={loading}
        className="w-full bg-primary text-primary-foreground py-2 rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50 font-medium"
      >
        {loading ? "Logging in..." : "Log in"}
      </button>

      <div className="text-center text-sm text-muted-foreground space-y-2">
        <button
          type="button"
          onClick={() => onNavigate("forgot-password")}
          className="text-primary hover:underline text-sm block w-full"
        >
          Forgot password?
        </button>
        <p>
          Don&apos;t have an account?{" "}
          <button
            type="button"
            onClick={() => onNavigate("register")}
            className="font-medium text-primary hover:underline"
          >
            Sign up
          </button>
        </p>
      </div>
    </form>
  );
}

/* -------------------------------------------------------------------------- */
/*  Demo Register form — mirrors features/auth/components/register-form.tsx    */
/* -------------------------------------------------------------------------- */

function DemoRegisterForm({
  onNavigate,
  onRegister,
}: {
  onNavigate: (view: View) => void;
  onRegister: (session: AuthSession) => void;
}) {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setLoading(true);

    const result = await apiRegister(name, email, password);

    if (result.error || !result.data) {
      setError(result.error || "Registration failed");
      setLoading(false);
      return;
    }

    setLoading(false);
    onRegister({ user: result.data.user, token: result.data.token });
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <h1 className="text-2xl font-bold text-center text-card-foreground">
        Create account
      </h1>

      {error && (
        <p className="text-sm text-destructive text-center">{error}</p>
      )}

      <div>
        <label
          htmlFor="demo-register-name"
          className="block text-sm font-medium mb-1 text-card-foreground"
        >
          Name
        </label>
        <input
          id="demo-register-name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
          className="w-full px-3 py-2 border border-input rounded-lg bg-background text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
      </div>

      <div>
        <label
          htmlFor="demo-register-email"
          className="block text-sm font-medium mb-1 text-card-foreground"
        >
          Email
        </label>
        <input
          id="demo-register-email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          className="w-full px-3 py-2 border border-input rounded-lg bg-background text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
      </div>

      <div>
        <label
          htmlFor="demo-register-password"
          className="block text-sm font-medium mb-1 text-card-foreground"
        >
          Password
        </label>
        <input
          id="demo-register-password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          minLength={8}
          className="w-full px-3 py-2 border border-input rounded-lg bg-background text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
      </div>

      <button
        type="submit"
        disabled={loading}
        className="w-full bg-primary text-primary-foreground py-2 rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50 font-medium"
      >
        {loading ? "Creating account..." : "Create account"}
      </button>

      <p className="text-center text-sm text-muted-foreground">
        Already have an account?{" "}
        <button
          type="button"
          onClick={() => onNavigate("login")}
          className="font-medium text-primary hover:underline"
        >
          Log in
        </button>
      </p>
    </form>
  );
}

/* -------------------------------------------------------------------------- */
/*  Demo Forgot Password — mirrors features/auth/components/forgot-password-form.tsx */
/* -------------------------------------------------------------------------- */

function DemoForgotPasswordForm({
  onNavigate,
}: {
  onNavigate: (view: View) => void;
}) {
  const [email, setEmail] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setSuccess(false);
    setLoading(true);

    try {
      const res = await fetch(`${API_URL}/auth/forgot-password`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email }),
      });

      await res.json().catch(() => null);

      // Always show success (backend doesn't leak user existence)
      setSuccess(true);
    } catch {
      setError("Network error — is the backend running?");
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <h1 className="text-2xl font-bold text-center text-card-foreground">
        Forgot password
      </h1>

      <p className="text-sm text-muted-foreground text-center">
        Enter your email and we&apos;ll send you a link to reset your password.
      </p>

      {error && (
        <p className="text-sm text-destructive text-center">{error}</p>
      )}

      {success && (
        <p className="text-sm text-success text-center">
          If an account exists with that email, we&apos;ve sent a reset link.
        </p>
      )}

      <div>
        <label
          htmlFor="demo-forgot-email"
          className="block text-sm font-medium mb-1 text-card-foreground"
        >
          Email
        </label>
        <input
          id="demo-forgot-email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          className="w-full px-3 py-2 border border-border rounded-lg bg-background text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
      </div>

      <button
        type="submit"
        disabled={loading}
        className="w-full bg-primary text-primary-foreground py-2 rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50 font-medium"
      >
        {loading ? "Sending..." : "Send reset link"}
      </button>

      <div className="text-center text-sm">
        <button
          type="button"
          onClick={() => onNavigate("login")}
          className="text-primary hover:underline"
        >
          Back to login
        </button>
      </div>
    </form>
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
/*  Demo Dashboard — shows real user data from the API                         */
/* -------------------------------------------------------------------------- */

function DemoDashboard({
  session,
  onLogout,
}: {
  session: AuthSession;
  onLogout: () => void;
}) {
  const [showToken, setShowToken] = useState(false);

  return (
    <div className="max-w-2xl">
      <h1 className="text-2xl font-bold text-foreground mb-6">Dashboard</h1>
      <div className="bg-card border border-border rounded-xl p-6 space-y-4">
        <div>
          <p className="text-sm text-muted-foreground">Welcome back,</p>
          <p className="text-lg font-semibold text-foreground">
            {session.user.name}
          </p>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
          <div>
            <p className="text-muted-foreground">Email</p>
            <p className="text-foreground">{session.user.email}</p>
          </div>
          <div>
            <p className="text-muted-foreground">Role</p>
            <p className="text-foreground capitalize">{session.user.role}</p>
          </div>
          <div>
            <p className="text-muted-foreground">User ID</p>
            <p className="text-foreground font-mono text-xs">
              {session.user.id}
            </p>
          </div>
          <div>
            <p className="text-muted-foreground">Email verified</p>
            <p className="text-foreground">
              {session.user.email_verified ? "Yes" : "No"}
            </p>
          </div>
        </div>

        {/* JWT token reveal */}
        <div className="pt-2 border-t border-border">
          <button
            onClick={() => setShowToken((v) => !v)}
            className="text-sm text-primary hover:underline"
          >
            {showToken ? "Hide JWT token" : "Show JWT token"}
          </button>
          {showToken && (
            <div className="mt-2 rounded-lg bg-muted p-3">
              <p className="text-xs font-mono text-muted-foreground break-all leading-relaxed">
                {session.token}
              </p>
            </div>
          )}
        </div>

        <div className="pt-2">
          <button
            onClick={onLogout}
            className="text-sm text-destructive hover:underline"
          >
            Log out
          </button>
        </div>
      </div>

      {/* Demo hint */}
      <div className="mt-4 rounded-lg bg-muted/50 border border-border p-3">
        <p className="text-xs text-muted-foreground text-center">
          <span className="font-medium text-foreground">Live demo</span> — This
          session uses a real JWT token from the backend. Click &ldquo;Log
          out&rdquo; to clear the session and return to the login form.
        </p>
      </div>
    </div>
  );
}

/* -------------------------------------------------------------------------- */
/*  Main demo component                                                        */
/* -------------------------------------------------------------------------- */

export function AuthDemo() {
  const [activeView, setActiveView] = useState<View>("login");
  const [session, setSession] = useState<AuthSession | null>(null);

  const handleLogin = useCallback((s: AuthSession) => {
    setSession(s);
  }, []);

  const handleRegister = useCallback((s: AuthSession) => {
    setSession(s);
  }, []);

  const handleLogout = useCallback(() => {
    setSession(null);
    setActiveView("login");
  }, []);

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
          {session
            ? `Logged in as ${session.user.name}. This is what the dashboard looks like after authentication.`
            : "Login, register, and password reset — all pre-built and ready to customise. Try the forms below with real data."}
        </p>
      </div>

      {/* Logged-in state: show dashboard */}
      {session ? (
        <DemoDashboard session={session} onLogout={handleLogout} />
      ) : (
        /* Logged-out state: form + features side by side */
        <div className="grid grid-cols-1 lg:grid-cols-5 gap-8 lg:gap-12 max-w-5xl">
          {/* Left: form card */}
          <div className="lg:col-span-3">
            <div className="bg-card border border-border rounded-xl p-6 sm:p-8 shadow-sm">
              {activeView === "login" && (
                <DemoLoginForm
                  onNavigate={setActiveView}
                  onLogin={handleLogin}
                />
              )}
              {activeView === "register" && (
                <DemoRegisterForm
                  onNavigate={setActiveView}
                  onRegister={handleRegister}
                />
              )}
              {activeView === "forgot-password" && (
                <DemoForgotPasswordForm onNavigate={setActiveView} />
              )}

              {/* Demo hint */}
              <div className="mt-6 rounded-lg bg-muted/50 border border-border p-3">
                <p className="text-xs text-muted-foreground text-center">
                  <span className="font-medium text-foreground">
                    Real registration
                  </span>{" "}
                  — Your account is created in our test database. Register, then
                  log in with your credentials.
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
      )}
    </div>
  );
}
