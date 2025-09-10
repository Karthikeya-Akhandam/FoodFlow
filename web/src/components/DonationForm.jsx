import React, { useState } from "react";

export default function DonationForm({ token, onCreated }) {
  const [title, setTitle] = useState("");
  const [qty, setQty] = useState(1);
  const [purpose, setPurpose] = useState("");
  const API_BASE = import.meta.env.VITE_API_URL || "http://localhost:8080";

  async function submit(e) {
    e.preventDefault();
    const res = await fetch(`${API_BASE}/donations`, {
      method: "POST",
      headers: { "Authorization": token, "Content-Type": "application/json" },
      body: JSON.stringify({ title, quantity: qty, purpose })
    });
    const j = await res.json();
    if (res.ok) {
      setTitle(""); setQty(1); setPurpose("");
      onCreated();
    } else {
      alert(j.error || "error");
    }
  }

  return (
    <form onSubmit={submit} className="p-4 bg-white rounded shadow">
      <h3 className="font-semibold mb-2">Create Donation (collaborator)</h3>
      <input value={title} onChange={e=>setTitle(e.target.value)} placeholder="Title" className="w-full p-2 border mb-2 rounded"/>
      <input type="number" value={qty} onChange={e=>setQty(parseInt(e.target.value) || 1)} className="w-1/4 p-2 border mb-2 rounded"/>
      <input value={purpose} onChange={e=>setPurpose(e.target.value)} placeholder="Purpose (optional)" className="w-full p-2 border mb-2 rounded"/>
      <button className="px-4 py-2 bg-green-600 text-white rounded">Create</button>
    </form>
  );
}