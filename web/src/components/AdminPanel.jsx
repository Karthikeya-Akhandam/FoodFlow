import React, { useState, useEffect } from "react";

export default function AdminPanel({ token, onUpdate }) {
  const [organizations, setOrganizations] = useState([]);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  const [assignProxyForm, setAssignProxyForm] = useState({
    orgId: "",
    proxyOrgId: "",
    extraCredits: 0
  });
  const API_BASE = import.meta.env.VITE_API_URL || "http://localhost:8080";

  const api = (path, opts={}) => fetch(`${API_BASE}${path}`, {
    ...opts,
    headers: {...(opts.headers||{}), "Authorization": token, "Content-Type":"application/json"}
  });

  async function loadOrganizations() {
    try {
      // Since we don't have a direct endpoint for organizations, we'll create a mock list
      // In a real app, you'd have an admin endpoint to get all organizations
      setOrganizations([
        { id: 1, name: "Sample Organization 1", credits: 50, peopleFed: 25 },
        { id: 2, name: "Sample Organization 2", credits: 30, peopleFed: 15 },
        { id: 3, name: "Sample Organization 3", credits: 75, peopleFed: 40 }
      ]);
    } catch (err) {
      console.error("Failed to load organizations:", err);
    }
  }

  async function resetCredits() {
    setLoading(true);
    setMessage("");
    try {
      const res = await api("/admin/reset-credits", { method: "POST" });
      const data = await res.json();
      if (res.ok) {
        setMessage("Credits reset successfully!");
        onUpdate();
      } else {
        setMessage("Failed to reset credits: " + (data.error || "Unknown error"));
      }
    } catch (err) {
      setMessage("Network error: " + err.message);
    } finally {
      setLoading(false);
    }
  }

  async function assignProxy() {
    if (!assignProxyForm.orgId || !assignProxyForm.proxyOrgId) {
      setMessage("Please select both organization and proxy organization");
      return;
    }

    setLoading(true);
    setMessage("");
    try {
      const res = await api("/admin/assign-proxy", {
        method: "POST",
        body: JSON.stringify({
          org_id: parseInt(assignProxyForm.orgId),
          proxy_org_id: parseInt(assignProxyForm.proxyOrgId),
          extra_credits: parseInt(assignProxyForm.extraCredits) || 0
        })
      });
      const data = await res.json();
      if (res.ok) {
        setMessage("Proxy assignment successful!");
        setAssignProxyForm({ orgId: "", proxyOrgId: "", extraCredits: 0 });
        onUpdate();
      } else {
        setMessage("Failed to assign proxy: " + (data.error || "Unknown error"));
      }
    } catch (err) {
      setMessage("Network error: " + err.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadOrganizations();
  }, []);

  return (
    <div className="space-y-8">
      <div className="bg-white rounded-lg shadow">
        <div className="px-6 py-4 border-b border-gray-200">
          <h3 className="text-lg font-medium text-gray-900">Admin Controls</h3>
          <p className="text-sm text-gray-600">Manage organizations and system settings</p>
        </div>
        <div className="p-6 space-y-6">
          {message && (
            <div className={`p-4 rounded-lg ${
              message.includes("success") || message.includes("Success")
                ? "bg-green-50 border border-green-200 text-green-700"
                : "bg-red-50 border border-red-200 text-red-700"
            }`}>
              {message}
            </div>
          )}

          {/* Reset Credits */}
          <div className="border border-gray-200 rounded-lg p-6">
            <h4 className="text-md font-medium text-gray-900 mb-4">🔄 Reset Monthly Credits</h4>
            <p className="text-sm text-gray-600 mb-4">
              Reset credits for all organizations based on their "people fed last month" data.
            </p>
            <button
              onClick={resetCredits}
              disabled={loading}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {loading ? "Processing..." : "Reset All Credits"}
            </button>
          </div>

          {/* Assign Proxy */}
          <div className="border border-gray-200 rounded-lg p-6">
            <h4 className="text-md font-medium text-gray-900 mb-4">🔗 Assign Proxy Organization</h4>
            <p className="text-sm text-gray-600 mb-4">
              Assign a proxy organization to handle donations for remote organizations.
            </p>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Organization</label>
                <select
                  value={assignProxyForm.orgId}
                  onChange={(e) => setAssignProxyForm({...assignProxyForm, orgId: e.target.value})}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  <option value="">Select Organization</option>
                  {organizations.map(org => (
                    <option key={org.id} value={org.id}>
                      {org.name} (Credits: {org.credits})
                    </option>
                  ))}
                </select>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Proxy Organization</label>
                <select
                  value={assignProxyForm.proxyOrgId}
                  onChange={(e) => setAssignProxyForm({...assignProxyForm, proxyOrgId: e.target.value})}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  <option value="">Select Proxy Organization</option>
                  {organizations.map(org => (
                    <option key={org.id} value={org.id}>
                      {org.name} (Credits: {org.credits})
                    </option>
                  ))}
                </select>
              </div>
            </div>
            
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 mb-1">Extra Credits for Proxy (Optional)</label>
              <input
                type="number"
                value={assignProxyForm.extraCredits}
                onChange={(e) => setAssignProxyForm({...assignProxyForm, extraCredits: e.target.value})}
                placeholder="0"
                min="0"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            
            <button
              onClick={assignProxy}
              disabled={loading}
              className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {loading ? "Assigning..." : "Assign Proxy"}
            </button>
          </div>
        </div>
      </div>

      {/* Organizations Overview */}
      <div className="bg-white rounded-lg shadow">
        <div className="px-6 py-4 border-b border-gray-200">
          <h3 className="text-lg font-medium text-gray-900">Organizations Overview</h3>
        </div>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Organization
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Credits
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  People Fed Last Month
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Status
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {organizations.map((org) => (
                <tr key={org.id}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                    {org.name}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {org.credits}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {org.peopleFed}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                      org.credits > 50 ? 'bg-green-100 text-green-800' :
                      org.credits > 20 ? 'bg-yellow-100 text-yellow-800' :
                      'bg-red-100 text-red-800'
                    }`}>
                      {org.credits > 50 ? 'Active' : org.credits > 20 ? 'Moderate' : 'Low Credits'}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
