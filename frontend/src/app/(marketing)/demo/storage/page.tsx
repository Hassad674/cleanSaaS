import type { Metadata } from "next";
import { StorageDemo } from "./storage-demo";

export const metadata: Metadata = {
  title: "File Storage Demo — CleanSaaS",
  description:
    "Interactive file storage demo with drag-and-drop uploads, file management, and storage usage tracking. No account required.",
};

export default function StorageDemoPage() {
  return <StorageDemo />;
}
