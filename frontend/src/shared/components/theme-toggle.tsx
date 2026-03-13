"use client";

import { useTheme } from "@/shared/components/theme-provider";

type Theme = "light" | "dark" | "system";

const CYCLE_ORDER: Theme[] = ["light", "dark", "system"];

const LABELS: Record<Theme, string> = {
  light: "Light mode",
  dark: "Dark mode",
  system: "System theme",
};

function SunIcon({ className }: { className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth={2}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      aria-hidden="true"
    >
      <circle cx={12} cy={12} r={5} />
      <line x1={12} y1={1} x2={12} y2={3} />
      <line x1={12} y1={21} x2={12} y2={23} />
      <line x1={4.22} y1={4.22} x2={5.64} y2={5.64} />
      <line x1={18.36} y1={18.36} x2={19.78} y2={19.78} />
      <line x1={1} y1={12} x2={3} y2={12} />
      <line x1={21} y1={12} x2={23} y2={12} />
      <line x1={4.22} y1={19.78} x2={5.64} y2={18.36} />
      <line x1={18.36} y1={5.64} x2={19.78} y2={4.22} />
    </svg>
  );
}

function MoonIcon({ className }: { className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth={2}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      aria-hidden="true"
    >
      <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
    </svg>
  );
}

function MonitorIcon({ className }: { className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth={2}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      aria-hidden="true"
    >
      <rect x={2} y={3} width={20} height={14} rx={2} ry={2} />
      <line x1={8} y1={21} x2={16} y2={21} />
      <line x1={12} y1={17} x2={12} y2={21} />
    </svg>
  );
}

const ICONS: Record<Theme, typeof SunIcon> = {
  light: SunIcon,
  dark: MoonIcon,
  system: MonitorIcon,
};

export function ThemeToggle() {
  const { theme, setTheme } = useTheme();

  function handleClick() {
    const currentIndex = CYCLE_ORDER.indexOf(theme);
    const nextIndex = (currentIndex + 1) % CYCLE_ORDER.length;
    setTheme(CYCLE_ORDER[nextIndex]);
  }

  const Icon = ICONS[theme];
  const label = LABELS[theme];

  return (
    <button
      type="button"
      onClick={handleClick}
      className="inline-flex items-center justify-center rounded-lg p-2 text-muted-foreground hover:text-foreground hover:bg-muted transition-colors"
      aria-label={label}
      title={label}
    >
      <Icon className="size-4" />
    </button>
  );
}
