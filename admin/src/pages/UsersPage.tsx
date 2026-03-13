import { useEffect, useState, useCallback } from "react";
import DataTable, { type Column } from "../components/DataTable.tsx";
import { fetchUsers, updateUserRole, updateUserStatus } from "../lib/api.ts";
import type { User } from "../types/index.ts";

const LIMIT = 20;

const ROLE_OPTIONS = ["user", "admin"] as const;
const STATUS_OPTIONS = ["active", "suspended"] as const;

export default function UsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [isLoading, setIsLoading] = useState(true);

  // Dropdown state: track which user + which dropdown is open
  const [openDropdown, setOpenDropdown] = useState<{
    userId: string;
    type: "role" | "status";
  } | null>(null);

  const loadUsers = useCallback(async () => {
    setIsLoading(true);
    try {
      const res = await fetchUsers({ page, limit: LIMIT, search: search || undefined });
      setUsers(res.data);
      setTotal(res.total);
    } catch {
      setUsers([]);
      setTotal(0);
    } finally {
      setIsLoading(false);
    }
  }, [page, search]);

  useEffect(() => {
    loadUsers();
  }, [loadUsers]);

  // Debounced search — reset page on new search
  function handleSearchChange(value: string) {
    setSearch(value);
    setPage(1);
  }

  async function handleRoleChange(userId: string, role: string) {
    try {
      await updateUserRole(userId, role);
      setOpenDropdown(null);
      await loadUsers();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Failed to update role");
    }
  }

  async function handleStatusChange(userId: string, status: string) {
    try {
      await updateUserStatus(userId, status);
      setOpenDropdown(null);
      await loadUsers();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Failed to update status");
    }
  }

  const totalPages = Math.max(1, Math.ceil(total / LIMIT));

  const columns: Column<User>[] = [
    {
      key: "name",
      header: "Name",
      render: (u) => <span className="font-medium">{u.name}</span>,
    },
    {
      key: "email",
      header: "Email",
      render: (u) => u.email,
    },
    {
      key: "role",
      header: "Role",
      render: (u) => (
        <div className="relative">
          <button
            onClick={() =>
              setOpenDropdown(
                openDropdown?.userId === u.id && openDropdown.type === "role"
                  ? null
                  : { userId: u.id, type: "role" }
              )
            }
            className={`rounded-full px-2.5 py-0.5 text-xs font-medium ${
              u.role === "admin"
                ? "bg-primary/10 text-primary"
                : "bg-muted text-muted-foreground"
            }`}
          >
            {u.role}
          </button>

          {openDropdown?.userId === u.id && openDropdown.type === "role" && (
            <div className="absolute left-0 top-full z-10 mt-1 rounded-lg border border-border bg-popover py-1 shadow-lg">
              {ROLE_OPTIONS.map((r) => (
                <button
                  key={r}
                  onClick={() => handleRoleChange(u.id, r)}
                  className="block w-full px-4 py-1.5 text-left text-sm text-popover-foreground hover:bg-muted"
                >
                  {r}
                </button>
              ))}
            </div>
          )}
        </div>
      ),
    },
    {
      key: "status",
      header: "Status",
      render: (u) => {
        const isVerified = u.email_verified;
        return (
          <div className="relative">
            <button
              onClick={() =>
                setOpenDropdown(
                  openDropdown?.userId === u.id && openDropdown.type === "status"
                    ? null
                    : { userId: u.id, type: "status" }
                )
              }
              className={`rounded-full px-2.5 py-0.5 text-xs font-medium ${
                isVerified
                  ? "bg-success/10 text-success"
                  : "bg-warning/10 text-warning"
              }`}
            >
              {isVerified ? "active" : "unverified"}
            </button>

            {openDropdown?.userId === u.id &&
              openDropdown.type === "status" && (
                <div className="absolute left-0 top-full z-10 mt-1 rounded-lg border border-border bg-popover py-1 shadow-lg">
                  {STATUS_OPTIONS.map((s) => (
                    <button
                      key={s}
                      onClick={() => handleStatusChange(u.id, s)}
                      className="block w-full px-4 py-1.5 text-left text-sm text-popover-foreground hover:bg-muted"
                    >
                      {s}
                    </button>
                  ))}
                </div>
              )}
          </div>
        );
      },
    },
    {
      key: "created_at",
      header: "Created",
      render: (u) =>
        new Date(u.created_at).toLocaleDateString("en-US", {
          year: "numeric",
          month: "short",
          day: "numeric",
        }),
    },
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-foreground">Users</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Manage users and their roles
          </p>
        </div>
      </div>

      {/* Search */}
      <div className="max-w-sm">
        <input
          type="text"
          value={search}
          onChange={(e) => handleSearchChange(e.target.value)}
          placeholder="Search by name or email..."
          className="w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
        />
      </div>

      {/* Table */}
      <DataTable
        columns={columns}
        data={users}
        keyExtractor={(u) => u.id}
        page={page}
        totalPages={totalPages}
        onPageChange={setPage}
        isLoading={isLoading}
      />
    </div>
  );
}
