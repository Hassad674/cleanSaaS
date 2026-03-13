import type { Metadata } from "next";
import { TeamsDemo } from "./teams-demo";

export const metadata: Metadata = {
  title: "Teams Demo — CleanSaaS",
  description:
    "Interactive teams and organizations demo with member management, role assignment, and invitations. No account required.",
};

export default function TeamsDemoPage() {
  return <TeamsDemo />;
}
