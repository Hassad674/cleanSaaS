import type { Metadata } from "next";
import { SettingsProfile } from "@/features/user/components/settings-profile";
import { SettingsPassword } from "@/features/user/components/settings-password";
import { SettingsDanger } from "@/features/user/components/settings-danger";

export const metadata: Metadata = {
  title: "Settings",
  robots: { index: false, follow: false },
};

export default function SettingsPage() {
  return (
    <div className="max-w-2xl space-y-6">
      <h1 className="text-2xl font-bold text-foreground">Settings</h1>
      <SettingsProfile />
      <SettingsPassword />
      <SettingsDanger />
    </div>
  );
}
