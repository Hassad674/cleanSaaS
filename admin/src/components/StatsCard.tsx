import type { ReactNode } from "react";

type StatsCardProps = {
  icon: ReactNode;
  label: string;
  value: string | number;
};

export default function StatsCard({ icon, label, value }: StatsCardProps) {
  return (
    <div className="rounded-xl border border-border bg-card p-6 shadow-sm">
      <div className="flex items-center gap-4">
        <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
          {icon}
        </div>

        <div>
          <p className="text-sm text-muted-foreground">{label}</p>
          <p className="text-2xl font-semibold text-card-foreground">{value}</p>
        </div>
      </div>
    </div>
  );
}
