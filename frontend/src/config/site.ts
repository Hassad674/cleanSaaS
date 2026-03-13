export const siteConfig = {
  name: "CleanSaaS",
  description: "Open-source SaaS boilerplate for modern applications",
  url: process.env.NEXT_PUBLIC_APP_URL || "http://localhost:3000",
  nav: {
    marketing: [
      { label: "Features", href: "/#features" },
      { label: "Pricing", href: "/pricing" },
      { label: "Blog", href: "/blog" },
      { label: "Docs", href: "/docs" },
    ],
    dashboard: [
      { label: "Dashboard", href: "/dashboard", icon: "home" },
      { label: "AI Chat", href: "/ai", icon: "bot" },
      { label: "Files", href: "/files", icon: "folder" },
      { label: "Notifications", href: "/notifications", icon: "bell" },
      { label: "Settings", href: "/settings", icon: "settings" },
    ],
  },
} as const;
