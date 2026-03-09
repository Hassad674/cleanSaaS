"use client";

import { useState, useEffect } from "react";

const API = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081";

interface TestItem {
  id: number;
  name: string;
  created_at: string;
}

export default function Home() {
  const [items, setItems] = useState<TestItem[]>([]);
  const [name, setName] = useState("");
  const [status, setStatus] = useState("loading...");

  async function fetchItems() {
    try {
      const res = await fetch(`${API}/test`);
      const data = await res.json();
      setItems(data || []);
      setStatus("connected");
    } catch {
      setStatus("backend unreachable");
    }
  }

  async function addItem() {
    if (!name.trim()) return;
    await fetch(`${API}/test`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name }),
    });
    setName("");
    fetchItems();
  }

  useEffect(() => {
    fetchItems();
  }, []);

  return (
    <div className="min-h-screen bg-zinc-950 text-white p-8">
      <h1 className="text-2xl font-bold mb-2">CleanSaaS — Test</h1>
      <p className="mb-6 text-zinc-400">
        Backend: <span className={status === "connected" ? "text-green-400" : "text-red-400"}>{status}</span>
      </p>

      <div className="flex gap-2 mb-6">
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && addItem()}
          placeholder="Enter a name..."
          className="bg-zinc-800 border border-zinc-700 rounded px-3 py-2 flex-1"
        />
        <button
          onClick={addItem}
          className="bg-white text-black px-4 py-2 rounded font-medium hover:bg-zinc-200"
        >
          Add
        </button>
      </div>

      <ul className="space-y-2">
        {items.map((item) => (
          <li key={item.id} className="bg-zinc-900 border border-zinc-800 rounded p-3 flex justify-between">
            <span>{item.name}</span>
            <span className="text-zinc-500 text-sm">{item.created_at}</span>
          </li>
        ))}
      </ul>
    </div>
  );
}
