import type { Metadata } from "next";
import { BlogDemo } from "./blog-demo";

export const metadata: Metadata = {
  title: "Blog Demo — CleanSaaS",
  description:
    "Interactive blog CMS demo with tag filtering, search, and full post view. No account required.",
};

export default function BlogDemoPage() {
  return <BlogDemo />;
}
