"use client";

import { useState, useEffect } from "react";
import { useAuth } from "@/features/auth/hooks/use-auth";
import { getInvoices } from "@/features/billing/actions/billing";
import { formatCurrency, formatDate } from "@/shared/lib/utils";
import { PAGINATION_DEFAULT_LIMIT } from "@/shared/lib/constants";
import type { Invoice } from "@/features/billing/types";

export function InvoiceList() {
  const { getToken } = useAuth({ required: true });

  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [total, setTotal] = useState(0);
  const [offset, setOffset] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const limit = PAGINATION_DEFAULT_LIMIT;

  useEffect(() => {
    const token = getToken();
    if (!token) return;

    setLoading(true);
    getInvoices(token, offset, limit).then((res) => {
      if (res.data) {
        setInvoices(res.data.invoices ?? []);
        setTotal(res.data.total);
      } else {
        setError(res.error ?? "Failed to load invoices");
      }
      setLoading(false);
    });
  }, [getToken, offset, limit]);

  if (loading) {
    return (
      <div className="bg-card border border-border rounded-xl p-6 shadow-sm animate-pulse h-40" />
    );
  }

  if (error) {
    return (
      <div className="bg-card border border-border rounded-xl p-6 shadow-sm">
        <h2 className="text-lg font-semibold text-foreground mb-2">Invoices</h2>
        <p className="text-sm text-destructive">{error}</p>
      </div>
    );
  }

  if (invoices.length === 0 && offset === 0) {
    return (
      <div className="bg-card border border-border rounded-xl p-6 shadow-sm">
        <h2 className="text-lg font-semibold text-foreground mb-2">Invoices</h2>
        <p className="text-muted-foreground text-sm">No invoices yet.</p>
      </div>
    );
  }

  const hasNext = offset + limit < total;
  const hasPrev = offset > 0;

  return (
    <div className="bg-card border border-border rounded-xl p-6 shadow-sm space-y-4">
      <h2 className="text-lg font-semibold text-foreground">Invoices</h2>

      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border text-left">
              <th className="pb-2 font-medium text-muted-foreground">Date</th>
              <th className="pb-2 font-medium text-muted-foreground">Amount</th>
              <th className="pb-2 font-medium text-muted-foreground">Status</th>
              <th className="pb-2 font-medium text-muted-foreground" />
            </tr>
          </thead>
          <tbody>
            {invoices.map((invoice) => (
              <tr key={invoice.id} className="border-b border-border last:border-0">
                <td className="py-3 text-foreground">
                  {formatDate(invoice.created_at)}
                </td>
                <td className="py-3 text-foreground">
                  {formatCurrency(invoice.amount_cents, invoice.currency)}
                </td>
                <td className="py-3">
                  <span
                    className={`text-xs font-medium px-2 py-0.5 rounded-full capitalize ${
                      invoice.status === "paid"
                        ? "bg-accent text-primary"
                        : "bg-muted text-muted-foreground"
                    }`}
                  >
                    {invoice.status}
                  </span>
                </td>
                <td className="py-3 text-right">
                  {invoice.invoice_url && (
                    <a
                      href={invoice.invoice_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-primary hover:underline text-xs"
                    >
                      View
                    </a>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {(hasPrev || hasNext) && (
        <div className="flex items-center justify-between pt-2">
          <button
            onClick={() => setOffset((prev) => Math.max(0, prev - limit))}
            disabled={!hasPrev}
            className="text-sm text-muted-foreground hover:text-foreground transition-colors disabled:opacity-50"
          >
            Previous
          </button>
          <span className="text-xs text-muted-foreground">
            {offset + 1}&ndash;{Math.min(offset + limit, total)} of {total}
          </span>
          <button
            onClick={() => setOffset((prev) => prev + limit)}
            disabled={!hasNext}
            className="text-sm text-muted-foreground hover:text-foreground transition-colors disabled:opacity-50"
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
}
