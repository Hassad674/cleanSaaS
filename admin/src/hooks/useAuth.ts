import {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  type ReactNode,
} from "react";
import { createElement } from "react";
import {
  login as apiLogin,
  fetchCurrentUser,
} from "../lib/api.ts";
import type { User } from "../types/index.ts";

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

type AuthState = {
  token: string | null;
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
};

// ---------------------------------------------------------------------------
// Context
// ---------------------------------------------------------------------------

const AuthContext = createContext<AuthState | null>(null);

// ---------------------------------------------------------------------------
// Provider
// ---------------------------------------------------------------------------

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(() =>
    localStorage.getItem("admin_token")
  );
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Validate token on mount
  useEffect(() => {
    if (!token) {
      setIsLoading(false);
      return;
    }

    fetchCurrentUser()
      .then((u) => {
        if (u.role !== "admin") {
          // Only admins are allowed
          localStorage.removeItem("admin_token");
          setToken(null);
          setUser(null);
        } else {
          setUser(u);
        }
      })
      .catch(() => {
        localStorage.removeItem("admin_token");
        setToken(null);
        setUser(null);
      })
      .finally(() => setIsLoading(false));
  }, [token]);

  const login = useCallback(async (email: string, password: string) => {
    const res = await apiLogin(email, password);

    if (res.user.role !== "admin") {
      throw new Error("Access denied. Admin role required.");
    }

    localStorage.setItem("admin_token", res.token);
    setToken(res.token);
    setUser(res.user);
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem("admin_token");
    setToken(null);
    setUser(null);
  }, []);

  const value: AuthState = {
    token,
    user,
    isAuthenticated: !!token && !!user,
    isLoading,
    login,
    logout,
  };

  return createElement(AuthContext.Provider, { value }, children);
}

// ---------------------------------------------------------------------------
// Hook
// ---------------------------------------------------------------------------

export function useAuth(): AuthState {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return ctx;
}
