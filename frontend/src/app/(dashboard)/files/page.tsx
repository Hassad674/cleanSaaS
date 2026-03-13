import type { Metadata } from "next";
import { StoragePage } from "@/features/storage/components/storage-page";

export const metadata: Metadata = {
  title: "Files",
  robots: { index: false, follow: false },
};

export default function FilesPage() {
  return <StoragePage />;
}
