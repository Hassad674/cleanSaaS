import type { Metadata } from "next";
import { TeamList } from "@/features/team/components/team-list";

export const metadata: Metadata = {
  title: "Teams",
  robots: { index: false, follow: false },
};

export default function TeamsPage() {
  return <TeamList />;
}
