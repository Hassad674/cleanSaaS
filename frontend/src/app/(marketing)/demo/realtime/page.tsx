import type { Metadata } from "next";
import { RealtimeDemo } from "./realtime-demo";

export const metadata: Metadata = {
  title: "Real-time Demo — CleanSaaS",
  description:
    "Interactive WebSocket real-time demo with live event feed, connection metrics, and architecture overview. No account required.",
};

export default function RealtimeDemoPage() {
  return <RealtimeDemo />;
}
